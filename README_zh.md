# xmux

[![Go 版本](https://img.shields.io/badge/go-1.18%2B-blue)](https://golang.org/)
[![许可证](https://img.shields.io/badge/license-MIT-green)](LICENSE)

**Go 语言框架无关的 HTTP 路由器**。一次编写业务逻辑，任意 Web 框架部署。

## 特性

- **框架无关** - 业务逻辑与 Web 框架完全解耦
- **类型安全** - 使用 Go 泛型在编译期检查类型
- **依赖注入** - 通过简单的 bind 函数实现 DI
- **零依赖** - xmux 本身没有任何外部依赖
- **路由组** - 使用共享服务组织相关路由

## 安装

```bash
go get github.com/Just-maple/xmux
```

## 快速开始

### 1. 定义业务逻辑

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
    // 纯业务逻辑 - 无框架依赖
    return &UserResponse{
        ID:    "user-123",
        Name:  req.Name,
        Email: req.Email,
    }, nil
}
```

### 2. 实现 Controller 接口

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

// Handle 实现 xmux.Controller 接口
func (c *Controller) Handle(method, path string, api xmux.Api, opts ...map[string]string) {
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

### 3. 组合使用

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

## 示例

`/examples` 目录包含完整的框架示例：

| 框架 | 目录 |
|------|------|
| net/http | `examples/nethttp` |
| Gin | `examples/gin` |
| Echo | `examples/echo` |
| Fiber | `examples/fiber` |
| Chi | `examples/chi` |
| Gorilla/mux | `examples/gorilla` |

运行示例：

```bash
cd examples/nethttp
go mod tidy
go run .
```

## API 参考

### 核心类型

```go
// Api - 类型安全的处理器接口
type Api interface {
    Invoke(ctx context.Context, bind func(params any) error) (any, error)
    Params() any
    Response() any
    Function() any
}

// Controller - 框架适配器接口（用户实现这个！）
type Controller interface {
    Handle(method string, path string, service any, api Api, options ...map[string]string)
}

// Router - 内部路由注册器
type Router interface {
    Register(method string, path string, api Api, options ...map[string]string)
}

// Binder - 依赖注入接口
type Binder interface {
    Bind(handler Controller, bind func(service any) error) (err error)
}
```

### 核心函数

```go
// Register - 注册业务逻辑函数为路由处理器
func Register[Params any, Response any](
    router Router,
    method string,
    path string,
    fn func(ctx context.Context, params *Params) (Response, error),
    options ...map[string]string,
)

// DefineGroup - 创建带有共享服务的路由组
func DefineGroup[Service any](
    fn func(router Router, handler Service),
    options ...map[string]string,
) Binder

// NewGroups - 创建路由组集合
func NewGroups(gs ...Binder) Groups
```

## 架构说明

```
┌─────────────────┐     ┌──────────────────┐     ┌─────────────────┐
│  业务逻辑代码   │────▶│   xmux 核心      │◀────│  Controller     │
│  (你的代码)     │     │   (框架)         │     │ (你的适配器)    │
└─────────────────┘     └──────────────────┘     └─────────────────┘
        │                       │                        │
        │                       │                        │
        ▼                       ▼                        ▼
  func(ctx, *Req)         Api 接口               Handle() 方法
  (Resp, error)           Router 接口            (net/http, Gin,
                          Binder 接口             Echo 等)
```

## 项目结构

```
your-project/
├── business/           # 纯业务逻辑
│   ├── user.go
│   └── order.go
├── controller/         # 框架控制器
│   ├── http.go        # net/http 控制器
│   └── gin.go         # Gin 控制器
└── main.go            # 应用入口
```

## 优势

1. **框架切换无成本** - 更换 Web 框架无需重写业务逻辑
2. **易于测试** - 业务逻辑可独立于 HTTP 框架测试
3. **代码清晰** - 业务逻辑与框架关注点分离
4. **类型安全** - 编译期发现类型错误

## 许可证

MIT License - 详见 [LICENSE](LICENSE)
