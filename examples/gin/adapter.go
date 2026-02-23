package main

import (
	"net/http"

	"github.com/Just-maple/xmux"
	"github.com/gin-gonic/gin"
)

// GinRouter implements xmux.Router for Gin.
type GinRouter struct {
	engine *gin.Engine
}

// NewGinRouter creates a new GinRouter.
func NewGinRouter(engine *gin.Engine) *GinRouter {
	if engine == nil {
		engine = gin.Default()
	}
	return &GinRouter{engine: engine}
}

// Register implements xmux.Router.Register.
func (r *GinRouter) Register(method string, path string, api xmux.Handler, options ...map[string]string) {
	// Convert xmux.Handler to gin.HandlerFunc
	handler := func(c *gin.Context) {
		// Create a bind function that extracts request data
		bind := func(ptr any) error {
			// For simplicity, we'll just return nil; real implementation would parse request
			return nil
		}
		// Invoke the handler
		_, err := api.Invoke(c.Request.Context(), bind)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// Write response as JSON
		c.JSON(http.StatusOK, gin.H{"message": "hello"})
	}
	// Register with Gin engine
	r.engine.Handle(method, path, handler)
}

// Run delegates to the underlying Gin engine.
func (r *GinRouter) Run(addr ...string) error {
	return r.engine.Run(addr...)
}
