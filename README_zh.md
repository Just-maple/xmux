# xmux

[![Go 版本](https://img.shields.io/badge/go-1.18%2B-blue)](https://golang.org/)
[![许可证](https://img.shields.io/badge/license-MIT-green)](LICENSE)

**为中大型 Go 项目设计的框架无关业务逻辑声明库。** 编写一次业务逻辑，随处部署。

xmux 使您能够定义类型安全、依赖注入的业务逻辑，完全独立于任何 Web 框架。专注于核心业务领域，同时保持适应任何 HTTP 框架或自动生成全面 API 文档的灵活性。

## 适合中大型项目

xmux 专为需要以下功能的项目设计：
- **框架无关的业务声明** - 定义一次 API，与任何 Web 框架配合使用
- **跨团队统一的 API 风格** - 通过类型安全的泛型强制执行一致的模式
- **自动化文档生成** - 利用类型信息生成 OpenAPI/Swagger 文档
- **清晰架构** - 将业务逻辑与框架依赖分离
- **面向未来的 API** - 无需重写业务逻辑即可切换框架

## 为什么选择 xmux？

### 问题
传统的 Go Web 开发将业务逻辑与特定框架 API 绑定。切换框架或适应不同的 API 需求意味着重写应用程序的大部分代码。

### 解决方案
xmux 通过以下方式将业务逻辑与 Web 框架依赖解耦：
- **类型安全的处理器**，提供编译时验证
- **依赖注入**，实现可测试、模块化的代码
- **框架无关的设计**，适用于任何 HTTP 路由器
- **可定制的适配器**，针对特定 API 需求进行优化

## 核心特性

- **框架无关的业务声明** - 完全独立于 Web 框架定义业务逻辑
- **统一的 API 风格强制执行** - 使用 Go 泛型在团队间强制执行一致的模式
- **自动化文档就绪** - 类型信息支持自动生成 OpenAPI/Swagger 文档
- **零依赖** - xmux 本身没有外部依赖
- **业务逻辑优先** - 编写一次，适配到任何框架
- **类型安全** - 请求/响应类型的编译时检查
- **依赖注入** - 通过简单的绑定函数实现灵活的 DI
- **框架无关** - 适用于 net/http、Gin、Echo、Fiber、Chi、Gorilla/mux 等
- **可定制的适配器** - 为您的 API 约定创建优化的适配器
- **路由分组** - 通过共享依赖组织相关路由
- **线程安全** - 支持并发注册和绑定

## 安装

```bash
go get github.com/Just-maple/xmux
```

## 快速开始

### 1. 定义业务逻辑

```go
package business

import "context"

// 请求类型（自动从 HTTP 请求绑定）
type CreateUserRequest struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}

// 响应类型（自动序列化为 HTTP 响应）
type UserResponse struct {
    ID    string `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

// 纯粹的业务逻辑函数
func CreateUser(ctx context.Context, req *CreateUserRequest) (*UserResponse, error) {
    // 您的业务逻辑在这里
    // 没有框架依赖！
    return &UserResponse{
        ID:    "user-123",
        Name:  req.Name,
        Email: req.Email,
    }, nil
}
```

### 2. 为框架创建自定义适配器

```go
package adapter

import (
    "encoding/json"
    "net/http"
    
    "github.com/Just-maple/xmux"
    "github.com/your-app/business"
)

// 针对 API 需求优化的自定义适配器
type CustomRouter struct {
    mux *http.ServeMux
}

func NewCustomRouter() *CustomRouter {
    return &CustomRouter{mux: http.NewServeMux()}
}

func (r *CustomRouter) Register(method, path string, api xmux.Handler, opts ...map[string]string) {
    r.mux.HandleFunc(path, func(w http.ResponseWriter, req *http.Request) {
        // 为您的 API 约定自定义请求解析
        bind := func(ptr any) error {
            // 为 POST/PUT 请求解析 JSON 体
            if method == http.MethodPost || method == http.MethodPut {
                return json.NewDecoder(req.Body).Decode(ptr)
            }
            // 为 GET 请求解析查询参数
            // 在此添加您的自定义解析逻辑
            return nil
        }
        
        // 执行业务逻辑
        result, err := api.Invoke(req.Context(), bind)
        if err != nil {
            // 您的自定义错误处理
            w.WriteHeader(http.StatusBadRequest)
            json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
            return
        }
        
        // 您的自定义成功响应
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(result)
    })
}
```

### 3. 连接所有组件

```go
package main

import (
    "net/http"
    
    "github.com/Just-maple/xmux"
    "github.com/your-app/adapter"
    "github.com/your-app/business"
)

func main() {
    // 创建自定义适配器
    router := adapter.NewCustomRouter()
    
    // 类型安全地注册业务逻辑
    xmux.Register(router, http.MethodPost, "/users", business.CreateUser)
    
    // 启动服务器
    http.ListenAndServe(":8080", router)
}
```

## 切面化路由器设计

xmux 的 `Router` 接口作为一个**统一的横切关注点**，支持跨所有端点的一致请求/响应处理。这种切面化方法使您能够：

### 1. **统一的请求/响应处理**
```go
// 强制执行一致 API 模式的自定义适配器
type StandardizedRouter struct {
    baseRouter xmux.Router
    config     RouterConfig  // 横切关注点配置
}

func NewStandardizedRouter(baseRouter xmux.Router, config RouterConfig) *StandardizedRouter {
    return &StandardizedRouter{
        baseRouter: baseRouter,
        config:     config,
    }
}

func (r *StandardizedRouter) Register(method, path string, api xmux.Handler, opts ...map[string]string) {
    // 在传递给基础路由器之前应用横切关注点
    wrappedHandler := r.applyCrossCuttingConcerns(api)
    r.baseRouter.Register(method, path, wrappedHandler, opts...)
}

func (r *StandardizedRouter) applyCrossCuttingConcerns(api xmux.Handler) xmux.Handler {
    return xmux.HandlerFunc(func(ctx context.Context, bind xmux.Bind) (any, error) {
        // 1. 统一的请求验证
        if err := r.validateRequest(ctx); err != nil {
            return nil, err
        }
        
        // 2. 统一的认证/授权
        if err := r.authenticate(ctx); err != nil {
            return nil, err
        }
        
        // 3. 执行业务逻辑
        result, err := api.Invoke(ctx, bind)
        
        // 4. 统一的错误处理
        if err != nil {
            return nil, r.normalizeError(err)
        }
        
        // 5. 统一的响应格式化
        return r.formatResponse(result), nil
    })
}
```

### 2. **一致的参数绑定**
```go
// 所有端点的标准化绑定函数
type StandardBindFunc struct {
    parsers []RequestParser  // 请求解析策略
}

func (b *StandardBindFunc) Bind(ptr any) error {
    // 应用一致的解析规则：
    
    // 1. JSON 体解析（用于 POST/PUT/PATCH）
    if err := b.parseJSONBody(ptr); err != nil {
        return err
    }
    
    // 2. 查询参数绑定
    if err := b.parseQueryParams(ptr); err != nil {
        return err
    }
    
    // 3. 路径参数绑定  
    if err := b.parsePathParams(ptr); err != nil {
        return err
    }
    
    // 4. 头部解析（例如，认证令牌）
    if err := b.parseHeaders(ptr); err != nil {
        return err
    }
    
    // 5. 验证（使用如 `validate:"required"` 的结构体标签）
    return b.validate(ptr)
}

// 具有统一绑定规则的请求结构体示例
type CreateOrderRequest struct {
    // JSON 体字段
    ProductID string `json:"product_id" validate:"required,uuid"`
    Quantity  int    `json:"quantity" validate:"required,min=1,max=100"`
    
    // 查询参数
    UserID string `query:"user_id" validate:"required"`
    
    // 路径参数（从路由中提取）
    StoreID string `path:"store_id"`
    
    // 头部值
    AuthToken string `header:"Authorization"`
}
```

### 3. **统一的响应格式化**
```go
// 所有端点的标准响应包装器
type APIResponse[T any] struct {
    Success bool         `json:"success"`
    Data    T            `json:"data,omitempty"`
    Error   *APIError    `json:"error,omitempty"`
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

// 应用于所有处理器的响应格式化器
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

### 4. **业务逻辑隔离**
```go
// 业务逻辑完全与 HTTP 关注点隔离
type OrderService interface {
    CreateOrder(ctx context.Context, req *CreateOrderRequest) (*OrderResponse, error)
    CancelOrder(ctx context.Context, orderID string) error
    GetOrder(ctx context.Context, orderID string) (*OrderResponse, error)
}

// 纯粹的业务实现
type orderServiceImpl struct {
    repo      OrderRepository
    inventory InventoryService
    payment   PaymentService
}

func (s *orderServiceImpl) CreateOrder(ctx context.Context, req *CreateOrderRequest) (*OrderResponse, error) {
    // 没有 HTTP 框架依赖！
    // 纯粹的业务逻辑：
    
    // 1. 检查库存
    available, err := s.inventory.CheckAvailability(ctx, req.ProductID, req.Quantity)
    if err != nil || !available {
        return nil, ErrInsufficientInventory
    }
    
    // 2. 在数据库中创建订单
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
    
    // 3. 处理支付
    paymentResult, err := s.payment.Process(ctx, order.ID, order.CalculateTotal())
    if err != nil {
        // 更新订单状态为失败
        order.Status = OrderStatusPaymentFailed
        s.repo.Save(ctx, order)
        return nil, err
    }
    
    // 4. 返回结果
    return &OrderResponse{
        OrderID:   order.ID,
        Status:    order.Status,
        Total:     order.CalculateTotal(),
        PaymentID: paymentResult.PaymentID,
    }, nil
}
```

### 5. **横切关注点示例**

```go
// 应用于所有路由的认证中间件
type AuthMiddleware struct {
    authService AuthService
}

func (m *AuthMiddleware) Wrap(api xmux.Handler) xmux.Handler {
    return xmux.HandlerFunc(func(ctx context.Context, bind xmux.Bind) (any, error) {
        // 提取并验证令牌
        token := extractTokenFromContext(ctx)
        user, err := m.authService.ValidateToken(ctx, token)
        if err != nil {
            return nil, ErrUnauthorized
        }
        
        // 为业务逻辑添加上下文
        ctx = context.WithValue(ctx, "user", user)
        
        // 继续执行业务逻辑
        return api.Invoke(ctx, bind)
    })
}

// 所有请求的日志中间件
type LoggingMiddleware struct {
    logger Logger
}

func (m *LoggingMiddleware) Wrap(api xmux.Handler) xmux.Handler {
    return xmux.HandlerFunc(func(ctx context.Context, bind xmux.Bind) (any, error) {
        start := time.Now()
        requestID := generateRequestID()
        
        m.logger.Info("请求开始",
            "request_id", requestID,
            "method", getMethodFromContext(ctx),
            "path", getPathFromContext(ctx))
        
        // 执行业务逻辑
        result, err := api.Invoke(ctx, bind)
        
        duration := time.Since(start)
        status := "成功"
        if err != nil {
            status = "错误"
        }
        
        m.logger.Info("请求完成",
            "request_id", requestID,
            "duration", duration,
            "status", status,
            "error", err)
        
        return result, err
    })
}

// 速率限制中间件
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

### 6. **可组合的中间件链**

```go
// 为横切关注点构建中间件链
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
    // 反向应用（最后添加的先执行）
    for i := len(c.middlewares) - 1; i >= 0; i-- {
        handler = c.middlewares[i].Wrap(handler)
    }
    return handler
}

// 在适配器中使用
func (r *StandardizedRouter) applyCrossCuttingConcerns(api xmux.Handler) xmux.Handler {
    chain := NewMiddlewareChain()
    return chain.Apply(api)
}
```

这种切面化设计使您能够：
- **跨所有端点强制执行一致的 API 模式**
- **将业务逻辑与 HTTP 框架关注点隔离**
- **统一应用横切关注点**（认证、日志、速率限制）
- **无需触碰业务逻辑即可切换框架**
- **从类型信息生成一致的文档**
- **维护基础设施和领域逻辑之间的清晰分离**

## 框架无关的业务声明

xmux 使您能够定义完全独立于任何 Web 框架的业务逻辑。这对于中大型项目特别有价值：

### 1. **业务逻辑作为一等公民**
```go
// 在没有框架依赖的情况下定义业务接口
type UserService interface {
    CreateUser(ctx context.Context, req *CreateUserRequest) (*UserResponse, error)
    GetUser(ctx context.Context, id string) (*UserResponse, error)
    ListUsers(ctx context.Context, filter *UserFilter) ([]*UserResponse, error)
}

// 纯粹实现业务逻辑
type userServiceImpl struct {
    db *database.DB
}

func (s *userServiceImpl) CreateUser(ctx context.Context, req *CreateUserRequest) (*UserResponse, error) {
    // 纯粹的业务逻辑 - 没有 HTTP 框架概念
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

### 2. **依赖注入实现可测试性**
```go
// 将依赖定义为接口
type UserRepository interface {
    Save(ctx context.Context, user *User) error
    FindByID(ctx context.Context, id string) (*User, error)
    FindAll(ctx context.Context, filter *UserFilter) ([]*User, error)
}

// 在运行时注入依赖
func NewUserService(repo UserRepository) UserService {
    return &userServiceImpl{repo: repo}
}

// 为测试模拟依赖
func TestUserService_CreateUser(t *testing.T) {
    mockRepo := &MockUserRepository{}
    service := NewUserService(mockRepo)
    
    // 在没有 HTTP 框架的情况下测试业务逻辑
    result, err := service.CreateUser(context.Background(), &CreateUserRequest{
        Name:  "张三",
        Email: "zhangsan@example.com",
    })
    
    assert.NoError(t, err)
    assert.NotNil(t, result)
}
```

### 3. **框架独立性**
```go
// 相同的业务逻辑适用于不同的框架
func main() {
    // 根据需求选择框架
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
    
    // 注册相同的业务逻辑
    userService := business.NewUserService(db)
    userGroup := xmux.DefineGroup(func(r xmux.Router, svc business.UserService) {
        xmux.Register(r, http.MethodPost, "/users", svc.CreateUser)
        xmux.Register(r, http.MethodGet, "/users/:id", svc.GetUser)
        xmux.Register(r, http.MethodGet, "/users", svc.ListUsers)
    })
    
    // 无论框架如何，业务逻辑保持不变
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

## 统一 API 风格与文档生成

xmux 使用 Go 泛ics 的类型安全方法支持自动 API 文档生成和一致的 API 风格强制执行。

### 1. **文档的类型信息**
```go
// 请求和响应类型携带文档元数据
type CreateUserRequest struct {
    Name  string `json:"name" validate:"required,min=2,max=100" doc:"用户全名"`
    Email string `json:"email" validate:"required,email" doc:"用户电子邮件地址"`
    Age   int    `json:"age" validate:"min=18" doc:"用户年龄（岁）"`
}

type UserResponse struct {
    ID        string    `json:"id" doc:"唯一用户标识符"`
    Name      string    `json:"name" doc:"用户全名"`
    Email     string    `json:"email" doc:"用户电子邮件地址"`
    CreatedAt time.Time `json:"created_at" doc:"账户创建时间戳"`
}

// 可以从注册的处理器中提取类型信息用于文档
func generateOpenAPISpec(router xmux.Router) openapi.Spec {
    spec := openapi.NewSpec()
    
    // 从注册的处理器中提取类型信息
    // 这支持自动 OpenAPI 生成
    return spec
}
```

### 2. **一致的 API 模式**
```go
// 在所有端点强制执行一致的错误响应
type APIError struct {
    Code    string `json:"code"`
    Message string `json:"message"`
    Details any    `json:"details,omitempty"`
}

// 标准化的响应包装器
type APIResponse[T any] struct {
    Success bool       `json:"success"`
    Data    T          `json:"data,omitempty"`
    Error   *APIError  `json:"error,omitempty"`
    Meta    map[string]any `json:"meta,omitempty"`
}

// 适配器强制执行一致的响应格式
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
            // 将任何错误转换为标准化格式
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

### 3. **自动文档工具**
```go
// 示例：从 xmux 路由生成 OpenAPI 文档
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
    // 为 OpenAPI 模式生成提取类型信息
    paramType := api.Params()
    responseType := api.Response()
    
    // 从类型信息生成 OpenAPI 操作
    operation := g.generateOperation(paramType, responseType)
    
    // 添加到 OpenAPI 规范
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

## 适配器定制指南

### 何时创建自定义适配器

`/examples` 目录中的示例是**仅作参考的实现**。在生产环境中，您应该创建针对以下方面定制的适配器：

1. **您的 API 约定**
   - 请求/响应格式（JSON、XML、Protobuf 等）
   - 错误处理模式
   - 认证/授权方案
   - 速率限制和配额

2. **您的框架选择**
   - 特定的中间件需求
   - 性能优化
   - 监控和指标集成
   - 部署环境约束

3. **您的业务需求**
   - 领域特定的验证
   - 审计日志
   - 合规性要求
   - 与现有系统的集成

### 关键定制领域

#### 1. 请求解析
```go
bind := func(ptr any) error {
    // 根据您的 API 定制：
    // - 带有验证的 JSON 体解析
    // - URL 参数提取  
    // - 查询字符串解析
    // - 表单数据处理
    // - 基于头部的认证
    // - 内容协商
    return yourCustomParser(req, ptr)
}
```

#### 2. 响应格式化
```go
// 自定义成功响应
w.Header().Set("Content-Type", "application/json")
json.NewEncoder(w).Encode(YourResponseFormat{
    Success: true,
    Data:    result,
    Meta:    yourMetadata,
})

// 自定义错误响应  
w.WriteHeader(yourErrorMapper(err))
json.NewEncoder(w).Encode(YourErrorFormat{
    Success: false,
    Error:   err.Error(),
    Code:    yourErrorCode(err),
})
```

#### 3. 中间件集成
```go
// 添加框架特定的中间件：
// - 认证/授权
// - 请求日志
// - 指标收集
// - CORS 处理
// - 速率限制
// - 压缩
// - 缓存控制
```

#### 4. 性能优化
```go
// 为您的用例优化：
// - 连接池
// - 响应缓存
// - 请求批处理
// - 异步处理
// - 负载削峰
```

## 参考示例

`/examples` 目录包含常见框架的简化实现：

```bash
# 这些仅是参考实现
# 为您的生产需求定制它们

examples/nethttp/     # 基础 net/http 适配器
examples/gin/         # Gin 框架适配器  
examples/echo/        # Echo 框架适配器
examples/fiber/       # Fiber 框架适配器
examples/chi/         # Chi 路由器适配器
examples/gorilla/     # Gorilla/mux 适配器
```

**重要**：这些示例演示了集成模式，但缺少生产就绪的功能，如：
- 适当的错误处理
- 请求验证
- 认证/授权
- 日志和监控
- 性能优化

## API 参考

### 核心类型

```go
// Handler - 所有业务逻辑处理器的接口
type Handler interface {
    Invoke(ctx context.Context, bind Bind) (any, error)
    Params() any      // 返回参数类型的零值
    Response() any    // 返回响应类型的零值
}

// Router - 框架适配器的接口
type Router interface {
    Register(method string, path string, api Handler, options ...map[string]string)
}

// Bind - 依赖注入函数类型
type Bind = func(ptr any) (err error)
```

### 关键函数

```go
// Register - 类型安全的路由注册
func Register[Params any, Response any](
    router Router,
    method string, 
    path string,
    fn func(ctx context.Context, params *Params) (Response, error),
    options ...map[string]string,
)

// DefineGroup - 创建具有共享依赖的路由组
func DefineGroup[Handler any](
    registerFn func(router Router, handler Handler),
    options ...map[string]string,
) Binder

// NewGroups - 线程安全的路由组集合
func NewGroups() *Groups
```

### 处理器类型

```go
// HandlerFunc - 业务逻辑函数的通用适配器
type HandlerFunc[Params any, Response any] func(ctx context.Context, bind Bind) (Response, error)
```

## 项目结构建议

```
your-project/
├── internal/
│   ├── business/          # 纯粹的业务逻辑（无框架依赖）
│   │   ├── user.go        # 用户相关业务函数
│   │   ├── product.go     # 产品相关业务函数
│   │   └── order.go       # 订单相关业务函数
│   │
│   ├── adapter/           # 框架适配器
│   │   ├── http/          # HTTP 框架适配器
│   │   │   ├── gin.go     # Gin 适配器（定制）
│   │   │   ├── echo.go    # Echo 适配器（定制）
│   │   │   └── common.go  # 共享适配器工具
│   │   └── cli/           # CLI 适配器（如需要）
│   │
│   └── types/             # 共享类型
│       ├── request.go     # 请求类型定义
│       └── response.go    # 响应类型定义
│
├── cmd/
│   └── server/           # 应用程序入口点
│       └── main.go       # 连接所有组件
│
└── go.mod
```

## 优势

1. **框架独立性**：无需重写业务逻辑即可切换 Web 框架
2. **可测试性**：在没有 HTTP 框架依赖的情况下测试业务逻辑
3. **可维护性**：清晰的关注点分离，更容易理解和修改
4. **可重用性**：在不同的接口（HTTP、CLI、gRPC 等）中使用相同的业务逻辑
5. **类型安全**：在编译时捕获错误，而不是运行时
6. **性能**：在不影响业务逻辑的情况下为特定用例优化适配器

## 何时使用 xmux

✅ **适合**：
- 需要支持多个 Web 框架的应用程序
- 希望首先关注业务逻辑的团队
- 需要高测试覆盖率的项目
- 未来可能需要切换框架的系统
- 具有复杂请求/响应格式的 API

❌ **不理想**：
- 简单的一次性脚本
- 紧密耦合到特定框架功能的项目
- 不熟悉依赖注入模式的团队

## 贡献

欢迎贡献！请参阅 [CONTRIBUTING.md](CONTRIBUTING.md) 了解指南。

## 许可证

MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。