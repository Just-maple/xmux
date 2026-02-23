package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/Just-maple/xmux"
)

type HelloParams struct {
	Name string `json:"name"`
}

type HelloResponse struct {
	Message string `json:"message"`
}

func main() {
	// Create a new router adapter
	router := NewChiRouter(nil)

	// Register a route using xmux.Register
	xmux.Register(router, http.MethodGet, "/hello",
		func(ctx context.Context, params *HelloParams) (*HelloResponse, error) {
			return &HelloResponse{
				Message: fmt.Sprintf("Hello, %s!", params.Name),
			}, nil
		})

	// Start the server
	fmt.Println("Server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
