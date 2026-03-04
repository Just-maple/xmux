# Web 应用示例

使用 **xmux**、**godi** 和 **Gin** 构建分层架构 Web 应用的完整大型项目示例。

## 项目结构

```
webapp/
├── cmd/
│   └── server/
│       └── main.go              # 应用入口
├── internal/                    # 业务逻辑层
│   ├── user/                    # 用户模块
│   │   ├── model/
│   │   │   └── model.go         # 用户模型
│   │   ├── repository/
│   │   │   └── user_repository.go  # 用户数据访问层（接口）
│   │   └── service/
│   │       └── user_service.go     # 用户业务逻辑层
│   ├── product/                 # 商品模块
│   │   ├── model/
│   │   │   └── model.go
│   │   ├── repository/
│   │   │   └── product_repository.go
│   │   └── service/
│   │       └── product_service.go
│   └── order/                   # 订单模块
│       ├── model/
│       │   └── model.go
│       ├── repository/
│       │   └── order_repository.go
│       └── service/
│           └── order_service.go
├── pkg/                         # 公共组件层
│   ├── app/
│   │   └── application.go       # 应用层（使用 xmux.Groups 组合路由）
│   ├── controller/
│   │   └── controller.go        # Gin 控制器（xmux.Controller 实现）
│   ├── di/
│   │   └── container.go         # godi 依赖注入容器配置
│   └── server/
│       └── server.go            # HTTP 服务器配置
└── go.mod
```

## 架构分层

```
┌─────────────────────────────────────────────────────────────┐
│                        cmd/server                            │
│                           main.go                            │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                      pkg/server                              │
│                   HTTP 服务器组件                             │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────┐     │
│  │   Config    │  │  HTTP Server│  │  Graceful Shutdown│   │
│  └─────────────┘  └─────────────┘  └─────────────────┘     │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                       pkg/di                                 │
│                   依赖注入容器                               │
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
│                   应用协调层                                 │
│  ┌─────────────────────────────────────────────────────┐    │
│  │  xmux.Groups 组合路由                                 │    │
│  │  - userGroup (用户路由组)                             │    │
│  │  - productGroup (商品路由组)                          │    │
│  │  - orderGroup (订单路由组)                            │    │
│  │  - groups.Bind(ctrl, bindService)                   │    │
│  └─────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                   pkg/controller                             │
│                  Gin 控制器组件                               │
│  ┌─────────────────────────────────────────────────────┐    │
│  │  xmux.Controller 实现                                │    │
│  │  - Handle(): 注册 Gin 路由                            │    │
│  │  - ServeHTTP(): http.Handler 接口                    │    │
│  └─────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    internal/*                                │
│                   业务逻辑层                                 │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │    Model     │  │  Repository  │  │   Service    │      │
│  │  (数据模型)  │  │  (数据访问)  │  │  (业务逻辑)  │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└─────────────────────────────────────────────────────────────┘
```

## 快速开始

```bash
cd examples/webapp
go mod tidy
go run ./cmd/server
```

## API 端点

### 用户管理

| 方法 | 路径 | 描述 | 请求体 |
|------|------|------|--------|
| POST | `/api/users` | 创建用户 | `{"name": "张三", "email": "zhangsan@example.com"}` |
| GET | `/api/users/:id` | 获取用户 | - |
| PUT | `/api/users/:id` | 更新用户 | `{"name": "张三更新"}` |
| DELETE | `/api/users/:id` | 删除用户 | - |

### 商品管理

| 方法 | 路径 | 描述 | 请求体 |
|------|------|------|--------|
| POST | `/api/products` | 创建商品 | `{"name": "iPhone", "price": 999.99}` |
| GET | `/api/products/:id` | 获取商品 | - |
| GET | `/api/products` | 获取商品列表 | - |
| PUT | `/api/products/:id` | 更新商品 | `{"name": "iPhone 15", "price": 1099.99}` |
| DELETE | `/api/products/:id` | 删除商品 | - |

### 订单管理

| 方法 | 路径 | 描述 | 请求体 |
|------|------|------|--------|
| POST | `/api/orders` | 创建订单 | `{"user_id": "user-123", "items": [...]}` |
| GET | `/api/orders/:id` | 获取订单 | - |

## 示例请求

### 创建用户

```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"name":"张三","email":"zhangsan@example.com"}'
```

响应：
```json
{
  "id": "user-1234567890",
  "name": "张三",
  "email": "zhangsan@example.com"
}
```

### 创建商品

```bash
curl -X POST http://localhost:8080/api/products \
  -H "Content-Type: application/json" \
  -d '{"name":"iPhone","description":"苹果手机","price":999.99}'
```

### 创建订单

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

### 获取商品列表

```bash
curl http://localhost:8080/api/products
```

## 核心代码

### 1. 依赖注入配置 (pkg/di/container.go)

```go
func BuildContainer() (*godi.Container, func(context.Context, bool), error) {
    c := &godi.Container{}
    
    c.MustAdd(
        // Repository 层
        godi.Build(func(c *godi.Container) (UserRepository, error) {
            return NewUserRepository(), nil
        }),
        
        // Service 层（依赖 Repository）
        godi.Build(func(repo UserRepository) (*UserService, error) {
            return NewUserService(repo), nil
        }),
        
        // OrderService 依赖多个 Repository
        godi.Build(func(c *godi.Container) (*OrderService, error) {
            orderRepo, _ := godi.Inject[OrderRepository](c)
            prodRepo, _ := godi.Inject[ProductRepository](c)
            return NewOrderService(orderRepo, prodRepo), nil
        }),
    )
    
    return c, shutdown.Iterate, nil
}
```

### 2. 使用 xmux.Groups 组合路由 (pkg/app/application.go)

```go
func (a *Application) RegisterRoutes(ctrl interface {
    Handle(string, string, any, xmux.Api, ...map[string]string)
}) {
    // 简化的 bind 函数，使用 godi.InjectAs
    bindService := func(ptr any) error {
        return godi.InjectAs(a.container, ptr)
    }
    
    // 定义路由组
    userGroup := xmux.DefineGroup(func(r xmux.Router, svc *UserService) {
        xmux.Register(r, http.MethodPost, "/api/users", svc.CreateUser)
        xmux.Register(r, http.MethodGet, "/api/users/:id", svc.GetUser)
    })
    
    productGroup := xmux.DefineGroup(func(r xmux.Router, svc *ProductService) {
        xmux.Register(r, http.MethodPost, "/api/products", svc.CreateProduct)
        xmux.Register(r, http.MethodGet, "/api/products", svc.ListProducts)
    })
    
    orderGroup := xmux.DefineGroup(func(r xmux.Router, svc *OrderService) {
        xmux.Register(r, http.MethodPost, "/api/orders", svc.CreateOrder)
        xmux.Register(r, http.MethodGet, "/api/orders/:id", svc.GetOrder)
    })
    
    // 使用 Groups 组合所有路由组
    groups := xmux.NewGroups(userGroup, productGroup, orderGroup)
    _ = groups.Bind(ctrl, bindService)
}
```

### 3. Gin 控制器实现 (pkg/controller/controller.go)

```go
type Controller struct {
    engine *gin.Engine
}

func (c *Controller) Handle(method, path string, service any, api xmux.Api, ...) {
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

### 4. 业务服务层 (internal/user/service/user_service.go)

```go
type UserService struct {
    repo UserRepository  // 依赖接口，实现依赖倒置
}

func NewUserService(repo UserRepository) *UserService {
    return &UserService{repo: repo}
}

func (s *UserService) CreateUser(ctx context.Context, req *CreateUserRequest) (*UserResponse, error) {
    // 纯业务逻辑，无框架依赖
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

## 核心设计原则

### 1. 依赖倒置原则

```go
// Repository 接口定义
type UserRepository interface {
    Create(ctx context.Context, user *User) error
    GetByID(ctx context.Context, id string) (*User, error)
}

// Service 依赖接口而非具体实现
type UserService struct {
    repo UserRepository
}
```

### 2. 分层架构

- **internal/**: 业务逻辑（Model/Repository/Service）
- **pkg/di/**: 依赖注入配置
- **pkg/app/**: 应用协调层（使用 xmux.Groups 组合路由）
- **pkg/controller/**: HTTP 适配器（Gin）
- **pkg/server/**: 服务器配置

### 3. 框架无关

业务逻辑（Service 层）完全不依赖 Gin 框架，只依赖 xmux 的类型安全接口。

### 4. 使用 godi.InjectAs 简化注入

```go
bindService := func(ptr any) error {
    return godi.InjectAs(a.container, ptr)
}
```

替代了繁琐的 type switch 写法。

### 5. 使用 xmux.Groups 组合路由

```go
groups := xmux.NewGroups(userGroup, productGroup, orderGroup)
_ = groups.Bind(ctrl, bindService)
```

一次性绑定所有路由组。

## 优势

1. **清晰的职责分离** - 每层有明确的职责
2. **易于测试** - 各层可独立测试
3. **依赖注入** - 使用 godi 进行类型安全的 DI
4. **Gin 框架** - 高性能 HTTP 框架
5. **框架无关** - 业务逻辑可移植到其他框架
6. **优雅关闭** - 完整的资源清理
7. **类型安全** - 编译期检查所有类型
8. **代码简洁** - 使用 `godi.InjectAs` 和 `xmux.Groups` 简化代码

## 依赖版本

- Go: 1.18+
- xmux: v1.0.0
- godi: v0.0.0-20260304015920-020362515ad7
- gin: v1.9.1+

## 许可证

MIT License
