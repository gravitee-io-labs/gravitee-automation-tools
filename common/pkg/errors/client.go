// Copyright (C) 2015 The Gravitee team (http://gravitee.io)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//         http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package errors is the SDK error vocabulary: client construction, auth config, and HTTP status.
package errors

import (
	"fmt"
)

// ClientError wraps a construction or config failure. It is not an HTTP status. Unwrap returns the cause.
type ClientError struct {
	err error
}

// NewClientError wraps err. A nil err still produces a ClientError whose Error is "unknown client error".
func NewClientError(err error) ClientError {
	return ClientError{err: err}
}

func (e ClientError) Error() string {
	if e.err == nil {
		return "unknown client error"
	}
	return fmt.Sprintf("client error: %s", e.err.Error())
}

// Unwrap returns the wrapped cause for errors.Is / errors.As.
func (e ClientError) Unwrap() error { return e.err }

// HttpError is a completed HTTP response that the caller treats as failure.
// Status is the status code; Body is the raw response body, also used as Error().
type HttpError struct {
	Status int
	Body   string
}

// StatusCode returns Status.
func (e HttpError) StatusCode() int {
	return e.Status
}

func (e HttpError) Error() string {
	return e.Body
}

// Has reports whether Status equals s.
func (e HttpError) Has(s int) bool {
	return e.Status == s
}
