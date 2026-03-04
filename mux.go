// Package xmux provides a type-safe, dependency injection friendly HTTP router
// with support for multiple web frameworks through adapters.
//
// Core design principles:
// 1. Framework agnostic - Business logic has no framework dependencies
// 2. Type safety - Compile-time request/response type checking with generics
// 3. Dependency injection - Flexible DI through simple bind functions
// 4. Adapter pattern - Support any HTTP framework through adapters
//
// Typical workflow:
// 1. Define pure business logic functions (func(ctx, *Params) (Response, error))
// 2. Create framework adapter (implement Router interface)
// 3. Register routes with Register
// 4. Organize route groups with ServiceGroup
// 5. Manage multiple groups with Groups
package xmux

import (
	"context"
	"reflect"
	"runtime"
	"sync"
)

// Api represents a type-safe handler interface.
// It encapsulates business logic with type information for request/response.
// The framework adapter calls Invoke to execute business logic,
// providing a bind function to populate request parameters.
type Api interface {
	// Invoke executes the business logic.
	// bind: function to populate the params struct from HTTP request
	// Returns: response data (typed) and error
	Invoke(ctx context.Context, bind func(params any) error) (any, error)

	// Params returns zero value of the parameter type.
	// Used for type introspection and documentation generation.
	Params() any

	// Response returns zero value of the response type.
	// Used for type introspection and documentation generation.
	Response() any

	// Function returns the underlying function.
	// Useful for reflection and advanced use cases.
	Function() any

	Name() string

	Service() (any, reflect.Type)
}

// Controller represents a framework-specific handler controller.
// Implemented by framework adapters to handle requests.
// The service parameter is the injected business logic instance.
type Controller interface {
	// Handle processes an HTTP request using the provided service and api.
	// method: HTTP method (GET, POST, PUT, DELETE, etc.)
	// path: URL path pattern
	// service: the business logic service instance (injected via Bind)
	// api: the type-safe handler to invoke
	// options: additional route configuration (middleware, metadata, etc.)
	Handle(method string, path string, api Api, options ...map[string]string)
}

// Router represents a route registrar.
// Implemented by framework adapters to register routes.
// This is the primary interface for connecting business logic to frameworks.
type Router interface {
	// Register adds a new route to the router.
	// method: HTTP method (GET, POST, PUT, DELETE, etc.)
	// path: URL path pattern (e.g., "/users/:id")
	// api: the type-safe handler to invoke
	// options: additional route configuration (middleware, metadata, etc.)
	Register(method string, path string, api Api, options ...map[string]string)
}

// Binder represents a bindable entity that can inject dependencies.
// Used for route groups that need dependency injection.
type Binder interface {
	// Bind injects dependencies and registers handlers.
	// handler: the controller that will handle requests
	// bind: function to inject service dependencies
	// Returns error if dependency injection fails
	Bind(handler Controller, bind func(service any) error) (err error)
}

// function is the internal implementation of Api for business logic functions.
// It wraps a function with signature func(context.Context, *Params) (Response, error)
// and provides type-safe invocation with automatic parameter binding.
type function[Params any, Response any] func(context.Context, *Params) (Response, error)

func (h function[Params, Response]) Service() (any, reflect.Type) {
	return nil, nil
}

// Invoke executes the business logic function.
// It first calls unmarshal to populate the params struct from the HTTP request,
// then calls the underlying function with the populated params.
func (h function[Params, Response]) Invoke(ctx context.Context, unmarshal func(params any) error) (ret any, err error) {
	var params Params
	if err = unmarshal(&params); err != nil {
		return
	}
	return h(ctx, &params)
}

// Function returns the underlying function instance.
// Useful for reflection and type introspection.
func (h function[Params, Response]) Function() any {
	return h
}

// Params returns a zero value of the Params type.
// Used for type introspection, documentation generation, and validation.
func (h function[Params, Response]) Params() any {
	var zero Params
	return zero
}

// Response returns a zero value of the Response type.
// Used for type introspection, documentation generation, and validation.
func (h function[Params, Response]) Response() any {
	var zero Response
	return zero
}

func (h function[Params, Response]) Name() string {
	return runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name()
}

// Register registers a business logic function as a route handler.
// This is the primary function for connecting business logic to a framework router.
//
// Type parameters:
//   - Params: the request parameter struct type
//   - Response: the response data struct type
//
// Parameters:
//   - router: the framework router (implements Router interface)
//   - method: HTTP method (GET, POST, PUT, DELETE, etc.)
//   - path: URL path pattern (e.g., "/users", "/users/:id")
//   - fn: the business logic function to execute
//   - options: optional route configuration (middleware, metadata, etc.)
//
// Example:
//
//	type CreateUserReq struct { Name string `json:"name"` }
//	type UserResp struct { ID string `json:"id"` }
//	func CreateUser(ctx context.Context, req *CreateUserReq) (*UserResp, error) { ... }
//	xmux.Register(router, http.MethodPost, "/users", CreateUser)
func Register[Params any, Response any](
	router Router,
	method string,
	path string,
	fn func(ctx context.Context, params *Params) (Response, error),
	options ...map[string]string,
) {
	router.Register(method, path, function[Params, Response](fn), options...)
}

// MergeOptions merges multiple option maps into a single map.
// Useful for combining route-level, group-level, and global options.
//
// Parameters:
//   - options: slice of option maps to merge
//   - desc: if true, merge in descending order (later options override earlier)
//     if false, merge in ascending order (earlier options override later)
//
// Returns:
//   - merged option map
//
// Example:
//
//	globalOpts := map[string]string{"timeout": "30s"}
//	groupOpts := map[string]string{"prefix": "/api"}
//	routeOpts := map[string]string{"middleware": "auth"}
//	merged := xmux.MergeOptions([]map[string]string{globalOpts, groupOpts, routeOpts}, true)
//	// Result: {"timeout": "30s", "prefix": "/api", "middleware": "auth"}
func MergeOptions(options []map[string]string, desc bool) map[string]string {
	opt := make(map[string]string)
	for i := 0; i < len(options); i++ {
		o := options[i]
		if desc {
			o = options[len(options)-1-i]
		}
		for k, v := range o {
			opt[k] = v
		}
	}
	return opt
}

// serviceGroup represents a group of routes that share a common service.
// It enables dependency injection for a specific service type and allows
// registering multiple handlers that use the same service instance.
//
// Type parameters:
//   - Service: the business logic service type (e.g., UserService, OrderService)
type serviceGroup[Service any] struct {
	// register is the function that defines routes for this service
	register func(router Router, handler Service)

	// options are route-level options that apply to all routes in this group
	options []map[string]string
}

// Bind injects the service dependency and registers all routes in the group.
// This implements the Binder interface.
//
// Parameters:
//   - controller: the framework controller that handles requests
//   - bind: function to inject the service dependency
//
// Returns:
//   - error if dependency injection or route registration fails
func (g serviceGroup[Service]) Bind(controller Controller, bind func(any) error) (err error) {
	var s Service
	if err = bind(&s); err != nil {
		return
	}
	g.register(registerFunc(func(method string, path string, api Api, options ...map[string]string) {
		controller.Handle(method, path, serviceApi[Service]{
			Api:  api,
			impl: s,
		}, append(g.options, options...)...)
	}), s)
	return
}

type serviceApi[Service any] struct {
	Api
	impl Service
}

func (api serviceApi[Service]) Service() (any, reflect.Type) {
	return api.impl, reflect.TypeOf((*Service)(nil)).Elem()
}

// registerFunc is a function type that implements the Router interface.
// It allows converting a function into a Router for flexible route registration.
type registerFunc func(method string, path string, api Api, options ...map[string]string)

// Register implements the Router interface for registerFunc.
func (fn registerFunc) Register(method string, path string, api Api, options ...map[string]string) {
	fn(method, path, api, options...)
}

// ServiceGroup creates a new route group for a specific service type.
// This is the primary way to organize related routes with shared dependencies.
//
// Type parameters:
//   - Service: the business logic service type
//
// Parameters:
//   - fn: function that defines routes for the service
//     It receives a Router and the service instance
//   - options: optional group-level options (applied to all routes in the group)
//
// Returns:
//   - Binder that can be used with Groups to register multiple groups
//
// Example:
//
//	userGroup := xmux.ServiceGroup(func(router xmux.Router, svc UserService) {
//	    xmux.Register(router, http.MethodGet, "/users", svc.GetUser)
//	    xmux.Register(router, http.MethodPost, "/users", svc.CreateUser)
//	    xmux.Register(router, http.MethodDelete, "/users/:id", svc.DeleteUser)
//	}, map[string]string{"prefix": "/api/v1"})
func ServiceGroup[Service any](fn func(router Router, handler Service), options ...map[string]string) Binder {
	return serviceGroup[Service]{
		options:  options,
		register: fn,
	}
}

// Groups represents a collection of route groups.
// It enables registering and binding multiple service groups together.
// This is useful for organizing large applications with multiple services.
type Groups interface {
	Binder
	// Register adds more groups to the collection
	Register(groups ...Binder) Groups
}

// groups is the internal implementation of Groups.
// It maintains a thread-safe slice of Binder instances.
type groups struct {
	mu     sync.Mutex
	groups []Binder
}

// NewGroups creates a new Groups instance with the provided initial groups.
//
// Parameters:
//   - gs: initial groups to register
//
// Returns:
//   - Groups instance for managing multiple route groups
//
// Example:
//
//	groups := xmux.NewGroups(userGroup, orderGroup, productGroup)
//	groups.Bind(router, func(ptr any) error {
//	    switch p := ptr.(type) {
//	    case *UserService:
//	        *p = NewUserService(db)
//	    case *OrderService:
//	        *p = NewOrderService(db)
//	    }
//	    return nil
//	})
func NewGroups(gs ...Binder) Groups {
	return &groups{groups: append(make([]Binder, 0, len(gs)), gs...)}
}

// Register adds more groups to the collection.
// This method is thread-safe and can be called concurrently.
//
// Parameters:
//   - groups: additional groups to register
//
// Returns:
//   - self for method chaining
func (g *groups) Register(groups ...Binder) Groups {
	g.mu.Lock()
	g.groups = append(g.groups, groups...)
	g.mu.Unlock()
	return g
}

// Bind injects dependencies and binds all registered groups.
// This method is thread-safe and can be called concurrently.
//
// Parameters:
//   - controller: the framework controller that handles requests
//   - bind: function to inject service dependencies
//     Typically uses a type switch to inject multiple services
//
// Returns:
//   - error if any group fails to bind
//
// Example:
//
//	groups := xmux.NewGroups(userGroup, orderGroup)
//	err := groups.Bind(router, func(ptr any) error {
//	    switch p := ptr.(type) {
//	    case *UserService:
//	        *p = NewUserService(db)
//	    case *OrderService:
//	        *p = NewOrderService(db)
//	    }
//	    return nil
//	})
func (g *groups) Bind(controller Controller, bind func(service any) error) (err error) {
	g.mu.Lock()
	gs := append(make([]Binder, 0, len(g.groups)), g.groups...)
	g.mu.Unlock()
	for _, group := range gs {
		if err = group.Bind(controller, bind); err != nil {
			return
		}
	}
	return
}
