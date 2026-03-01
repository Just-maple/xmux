package adapter

import (
	"context"
	"net/http"

	"github.com/Just-maple/godi"
	"github.com/google/uuid"

	"github.com/Just-maple/xmux"
	"github.com/Just-maple/xmux/examples/gin-business/internal/business"
	"github.com/Just-maple/xmux/examples/gin-business/internal/types"
)

type SuccessResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

func RegisterUserRoutes() *xmux.Groups {
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
		types.ListUsersRequest
	}

	publicGroup := xmux.DefineGroup(func(r xmux.Router, svc business.UserService) {
		xmux.Register(r, http.MethodPost, "/users/register",
			func(ctx context.Context, params *RegisterParams) (*types.UserResponse, error) {
				return svc.CreateUser(ctx, params.CreateUserRequest)
			})

		xmux.Register(r, http.MethodPost, "/users/login",
			func(ctx context.Context, params *LoginParams) (*types.LoginResponse, error) {
				return svc.Login(ctx, params.LoginRequest)
			})
	}, map[string]string{"public": "true"})

	protectedGroup := xmux.DefineGroup(func(r xmux.Router, svc business.UserService) {
		xmux.Register(r, http.MethodGet, "/users/me",
			func(ctx context.Context, params *GetUserProfileParams) (*types.UserResponse, error) {
				return svc.GetProfile(ctx, params.UserID)
			})

		xmux.Register(r, http.MethodGet, "/users/:id",
			func(ctx context.Context, params *GetUserByIDParams) (*types.UserResponse, error) {
				return svc.GetUser(ctx, params.UserID)
			})

		xmux.Register(r, http.MethodPut, "/users/:id",
			func(ctx context.Context, params *UpdateUserParams) (*types.UserResponse, error) {
				return svc.UpdateUser(ctx, params.UserID, params.Req)
			})

		xmux.Register(r, http.MethodDelete, "/users/:id",
			func(ctx context.Context, params *DeleteUserParams) (*SuccessResponse, error) {
				if err := svc.DeleteUser(ctx, params.UserID); err != nil {
					return nil, err
				}
				return &SuccessResponse{Success: true, Message: "User deleted successfully"}, nil
			})

		xmux.Register(r, http.MethodPost, "/users/:id/change-password",
			func(ctx context.Context, params *ChangePasswordParams) (*SuccessResponse, error) {
				if err := svc.ChangePassword(ctx, params.UserID, params.Req); err != nil {
					return nil, err
				}
				return &SuccessResponse{Success: true, Message: "Password changed successfully"}, nil
			})

		xmux.Register(r, http.MethodGet, "/users",
			func(ctx context.Context, params *ListUsersParams) (*types.ListUsersResponse, error) {
				return svc.ListUsers(ctx, &params.ListUsersRequest)
			})
	}, map[string]string{"protected": "true", "prefix": "/api/v1"})

	return xmux.NewGroups().Register(publicGroup, protectedGroup)
}

func RegisterAllRoutes(router xmux.Router, container *godi.Container) (err error) {
	return RegisterUserRoutes().Bind(router, func(ptr any) error {
		e := godi.InjectAs(ptr, container)
		return e
	})
}
