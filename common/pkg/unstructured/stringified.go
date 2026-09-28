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

package unstructured

import (
	"bytes"
	"encoding/json"

	"gopkg.in/yaml.v3"
)

// Stringified is an Unstructured serialized as a JSON-encoded string:
// {"data":"{\"key\":\"value\"}"}. It also reads a plain object tree, so it can
// be mapped from an Unstructured. An empty string decodes to a nil object.
type Stringified struct {
	Unstructured
}

// StringifiedFrom wraps obj without copying it.
func StringifiedFrom(obj map[string]any) *Stringified {
	return &Stringified{Unstructured{Object: obj}}
}

func (in Stringified) MarshalJSON() ([]byte, error) {
	if in.Object == nil {
		return []byte("null"), nil
	}
	s, err := in.encode()
	if err != nil {
		return nil, err
	}
	return json.Marshal(s)
}

func (in *Stringified) UnmarshalJSON(data []byte) error {
	if isJSONObject(data) {
		return in.Unstructured.UnmarshalJSON(data)
	}
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	return in.decode(s)
}

func (in Stringified) MarshalYAML() (any, error) {
	if in.Object == nil {
		return nil, nil
	}
	return in.encode()
}

func (in *Stringified) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind == yaml.MappingNode {
		return in.Unstructured.UnmarshalYAML(node)
	}
	var s string
	if err := node.Decode(&s); err != nil {
		return err
	}
	return in.decode(s)
}

func (in *Stringified) encode() (string, error) {
	b, err := json.Marshal(in.Object)
	return string(b), err
}

func (in *Stringified) decode(s string) error {
	if s == "" {
		in.Object = nil
		return nil
	}
	return in.Unstructured.UnmarshalJSON([]byte(s))
}

func isJSONObject(data []byte) bool {
	return bytes.HasPrefix(bytes.TrimSpace(data), []byte("{"))
}
