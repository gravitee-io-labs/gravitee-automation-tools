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

package response

import (
	"net/http"
	"reflect"
)

type StatusCoder interface {
	StatusCode() int
}

func IsNetworkError(_ StatusCoder, err error) bool {
	return err != nil
}

func IsServerError(sc StatusCoder, _ error) bool {
	c := statusCode(sc)
	return c >= http.StatusInternalServerError && c <= 599
}

func IsNotFound(sc StatusCoder, _ error) bool {
	return statusCode(sc) == http.StatusNotFound
}

func IsForbidden(sc StatusCoder, _ error) bool {
	return statusCode(sc) == http.StatusForbidden
}

func IsUnauthorized(sc StatusCoder, _ error) bool {
	return statusCode(sc) == http.StatusUnauthorized
}

func IsRetryable(sc StatusCoder, err error) bool {
	return IsNetworkError(sc, err) || IsServerError(sc, err) || IsForbidden(sc, err)
}

// Payload returns the JSON200 value when it is assignable to T.
func Payload[T any](resp any) (T, bool) {
	var zero T
	v, ok := structValue(resp)
	if !ok {
		return zero, false
	}
	f := v.FieldByName("JSON200")
	if !f.IsValid() {
		return zero, false
	}
	want := reflect.TypeFor[T]()
	if val, ok := assignable[T](f, want); ok {
		return val, true
	}
	if f.Kind() == reflect.Pointer && !f.IsNil() {
		return assignable[T](f.Elem(), want)
	}
	return zero, false
}

func statusCode(sc StatusCoder) int {
	if sc == nil {
		return 0
	}
	v := reflect.ValueOf(sc)
	if v.Kind() == reflect.Pointer && v.IsNil() {
		return 0
	}
	return sc.StatusCode()
}

func structValue(resp any) (reflect.Value, bool) {
	if resp == nil {
		return reflect.Value{}, false
	}
	v := reflect.ValueOf(resp)
	for v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return reflect.Value{}, false
		}
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return reflect.Value{}, false
	}
	return v, true
}

func assignable[T any](v reflect.Value, want reflect.Type) (T, bool) {
	var zero T
	if !v.IsValid() || !v.Type().AssignableTo(want) {
		return zero, false
	}
	return v.Interface().(T), true
}
