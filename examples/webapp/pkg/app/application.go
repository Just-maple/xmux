package app

import (
	"context"
	"log"
	"net/http"

	"github.com/Just-maple/godi"

	"github.com/Just-maple/xmux"
	orderService "github.com/Just-maple/xmux/examples/webapp/internal/order/service"
	productModel "github.com/Just-maple/xmux/examples/webapp/internal/product/model"
	productService "github.com/Just-maple/xmux/examples/webapp/internal/product/service"
	userModel "github.com/Just-maple/xmux/examples/webapp/internal/user/model"
	userService "github.com/Just-maple/xmux/examples/webapp/internal/user/service"
)

type Application struct {
	container *godi.Container
}

func NewApplication(container *godi.Container) *Application {
	return &Application{
		container: container,
	}
}

func (a *Application) RegisterRoutes(ctrl xmux.Controller) {
	bindService := func(ptr any) error {
		return godi.InjectAs(a.container, ptr)
	}

	userGroup := xmux.DefineGroup(func(r xmux.Router, svc *userService.UserService) {
		log.Println("Registering user routes")
		xmux.Register(r, http.MethodPost, "/api/users", svc.CreateUser)
		xmux.Register(r, http.MethodGet, "/api/users/:id", svc.GetUser)
		xmux.Register(r, http.MethodPut, "/api/users/:id", svc.UpdateUser)
		xmux.Register(r, http.MethodDelete, "/api/users/:id", func(ctx context.Context, req *userModel.DeleteUserRequest) (any, error) {
			return nil, svc.DeleteUser(ctx, req)
		})
	})

	productGroup := xmux.DefineGroup(func(r xmux.Router, svc *productService.ProductService) {
		log.Println("Registering product routes")
		xmux.Register(r, http.MethodPost, "/api/products", svc.CreateProduct)
		xmux.Register(r, http.MethodGet, "/api/products/:id", svc.GetProduct)
		xmux.Register(r, http.MethodGet, "/api/products", func(ctx context.Context, _ *struct{}) ([]*productModel.ProductResponse, error) {
			return svc.ListProducts(ctx)
		})
		xmux.Register(r, http.MethodPut, "/api/products/:id", svc.UpdateProduct)
		xmux.Register(r, http.MethodDelete, "/api/products/:id", func(ctx context.Context, req *productModel.DeleteProductRequest) (any, error) {
			return nil, svc.DeleteProduct(ctx, req)
		})
	})

	orderGroup := xmux.DefineGroup(func(r xmux.Router, svc *orderService.OrderService) {
		log.Println("Registering order routes")
		xmux.Register(r, http.MethodPost, "/api/orders", svc.CreateOrder)
		xmux.Register(r, http.MethodGet, "/api/orders/:id", svc.GetOrder)
	})

	groups := xmux.NewGroups(userGroup, productGroup, orderGroup)

	if err := groups.Bind(ctrl, bindService); err != nil {
		log.Printf("Error binding routes: %v", err)
	} else {
		log.Println("All routes registered successfully")
	}
}
