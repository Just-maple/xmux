package main

import (
	"net/http"

	"github.com/Just-maple/xmux"
	"github.com/gorilla/mux"
)

// GorillaRouter implements xmux.Router for Gorilla/mux.
type GorillaRouter struct {
	router *mux.Router
}

// NewGorillaRouter creates a new GorillaRouter.
func NewGorillaRouter(router *mux.Router) *GorillaRouter {
	if router == nil {
		router = mux.NewRouter()
	}
	return &GorillaRouter{router: router}
}

// Register implements xmux.Router.Register.
func (r *GorillaRouter) Register(method string, path string, api xmux.Handler, options ...map[string]string) {
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
	// Register with Gorilla router
	r.router.HandleFunc(path, handler).Methods(method)
}

// ServeHTTP delegates to the underlying Gorilla router.
func (r *GorillaRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.router.ServeHTTP(w, req)
}
