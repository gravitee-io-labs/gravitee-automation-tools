package errors

import (
	"fmt"
)

type ClientError struct {
	err error
}

func NewClientError(err error) ClientError {
	return ClientError{err: err}
}

func (e ClientError) Error() string {
	if e.err == nil {
		return "unknown client error"
	}
	return fmt.Sprintf("client error: %s", e.err.Error())
}

func (e ClientError) Unwrap() error { return e.err }

type HttpError struct {
	Status int
	Body   string
}

func (e HttpError) StatusCode() int {
	return e.Status
}

func (e HttpError) Error() string {
	return e.Body
}

func (e HttpError) Has(s int) bool {
	return e.Status == s
}
