package controller

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Just-maple/xmux"
)

type Controller struct {
	engine *gin.Engine
}

func NewController() *Controller {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(gin.Recovery())
	return &Controller{
		engine: engine,
	}
}

func (c *Controller) Handle(method, path string, api xmux.Api, options ...map[string]string) {
	c.engine.Handle(method, path, func(ctx *gin.Context) {
		bind := func(ptr any) error {
			if ctx.Request.Body == nil {
				return nil
			}
			if err := ctx.ShouldBindJSON(ptr); err != nil {
				return ctx.ShouldBindQuery(ptr)
			}
			return nil
		}

		result, err := api.Invoke(ctx.Request.Context(), bind)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, result)
	})
}

func (c *Controller) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c.engine.ServeHTTP(w, req)
}

func (c *Controller) Shutdown(ctx context.Context) error {
	return nil
}
