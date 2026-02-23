package main

import (
	"github.com/Just-maple/xmux"
	"github.com/gofiber/fiber/v2"
)

// FiberRouter implements xmux.Router for Fiber.
type FiberRouter struct {
	app *fiber.App
}

// NewFiberRouter creates a new FiberRouter.
func NewFiberRouter(app *fiber.App) *FiberRouter {
	if app == nil {
		app = fiber.New()
	}
	return &FiberRouter{app: app}
}

// Register implements xmux.Router.Register.
func (r *FiberRouter) Register(method string, path string, api xmux.Handler, options ...map[string]string) {
	// Convert xmux.Handler to fiber.Handler
	handler := func(c *fiber.Ctx) error {
		// Create a bind function that extracts request data
		bind := func(ptr any) error {
			// For simplicity, we'll just return nil; real implementation would parse request
			return nil
		}
		// Invoke the handler
		_, err := api.Invoke(c.Context(), bind)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		// Write response as JSON
		return c.JSON(fiber.Map{"message": "hello"})
	}
	// Register with Fiber app
	r.app.Add(method, path, handler)
}

// Listen delegates to the underlying Fiber app.
func (r *FiberRouter) Listen(addr string) error {
	return r.app.Listen(addr)
}
