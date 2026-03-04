package main

import (
	"net/http"

	"github.com/Just-maple/xmux"
	"github.com/labstack/echo/v4"
)

// Controller adapts Echo to xmux.Controller interface.
type Controller struct {
	engine *echo.Echo
}

// NewController creates a new Echo controller.
func NewController() *Controller {
	return &Controller{
		engine: echo.New(),
	}
}

// Handle implements xmux.Controller interface.
func (c *Controller) Handle(method, path string, service any, api xmux.Api, opts ...map[string]string) {
	c.engine.Add(method, path, func(ctx echo.Context) error {
		// Create bind function to parse request
		bind := func(ptr any) error {
			return ctx.Bind(ptr)
		}

		// Execute business logic
		result, err := api.Invoke(ctx.Request().Context(), bind)
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}

		// Send response
		return ctx.JSON(http.StatusOK, result)
	})
}

// ServeHTTP implements http.Handler interface.
func (c *Controller) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c.engine.ServeHTTP(w, req)
}

// Use adds middleware to the controller.
func (c *Controller) Use(middleware ...echo.MiddlewareFunc) {
	c.engine.Use(middleware...)
}
