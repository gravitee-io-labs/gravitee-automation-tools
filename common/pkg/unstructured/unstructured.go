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

// Package unstructured holds free-form JSON objects, mirroring the shape of
// Kubernetes' unstructured.Unstructured without depending on apimachinery.
package unstructured

import (
	"encoding/json"

	"gopkg.in/yaml.v3"
)

// Unstructured is a free-form object serialized as a JSON/YAML tree:
// {"data":{"key":"value"}}.
type Unstructured struct {
	Object map[string]any `json:"-"`
}

// New returns an empty object.
func New() *Unstructured {
	return &Unstructured{Object: map[string]any{}}
}

// From wraps obj without copying it.
func From(obj map[string]any) *Unstructured {
	return &Unstructured{Object: obj}
}

func (in *Unstructured) Put(key string, value any) *Unstructured {
	if in.Object == nil {
		in.Object = map[string]any{}
	}
	in.Object[key] = value
	return in
}

func (in *Unstructured) Get(key string) any {
	return in.Object[key]
}

func (in *Unstructured) Remove(key string) {
	delete(in.Object, key)
}

func (in *Unstructured) GetString(key string) string {
	s, _ := in.Object[key].(string)
	return s
}

func (in *Unstructured) GetBool(key string) bool {
	b, _ := in.Object[key].(bool)
	return b
}

func (in *Unstructured) GetSlice(key string) []any {
	s, _ := in.Object[key].([]any)
	return s
}

// IsZero lets `omitzero` (JSON) and `omitempty` (YAML) skip empty objects.
func (in Unstructured) IsZero() bool {
	return len(in.Object) == 0
}

func (in Unstructured) MarshalJSON() ([]byte, error) {
	return json.Marshal(in.Object)
}

func (in *Unstructured) UnmarshalJSON(data []byte) error {
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	in.Object = m
	return nil
}

func (in Unstructured) MarshalYAML() (any, error) {
	return in.Object, nil
}

func (in *Unstructured) UnmarshalYAML(node *yaml.Node) error {
	var m map[string]any
	if err := node.Decode(&m); err != nil {
		return err
	}
	in.Object = m
	return nil
}
