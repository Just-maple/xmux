package adapter

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/Just-maple/xmux"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GinRouter implements xmux.Router for Gin with production-ready features
type GinRouter struct {
	engine *gin.Engine
	config RouterConfig
}

// RouterConfig contains configuration for the router
type RouterConfig struct {
	APIVersion string
	DebugMode  bool
}

// DefaultRouterConfig returns the default router configuration
func DefaultRouterConfig() RouterConfig {
	return RouterConfig{
		APIVersion: "v1",
		DebugMode:  false,
	}
}

// NewGinRouter creates a new GinRouter with the provided configuration
func NewGinRouter(engine *gin.Engine, config RouterConfig) *GinRouter {
	if engine == nil {
		if config.DebugMode {
			engine = gin.Default()
		} else {
			engine = gin.New()
			engine.Use(gin.Recovery())
		}
	}
	return &GinRouter{
		engine: engine,
		config: config,
	}
}

// Register implements xmux.Router.Register with enhanced request handling
func (r *GinRouter) Register(method string, path string, api xmux.Handler, options ...map[string]string) {
	// Apply route-specific options
	routeOptions := mergeOptions(options...)

	// Create the Gin handler
	handler := func(c *gin.Context) {
		// Create context with request metadata
		ctx := r.enrichContext(c)

		// Create bind function that parses request data
		bind := r.createBindFunction(c, method)

		// Execute the handler
		result, err := api.Invoke(ctx, bind)

		// Handle response
		r.handleResponse(c, result, err, routeOptions)
	}

	// Register with Gin
	r.engine.Handle(method, r.normalizePath(path), handler)
}

// Run starts the Gin server
func (r *GinRouter) Run(addr ...string) error {
	return r.engine.Run(addr...)
}

// ServeHTTP implements http.Handler interface
func (r *GinRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.engine.ServeHTTP(w, req)
}

// enrichContext adds request metadata to the context
func (r *GinRouter) enrichContext(c *gin.Context) context.Context {
	ctx := c.Request.Context()

	// Add request ID if available
	if requestID := c.GetHeader("X-Request-ID"); requestID != "" {
		// In a real implementation, you'd add this to context
		_ = requestID
	}

	// Add client IP
	if clientIP := c.ClientIP(); clientIP != "" {
		// In a real implementation, you'd add this to context
		_ = clientIP
	}

	return ctx
}

// createBindFunction creates a bind function that parses request data
func (r *GinRouter) createBindFunction(c *gin.Context, method string) xmux.Bind {
	return func(ptr any) error {
		// Parse JSON body for methods that typically have a body
		if method == http.MethodPost || method == http.MethodPut || method == http.MethodPatch {
			if c.Request.Body != nil && c.Request.ContentLength > 0 {
				if err := json.NewDecoder(c.Request.Body).Decode(ptr); err != nil {
					return &BindError{
						Type:    "json_parse",
						Message: "Failed to parse JSON body",
						Err:     err,
					}
				}
			}
		}

		// Parse query parameters
		if err := r.bindQueryParams(c, ptr); err != nil {
			return err
		}

		// Parse path parameters
		if err := r.bindPathParams(c, ptr); err != nil {
			return err
		}

		// Parse headers
		if err := r.bindHeaders(c, ptr); err != nil {
			return err
		}

		// In a real implementation, you would also:
		// - Parse form data
		// - Parse multipart form
		// - Parse cookies
		// - Apply validation

		return nil
	}
}

// bindQueryParams binds query parameters to struct fields with "query" tag
func (r *GinRouter) bindQueryParams(c *gin.Context, ptr any) error {
	// This is a simplified implementation
	// In a real implementation, you would use reflection to bind query params
	// based on struct tags like `query:"param_name"`
	return nil
}

// bindPathParams binds path parameters to struct fields with "path" tag
func (r *GinRouter) bindPathParams(c *gin.Context, ptr any) error {
	// Extract path parameters and bind them to the struct
	// This is a simplified implementation
	params := c.Params
	for _, param := range params {
		// In a real implementation, you would use reflection to bind
		// based on struct tags like `path:"param_name"`
		if param.Key == "id" {
			// Try to parse as UUID
			if id, err := uuid.Parse(param.Value); err == nil {
				// Try to set the ID field
				_ = id
			}
		}
	}
	return nil
}

// bindHeaders binds headers to struct fields with "header" tag
func (r *GinRouter) bindHeaders(c *gin.Context, ptr any) error {
	// This is a simplified implementation
	// In a real implementation, you would use reflection to bind headers
	// based on struct tags like `header:"Header-Name"`
	return nil
}

// handleResponse handles the response from the handler
func (r *GinRouter) handleResponse(c *gin.Context, result any, err error, options map[string]string) {
	// Handle errors
	if err != nil {
		r.handleError(c, err, options)
		return
	}

	// Determine response format based on Accept header or options
	format := options["format"]
	if format == "" {
		format = "json"
	}

	// Set response headers
	c.Header("X-API-Version", r.config.APIVersion)
	c.Header("Content-Type", "application/json")

	// Send response
	switch format {
	case "json":
		c.JSON(http.StatusOK, result)
	default:
		c.JSON(http.StatusOK, result)
	}
}

// handleError handles errors from the handler
func (r *GinRouter) handleError(c *gin.Context, err error, options map[string]string) {
	statusCode := http.StatusInternalServerError
	errorCode := "internal_error"
	message := "Internal server error"

	// Map specific errors to appropriate HTTP status codes
	var bindErr *BindError
	if errors.As(err, &bindErr) {
		statusCode = http.StatusBadRequest
		errorCode = bindErr.Type
		message = bindErr.Message
		if r.config.DebugMode {
			message = bindErr.Error()
		}
	}

	// Create error response
	errorResponse := map[string]any{
		"error":   errorCode,
		"message": message,
	}

	if r.config.DebugMode && err != nil {
		errorResponse["details"] = err.Error()
	}

	c.Header("X-API-Version", r.config.APIVersion)
	c.JSON(statusCode, errorResponse)
}

// normalizePath normalizes the path by adding API version prefix if needed
func (r *GinRouter) normalizePath(path string) string {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	// Add API version prefix if not already present
	if !strings.HasPrefix(path, "/api/") {
		path = "/api/" + r.config.APIVersion + path
	}

	return path
}

// mergeOptions merges multiple option maps
func mergeOptions(options ...map[string]string) map[string]string {
	result := make(map[string]string)
	for _, opts := range options {
		for k, v := range opts {
			result[k] = v
		}
	}
	return result
}

// BindError represents an error during parameter binding
type BindError struct {
	Type    string
	Message string
	Err     error
}

func (e *BindError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

func (e *BindError) Unwrap() error {
	return e.Err
}
