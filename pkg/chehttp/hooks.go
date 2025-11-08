package chehttp

import (
	"net/http"
	"time"
)

// HookContext contains information passed to hooks
type HookContext struct {
	// Request information
	Method  string
	URL     string
	Headers http.Header

	// Response information (available in post-request hooks)
	StatusCode int
	Response   Response

	// Error information (available if error occurred)
	Error error

	// Timing information
	StartTime time.Time
	Duration  time.Duration
}

// PreRequestHook is called before sending the request
// It can modify the request or cancel it by returning an error
type PreRequestHook func(ctx *HookContext) error

// PostRequestHook is called after receiving the response
type PostRequestHook func(ctx *HookContext)

// SuccessHook is called when the response status is 2xx
type SuccessHook func(ctx *HookContext)

// ErrorHook is called when the response status is 4xx or 5xx
type ErrorHook func(ctx *HookContext)

// CompleteHook is called after the request completes (success or failure)
type CompleteHook func(ctx *HookContext)

// Hooks contains all the lifecycle hooks for HTTP requests
type Hooks struct {
	PreRequest  []PreRequestHook
	PostRequest []PostRequestHook
	OnSuccess   []SuccessHook
	OnError     []ErrorHook
	OnComplete  []CompleteHook
}

// Clone creates a copy of the hooks
func (h *Hooks) Clone() *Hooks {
	if h == nil {
		return &Hooks{}
	}

	clone := &Hooks{
		PreRequest:  make([]PreRequestHook, len(h.PreRequest)),
		PostRequest: make([]PostRequestHook, len(h.PostRequest)),
		OnSuccess:   make([]SuccessHook, len(h.OnSuccess)),
		OnError:     make([]ErrorHook, len(h.OnError)),
		OnComplete:  make([]CompleteHook, len(h.OnComplete)),
	}

	copy(clone.PreRequest, h.PreRequest)
	copy(clone.PostRequest, h.PostRequest)
	copy(clone.OnSuccess, h.OnSuccess)
	copy(clone.OnError, h.OnError)
	copy(clone.OnComplete, h.OnComplete)

	return clone
}
