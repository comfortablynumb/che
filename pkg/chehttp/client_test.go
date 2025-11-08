package chehttp

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
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
		json.NewEncoder(w).Encode(testResponse{Message: "success", Code: 200})
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
		json.NewDecoder(r.Body).Decode(&body)
		chetest.RequireEqual(t, body.Message, "test")

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(testResponse{Message: "created", Code: 201})
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
		json.NewEncoder(w).Encode(testResponse{Message: "success", Code: 200})
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
		json.NewEncoder(w).Encode(testError{Error: "bad request", Details: "invalid input"})
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
		w.Write([]byte("test response"))
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
		json.NewDecoder(r.Body).Decode(&body)

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
