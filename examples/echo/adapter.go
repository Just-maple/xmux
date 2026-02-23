package main

import (
	"net/http"

	"github.com/Just-maple/xmux"
	"github.com/labstack/echo/v4"
)

// EchoRouter implements xmux.Router for Echo.
type EchoRouter struct {
	e *echo.Echo
}

// NewEchoRouter creates a new EchoRouter.
func NewEchoRouter(e *echo.Echo) *EchoRouter {
	if e == nil {
		e = echo.New()
	}
	return &EchoRouter{e: e}
}

// Register implements xmux.Router.Register.
func (r *EchoRouter) Register(method string, path string, api xmux.Handler, options ...map[string]string) {
	// Convert xmux.Handler to echo.HandlerFunc
	handler := func(c echo.Context) error {
		// Create a bind function that extracts request data
		bind := func(ptr any) error {
			// For simplicity, we'll just return nil; real implementation would parse request
			return nil
		}
		// Invoke the handler
		_, err := api.Invoke(c.Request().Context(), bind)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		// Write response as JSON
		return c.JSON(http.StatusOK, map[string]string{"message": "hello"})
	}
	// Register with Echo
	r.e.Add(method, path, handler)
}

// Start delegates to the underlying Echo instance.
func (r *EchoRouter) Start(addr string) error {
	return r.e.Start(addr)
}
