package response

import (
	"errors"
	"fmt"
	"testing"

	pkgerrors "github.com/gravitee-io-labs/gravitee-automation-sdks/common/pkg/errors"
)

type apiError struct {
	Message *string
}

type getResp struct {
	status      int
	Body        []byte
	JSON200     *string
	JSON403     *apiError
	JSONDefault *apiError
}

func (r getResp) StatusCode() int { return r.status }
func (r getResp) GetBody() []byte { return r.Body }

type deleteResp struct {
	status      int
	Body        []byte
	JSON403     *apiError
	JSONDefault *apiError
}

func (r deleteResp) StatusCode() int { return r.status }
func (r deleteResp) GetBody() []byte { return r.Body }

func TestRespond_SuccessJSON200(t *testing.T) {
	value := "example-domain"
	got, err := Respond[string](getResp{status: 200, JSON200: &value}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == nil || *got != value {
		t.Fatalf("got %#v, want %q", got, value)
	}
}

func TestRespond_NoContent(t *testing.T) {
	got, err := Respond[struct{}](deleteResp{status: 204}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != nil {
		t.Fatalf("got %#v, want nil", got)
	}
}

func TestRespond_StatusErrorMessage(t *testing.T) {
	msg := "permission denied"
	_, err := Respond[string](getResp{status: 403, JSON403: &apiError{Message: &msg}}, nil)
	var httpErr pkgerrors.HttpError
	if !errors.As(err, &httpErr) {
		t.Fatalf("got %T %v, want HttpError", err, err)
	}
	if httpErr.Status != 403 || httpErr.Body != msg {
		t.Fatalf("got %+v, want status 403 body %q", httpErr, msg)
	}
}

func TestRespond_DefaultErrorThenBody(t *testing.T) {
	msg := "boom"
	_, err := Respond[string](getResp{status: 500, Body: []byte("raw"), JSONDefault: &apiError{Message: &msg}}, nil)
	var httpErr pkgerrors.HttpError
	if !errors.As(err, &httpErr) {
		t.Fatalf("got %T %v, want HttpError", err, err)
	}
	if httpErr.Status != 500 || httpErr.Body != msg {
		t.Fatalf("got %+v, want status 500 body %q", httpErr, msg)
	}

	_, err = Respond[string](getResp{status: 404, Body: []byte("missing")}, nil)
	if !errors.As(err, &httpErr) {
		t.Fatalf("got %T %v, want HttpError", err, err)
	}
	if httpErr.Status != 404 || httpErr.Body != "missing" {
		t.Fatalf("got %+v, want status 404 body missing", httpErr)
	}
}

func TestRespond_TransportError(t *testing.T) {
	_, err := Respond[string](getResp{}, fmt.Errorf("dial tcp"))
	var clientErr pkgerrors.ClientError
	if !errors.As(err, &clientErr) {
		t.Fatalf("got %T %v, want ClientError", err, err)
	}
}

func TestRespond_SuccessTypeMismatch(t *testing.T) {
	value := "example-domain"
	_, err := Respond[int](getResp{status: 200, JSON200: &value}, nil)
	if err == nil {
		t.Fatal("expected type mismatch error")
	}
}
