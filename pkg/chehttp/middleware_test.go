package chehttp

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/comfortablynumb/che/pkg/chetest"
)

func TestRequestIDMiddleware_WithExistingID(t *testing.T) {
	middleware := RequestIDMiddleware(nil)

	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check that request ID is in context
		requestID := GetRequestID(r)
		chetest.RequireEqual(t, requestID, "existing-id-123")

		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("X-Request-ID", "existing-id-123")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	// Check response header
	chetest.RequireEqual(t, rec.Header().Get("X-Request-ID"), "existing-id-123")
	chetest.RequireEqual(t, rec.Code, http.StatusOK)
}

func TestRequestIDMiddleware_GenerateNewID(t *testing.T) {
	middleware := RequestIDMiddleware(nil)

	var capturedRequestID string
	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check that request ID is generated and in context
		requestID := GetRequestID(r)
		chetest.RequireEqual(t, requestID != "", true)
		chetest.RequireEqual(t, len(requestID) > 0, true)
		capturedRequestID = requestID

		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	// Check response header contains generated ID
	chetest.RequireEqual(t, rec.Header().Get("X-Request-ID"), capturedRequestID)
	chetest.RequireEqual(t, rec.Code, http.StatusOK)
}

func TestRequestIDMiddleware_CustomHeaderName(t *testing.T) {
	config := &RequestIDConfig{
		HeaderName: "X-Custom-Request-ID",
	}
	middleware := RequestIDMiddleware(config)

	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check that request ID is in context
		requestID := GetRequestID(r)
		chetest.RequireEqual(t, requestID, "custom-id-456")

		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("X-Custom-Request-ID", "custom-id-456")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	// Check custom header in response
	chetest.RequireEqual(t, rec.Header().Get("X-Custom-Request-ID"), "custom-id-456")
	chetest.RequireEqual(t, rec.Code, http.StatusOK)
}

func TestRequestIDMiddleware_CustomGenerator(t *testing.T) {
	callCount := 0
	config := &RequestIDConfig{
		Generator: func() string {
			callCount++
			return "generated-id-789"
		},
	}
	middleware := RequestIDMiddleware(config)

	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := GetRequestID(r)
		chetest.RequireEqual(t, requestID, "generated-id-789")

		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	// Generator should have been called once
	chetest.RequireEqual(t, callCount, 1)
	chetest.RequireEqual(t, rec.Header().Get("X-Request-ID"), "generated-id-789")
}

func TestGetRequestID_NotFound(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/test", nil)

	// Should return empty string when not found
	requestID := GetRequestID(req)
	chetest.RequireEqual(t, requestID, "")
}

func TestGetRequestID_Found(t *testing.T) {
	middleware := RequestIDMiddleware(nil)

	var retrievedID string
	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		retrievedID = GetRequestID(r)
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("X-Request-ID", "test-id-999")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	chetest.RequireEqual(t, retrievedID, "test-id-999")
}

func TestGenerateRequestID(t *testing.T) {
	// Generate multiple IDs and ensure they're unique
	id1 := generateRequestID()
	id2 := generateRequestID()
	id3 := generateRequestID()

	chetest.RequireEqual(t, id1 != "", true)
	chetest.RequireEqual(t, id2 != "", true)
	chetest.RequireEqual(t, id3 != "", true)

	// IDs should be different
	chetest.RequireEqual(t, id1 != id2, true)
	chetest.RequireEqual(t, id2 != id3, true)
	chetest.RequireEqual(t, id1 != id3, true)

	// IDs should be hex encoded (32 characters for 16 bytes)
	chetest.RequireEqual(t, len(id1), 32)
	chetest.RequireEqual(t, len(id2), 32)
	chetest.RequireEqual(t, len(id3), 32)
}

func TestRequestIDMiddleware_ChainedMiddlewares(t *testing.T) {
	requestIDMiddleware := RequestIDMiddleware(nil)

	// Simulate another middleware that uses the request ID
	loggingMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// This middleware should have access to the request ID
			requestID := GetRequestID(r)
			chetest.RequireEqual(t, requestID != "", true)

			next.ServeHTTP(w, r)
		})
	}

	handler := requestIDMiddleware(loggingMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	chetest.RequireEqual(t, rec.Code, http.StatusOK)
	chetest.RequireEqual(t, rec.Header().Get("X-Request-ID") != "", true)
}
