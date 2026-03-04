# xmux

[![Go Version](https://img.shields.io/badge/go-1.18%2B-blue)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

**Framework-agnostic HTTP router for Go.** Write business logic once, deploy with any web framework.

## Features

- **Framework Agnostic** - Write business logic once, use with net/http, Gin, Echo, Fiber, Chi, Gorilla/mux
- **Type Safe** - Compile-time type checking with Go generics
- **Dependency Injection** - Clean DI through simple bind functions
- **Zero Dependencies** - xmux itself has no external dependencies
- **Route Groups** - Organize related routes with shared services

## Installation

```bash
go get github.com/Just-maple/xmux
```

## Quick Start

### 1. Define Business Logic

```go
package business

import "context"

type CreateUserRequest struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}

type UserResponse struct {
    ID    string `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

func CreateUser(ctx context.Context, req *CreateUserRequest) (*UserResponse, error) {
    // Pure business logic - no framework dependencies
    return &UserResponse{
        ID:    "user-123",
        Name:  req.Name,
        Email: req.Email,
    }, nil
}
```

### 2. Implement Controller Interface

```go
package main

import (
    "encoding/json"
    "net/http"
    "github.com/Just-maple/xmux"
)

type Controller struct {
    mux *http.ServeMux
}

func NewController() *Controller {
    return &Controller{mux: http.NewServeMux()}
}

// Handle implements xmux.Controller interface
func (c *Controller) Handle(method, path string, service any, api xmux.Api, opts ...map[string]string) {
    c.mux.HandleFunc(path, func(w http.ResponseWriter, req *http.Request) {
        if req.Method != method {
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
            return
        }

        bind := func(ptr any) error {
            return json.NewDecoder(req.Body).Decode(ptr)
        }

        result, err := api.Invoke(req.Context(), bind)
        if err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(result)
    })
}

func (c *Controller) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    c.mux.ServeHTTP(w, req)
}
```

### 3. Wire Everything Together

```go
package main

import (
    "log"
    "net/http"
    "github.com/Just-maple/xmux"
    "your-app/business"
)

func main() {
    controller := NewController()
    userService := business.NewUserService()

    userGroup := xmux.DefineGroup(func(r xmux.Router, svc *business.UserService) {
        xmux.Register(r, http.MethodPost, "/users", svc.CreateUser)
        xmux.Register(r, http.MethodGet, "/users", svc.ListUsers)
        xmux.Register(r, http.MethodGet, "/user", svc.GetUser)
    })

    err := userGroup.Bind(controller, func(ptr any) error {
        switch p := ptr.(type) {
        case **business.UserService:
            *p = userService
        }
        return nil
    })
    if err != nil {
        log.Fatal(err)
    }

    log.Fatal(http.ListenAndServe(":8080", controller))
}
```

## Examples

See the `/examples` directory for complete implementations:

| Framework | Directory |
|-----------|-----------|
| net/http | `examples/nethttp` |
| Gin | `examples/gin` |
| Echo | `examples/echo` |
| Fiber | `examples/fiber` |
| Chi | `examples/chi` |
| Gorilla/mux | `examples/gorilla` |

Run an example:

```bash
cd examples/nethttp
go mod tidy
go run .
```

## API Reference

### Core Types

```go
// Api - Type-safe handler interface
type Api interface {
    Invoke(ctx context.Context, bind func(params any) error) (any, error)
    Params() any
    Response() any
    Function() any
}

// Controller - Framework adapter interface (implement this!)
type Controller interface {
    Handle(method string, path string, service any, api Api, options ...map[string]string)
}

// Router - Internal route registrar
type Router interface {
    Register(method string, path string, api Api, options ...map[string]string)
}

// Binder - Dependency injection interface
type Binder interface {
    Bind(handler Controller, bind func(service any) error) (err error)
}
```

### Key Functions

```go
// Register - Register a business logic function as a route handler
func Register[Params any, Response any](
    router Router,
    method string,
    path string,
    fn func(ctx context.Context, params *Params) (Response, error),
    options ...map[string]string,
)

// DefineGroup - Create a route group with shared service
func DefineGroup[Service any](
    fn func(router Router, handler Service),
    options ...map[string]string,
) Binder

// NewGroups - Create a collection of route groups
func NewGroups(gs ...Binder) Groups
```

## Architecture

```
┌─────────────────┐     ┌──────────────────┐     ┌─────────────────┐
│  Business Logic │────▶│   xmux.Core      │◀────│   Controller    │
│  (Your Code)    │     │  (Framework)     │     │ (Your Adapter)  │
└─────────────────┘     └──────────────────┘     └─────────────────┘
        │                       │                        │
        │                       │                        │
        ▼                       ▼                        ▼
  func(ctx, *Req)         Api interface           Handle() method
  (Resp, error)           Router interface        (net/http, Gin,
                          Binder interface         Echo, etc.)
```

## Project Structure

```
your-project/
├── business/           # Pure business logic
│   ├── user.go
│   └── order.go
├── controller/         # Framework controllers
│   ├── http.go        # net/http controller
│   └── gin.go         # Gin controller
└── main.go            # Application entry point
```

## License

MIT License - see [LICENSE](LICENSE) for details.
