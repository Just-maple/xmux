package main

import (
	"net/http"

	"github.com/Just-maple/xmux"
	"github.com/go-chi/chi/v5"
)

// ChiRouter implements xmux.Router for Chi.
type ChiRouter struct {
	router *chi.Mux
}

// NewChiRouter creates a new ChiRouter.
func NewChiRouter(router *chi.Mux) *ChiRouter {
	if router == nil {
		router = chi.NewRouter()
	}
	return &ChiRouter{router: router}
}

// Register implements xmux.Router.Register.
func (r *ChiRouter) Register(method string, path string, api xmux.Handler, options ...map[string]string) {
	// Convert xmux.Handler to http.HandlerFunc
	handler := func(w http.ResponseWriter, req *http.Request) {
		// Create a bind function that extracts request data
		bind := func(ptr any) error {
			// For simplicity, we'll just return nil; real implementation would parse request
			return nil
		}
		// Invoke the handler
		_, err := api.Invoke(req.Context(), bind)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Write response as JSON
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message":"hello"}`))
	}
	// Register with Chi router
	r.router.MethodFunc(method, path, handler)
}

// ServeHTTP delegates to the underlying Chi router.
func (r *ChiRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.router.ServeHTTP(w, req)
}
