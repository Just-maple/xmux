package main

import (
	"log"
	"os"

	"github.com/Just-maple/xmux/examples/gin-business/internal/adapter"
)

func main() {
	deps := adapter.NewContainer()
	adapter.InitSampleData(deps)

	routerConfig := adapter.DefaultRouterConfig()
	routerConfig.DebugMode = os.Getenv("DEBUG") == "false"

	router := adapter.NewGinRouter(nil, routerConfig)

	if err := adapter.RegisterAllRoutes(router, deps); err != nil {
		log.Fatal("Failed to register routes:", err)
	}

	port := getPort()
	log.Printf("Server starting on port %s", port)
	log.Printf("API version: %s", routerConfig.APIVersion)
	log.Printf("Debug mode: %v", routerConfig.DebugMode)

	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func getPort() string {
	if port := os.Getenv("PORT"); port != "" {
		return port
	}
	return "8080"
}
