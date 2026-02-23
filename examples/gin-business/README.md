# Gin Business Example - Complete User Management System

This is a complete business application example built with **xmux** and **Gin**. It demonstrates how to build a production-ready user management system with clean architecture, dependency injection, and framework-agnostic business logic.

## Project Structure

```
gin-business/
├── cmd/
│   └── server/                # Application entry point
│       └── main.go           # Server setup and dependency wiring
├── internal/
│   ├── adapter/              # Framework adapters
│   │   ├── gin_adapter.go   # Production-ready Gin adapter
│   │   └── routes.go        # Route definitions using xmux
│   ├── business/            # Business logic layer
│   │   └── user_service.go  # User business logic (framework-agnostic)
│   ├── models/              # Domain models
│   │   └── user.go          # User domain model
│   ├── repository/          # Data access layer
│   │   └── user_repository.go # User repository interface & implementation
│   └── types/               # Request/Response types
│       └── user.go          # API types for user operations
├── go.mod
└── README.md
```

## Key Features Demonstrated

### 1. **Framework-Agnostic Business Logic**
- `internal/business/user_service.go` contains pure business logic
- No dependency on any web framework
- Easy to test without HTTP context

### 2. **Clean Architecture**
- Clear separation of concerns
- Dependency injection through interfaces
- Business logic isolated from infrastructure concerns

### 3. **Production-Ready Gin Adapter**
- `internal/adapter/gin_adapter.go` shows a realistic adapter
- Request parsing with error handling
- Response formatting and error mapping
- Context enrichment and middleware support

### 4. **Type-Safe Routing with xmux**
- Route groups with shared dependencies
- Compile-time type checking
- Dependency injection for route handlers

### 5. **Complete User Management**
- User registration and authentication
- Profile management
- Password management
- User listing with pagination
- Role-based access control

## Running the Example

### Prerequisites
- Go 1.18 or higher

### Installation
```bash
cd examples/gin-business
go mod tidy
```

### Running the Server
```bash
go run cmd/server/main.go
```

### Environment Variables
```bash
# Server port (default: 8080)
export PORT=3000

# Enable debug mode
export DEBUG=true
```

## API Endpoints

### Public Routes (No Authentication Required)

#### Register User
```http
POST /api/v1/users/register
Content-Type: application/json

{
  "username": "john_doe",
  "email": "john@example.com",
  "password": "SecurePass123!",
  "full_name": "John Doe",
  "role": "user"
}
```

#### Login
```http
POST /api/v1/users/login
Content-Type: application/json

{
  "username": "john_doe",
  "password": "SecurePass123!"
}
```

### Protected Routes (Authentication Required)

#### Get User Profile
```http
GET /api/v1/users/me
Authorization: Bearer <token>
```

#### Get User by ID
```http
GET /api/v1/users/{id}
Authorization: Bearer <token>
```

#### Update User
```http
PUT /api/v1/users/{id}
Authorization: Bearer <token>
Content-Type: application/json

{
  "email": "newemail@example.com",
  "full_name": "John Smith"
}
```

#### Delete User
```http
DELETE /api/v1/users/{id}
Authorization: Bearer <token>
```

#### Change Password
```http
POST /api/v1/users/{id}/change-password
Authorization: Bearer <token>
Content-Type: application/json

{
  "old_password": "SecurePass123!",
  "new_password": "NewSecurePass456!"
}
```

#### List Users (Admin Only)
```http
GET /api/v1/users?limit=10&offset=0
Authorization: Bearer <token>
```

## Architecture Highlights

### 1. **Business Logic Isolation**
```go
// Pure business logic - no framework dependencies
func (s *userServiceImpl) CreateUser(ctx context.Context, req *types.CreateUserRequest) (*types.UserResponse, error) {
    // Business rules and validation
    // Database operations through repository interface
    // No HTTP framework concepts
}
```

### 2. **Dependency Injection**
```go
// Define dependencies as interfaces
type UserRepository interface {
    Create(ctx context.Context, user *models.User) error
    FindByID(ctx context.Context, id uuid.UUID) (*models.User, error)
    // ...
}

// Inject at runtime
func NewUserService(userRepo UserRepository) UserService {
    return &userServiceImpl{userRepo: userRepo}
}
```

### 3. **Type-Safe Routing**
```go
// Route definitions with type safety
type RegisterParams struct {
    *types.CreateUserRequest
}

xmux.Register(r, http.MethodPost, "/users/register",
    func(ctx context.Context, params *RegisterParams) (*types.UserResponse, error) {
        return svc.CreateUser(ctx, params.CreateUserRequest)
    })
```

### 4. **Production Adapter Features**
- **Request Parsing**: JSON body, query params, path params
- **Error Handling**: Structured error responses
- **Response Formatting**: Consistent API responses
- **Context Enrichment**: Request metadata in context
- **Middleware Support**: Auth, logging, rate limiting

## Testing Business Logic

```go
func TestUserService_CreateUser(t *testing.T) {
    // Create mock repository
    mockRepo := &MockUserRepository{}
    
    // Create service with mock
    service := business.NewUserService(mockRepo)
    
    // Test business logic without HTTP
    result, err := service.CreateUser(context.Background(), &types.CreateUserRequest{
        Username: "testuser",
        Email:    "test@example.com",
        Password: "password123",
        FullName: "Test User",
        Role:     "user",
    })
    
    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.Equal(t, "testuser", result.Username)
}
```

## Extending the Example

### Adding New Features
1. **Add new business service** in `internal/business/`
2. **Define API types** in `internal/types/`
3. **Add repository interface** in `internal/repository/`
4. **Register routes** in `internal/adapter/routes.go`

### Switching Frameworks
To switch from Gin to another framework (e.g., Echo, Fiber):
1. Create a new adapter in `internal/adapter/`
2. Update `main.go` to use the new adapter
3. Business logic remains unchanged

### Adding Middleware
```go
// Example: Authentication middleware
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
        
        // Add user to context
        ctx = context.WithValue(ctx, "user", user)
        
        // Proceed to business logic
        return api.Invoke(ctx, bind)
    })
}
```

## Benefits of This Architecture

1. **Framework Independence**: Business logic works with any HTTP framework
2. **Testability**: Easy to test without HTTP context
3. **Maintainability**: Clear separation of concerns
4. **Flexibility**: Easy to add new features or change infrastructure
5. **Consistency**: Enforced API patterns and error handling
6. **Scalability**: Clean architecture supports growth

## Next Steps

1. **Add Database Integration**: Replace in-memory repository with SQL/NoSQL
2. **Implement Authentication**: Add JWT token generation and validation
3. **Add Logging and Metrics**: Integrate structured logging and monitoring
4. **Add Caching**: Implement caching layer for performance
5. **Add Validation**: Enhance request validation with custom rules
6. **Add Documentation**: Generate OpenAPI documentation from types

This example demonstrates how xmux enables you to focus on business logic while maintaining flexibility in infrastructure choices.