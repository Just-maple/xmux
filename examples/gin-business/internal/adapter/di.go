package adapter

import (
	"context"

	"github.com/Just-maple/godi"
	"github.com/Just-maple/xmux/examples/gin-business/internal/business"
	"github.com/Just-maple/xmux/examples/gin-business/internal/repository"
	"github.com/Just-maple/xmux/examples/gin-business/internal/types"
)

func NewContainer() *godi.Container {
	c := &godi.Container{}

	c.MustAdd(
		godi.Build(func(c *godi.Container) (repository.UserRepository, error) {
			return repository.NewInMemoryUserRepository(), nil
		}),
		godi.Build(func(c *godi.Container) (business.UserService, error) {
			repo, err := godi.Inject[repository.UserRepository](c)
			if err != nil {
				return nil, err
			}
			return business.NewUserService(repo), nil
		}),
	)

	return c
}

func InitSampleData(c *godi.Container) {
	userService, err := godi.Inject[business.UserService](c)
	if err != nil {
		return
	}

	ctx := context.Background()

	userService.CreateUser(ctx, &types.CreateUserRequest{
		Username: "admin",
		Email:    "admin@example.com",
		Password: "Admin123!",
		FullName: "System Administrator",
		Role:     "admin",
	})

	userService.CreateUser(ctx, &types.CreateUserRequest{
		Username: "user",
		Email:    "user@example.com",
		Password: "User123!",
		FullName: "Regular User",
		Role:     "user",
	})
}
