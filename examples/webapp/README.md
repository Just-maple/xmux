# Web Application Example

Complete large-scale project example demonstrating how to build a layered web application with **xmux**, **godi**, and **Gin**.

## Project Structure

```
webapp/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/                    # Business logic layer
│   ├── user/
│   │   ├── model/
│   │   │   └── model.go         # User model
│   │   ├── repository/
│   │   │   └── user_repository.go  # User data access layer (interface)
│   │   └── service/
│   │       └── user_service.go     # User business logic layer
│   ├── product/
│   │   ├── model/
│   │   │   └── model.go
│   │   ├── repository/
│   │   │   └── product_repository.go
│   │   └── service/
│   │       └── product_service.go
│   └── order/
│       ├── model/
│       │   └── model.go
│       ├── repository/
│       │   └── order_repository.go
│       └── service/
│           └── order_service.go
├── pkg/                       # Common components layer
│   ├── app/
│   │   └── application.go     # Application layer (using xmux.Groups to combine routes)
│   ├── controller/
│   │   └── controller.go      # Gin controller (xmux.Controller implementation)
│   ├── di/
│   │   └── container.go       # godi dependency injection container configuration
│   └── server/
│       └── server.go          # HTTP server configuration
└── go.mod
```

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                        cmd/server                            │
│                           main.go                            │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                      pkg/server                              │
│                      Server Component                        │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────┐     │
│  │   Config    │  │  HTTP Server│  │  Graceful Shutdown│   │
│  └─────────────┘  └─────────────┘  └─────────────────┘     │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                       pkg/di                                 │
│                   DI Container                               │
│  ┌─────────────────────────────────────────────────────┐    │
│  │  godi.Container                                     │    │
│  │  - UserRepository                                   │    │
│  │  - UserService                                      │    │
│  │  - ProductRepository                                │    │
│  │  - ProductService                                   │    │
│  │  - OrderRepository                                  │    │
│  │  - OrderService                                     │    │
│  └─────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                      pkg/app                                 │
│                   Application Component                      │
│  ┌─────────────────────────────────────────────────────┐    │
│  │  xmux.Groups combine routes                          │    │
│  │  - userGroup                                         │    │
│  │  - productGroup                                      │    │
│  │  - orderGroup                                        │    │
│  │  - groups.Bind(ctrl, bindService)                    │    │
│  └─────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                   pkg/controller                             │
│                  Gin Controller Component                    │
│  ┌─────────────────────────────────────────────────────┐    │
│  │  xmux.Controller implementation                      │    │
│  │  - Handle(): Register Gin routes                     │    │
│  │  - ServeHTTP(): http.Handler interface               │    │
│  └─────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    internal/*                                │
│                   Business Logic Layer                       │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │    Model     │  │  Repository  │  │   Service    │      │
│  │  (Data)      │  │  (Data Access)│ │  (Business)  │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└─────────────────────────────────────────────────────────────┘
```

## Quick Start

```bash
cd examples/webapp
go mod tidy
go run ./cmd/server
```

## API Endpoints

### Users

| Method | Path | Description | Request Body |
|--------|------|-------------|--------------|
| POST | `/api/users` | Create user | `{"name": "John", "email": "john@example.com"}` |
| GET | `/api/users/:id` | Get user | - |
| PUT | `/api/users/:id` | Update user | `{"name": "John Updated"}` |
| DELETE | `/api/users/:id` | Delete user | - |

### Products

| Method | Path | Description | Request Body |
|--------|------|-------------|--------------|
| POST | `/api/products` | Create product | `{"name": "iPhone", "price": 999.99}` |
| GET | `/api/products/:id` | Get product | - |
| GET | `/api/products` | List products | - |
| PUT | `/api/products/:id` | Update product | `{"name": "iPhone 15", "price": 1099.99}` |
| DELETE | `/api/products/:id` | Delete product | - |

### Orders

| Method | Path | Description | Request Body |
|--------|------|-------------|--------------|
| POST | `/api/orders` | Create order | `{"user_id": "user-123", "items": [...]}` |
| GET | `/api/orders/:id` | Get order | - |

## Example Requests

### Create User

```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"name":"John","email":"john@example.com"}'
```

Response:
```json
{
  "id": "user-1234567890",
  "name": "John",
  "email": "john@example.com"
}
```

### Create Product

```bash
curl -X POST http://localhost:8080/api/products \
  -H "Content-Type: application/json" \
  -d '{"name":"iPhone","description":"Apple Phone","price":999.99}'
```

### Create Order

```bash
curl -X POST http://localhost:8080/api/orders \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user-1234567890",
    "items": [
      {"product_id": "product-123", "quantity": 2}
    ]
  }'
```

### List Products

```bash
curl http://localhost:8080/api/products
```

## Core Code

### 1. Dependency Injection Configuration (pkg/di/container.go)

```go
func BuildContainer() (*godi.Container, func(context.Context, bool), error) {
    c := &godi.Container{}
    
    c.MustAdd(
        // Repository layer
        godi.Build(func(c *godi.Container) (UserRepository, error) {
            return NewUserRepository(), nil
        }),
        
        // Service layer (depends on Repository)
        godi.Build(func(repo UserRepository) (*UserService, error) {
            return NewUserService(repo), nil
        }),
        
        // OrderService depends on multiple Repositories
        godi.Build(func(c *godi.Container) (*OrderService, error) {
            orderRepo, _ := godi.Inject[OrderRepository](c)
            prodRepo, _ := godi.Inject[ProductRepository](c)
            return NewOrderService(orderRepo, prodRepo), nil
        }),
    )
    
    return c, shutdown.Iterate, nil
}
```

### 2. Using xmux.Groups to Combine Routes (pkg/app/application.go)

```go
func (a *Application) RegisterRoutes(ctrl interface {
    Handle(string, string, any, xmux.Api, ...map[string]string)
}) {
    // Simplified bind function using godi.InjectAs
    bindService := func(ptr any) error {
        return godi.InjectAs(a.container, ptr)
    }
    
    // Define route groups
    userGroup := xmux.ServiceGroup(func(r xmux.Router, svc *UserService) {
        xmux.Register(r, http.MethodPost, "/api/users", svc.CreateUser)
        xmux.Register(r, http.MethodGet, "/api/users/:id", svc.GetUser)
    })
    
    productGroup := xmux.ServiceGroup(func(r xmux.Router, svc *ProductService) {
        xmux.Register(r, http.MethodPost, "/api/products", svc.CreateProduct)
        xmux.Register(r, http.MethodGet, "/api/products", svc.ListProducts)
    })
    
    orderGroup := xmux.ServiceGroup(func(r xmux.Router, svc *OrderService) {
        xmux.Register(r, http.MethodPost, "/api/orders", svc.CreateOrder)
        xmux.Register(r, http.MethodGet, "/api/orders/:id", svc.GetOrder)
    })
    
    // Use Groups to combine all route groups
    groups := xmux.NewGroups(userGroup, productGroup, orderGroup)
    _ = groups.Bind(ctrl, bindService)
}
```

### 3. Gin Controller Implementation (pkg/controller/controller.go)

```go
type Controller struct {
    engine *gin.Engine
}

func (c *Controller) Handle(method, path string, api xmux.Api, ...) {
    c.engine.Handle(method, path, func(ctx *gin.Context) {
        bind := func(ptr any) error {
            if err := ctx.ShouldBindJSON(ptr); err != nil {
                return ctx.ShouldBindQuery(ptr)
            }
            return nil
        }
        
        result, err := api.Invoke(ctx.Request.Context(), bind)
        if err != nil {
            ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        
        ctx.JSON(http.StatusOK, result)
    })
}
```

### 4. Business Service Layer (internal/user/service/user_service.go)

```go
type UserService struct {
    repo UserRepository  // Depends on interface, implementing dependency inversion
}

func NewUserService(repo UserRepository) *UserService {
    return &UserService{repo: repo}
}

func (s *UserService) CreateUser(ctx context.Context, req *CreateUserRequest) (*UserResponse, error) {
    // Pure business logic, no framework dependencies
    user := &model.User{
        ID:    fmt.Sprintf("user-%d", time.Now().UnixNano()),
        Name:  req.Name,
        Email: req.Email,
    }
    
    if err := s.repo.Create(ctx, user); err != nil {
        return nil, err
    }
    
    return &UserResponse{ID: user.ID, Name: user.Name, Email: user.Email}, nil
}
```

## Core Design Principles

### 1. Dependency Inversion Principle

```go
// Repository interface definition
type UserRepository interface {
    Create(ctx context.Context, user *User) error
    GetByID(ctx context.Context, id string) (*User, error)
}

// Service depends on interface, not concrete implementation
type UserService struct {
    repo UserRepository
}
```

### 2. Layered Architecture

- **internal/**: Business logic (Model/Repository/Service)
- **pkg/di/**: Dependency injection configuration
- **pkg/app/**: Application coordination layer (using xmux.Groups to combine routes)
- **pkg/controller/**: HTTP adapter (Gin)
- **pkg/server/**: Server configuration

### 3. Framework Agnostic

Business logic (Service layer) has no Gin framework dependencies, only depends on xmux's type-safe interface.

### 4. Using godi.InjectAs for Simplified Injection

```go
bindService := func(ptr any) error {
    return godi.InjectAs(a.container, ptr)
}
```

Replaces verbose type switch approach.

### 5. Using xmux.Groups to Combine Routes

```go
groups := xmux.NewGroups(userGroup, productGroup, orderGroup)
_ = groups.Bind(ctrl, bindService)
```

Binds all route groups at once.

## Benefits

1. **Clear Separation of Concerns** - Each layer has well-defined responsibilities
2. **Easy to Test** - Each layer can be tested independently
3. **Dependency Injection** - Type-safe DI using godi
4. **Gin Framework** - High-performance HTTP framework
5. **Framework Agnostic** - Business logic is portable to other frameworks
6. **Graceful Shutdown** - Complete resource cleanup
7. **Type Safety** - Compile-time type checking for all types
8. **Concise Code** - Using `godi.InjectAs` and `xmux.Groups` to simplify code

## Dependencies

- Go: 1.18+
- xmux: v1.0.0
- godi: v0.0.0-20260304015920-020362515ad7
- gin: v1.9.1+

## License

MIT License
