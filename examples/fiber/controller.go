package main

import (
	"net/http"

	"github.com/Just-maple/xmux"
	"github.com/gofiber/fiber/v2"
)

// Controller adapts Fiber to xmux.Controller interface.
type Controller struct {
	app *fiber.App
}

// NewController creates a new Fiber controller.
func NewController() *Controller {
	return &Controller{
		app: fiber.New(),
	}
}

// Handle implements xmux.Controller interface.
func (c *Controller) Handle(method, path string, api xmux.Api, opts ...map[string]string) {
	c.app.Add(method, path, func(ctx *fiber.Ctx) error {
		// Create bind function to parse request
		bind := func(ptr any) error {
			return ctx.BodyParser(ptr)
		}

		// Execute business logic
		result, err := api.Invoke(ctx.Context(), bind)
		if err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(map[string]string{"error": err.Error()})
		}

		// Send response
		return ctx.JSON(result)
	})
}

// ServeHTTP implements http.Handler interface.
func (c *Controller) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// Fiber doesn't directly support http.Handler interface
	// This is a simplified implementation for testing
	c.app.Test(req, -1)
}
