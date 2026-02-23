package main

import (
	"context"
	"fmt"
	"log"

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
	router := NewEchoRouter(nil)

	// Register a route using xmux.Register
	xmux.Register(router, "GET", "/hello",
		func(ctx context.Context, params *HelloParams) (*HelloResponse, error) {
			return &HelloResponse{
				Message: fmt.Sprintf("Hello, %s!", params.Name),
			}, nil
		})

	// Start the server
	fmt.Println("Server listening on :8080")
	if err := router.Start(":8080"); err != nil {
		log.Fatal(err)
	}
}
