package di

import (
	"context"
	"github.com/Just-maple/godi"
	"github.com/Just-maple/xmux/examples/webapp/internal/order/repository"
	"github.com/Just-maple/xmux/examples/webapp/internal/order/service"
	productRepo "github.com/Just-maple/xmux/examples/webapp/internal/product/repository"
	productService "github.com/Just-maple/xmux/examples/webapp/internal/product/service"
	userRepo "github.com/Just-maple/xmux/examples/webapp/internal/user/repository"
	userService "github.com/Just-maple/xmux/examples/webapp/internal/user/service"
)

func BuildContainer() (*godi.Container, func(context.Context, bool), error) {
	c := &godi.Container{}

	shutdown := c.HookOnce("shutdown", func(v any) func(context.Context) {
		return func(ctx context.Context) {
			if closer, ok := v.(interface{ Close() error }); ok {
				closer.Close()
			}
		}
	})

	c.MustAdd(
		godi.Build(func(c *godi.Container) (userRepo.UserRepository, error) {
			return userRepo.NewUserRepository(), nil
		}),

		godi.Build(func(repo userRepo.UserRepository) (*userService.UserService, error) {
			return userService.NewUserService(repo), nil
		}),

		godi.Build(func(c *godi.Container) (productRepo.ProductRepository, error) {
			return productRepo.NewProductRepository(), nil
		}),

		godi.Build(func(repo productRepo.ProductRepository) (*productService.ProductService, error) {
			return productService.NewProductService(repo), nil
		}),

		godi.Build(func(c *godi.Container) (repository.OrderRepository, error) {
			return repository.NewOrderRepository(), nil
		}),

		godi.Build(func(c *godi.Container) (*service.OrderService, error) {
			orderRepo, _ := godi.Inject[repository.OrderRepository](c)
			prodRepo, _ := godi.Inject[productRepo.ProductRepository](c)
			return service.NewOrderService(orderRepo, prodRepo), nil
		}),
	)

	return c, shutdown.Iterate, nil
}
