package response

import (
	"fmt"
	"net/http"
	"testing"
)

type code int

func (c code) StatusCode() int { return int(c) }

type getResp struct {
	status  int
	JSON200 *string
}

func (r getResp) StatusCode() int { return r.status }

func TestIsNetworkError(t *testing.T) {
	if !IsNetworkError(nil, fmt.Errorf("dial tcp")) {
		t.Fatal("want true for transport error")
	}
	if IsNetworkError(code(http.StatusOK), nil) {
		t.Fatal("want false when err is nil")
	}
}

func TestIsHTTPStatus(t *testing.T) {
	tests := []struct {
		name string
		sc   StatusCoder
		fn   func(StatusCoder, error) bool
		want bool
	}{
		{name: "server 500", sc: code(http.StatusInternalServerError), fn: IsServerError, want: true},
		{name: "server 503", sc: code(http.StatusServiceUnavailable), fn: IsServerError, want: true},
		{name: "server 404", sc: code(http.StatusNotFound), fn: IsServerError, want: false},
		{name: "not found", sc: code(http.StatusNotFound), fn: IsNotFound, want: true},
		{name: "not found 200", sc: code(http.StatusOK), fn: IsNotFound, want: false},
		{name: "forbidden", sc: code(http.StatusForbidden), fn: IsForbidden, want: true},
		{name: "unauthorized", sc: code(http.StatusUnauthorized), fn: IsUnauthorized, want: true},
		{name: "nil coder", sc: nil, fn: IsNotFound, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fn(tt.sc, nil); got != tt.want {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsRetryable(t *testing.T) {
	if !IsRetryable(nil, fmt.Errorf("dial tcp")) {
		t.Fatal("network should be retryable")
	}
	if !IsRetryable(code(http.StatusInternalServerError), nil) {
		t.Fatal("server error should be retryable")
	}
	if !IsRetryable(code(http.StatusForbidden), nil) {
		t.Fatal("forbidden should be retryable")
	}
	if IsRetryable(code(http.StatusNotFound), nil) {
		t.Fatal("not found should not be retryable")
	}
	if IsRetryable(code(http.StatusUnauthorized), nil) {
		t.Fatal("unauthorized should not be retryable")
	}
}

func TestPayload(t *testing.T) {
	value := "example-domain"
	got, ok := Payload[string](getResp{status: 200, JSON200: &value})
	if !ok || got != value {
		t.Fatalf("got %q %v, want %q true", got, ok, value)
	}

	ptr, ok := Payload[*string](getResp{status: 200, JSON200: &value})
	if !ok || ptr == nil || *ptr != value {
		t.Fatalf("got %#v %v, want pointer to %q", ptr, ok, value)
	}

	if _, ok := Payload[int](getResp{status: 200, JSON200: &value}); ok {
		t.Fatal("want false on type mismatch")
	}
	if _, ok := Payload[string](getResp{status: 200}); ok {
		t.Fatal("want false on nil JSON200")
	}
	if _, ok := Payload[string](struct{ status int }{status: 204}); ok {
		t.Fatal("want false when JSON200 is missing")
	}
}
