package main

import (
	"net/http"

	"github.com/Just-maple/xmux"
	"github.com/gin-gonic/gin"
)

// Controller adapts Gin to xmux.Controller interface.
type Controller struct {
	engine *gin.Engine
}

// NewController creates a new Gin controller.
func NewController() *Controller {
	return &Controller{
		engine: gin.Default(),
	}
}

// Handle implements xmux.Controller interface.
func (c *Controller) Handle(method, path string, service any, api xmux.Api, opts ...map[string]string) {
	c.engine.Handle(method, path, func(ctx *gin.Context) {
		// Create bind function to parse request
		bind := func(ptr any) error {
			if err := ctx.ShouldBindJSON(ptr); err != nil {
				return ctx.ShouldBindQuery(ptr)
			}
			return nil
		}

		// Execute business logic
		result, err := api.Invoke(ctx.Request.Context(), bind)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Send response
		ctx.JSON(http.StatusOK, result)
	})
}

// ServeHTTP implements http.Handler interface.
func (c *Controller) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c.engine.ServeHTTP(w, req)
}

// Use adds middleware to the controller.
func (c *Controller) Use(middleware ...gin.HandlerFunc) {
	c.engine.Use(middleware...)
}
