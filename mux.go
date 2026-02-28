// Package xmux provides a type-safe, dependency injection friendly HTTP router
// with support for multiple web frameworks through adapters.
package xmux

import (
	"context"
	"sync"
)

// ============================================================================
// Core Types
// ============================================================================

// Bind is a function type that binds an interface to its concrete implementation.
// It's used for dependency injection, allowing different implementations to be
// provided at runtime (e.g., mock implementations for testing).
//
// Example:
//
//	bind := func(ptr any) error {
//	    switch p := ptr.(type) {
//	    case *IUser:
//	        *p = &UserService{}  // Inject concrete implementation
//	    }
//	    return nil
//	}
type Bind = func(ptr any) (err error)

// Handler is the core interface that all route handlers must implement.
// It provides methods for invoking the handler and retrieving type information
// about its parameters and response.
type Handler interface {
	// Invoke executes the handler with the given context and bind function.
	// The bind function is used to parse request data into the handler's parameter type.
	Invoke(ctx context.Context, bind Bind) (any, error)

	// Params returns a zero value of the handler's parameter type.
	// This is useful for runtime type reflection and documentation generation.
	Params() any

	// Response returns a zero value of the handler's response type.
	// This is useful for runtime type reflection and documentation generation.
	Response() any
}

// Router defines the interface that must be implemented by framework adapters.
// It allows registering routes with their handlers and optional configuration.
type Router interface {
	// Register adds a new route with the given HTTP method, path pattern, handler,
	// and optional configuration options.
	Register(method string, path string, api Handler, options ...map[string]string)
}

// Binder is the interface that wraps the Bind method.
// It represents a routable group that can resolve its dependencies and register
// its routes with a router.
type Binder interface {
	// Bind resolves the handler dependencies through the bind function and
	// registers all routes with the provided router.
	Bind(router Router, bind Bind) (err error)
}

// ============================================================================
// Handler Implementation
// ============================================================================

// HandlerFunc is a generic adapter that allows ordinary functions to be used as
// Handler implementations. The generic types Params and Response define the
// handler's input and output types, providing compile-time type safety.
type HandlerFunc[Params any, Response any] func(ctx context.Context, bind Bind) (Response, error)

// Invoke implements Handler.Invoke by calling the underlying function.
func (h HandlerFunc[Params, Response]) Invoke(ctx context.Context, bind Bind) (any, error) {
	return h(ctx, bind)
}

// Params implements Handler.Params by returning a zero value of the parameter type.
func (h HandlerFunc[Params, Response]) Params() any {
	var zero Params
	return zero
}

// Response implements Handler.Response by returning a zero value of the response type.
func (h HandlerFunc[Params, Response]) Response() any {
	var zero Response
	return zero
}

// ============================================================================
// Route Registration
// ============================================================================

// Register is a generic helper function that simplifies route registration.
// It creates a HandlerFunc from the provided function, automatically handling
// parameter binding and type conversion.
//
// Type Parameters:
//   - Params: The request parameter type (will be automatically bound)
//   - Response: The response type (must match the function's return type)
//
// Example:
//
//	type CreateUserParams struct {
//	    Name  string `json:"name"`
//	    Email string `json:"email"`
//	}
//
//	Register(router, http.MethodPost, "/users",
//	    func(ctx context.Context, params *CreateUserParams) (*User, error) {
//	        return userService.Create(ctx, params)
//	    })
func Register[Params any, Response any](
	router Router,
	method string,
	path string,
	fn func(ctx context.Context, params *Params) (Response, error),
	options ...map[string]string,
) {
	router.Register(method, path, HandlerFunc[Params, Response](func(ctx context.Context, bind Bind) (resp Response, err error) {
		var params Params
		if err = bind(&params); err != nil {
			return
		}
		return fn(ctx, &params)
	}), options...)
}

// ============================================================================
// Route Groups
// ============================================================================

// optionsRouter is a decorator that adds group-level options to all routes
// registered through it. It implements the Router interface.
type optionsRouter struct {
	router  Router              // The underlying router being decorated
	options []map[string]string // Options to be applied to all routes
}

// Register implements Router.Register by appending group options to the route-specific options.
func (o optionsRouter) Register(method string, path string, api Handler, options ...map[string]string) {
	o.router.Register(method, path, api, append(o.options, options...)...)
}

// RouterGroup represents a group of routes that share the same handler interface
// and can be registered together. It implements the Binder interface.
type RouterGroup[Handler any] struct {
	// register is the function that actually registers all routes in this group
	// with a concrete router and handler implementation.
	register func(router Router, handler Handler)
}

// Bind implements Binder.Bind by resolving the handler implementation through the
// provided bind function and then registering all routes with the given router.
func (g RouterGroup[Handler]) Bind(router Router, bind Bind) (err error) {
	var handler Handler
	if err = bind(&handler); err != nil {
		return
	}
	g.register(router, handler)
	return
}

// DefineGroup creates a new route group binder. It captures the group-level options
// and returns a Binder that can be registered with a Groups instance.
// The generic type Handler represents the interface type that this group's routes depend on.
//
// Example:
//
//	userGroup := DefineGroup(func(router Router, user IUser) {
//	    Register(router, http.MethodGet, "/users", user.List)
//	}, map[string]string{"prefix": "/api/v1"})
func DefineGroup[Handler any](
	registerFn func(router Router, handler Handler),
	options ...map[string]string,
) Binder {
	return RouterGroup[Handler]{
		register: func(router Router, handler Handler) {
			registerFn(&optionsRouter{
				router:  router,
				options: options,
			}, handler)
		},
	}
}

// ============================================================================
// Groups Collection
// ============================================================================

// Groups is a thread-safe collection of route groups that can be bound to a router
// as a single unit. It implements a registry pattern for route groups.
type Groups struct {
	mu     sync.Mutex
	groups []Binder // Registered route groups
}

// NewGroups creates a new empty Groups collection.
func NewGroups() *Groups {
	return &Groups{
		groups: make([]Binder, 0),
	}
}

// Register adds one or more route groups to the collection in a thread-safe manner.
// This can be called at any time, even after the server has started.
func (g *Groups) Register(groups ...Binder) *Groups {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.groups = append(g.groups, groups...)
	return g
}

// Bind applies all registered groups to the provided router, resolving their
// handler dependencies through the bind function. Returns the first error encountered,
// if any.
func (g *Groups) Bind(router Router, bind Bind) (err error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	for _, group := range g.groups {
		if err = group.Bind(router, bind); err != nil {
			return
		}
	}
	return
}
