package main

import (
	"context"
	"log"
	"os"

	"github.com/Just-maple/xmux/examples/gin-business/internal/adapter"
	"github.com/Just-maple/xmux/examples/gin-business/internal/business"
	"github.com/Just-maple/xmux/examples/gin-business/internal/repository"
	"github.com/Just-maple/xmux/examples/gin-business/internal/types"
)

func main() {
	// Initialize dependencies
	deps := setupDependencies()

	// Create router with configuration
	routerConfig := adapter.DefaultRouterConfig()
	routerConfig.DebugMode = os.Getenv("DEBUG") == "true"

	router := adapter.NewGinRouter(nil, routerConfig)

	// Register all routes
	adapter.RegisterAllRoutes(router, deps)

	// Start server
	port := getPort()
	log.Printf("Server starting on port %s", port)
	log.Printf("API version: %s", routerConfig.APIVersion)
	log.Printf("Debug mode: %v", routerConfig.DebugMode)

	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

// setupDependencies initializes all application dependencies
func setupDependencies() adapter.Dependencies {
	// Create repository
	userRepo := repository.NewInMemoryUserRepository()

	// Create business services
	userService := business.NewUserService(userRepo)

	// Initialize with some sample data
	initSampleData(userService)

	return adapter.Dependencies{
		UserService: userService,
	}
}

// initSampleData initializes the application with some sample data
func initSampleData(userService business.UserService) {
	ctx := context.Background()

	// Create admin user
	_, err := userService.CreateUser(ctx, &types.CreateUserRequest{
		Username: "admin",
		Email:    "admin@example.com",
		Password: "Admin123!",
		FullName: "System Administrator",
		Role:     "admin",
	})
	if err != nil && err.Error() != "user already exists" {
		log.Printf("Failed to create admin user: %v", err)
	}

	// Create regular user
	_, err = userService.CreateUser(ctx, &types.CreateUserRequest{
		Username: "user",
		Email:    "user@example.com",
		Password: "User123!",
		FullName: "Regular User",
		Role:     "user",
	})
	if err != nil && err.Error() != "user already exists" {
		log.Printf("Failed to create regular user: %v", err)
	}
}

// getPort gets the port from environment or defaults to 8080
func getPort() string {
	if port := os.Getenv("PORT"); port != "" {
		return port
	}
	return "8080"
}
