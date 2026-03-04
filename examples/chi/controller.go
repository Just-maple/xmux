package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/Just-maple/xmux"
)

// Controller adapts Chi to xmux.Controller interface.
type Controller struct {
	mux *chi.Mux
}

// NewController creates a new Chi controller.
func NewController() *Controller {
	return &Controller{
		mux: chi.NewMux(),
	}
}

// Handle implements xmux.Controller interface.
func (c *Controller) Handle(method, path string, service any, api xmux.Api, opts ...map[string]string) {
	c.mux.Method(method, path, http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// Create bind function to parse request
		bind := func(ptr any) error {
			if req.Body == nil {
				return nil
			}
			return json.NewDecoder(req.Body).Decode(ptr)
		}

		// Execute business logic
		result, err := api.Invoke(req.Context(), bind)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Send response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	}))
}

// ServeHTTP implements http.Handler interface.
func (c *Controller) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c.mux.ServeHTTP(w, req)
}

// Use adds middleware to the controller.
func (c *Controller) Use(middleware ...func(http.Handler) http.Handler) {
	for _, m := range middleware {
		c.mux.Use(m)
	}
}
