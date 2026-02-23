# xmux

[![Go Version](https://img.shields.io/badge/go-1.18%2B-blue)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

**Framework-agnostic business logic declaration for medium to large Go projects.** Write business logic once, deploy anywhere.

xmux enables you to define type-safe, dependency-injected business logic completely independent of web frameworks. Focus on your core domain while maintaining the flexibility to adapt to any HTTP framework or generate comprehensive API documentation automatically.

## Ideal for Medium to Large Projects

xmux is designed specifically for projects that need:
- **Framework-agnostic business declarations** - Define APIs once, use with any web framework
- **Unified API style across teams** - Enforce consistent patterns through type-safe generics
- **Automated documentation generation** - Leverage type information for OpenAPI/Swagger generation
- **Clean architecture** - Separate business logic from framework dependencies
- **Future-proof APIs** - Switch frameworks without rewriting business logic

## Why xmux?

### Problem
Traditional Go web development ties your business logic to specific framework APIs. Switching frameworks or adapting to different API requirements means rewriting substantial portions of your application.

### Solution  
xmux decouples your business logic from web framework dependencies through:
- **Type-safe handlers** with compile-time validation
- **Dependency injection** for testable, modular code  
- **Framework-agnostic design** that works with any HTTP router
- **Customizable adapters** tailored to your specific API requirements

## Core Features

- **Framework-Agnostic Business Declarations**: Define business logic completely independent of web frameworks
- **Unified API Style Enforcement**: Enforce consistent patterns across teams using Go generics
- **Automated Documentation Ready**: Type information enables automatic OpenAPI/Swagger generation
- **Zero Dependencies**: xmux itself has no external dependencies
- **Business Logic Focus**: Write once, adapt to any framework
- **Type Safety**: Compile-time checking for request/response types
- **Dependency Injection**: Flexible DI through simple bind functions
- **Framework Agnostic**: Works with net/http, Gin, Echo, Fiber, Chi, Gorilla/mux, etc.
- **Customizable Adapters**: Create adapters optimized for your API conventions
- **Route Groups**: Organize related routes with shared dependencies
- **Thread-Safe**: Concurrent registration and binding support

## Installation

```bash
go get github.com/Just-maple/xmux
```

## Quick Start

### 1. Define Your Business Logic

```go
package business

import "context"

// Request type (automatically bound from HTTP request)
type CreateUserRequest struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}

// Response type (automatically serialized to HTTP response)
type UserResponse struct {
    ID    string `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

// Pure business logic function
func CreateUser(ctx context.Context, req *CreateUserRequest) (*UserResponse, error) {
    // Your business logic here
    // No framework dependencies!
    return &UserResponse{
        ID:    "user-123",
        Name:  req.Name,
        Email: req.Email,
    }, nil
}
```

### 2. Create a Custom Adapter for Your Framework

```go
package adapter

import (
    "encoding/json"
    "net/http"
    
    "github.com/Just-maple/xmux"
    "github.com/your-app/business"
)

// Custom adapter optimized for your API requirements
type CustomRouter struct {
    mux *http.ServeMux
}

func NewCustomRouter() *CustomRouter {
    return &CustomRouter{mux: http.NewServeMux()}
}

func (r *CustomRouter) Register(method, path string, api xmux.Handler, opts ...map[string]string) {
    r.mux.HandleFunc(path, func(w http.ResponseWriter, req *http.Request) {
        // Custom request parsing for your API conventions
        bind := func(ptr any) error {
            // Parse JSON body for POST/PUT requests
            if method == http.MethodPost || method == http.MethodPut {
                return json.NewDecoder(req.Body).Decode(ptr)
            }
            // Parse query params for GET requests
            // Add your custom parsing logic here
            return nil
        }
        
        // Execute business logic
        result, err := api.Invoke(req.Context(), bind)
        if err != nil {
            // Your custom error handling
            w.WriteHeader(http.StatusBadRequest)
            json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
            return
        }
        
        // Your custom success response
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(result)
    })
}
```

### 3. Wire Everything Together

```go
package main

import (
    "net/http"
    
    "github.com/Just-maple/xmux"
    "github.com/your-app/adapter"
    "github.com/your-app/business"
)

func main() {
    // Create your custom adapter
    router := adapter.NewCustomRouter()
    
    // Register business logic with type safety
    xmux.Register(router, http.MethodPost, "/users", business.CreateUser)
    
    // Start server
    http.ListenAndServe(":8080", router)
}
```

## Aspect-Oriented Router Design

xmux's `Router` interface serves as a **unified cross-cutting concern** that enables consistent request/response handling across all endpoints. This aspect-oriented approach allows you to:

### 1. **Unified Request/Response Processing**
```go
// Custom adapter that enforces consistent API patterns
type StandardizedRouter struct {
    baseRouter xmux.Router
    // Cross-cutting concerns configuration
    config     RouterConfig
}

func NewStandardizedRouter(baseRouter xmux.Router, config RouterConfig) *StandardizedRouter {
    return &StandardizedRouter{
        baseRouter: baseRouter,
        config:     config,
    }
}

func (r *StandardizedRouter) Register(method, path string, api xmux.Handler, opts ...map[string]string) {
    // Apply cross-cutting concerns before passing to base router
    wrappedHandler := r.applyCrossCuttingConcerns(api)
    r.baseRouter.Register(method, path, wrappedHandler, opts...)
}

func (r *StandardizedRouter) applyCrossCuttingConcerns(api xmux.Handler) xmux.Handler {
    return xmux.HandlerFunc(func(ctx context.Context, bind xmux.Bind) (any, error) {
        // 1. Unified request validation
        if err := r.validateRequest(ctx); err != nil {
            return nil, err
        }
        
        // 2. Unified authentication/authorization
        if err := r.authenticate(ctx); err != nil {
            return nil, err
        }
        
        // 3. Execute business logic
        result, err := api.Invoke(ctx, bind)
        
        // 4. Unified error handling
        if err != nil {
            return nil, r.normalizeError(err)
        }
        
        // 5. Unified response formatting
        return r.formatResponse(result), nil
    })
}
```

### 2. **Consistent Parameter Binding**
```go
// Standardized bind function for all endpoints
type StandardBindFunc struct {
    // Request parsing strategies
    parsers []RequestParser
}

func (b *StandardBindFunc) Bind(ptr any) error {
    // Apply consistent parsing rules:
    
    // 1. JSON body parsing (for POST/PUT/PATCH)
    if err := b.parseJSONBody(ptr); err != nil {
        return err
    }
    
    // 2. Query parameter binding
    if err := b.parseQueryParams(ptr); err != nil {
        return err
    }
    
    // 3. Path parameter binding  
    if err := b.parsePathParams(ptr); err != nil {
        return err
    }
    
    // 4. Header parsing (e.g., authentication tokens)
    if err := b.parseHeaders(ptr); err != nil {
        return err
    }
    
    // 5. Validation (using struct tags like `validate:"required"`)
    return b.validate(ptr)
}

// Example request struct with unified binding rules
type CreateOrderRequest struct {
    // JSON body field
    ProductID string `json:"product_id" validate:"required,uuid"`
    Quantity  int    `json:"quantity" validate:"required,min=1,max=100"`
    
    // Query parameter
    UserID string `query:"user_id" validate:"required"`
    
    // Path parameter (extracted from route)
    StoreID string `path:"store_id"`
    
    // Header value
    AuthToken string `header:"Authorization"`
}
```

### 3. **Unified Response Formatting**
```go
// Standard response wrapper for all endpoints
type APIResponse[T any] struct {
    Success bool       `json:"success"`
    Data    T          `json:"data,omitempty"`
    Error   *APIError  `json:"error,omitempty"`
    Meta    ResponseMeta `json:"meta,omitempty"`
}

type APIError struct {
    Code    string `json:"code"`
    Message string `json:"message"`
    Details any    `json:"details,omitempty"`
}

type ResponseMeta struct {
    RequestID string    `json:"request_id"`
    Timestamp time.Time `json:"timestamp"`
    Version   string    `json:"version"`
}

// Response formatter that applies to all handlers
type ResponseFormatter struct {
    config FormatConfig
}

func (f *ResponseFormatter) Format(result any, err error) any {
    if err != nil {
        return APIResponse[any]{
            Success: false,
            Error:   f.mapError(err),
            Meta: ResponseMeta{
                RequestID: f.getRequestID(),
                Timestamp: time.Now(),
                Version:   f.config.APIVersion,
            },
        }
    }
    
    return APIResponse[any]{
        Success: true,
        Data:    result,
        Meta: ResponseMeta{
            RequestID: f.getRequestID(),
            Timestamp: time.Now(),
            Version:   f.config.APIVersion,
        },
    }
}
```

### 4. **Business Logic Isolation**
```go
// Business logic completely isolated from HTTP concerns
type OrderService interface {
    CreateOrder(ctx context.Context, req *CreateOrderRequest) (*OrderResponse, error)
    CancelOrder(ctx context.Context, orderID string) error
    GetOrder(ctx context.Context, orderID string) (*OrderResponse, error)
}

// Pure business implementation
type orderServiceImpl struct {
    repo      OrderRepository
    inventory InventoryService
    payment   PaymentService
}

func (s *orderServiceImpl) CreateOrder(ctx context.Context, req *CreateOrderRequest) (*OrderResponse, error) {
    // No HTTP framework dependencies!
    // Pure business logic:
    
    // 1. Check inventory
    available, err := s.inventory.CheckAvailability(ctx, req.ProductID, req.Quantity)
    if err != nil || !available {
        return nil, ErrInsufficientInventory
    }
    
    // 2. Create order in database
    order := &Order{
        ID:        uuid.New().String(),
        ProductID: req.ProductID,
        Quantity:  req.Quantity,
        Status:    OrderStatusPending,
        CreatedAt: time.Now(),
    }
    
    if err := s.repo.Save(ctx, order); err != nil {
        return nil, err
    }
    
    // 3. Process payment
    paymentResult, err := s.payment.Process(ctx, order.ID, order.CalculateTotal())
    if err != nil {
        // Update order status to failed
        order.Status = OrderStatusPaymentFailed
        s.repo.Save(ctx, order)
        return nil, err
    }
    
    // 4. Return result
    return &OrderResponse{
        OrderID:   order.ID,
        Status:    order.Status,
        Total:     order.CalculateTotal(),
        PaymentID: paymentResult.PaymentID,
    }, nil
}
```

### 5. **Cross-Cutting Concern Examples**

```go
// Authentication middleware applied to all routes
type AuthMiddleware struct {
    authService AuthService
}

func (m *AuthMiddleware) Wrap(api xmux.Handler) xmux.Handler {
    return xmux.HandlerFunc(func(ctx context.Context, bind xmux.Bind) (any, error) {
        // Extract and validate token
        token := extractTokenFromContext(ctx)
        user, err := m.authService.ValidateToken(ctx, token)
        if err != nil {
            return nil, ErrUnauthorized
        }
        
        // Add user to context for business logic
        ctx = context.WithValue(ctx, "user", user)
        
        // Proceed to business logic
        return api.Invoke(ctx, bind)
    })
}

// Logging middleware for all requests
type LoggingMiddleware struct {
    logger Logger
}

func (m *LoggingMiddleware) Wrap(api xmux.Handler) xmux.Handler {
    return xmux.HandlerFunc(func(ctx context.Context, bind xmux.Bind) (any, error) {
        start := time.Now()
        requestID := generateRequestID()
        
        m.logger.Info("Request started",
            "request_id", requestID,
            "method", getMethodFromContext(ctx),
            "path", getPathFromContext(ctx))
        
        // Execute business logic
        result, err := api.Invoke(ctx, bind)
        
        duration := time.Since(start)
        status := "success"
        if err != nil {
            status = "error"
        }
        
        m.logger.Info("Request completed",
            "request_id", requestID,
            "duration", duration,
            "status", status,
            "error", err)
        
        return result, err
    })
}

// Rate limiting middleware
type RateLimitMiddleware struct {
    limiter RateLimiter
}

func (m *RateLimitMiddleware) Wrap(api xmux.Handler) xmux.Handler {
    return xmux.HandlerFunc(func(ctx context.Context, bind xmux.Bind) (any, error) {
        clientIP := getClientIPFromContext(ctx)
        
        if !m.limiter.Allow(clientIP) {
            return nil, ErrRateLimitExceeded
        }
        
        return api.Invoke(ctx, bind)
    })
}
```

### 6. **Composable Middleware Chain**

```go
// Build middleware chain for cross-cutting concerns
type MiddlewareChain struct {
    middlewares []Middleware
}

func NewMiddlewareChain() *MiddlewareChain {
    return &MiddlewareChain{
        middlewares: []Middleware{
            &LoggingMiddleware{},
            &AuthMiddleware{},
            &RateLimitMiddleware{},
            &ValidationMiddleware{},
            &MetricsMiddleware{},
        },
    }
}

func (c *MiddlewareChain) Apply(api xmux.Handler) xmux.Handler {
    handler := api
    // Apply in reverse order (last added executes first)
    for i := len(c.middlewares) - 1; i >= 0; i-- {
        handler = c.middlewares[i].Wrap(handler)
    }
    return handler
}

// Usage in adapter
func (r *StandardizedRouter) applyCrossCuttingConcerns(api xmux.Handler) xmux.Handler {
    chain := NewMiddlewareChain()
    return chain.Apply(api)
}
```

This aspect-oriented design enables you to:
- **Enforce consistent API patterns** across all endpoints
- **Isolate business logic** from HTTP framework concerns  
- **Apply cross-cutting concerns** (auth, logging, rate limiting) uniformly
- **Switch frameworks** without touching business logic
- **Generate consistent documentation** from type information
- **Maintain clean separation** between infrastructure and domain logic

## Framework-Agnostic Business Declarations

xmux enables you to define business logic that is completely independent of any web framework. This is particularly valuable for medium to large projects where:

### 1. **Business Logic as First-Class Citizen**
```go
// Define business interfaces without framework dependencies
type UserService interface {
    CreateUser(ctx context.Context, req *CreateUserRequest) (*UserResponse, error)
    GetUser(ctx context.Context, id string) (*UserResponse, error)
    ListUsers(ctx context.Context, filter *UserFilter) ([]*UserResponse, error)
}

// Implement business logic purely
type userServiceImpl struct {
    db *database.DB
}

func (s *userServiceImpl) CreateUser(ctx context.Context, req *CreateUserRequest) (*UserResponse, error) {
    // Pure business logic - no HTTP framework concepts
    user := &User{
        ID:    uuid.New().String(),
        Name:  req.Name,
        Email: req.Email,
    }
    if err := s.db.SaveUser(ctx, user); err != nil {
        return nil, err
    }
    return &UserResponse{
        ID:    user.ID,
        Name:  user.Name,
        Email: user.Email,
    }, nil
}
```

### 2. **Dependency Injection for Testability**
```go
// Define dependencies as interfaces
type UserRepository interface {
    Save(ctx context.Context, user *User) error
    FindByID(ctx context.Context, id string) (*User, error)
    FindAll(ctx context.Context, filter *UserFilter) ([]*User, error)
}

// Inject dependencies at runtime
func NewUserService(repo UserRepository) UserService {
    return &userServiceImpl{repo: repo}
}

// Mock dependencies for testing
func TestUserService_CreateUser(t *testing.T) {
    mockRepo := &MockUserRepository{}
    service := NewUserService(mockRepo)
    
    // Test business logic without HTTP framework
    result, err := service.CreateUser(context.Background(), &CreateUserRequest{
        Name:  "John",
        Email: "john@example.com",
    })
    
    assert.NoError(t, err)
    assert.NotNil(t, result)
}
```

### 3. **Framework Independence**
```go
// Same business logic works with different frameworks
func main() {
    // Choose framework based on requirements
    framework := os.Getenv("HTTP_FRAMEWORK")
    
    var router xmux.Router
    switch framework {
    case "gin":
        router = ginadapter.NewRouter()
    case "echo":
        router = echoadapter.NewRouter()
    case "fiber":
        router = fiberadapter.NewRouter()
    default:
        router = nethttpadapter.NewRouter()
    }
    
    // Register same business logic
    userService := business.NewUserService(db)
    userGroup := xmux.DefineGroup(func(r xmux.Router, svc business.UserService) {
        xmux.Register(r, http.MethodPost, "/users", svc.CreateUser)
        xmux.Register(r, http.MethodGet, "/users/:id", svc.GetUser)
        xmux.Register(r, http.MethodGet, "/users", svc.ListUsers)
    })
    
    // Business logic remains unchanged regardless of framework
    groups := xmux.NewGroups()
    groups.Register(userGroup)
    groups.Bind(router, func(ptr any) error {
        switch p := ptr.(type) {
        case *business.UserService:
            *p = userService
        }
        return nil
    })
}
```

## Unified API Style & Documentation Generation

xmux's type-safe approach using Go generics enables automatic API documentation generation and consistent API style enforcement.

### 1. **Type Information for Documentation**
```go
// Request and response types carry metadata for documentation
type CreateUserRequest struct {
    Name  string `json:"name" validate:"required,min=2,max=100" doc:"User's full name"`
    Email string `json:"email" validate:"required,email" doc:"User's email address"`
    Age   int    `json:"age" validate:"min=18" doc:"User's age in years"`
}

type UserResponse struct {
    ID        string    `json:"id" doc:"Unique user identifier"`
    Name      string    `json:"name" doc:"User's full name"`
    Email     string    `json:"email" doc:"User's email address"`
    CreatedAt time.Time `json:"created_at" doc:"Account creation timestamp"`
}

// Handler type information can be extracted for documentation
func generateOpenAPISpec(router xmux.Router) openapi.Spec {
    spec := openapi.NewSpec()
    
    // Extract type information from registered handlers
    // This enables automatic OpenAPI generation
    return spec
}
```

### 2. **Consistent API Patterns**
```go
// Enforce consistent error responses across all endpoints
type APIError struct {
    Code    string `json:"code"`
    Message string `json:"message"`
    Details any    `json:"details,omitempty"`
}

// Standardized response wrapper
type APIResponse[T any] struct {
    Success bool   `json:"success"`
    Data    T      `json:"data,omitempty"`
    Error   *APIError `json:"error,omitempty"`
    Meta    map[string]any `json:"meta,omitempty"`
}

// Adapter enforces consistent response format
type StandardizedAdapter struct {
    baseRouter xmux.Router
}

func (a *StandardizedAdapter) Register(method, path string, api xmux.Handler, opts ...map[string]string) {
    a.baseRouter.Register(method, path, a.wrapHandler(api), opts...)
}

func (a *StandardizedAdapter) wrapHandler(api xmux.Handler) xmux.Handler {
    return xmux.HandlerFunc(func(ctx context.Context, bind xmux.Bind) (any, error) {
        result, err := api.Invoke(ctx, bind)
        if err != nil {
            // Convert any error to standardized format
            return APIResponse[any]{
                Success: false,
                Error:   a.mapError(err),
            }, nil
        }
        
        return APIResponse[any]{
            Success: true,
            Data:    result,
        }, nil
    })
}
```

### 3. **Automatic Documentation Tools**
```go
// Example: Generate OpenAPI documentation from xmux routes
package docs

import (
    "encoding/json"
    "os"
    
    "github.com/Just-maple/xmux"
    "github.com/getkin/kin-openapi/openapi3"
)

type DocGenerator struct {
    spec *openapi3.T
}

func NewDocGenerator(title, version string) *DocGenerator {
    return &DocGenerator{
        spec: &openapi3.T{
            OpenAPI: "3.0.0",
            Info: &openapi3.Info{
                Title:   title,
                Version: version,
            },
            Paths: openapi3.Paths{},
        },
    }
}

func (g *DocGenerator) AddRoute(router xmux.Router, method, path string, api xmux.Handler) {
    // Extract type information for OpenAPI schema generation
    paramType := api.Params()
    responseType := api.Response()
    
    // Generate OpenAPI operation from type information
    operation := g.generateOperation(paramType, responseType)
    
    // Add to OpenAPI spec
    if g.spec.Paths[path] == nil {
        g.spec.Paths[path] = &openapi3.PathItem{}
    }
    
    switch method {
    case "GET":
        g.spec.Paths[path].Get = operation
    case "POST":
        g.spec.Paths[path].Post = operation
    case "PUT":
        g.spec.Paths[path].Put = operation
    case "DELETE":
        g.spec.Paths[path].Delete = operation
    }
}

func (g *DocGenerator) SaveToFile(filename string) error {
    data, err := json.MarshalIndent(g.spec, "", "  ")
    if err != nil {
        return err
    }
    return os.WriteFile(filename, data, 0644)
}
```

## Adapter Customization Guide

### When to Create Custom Adapters

The examples in `/examples` are **reference implementations only**. In production, you should create adapters tailored to:

1. **Your API Conventions**
   - Request/response formats (JSON, XML, Protobuf, etc.)
   - Error handling patterns
   - Authentication/authorization schemes
   - Rate limiting and quotas

2. **Your Framework Choices**
   - Specific middleware requirements
   - Performance optimizations
   - Monitoring and metrics integration
   - Deployment environment constraints

3. **Your Business Requirements**
   - Domain-specific validation
   - Audit logging
   - Compliance requirements
   - Integration with existing systems

### Key Areas to Customize

#### 1. Request Parsing
```go
bind := func(ptr any) error {
    // Customize based on your API:
    // - JSON body parsing with validation
    // - URL parameter extraction  
    // - Query string parsing
    // - Form data handling
    // - Header-based authentication
    // - Content negotiation
    return yourCustomParser(req, ptr)
}
```

#### 2. Response Formatting
```go
// Custom success responses
w.Header().Set("Content-Type", "application/json")
json.NewEncoder(w).Encode(YourResponseFormat{
    Success: true,
    Data:    result,
    Meta:    yourMetadata,
})

// Custom error responses  
w.WriteHeader(yourErrorMapper(err))
json.NewEncoder(w).Encode(YourErrorFormat{
    Success: false,
    Error:   err.Error(),
    Code:    yourErrorCode(err),
})
```

#### 3. Middleware Integration
```go
// Add framework-specific middleware:
// - Authentication/authorization
// - Request logging
// - Metrics collection
// - CORS handling
// - Rate limiting
// - Compression
// - Cache control
```

#### 4. Performance Optimizations
```go
// Optimize for your use case:
// - Connection pooling
// - Response caching
// - Request batching
// - Async processing
// - Load shedding
```

## Reference Examples

The `/examples` directory contains simplified implementations for common frameworks:

```bash
# These are REFERENCE implementations only
# Customize them for your production needs

examples/nethttp/     # Basic net/http adapter
examples/gin/         # Gin framework adapter  
examples/echo/        # Echo framework adapter
examples/fiber/       # Fiber framework adapter
examples/chi/         # Chi router adapter
examples/gorilla/     # Gorilla/mux adapter
```

**Important**: These examples demonstrate the integration pattern but lack production-ready features like:
- Proper error handling
- Request validation
- Authentication/authorization
- Logging and monitoring
- Performance optimizations

## API Reference

### Core Types

```go
// Handler - Interface for all business logic handlers
type Handler interface {
    Invoke(ctx context.Context, bind Bind) (any, error)
    Params() any      // Returns zero value of parameter type
    Response() any    // Returns zero value of response type
}

// Router - Interface for framework adapters
type Router interface {
    Register(method string, path string, api Handler, options ...map[string]string)
}

// Bind - Dependency injection function type
type Bind = func(ptr any) (err error)
```

### Key Functions

```go
// Register - Type-safe route registration
func Register[Params any, Response any](
    router Router,
    method string, 
    path string,
    fn func(ctx context.Context, params *Params) (Response, error),
    options ...map[string]string,
)

// DefineGroup - Create route groups with shared dependencies
func DefineGroup[Handler any](
    registerFn func(router Router, handler Handler),
    options ...map[string]string,
) Binder

// NewGroups - Thread-safe route group collection
func NewGroups() *Groups
```

### Handler Types

```go
// HandlerFunc - Generic adapter for business logic functions
type HandlerFunc[Params any, Response any] func(ctx context.Context, bind Bind) (Response, error)
```

## Project Structure Recommendations

```
your-project/
├── internal/
│   ├── business/          # Pure business logic (no framework deps)
│   │   ├── user.go        # User-related business functions
│   │   ├── product.go     # Product-related business functions
│   │   └── order.go       # Order-related business functions
│   │
│   ├── adapter/           # Framework adapters
│   │   ├── http/          # HTTP framework adapters
│   │   │   ├── gin.go     # Gin adapter (customized)
│   │   │   ├── echo.go    # Echo adapter (customized)
│   │   │   └── common.go  # Shared adapter utilities
│   │   └── cli/           # CLI adapters (if needed)
│   │
│   └── types/             # Shared types
│       ├── request.go     # Request type definitions
│       └── response.go    # Response type definitions
│
├── cmd/
│   └── server/           # Application entry point
│       └── main.go       # Wire everything together
│
└── go.mod
```

## Benefits

1. **Framework Independence**: Switch web frameworks without rewriting business logic
2. **Testability**: Test business logic without HTTP framework dependencies  
3. **Maintainability**: Clear separation of concerns, easier to understand and modify
4. **Reusability**: Use the same business logic across different interfaces (HTTP, CLI, gRPC, etc.)
5. **Type Safety**: Catch errors at compile time, not runtime
6. **Performance**: Optimize adapters for specific use cases without affecting business logic

## When to Use xmux

✅ **Good for**:
- Applications that need to support multiple web frameworks
- Teams that want to focus on business logic first
- Projects requiring high test coverage
- Systems that may need to switch frameworks in the future
- APIs with complex request/response formats

❌ **Not ideal for**:
- Simple one-off scripts
- Projects tightly coupled to a specific framework's features
- Teams not comfortable with dependency injection patterns

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

MIT License - see [LICENSE](LICENSE) for details.