package main

import (
	"encoding/json"
	"net/http"

	"github.com/Just-maple/xmux"
)

// Controller adapts net/http.ServeMux to xmux.Controller interface.
type Controller struct {
	mux *http.ServeMux
}

// NewController creates a new net/http controller.
func NewController() *Controller {
	return &Controller{
		mux: http.NewServeMux(),
	}
}

// Handle implements xmux.Controller interface.
func (c *Controller) Handle(method, path string, service any, api xmux.Api, opts ...map[string]string) {
	c.mux.HandleFunc(path, func(w http.ResponseWriter, req *http.Request) {
		// Check HTTP method
		if req.Method != method {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

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
	})
}

// ServeHTTP implements http.Handler interface.
func (c *Controller) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c.mux.ServeHTTP(w, req)
}
