package chehttp

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/comfortablynumb/che/pkg/chetest"
)

type testResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type testError struct {
	Error   string `json:"error"`
	Details string `json:"details"`
}

func TestBuilder_Build(t *testing.T) {
	client := NewBuilder().
		WithBaseURL("https://api.example.com").
		WithDefaultHeader("User-Agent", "test").
		WithDefaultTimeout(30 * time.Second).
		Build()

	chetest.RequireEqual(t, client != nil, true)
}

func TestClient_Get(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		chetest.RequireEqual(t, r.Method, http.MethodGet)
		chetest.RequireEqual(t, r.URL.Path, "/test")

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(testResponse{Message: "success", Code: 200})
	}))
	defer server.Close()

	client := NewBuilder().WithBaseURL(server.URL).Build()

	resp, err := client.Get("/test")
	chetest.RequireEqual(t, err, nil)
	chetest.RequireEqual(t, resp.StatusCode(), http.StatusOK)
	chetest.RequireEqual(t, resp.IsSuccess(), true)
}

func TestClient_Post(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		chetest.RequireEqual(t, r.Method, http.MethodPost)
		chetest.RequireEqual(t, r.Header.Get("Content-Type"), "application/json")

		var body testResponse
		_ = json.NewDecoder(r.Body).Decode(&body)
		chetest.RequireEqual(t, body.Message, "test")

		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(testResponse{Message: "created", Code: 201})
	}))
	defer server.Close()

	client := NewBuilder().WithBaseURL(server.URL).Build()

	resp, err := client.Post("/test", WithJSONBody(testResponse{Message: "test"}))
	chetest.RequireEqual(t, err, nil)
	chetest.RequireEqual(t, resp.StatusCode(), http.StatusCreated)
}

func TestClient_Put(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		chetest.RequireEqual(t, r.Method, http.MethodPut)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewBuilder().WithBaseURL(server.URL).Build()

	resp, err := client.Put("/test")
	chetest.RequireEqual(t, err, nil)
	chetest.RequireEqual(t, resp.StatusCode(), http.StatusOK)
}

func TestClient_Patch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		chetest.RequireEqual(t, r.Method, http.MethodPatch)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewBuilder().WithBaseURL(server.URL).Build()

	resp, err := client.Patch("/test")
	chetest.RequireEqual(t, err, nil)
	chetest.RequireEqual(t, resp.StatusCode(), http.StatusOK)
}

func TestClient_Delete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		chetest.RequireEqual(t, r.Method, http.MethodDelete)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := NewBuilder().WithBaseURL(server.URL).Build()

	resp, err := client.Delete("/test")
	chetest.RequireEqual(t, err, nil)
	chetest.RequireEqual(t, resp.StatusCode(), http.StatusNoContent)
}

func TestClient_WithHeaders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		chetest.RequireEqual(t, r.Header.Get("X-Custom-Header"), "custom-value")
		chetest.RequireEqual(t, r.Header.Get("Authorization"), "Bearer token")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewBuilder().WithBaseURL(server.URL).Build()

	resp, err := client.Get("/test",
		WithHeader("X-Custom-Header", "custom-value"),
		WithHeader("Authorization", "Bearer token"),
	)
	chetest.RequireEqual(t, err, nil)
	chetest.RequireEqual(t, resp.StatusCode(), http.StatusOK)
}

func TestClient_WithDefaultHeaders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		chetest.RequireEqual(t, r.Header.Get("User-Agent"), "test-agent")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewBuilder().
		WithBaseURL(server.URL).
		WithDefaultHeader("User-Agent", "test-agent").
		Build()

	resp, err := client.Get("/test")
	chetest.RequireEqual(t, err, nil)
	chetest.RequireEqual(t, resp.StatusCode(), http.StatusOK)
}

func TestClient_AutoUnmarshalSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(testResponse{Message: "success", Code: 200})
	}))
	defer server.Close()

	client := NewBuilder().WithBaseURL(server.URL).Build()

	var successResp testResponse
	resp, err := client.Get("/test", WithSuccess(&successResp))

	chetest.RequireEqual(t, err, nil)
	chetest.RequireEqual(t, resp.StatusCode(), http.StatusOK)
	chetest.RequireEqual(t, successResp.Message, "success")
	chetest.RequireEqual(t, successResp.Code, 200)
}

func TestClient_AutoUnmarshalError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(testError{Error: "bad request", Details: "invalid input"})
	}))
	defer server.Close()

	client := NewBuilder().WithBaseURL(server.URL).Build()

	var errorResp testError
	resp, err := client.Get("/test", WithError(&errorResp))

	chetest.RequireEqual(t, err, nil)
	chetest.RequireEqual(t, resp.StatusCode(), http.StatusBadRequest)
	chetest.RequireEqual(t, resp.IsError(), true)
	chetest.RequireEqual(t, errorResp.Error, "bad request")
	chetest.RequireEqual(t, errorResp.Details, "invalid input")
}

func TestClient_Timeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewBuilder().WithBaseURL(server.URL).Build()

	_, err := client.Get("/test", WithTimeout(10*time.Millisecond))
	chetest.RequireEqual(t, err != nil, true)
}

func TestResponse_String(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("test response"))
	}))
	defer server.Close()

	client := NewBuilder().WithBaseURL(server.URL).Build()

	resp, err := client.Get("/test")
	chetest.RequireEqual(t, err, nil)
	chetest.RequireEqual(t, resp.String(), "test response")
}

func TestResponse_IsSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewBuilder().WithBaseURL(server.URL).Build()

	resp, err := client.Get("/test")
	chetest.RequireEqual(t, err, nil)
	chetest.RequireEqual(t, resp.IsSuccess(), true)
	chetest.RequireEqual(t, resp.IsError(), false)
}

func TestResponse_IsError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := NewBuilder().WithBaseURL(server.URL).Build()

	resp, err := client.Get("/test")
	chetest.RequireEqual(t, err, nil)
	chetest.RequireEqual(t, resp.IsSuccess(), false)
	chetest.RequireEqual(t, resp.IsError(), true)
}

func TestResponse_Headers(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Custom-Header", "custom-value")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewBuilder().WithBaseURL(server.URL).Build()

	resp, err := client.Get("/test")
	chetest.RequireEqual(t, err, nil)
	chetest.RequireEqual(t, resp.Headers().Get("X-Custom-Header"), "custom-value")
}

func TestBuilder_WithMultipleHeaders(t *testing.T) {
	headers := map[string]string{
		"X-Header-1": "value1",
		"X-Header-2": "value2",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		chetest.RequireEqual(t, r.Header.Get("X-Header-1"), "value1")
		chetest.RequireEqual(t, r.Header.Get("X-Header-2"), "value2")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewBuilder().
		WithBaseURL(server.URL).
		WithDefaultHeaders(headers).
		Build()

	resp, err := client.Get("/test")
	chetest.RequireEqual(t, err, nil)
	chetest.RequireEqual(t, resp.StatusCode(), http.StatusOK)
}

func TestClient_WithJSONBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body testResponse
		_ = json.NewDecoder(r.Body).Decode(&body)

		chetest.RequireEqual(t, body.Message, "test message")
		chetest.RequireEqual(t, body.Code, 123)

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewBuilder().WithBaseURL(server.URL).Build()

	resp, err := client.Post("/test", WithJSONBody(testResponse{
		Message: "test message",
		Code:    123,
	}))

	chetest.RequireEqual(t, err, nil)
	chetest.RequireEqual(t, resp.StatusCode(), http.StatusOK)
}

func TestClient_PreRequestHook(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	var hookCalled bool
	var hookMethod string
	var hookURL string

	client := NewBuilder().
		WithBaseURL(server.URL).
		WithPreRequestHook(func(ctx *HookContext) error {
			hookCalled = true
			hookMethod = ctx.Method
			hookURL = ctx.URL
			return nil
		}).
		Build()

	resp, err := client.Get("/test")
	chetest.RequireEqual(t, err, nil)
	chetest.RequireEqual(t, resp.StatusCode(), http.StatusOK)
	chetest.RequireEqual(t, hookCalled, true)
	chetest.RequireEqual(t, hookMethod, http.MethodGet)
	chetest.RequireEqual(t, hookURL, server.URL+"/test")
}

func TestClient_PostRequestHook(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	var hookCalled bool
	var hookStatusCode int

	client := NewBuilder().
		WithBaseURL(server.URL).
		WithPostRequestHook(func(ctx *HookContext) {
			hookCalled = true
			hookStatusCode = ctx.StatusCode
		}).
		Build()

	resp, err := client.Get("/test")
	chetest.RequireEqual(t, err, nil)
	chetest.RequireEqual(t, resp.StatusCode(), http.StatusOK)
	chetest.RequireEqual(t, hookCalled, true)
	chetest.RequireEqual(t, hookStatusCode, http.StatusOK)
}

func TestClient_SuccessHook(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	var successCalled bool
	var errorCalled bool

	client := NewBuilder().
		WithBaseURL(server.URL).
		WithSuccessHook(func(ctx *HookContext) {
			successCalled = true
		}).
		WithErrorHook(func(ctx *HookContext) {
			errorCalled = true
		}).
		Build()

	resp, err := client.Get("/test")
	chetest.RequireEqual(t, err, nil)
	chetest.RequireEqual(t, resp.StatusCode(), http.StatusOK)
	chetest.RequireEqual(t, successCalled, true)
	chetest.RequireEqual(t, errorCalled, false)
}

func TestClient_ErrorHook(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer server.Close()

	var successCalled bool
	var errorCalled bool

	client := NewBuilder().
		WithBaseURL(server.URL).
		WithSuccessHook(func(ctx *HookContext) {
			successCalled = true
		}).
		WithErrorHook(func(ctx *HookContext) {
			errorCalled = true
		}).
		Build()

	resp, err := client.Get("/test")
	chetest.RequireEqual(t, err, nil)
	chetest.RequireEqual(t, resp.StatusCode(), http.StatusBadRequest)
	chetest.RequireEqual(t, successCalled, false)
	chetest.RequireEqual(t, errorCalled, true)
}

func TestClient_CompleteHook(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	var completeCalled bool
	var duration time.Duration

	client := NewBuilder().
		WithBaseURL(server.URL).
		WithCompleteHook(func(ctx *HookContext) {
			completeCalled = true
			duration = ctx.Duration
		}).
		Build()

	resp, err := client.Get("/test")
	chetest.RequireEqual(t, err, nil)
	chetest.RequireEqual(t, resp.StatusCode(), http.StatusOK)
	chetest.RequireEqual(t, completeCalled, true)
	chetest.RequireEqual(t, duration > 0, true)
}

func TestClient_PreRequestHookError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	var completeCalled bool

	client := NewBuilder().
		WithBaseURL(server.URL).
		WithPreRequestHook(func(ctx *HookContext) error {
			return fmt.Errorf("hook error")
		}).
		WithCompleteHook(func(ctx *HookContext) {
			completeCalled = true
		}).
		Build()

	_, err := client.Get("/test")
	chetest.RequireEqual(t, err != nil, true)
	chetest.RequireEqual(t, completeCalled, true)
}

func TestClient_RequestTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewBuilder().
		WithBaseURL(server.URL).
		WithRequestTimeout(10 * time.Millisecond).
		Build()

	_, err := client.Get("/test")
	chetest.RequireEqual(t, err != nil, true)
}

func TestClient_MultipleHooks(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	var hook1Called bool
	var hook2Called bool

	client := NewBuilder().
		WithBaseURL(server.URL).
		WithPreRequestHook(func(ctx *HookContext) error {
			hook1Called = true
			return nil
		}).
		WithPreRequestHook(func(ctx *HookContext) error {
			hook2Called = true
			return nil
		}).
		Build()

	resp, err := client.Get("/test")
	chetest.RequireEqual(t, err, nil)
	chetest.RequireEqual(t, resp.StatusCode(), http.StatusOK)
	chetest.RequireEqual(t, hook1Called, true)
	chetest.RequireEqual(t, hook2Called, true)
}

// Test retry with exponential backoff
func TestClient_RetryWithExponentialBackoff(t *testing.T) {
	var attemptCount int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count := atomic.AddInt32(&attemptCount, 1)
		if count < 3 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(testResponse{Message: "success", Code: 200})
	}))
	defer server.Close()

	client := NewBuilder().
		WithBaseURL(server.URL).
		WithRetries(3).
		WithRetryBackoff(ExponentialBackoff{
			BaseDelay:  10 * time.Millisecond,
			Multiplier: 2.0,
			MaxDelay:   1 * time.Second,
		}).
		Build()

	start := time.Now()
	resp, err := client.Get("/test")
	duration := time.Since(start)

	chetest.RequireEqual(t, err, nil)
	chetest.RequireEqual(t, resp.StatusCode(), http.StatusOK)
	chetest.RequireEqual(t, atomic.LoadInt32(&attemptCount), int32(3))
	// Should have some backoff delay (at least 10ms + 20ms = 30ms)
	chetest.RequireEqual(t, duration > 30*time.Millisecond, true)
}

// Test retry with fixed backoff
func TestClient_RetryWithFixedBackoff(t *testing.T) {
	var attemptCount int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count := atomic.AddInt32(&attemptCount, 1)
		if count < 2 {
			w.WriteHeader(http.StatusBadGateway)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewBuilder().
		WithBaseURL(server.URL).
		WithRetries(2).
		WithRetryBackoff(FixedBackoff{Delay: 50 * time.Millisecond}).
		Build()

	resp, err := client.Get("/test")
	chetest.RequireEqual(t, err, nil)
	chetest.RequireEqual(t, resp.StatusCode(), http.StatusOK)
	chetest.RequireEqual(t, atomic.LoadInt32(&attemptCount), int32(2))
}

// Test retry exhausted
func TestClient_RetryExhausted(t *testing.T) {
	var attemptCount int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&attemptCount, 1)
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := NewBuilder().
		WithBaseURL(server.URL).
		WithRetries(2).
		WithRetryBackoff(FixedBackoff{Delay: 10 * time.Millisecond}).
		Build()

	resp, err := client.Get("/test")
	chetest.RequireEqual(t, err, nil)
	chetest.RequireEqual(t, resp.StatusCode(), http.StatusInternalServerError)
	chetest.RequireEqual(t, atomic.LoadInt32(&attemptCount), int32(3)) // initial + 2 retries
}

// Test non-retryable status codes
func TestClient_NonRetryableStatusCode(t *testing.T) {
	var attemptCount int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&attemptCount, 1)
		w.WriteHeader(http.StatusNotImplemented) // 501 is non-retryable by default
	}))
	defer server.Close()

	client := NewBuilder().
		WithBaseURL(server.URL).
		WithRetries(3).
		WithRetryBackoff(FixedBackoff{Delay: 10 * time.Millisecond}).
		Build()

	resp, err := client.Get("/test")
	chetest.RequireEqual(t, err, nil)
	chetest.RequireEqual(t, resp.StatusCode(), http.StatusNotImplemented)
	chetest.RequireEqual(t, atomic.LoadInt32(&attemptCount), int32(1)) // No retries
}

// Test custom retryable status codes
func TestClient_CustomRetryableStatusCodes(t *testing.T) {
	var attemptCount int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count := atomic.AddInt32(&attemptCount, 1)
		if count < 2 {
			w.WriteHeader(http.StatusBadRequest) // 400
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewBuilder().
		WithBaseURL(server.URL).
		WithRetries(2).
		WithRetryableStatusCodes(
			StatusCodeRange{Min: 400, Max: 400}, // Retry on 400
		).
		WithRetryBackoff(FixedBackoff{Delay: 10 * time.Millisecond}).
		Build()

	resp, err := client.Get("/test")
	chetest.RequireEqual(t, err, nil)
	chetest.RequireEqual(t, resp.StatusCode(), http.StatusOK)
	chetest.RequireEqual(t, atomic.LoadInt32(&attemptCount), int32(2))
}

// Test context cancellation with GetWithCtx
func TestClient_ContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewBuilder().WithBaseURL(server.URL).Build()

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	_, err := client.GetWithCtx(ctx, "/test")
	chetest.RequireEqual(t, err != nil, true)
}

// Test PostWithCtx
func TestClient_PostWithCtx(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		chetest.RequireEqual(t, r.Method, http.MethodPost)
		w.WriteHeader(http.StatusCreated)
	}))
	defer server.Close()

	client := NewBuilder().WithBaseURL(server.URL).Build()

	resp, err := client.PostWithCtx(context.Background(), "/test")
	chetest.RequireEqual(t, err, nil)
	chetest.RequireEqual(t, resp.StatusCode(), http.StatusCreated)
}

// Test PutWithCtx
func TestClient_PutWithCtx(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		chetest.RequireEqual(t, r.Method, http.MethodPut)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewBuilder().WithBaseURL(server.URL).Build()

	resp, err := client.PutWithCtx(context.Background(), "/test")
	chetest.RequireEqual(t, err, nil)
	chetest.RequireEqual(t, resp.StatusCode(), http.StatusOK)
}

// Test DeleteWithCtx
func TestClient_DeleteWithCtx(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		chetest.RequireEqual(t, r.Method, http.MethodDelete)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := NewBuilder().WithBaseURL(server.URL).Build()

	resp, err := client.DeleteWithCtx(context.Background(), "/test")
	chetest.RequireEqual(t, err, nil)
	chetest.RequireEqual(t, resp.StatusCode(), http.StatusNoContent)
}

// Test response body reader
func TestResponse_BodyReader(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("streaming response"))
	}))
	defer server.Close()

	client := NewBuilder().WithBaseURL(server.URL).Build()

	resp, err := client.Get("/test")
	chetest.RequireEqual(t, err, nil)

	// Use body reader directly
	bodyReader := resp.BodyReader()
	chetest.RequireEqual(t, bodyReader != nil, true)

	data, err := io.ReadAll(bodyReader)
	chetest.RequireEqual(t, err, nil)
	chetest.RequireEqual(t, string(data), "streaming response")
	_ = bodyReader.Close()
}

// Test body reader is not read when not using auto-unmarshal
func TestResponse_BodyReaderNotAutoRead(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("test data"))
	}))
	defer server.Close()

	client := NewBuilder().WithBaseURL(server.URL).Build()

	resp, err := client.Get("/test")
	chetest.RequireEqual(t, err, nil)

	// Body should not be read automatically
	// We can verify by using BodyReader
	bodyReader := resp.BodyReader()
	chetest.RequireEqual(t, bodyReader != nil, true)
}

// Test linear backoff
func TestClient_LinearBackoff(t *testing.T) {
	var attemptCount int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count := atomic.AddInt32(&attemptCount, 1)
		if count < 3 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewBuilder().
		WithBaseURL(server.URL).
		WithRetries(3).
		WithRetryBackoff(LinearBackoff{BaseDelay: 20 * time.Millisecond}).
		Build()

	start := time.Now()
	resp, err := client.Get("/test")
	duration := time.Since(start)

	chetest.RequireEqual(t, err, nil)
	chetest.RequireEqual(t, resp.StatusCode(), http.StatusOK)
	chetest.RequireEqual(t, atomic.LoadInt32(&attemptCount), int32(3))
	// Linear backoff: 20ms + 40ms = 60ms minimum
	chetest.RequireEqual(t, duration > 60*time.Millisecond, true)
}

// Test retry with context cancellation
func TestClient_RetryWithContextCancellation(t *testing.T) {
	var attemptCount int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&attemptCount, 1)
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := NewBuilder().
		WithBaseURL(server.URL).
		WithRetries(5).
		WithRetryBackoff(FixedBackoff{Delay: 100 * time.Millisecond}).
		Build()

	ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
	defer cancel()

	_, err := client.GetWithCtx(ctx, "/test")
	chetest.RequireEqual(t, err != nil, true)
	// Should have attempted at least once, maybe twice depending on timing
	chetest.RequireEqual(t, atomic.LoadInt32(&attemptCount) >= 1, true)
}
