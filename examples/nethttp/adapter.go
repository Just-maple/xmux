package main

import (
	"net/http"

	"github.com/Just-maple/xmux"
)

// NetHTTPRouter implements xmux.Router for net/http.
type NetHTTPRouter struct {
	mux *http.ServeMux
}

// NewNetHTTPRouter creates a new NetHTTPRouter.
func NewNetHTTPRouter(mux *http.ServeMux) *NetHTTPRouter {
	if mux == nil {
		mux = http.NewServeMux()
	}
	return &NetHTTPRouter{mux: mux}
}

// Register implements xmux.Router.Register.
func (r *NetHTTPRouter) Register(method string, path string, api xmux.Handler, options ...map[string]string) {
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
		// For now just write a placeholder
		w.Write([]byte(`{"message":"hello"}`))
	}
	// Register with the ServeMux (note: ServeMux doesn't distinguish methods)
	r.mux.HandleFunc(path, handler)
}

// ServeHTTP delegates to the underlying ServeMux.
func (r *NetHTTPRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}
