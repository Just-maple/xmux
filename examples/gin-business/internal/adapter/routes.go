package adapter

import (
	"context"
	"net/http"

	"github.com/Just-maple/xmux"
	"github.com/Just-maple/xmux/examples/gin-business/internal/business"
	"github.com/Just-maple/xmux/examples/gin-business/internal/types"
	"github.com/google/uuid"
)

// SuccessResponse represents a generic success response
type SuccessResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

// RegisterUserRoutes registers all user-related routes
func RegisterUserRoutes(router xmux.Router, userService business.UserService) {
	// Define parameter types for each route
	type RegisterParams struct {
		*types.CreateUserRequest
	}

	type LoginParams struct {
		*types.LoginRequest
	}

	type GetUserProfileParams struct {
		UserID uuid.UUID `path:"user_id"`
	}

	type GetUserByIDParams struct {
		UserID uuid.UUID `path:"id"`
	}

	type UpdateUserParams struct {
		UserID uuid.UUID                `path:"id"`
		Req    *types.UpdateUserRequest `json:"req"`
	}

	type DeleteUserParams struct {
		UserID uuid.UUID `path:"id"`
	}

	type ChangePasswordParams struct {
		UserID uuid.UUID                    `path:"id"`
		Req    *types.ChangePasswordRequest `json:"req"`
	}

	type ListUsersParams struct {
		*types.ListUsersRequest
	}

	// Public routes (no authentication required)
	publicGroup := xmux.DefineGroup(func(r xmux.Router, svc business.UserService) {
		// Register user
		xmux.Register(r, http.MethodPost, "/users/register",
			func(ctx context.Context, params *RegisterParams) (*types.UserResponse, error) {
				return svc.CreateUser(ctx, params.CreateUserRequest)
			})

		// Login
		xmux.Register(r, http.MethodPost, "/users/login",
			func(ctx context.Context, params *LoginParams) (*types.LoginResponse, error) {
				return svc.Login(ctx, params.LoginRequest)
			})
	}, map[string]string{"public": "true"})

	// Protected routes (authentication required)
	protectedGroup := xmux.DefineGroup(func(r xmux.Router, svc business.UserService) {
		// Get user profile
		xmux.Register(r, http.MethodGet, "/users/me",
			func(ctx context.Context, params *GetUserProfileParams) (*types.UserResponse, error) {
				return svc.GetProfile(ctx, params.UserID)
			})

		// Get user by ID
		xmux.Register(r, http.MethodGet, "/users/:id",
			func(ctx context.Context, params *GetUserByIDParams) (*types.UserResponse, error) {
				return svc.GetUser(ctx, params.UserID)
			})

		// Update user
		xmux.Register(r, http.MethodPut, "/users/:id",
			func(ctx context.Context, params *UpdateUserParams) (*types.UserResponse, error) {
				return svc.UpdateUser(ctx, params.UserID, params.Req)
			})

		// Delete user
		xmux.Register(r, http.MethodDelete, "/users/:id",
			func(ctx context.Context, params *DeleteUserParams) (*SuccessResponse, error) {
				if err := svc.DeleteUser(ctx, params.UserID); err != nil {
					return nil, err
				}
				return &SuccessResponse{Success: true, Message: "User deleted successfully"}, nil
			})

		// Change password
		xmux.Register(r, http.MethodPost, "/users/:id/change-password",
			func(ctx context.Context, params *ChangePasswordParams) (*SuccessResponse, error) {
				if err := svc.ChangePassword(ctx, params.UserID, params.Req); err != nil {
					return nil, err
				}
				return &SuccessResponse{Success: true, Message: "Password changed successfully"}, nil
			})

		// List users (admin only)
		xmux.Register(r, http.MethodGet, "/users",
			func(ctx context.Context, params *ListUsersParams) (*types.ListUsersResponse, error) {
				return svc.ListUsers(ctx, params.ListUsersRequest)
			})
	}, map[string]string{"protected": "true", "prefix": "/api/v1"})

	// Register groups
	groups := xmux.NewGroups()
	groups.Register(publicGroup, protectedGroup)

	// Bind the user service to the groups
	groups.Bind(router, func(ptr any) error {
		switch p := ptr.(type) {
		case *business.UserService:
			*p = userService
		default:
			// In a real implementation, you might have other dependencies
			// like authentication service, logger, etc.
		}
		return nil
	})
}

// RegisterAllRoutes registers all application routes
func RegisterAllRoutes(router xmux.Router, dependencies Dependencies) {
	// Register user routes
	RegisterUserRoutes(router, dependencies.UserService)

	// In a real application, you would also register:
	// - Product routes
	// - Order routes
	// - Payment routes
	// - Health check routes
	// - Metrics routes
	// - Admin routes
}

// Dependencies contains all application dependencies
type Dependencies struct {
	UserService business.UserService
	// Other dependencies would be added here:
	// ProductService  business.ProductService
	// OrderService    business.OrderService
	// AuthService     auth.Service
	// Logger          log.Logger
}
