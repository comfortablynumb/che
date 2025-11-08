package chehttp

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"

	"github.com/comfortablynumb/che/pkg/chectx"
)

// RequestIDKey is the context key for the request ID
var RequestIDKey = chectx.Key[string]("X-Request-ID")

// RequestIDConfig configures the request ID middleware
type RequestIDConfig struct {
	// HeaderName is the name of the header to read/write the request ID
	// Defaults to "X-Request-ID"
	HeaderName string

	// Generator is a function that generates a new request ID
	// Defaults to generateRequestID
	Generator func() string
}

// DefaultRequestIDConfig returns a sensible default configuration
func DefaultRequestIDConfig() *RequestIDConfig {
	return &RequestIDConfig{
		HeaderName: "X-Request-ID",
		Generator:  generateRequestID,
	}
}

// RequestIDMiddleware returns a middleware that extracts or generates a request ID
// and stores it in the context. It also sets the request ID in the response header.
func RequestIDMiddleware(config *RequestIDConfig) func(http.Handler) http.Handler {
	if config == nil {
		config = DefaultRequestIDConfig()
	}

	if config.HeaderName == "" {
		config.HeaderName = "X-Request-ID"
	}

	if config.Generator == nil {
		config.Generator = generateRequestID
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Try to get request ID from header
			requestID := r.Header.Get(config.HeaderName)

			// Generate a new one if not present
			if requestID == "" {
				requestID = config.Generator()
			}

			// Store in context
			ctx := chectx.WithValue(r.Context(), RequestIDKey, requestID)

			// Set response header
			w.Header().Set(config.HeaderName, requestID)

			// Call next handler with updated context
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetRequestID retrieves the request ID from the context
// Returns empty string if not found
func GetRequestID(r *http.Request) string {
	return chectx.GetOrDefault(r.Context(), RequestIDKey, "")
}

// generateRequestID generates a random request ID
func generateRequestID() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		// Fallback to a simple counter if random fails
		return "fallback-id"
	}
	return hex.EncodeToString(bytes)
}
