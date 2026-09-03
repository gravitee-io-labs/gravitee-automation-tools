package response

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"

	"github.com/gravitee-io-labs/gravitee-automation-sdks/common/pkg/errors"
)

var jsonStatusField = regexp.MustCompile(`^JSON(\d{3})$`)

type statusCoder interface {
	StatusCode() int
}

type bodyGetter interface {
	GetBody() []byte
}

// Respond maps a generated *WithResponse result to T.
//
// On 2xx it returns the JSON2xx field whose type is *T.
// On 204 (no JSON2xx field) it returns (nil, nil).
// On error status it reads JSON{status}, then JSONDefault, then the raw body.
func Respond[T any](resp any, err error) (*T, error) {
	if err != nil {
		return nil, errors.NewClientError(err)
	}
	if resp == nil {
		return nil, errors.NewClientError(fmt.Errorf("nil response"))
	}

	sc, ok := resp.(statusCoder)
	if !ok {
		return nil, errors.NewClientError(fmt.Errorf("response %T has no StatusCode", resp))
	}

	v, err := structValue(resp)
	if err != nil {
		return nil, err
	}

	status := sc.StatusCode()
	if status < 300 {
		return successEntity[T](v)
	}
	return nil, httpError(v, status, rawBody(resp))
}

func successEntity[T any](v reflect.Value) (*T, error) {
	want := reflect.TypeFor[*T]()
	found2xx := false
	for _, f := range jsonFields(v) {
		if f.code < 200 || f.code >= 300 {
			continue
		}
		found2xx = true
		if !f.value.Type().AssignableTo(want) {
			continue
		}
		if f.value.IsNil() {
			return nil, nil
		}
		return f.value.Interface().(*T), nil
	}
	if !found2xx {
		return nil, nil
	}
	return nil, errors.NewClientError(fmt.Errorf("no JSON2xx field of type %s", want))
}

func httpError(v reflect.Value, status int, body []byte) error {
	msg := string(body)
	if f, ok := fieldForStatus(v, status); ok && !f.IsNil() {
		if m := messageOf(f); m != "" {
			msg = m
		}
	} else if d := v.FieldByName("JSONDefault"); d.IsValid() && d.Kind() == reflect.Pointer && !d.IsNil() {
		if m := messageOf(d); m != "" {
			msg = m
		}
	}
	return errors.HttpError{Status: status, Body: msg}
}

type jsonField struct {
	code  int
	value reflect.Value
}

func jsonFields(v reflect.Value) []jsonField {
	t := v.Type()
	fields := make([]jsonField, 0, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		matches := jsonStatusField.FindStringSubmatch(t.Field(i).Name)
		if matches == nil {
			continue
		}
		code, _ := strconv.Atoi(matches[1])
		fields = append(fields, jsonField{code: code, value: v.Field(i)})
	}
	return fields
}

func fieldForStatus(v reflect.Value, status int) (reflect.Value, bool) {
	f := v.FieldByName("JSON" + strconv.Itoa(status))
	if !f.IsValid() || f.Kind() != reflect.Pointer {
		return reflect.Value{}, false
	}
	return f, true
}

func messageOf(errVal reflect.Value) string {
	if errVal.Kind() == reflect.Pointer {
		if errVal.IsNil() {
			return ""
		}
		errVal = errVal.Elem()
	}
	if errVal.Kind() != reflect.Struct {
		return ""
	}
	m := errVal.FieldByName("Message")
	if !m.IsValid() {
		return ""
	}
	if m.Kind() == reflect.Pointer {
		if m.IsNil() {
			return ""
		}
		m = m.Elem()
	}
	if m.Kind() != reflect.String {
		return ""
	}
	return m.String()
}

func structValue(resp any) (reflect.Value, error) {
	v := reflect.ValueOf(resp)
	for v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return reflect.Value{}, errors.NewClientError(fmt.Errorf("nil response"))
		}
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return reflect.Value{}, errors.NewClientError(fmt.Errorf("response %T is not a struct", resp))
	}
	return v, nil
}

func rawBody(resp any) []byte {
	if b, ok := resp.(bodyGetter); ok {
		return b.GetBody()
	}
	return nil
}
