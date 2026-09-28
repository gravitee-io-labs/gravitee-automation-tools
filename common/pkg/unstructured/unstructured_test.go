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
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

type treeDTO struct {
	Name string       `json:"name" yaml:"name"`
	Data Unstructured `json:"data" yaml:"data"`
}

type treeDTOPtr struct {
	Name string        `json:"name" yaml:"name"`
	Data *Unstructured `json:"data,omitempty" yaml:"data,omitempty"`
}

type treeDTOOmit struct {
	Name string       `json:"name" yaml:"name"`
	Data Unstructured `json:"data,omitzero" yaml:"data,omitempty"`
}

func TestUnstructuredJSON(t *testing.T) {
	t.Run("marshals as a tree", func(t *testing.T) {
		got, err := json.Marshal(treeDTO{Name: "test", Data: *From(map[string]any{"key": "value"})})
		require.NoError(t, err)
		assert.Equal(t, `{"name":"test","data":{"key":"value"}}`, string(got))
	})

	t.Run("marshals through a pointer field", func(t *testing.T) {
		got, err := json.Marshal(treeDTOPtr{Name: "test", Data: From(map[string]any{"key": "value"})})
		require.NoError(t, err)
		assert.Equal(t, `{"name":"test","data":{"key":"value"}}`, string(got))
	})

	t.Run("marshals nested values", func(t *testing.T) {
		data := From(map[string]any{
			"nested": map[string]any{"list": []any{1, "two", true, nil}},
		})
		got, err := json.Marshal(treeDTO{Name: "test", Data: *data})
		require.NoError(t, err)
		assert.JSONEq(t, `{"name":"test","data":{"nested":{"list":[1,"two",true,null]}}}`, string(got))
	})

	t.Run("marshals nil object as null", func(t *testing.T) {
		got, err := json.Marshal(treeDTO{Name: "test"})
		require.NoError(t, err)
		assert.Equal(t, `{"name":"test","data":null}`, string(got))
	})

	t.Run("omitzero skips nil and empty objects", func(t *testing.T) {
		for _, data := range []Unstructured{{}, *New()} {
			got, err := json.Marshal(treeDTOOmit{Name: "test", Data: data})
			require.NoError(t, err)
			assert.Equal(t, `{"name":"test"}`, string(got))
		}
	})

	t.Run("unmarshals a tree", func(t *testing.T) {
		var got treeDTO
		require.NoError(t, json.Unmarshal([]byte(`{"name":"test","data":{"key":"value","n":1.5,"obj":{"a":[1]}}}`), &got))
		assert.Equal(t, "test", got.Name)
		assert.Equal(t, map[string]any{"key": "value", "n": 1.5, "obj": map[string]any{"a": []any{float64(1)}}}, got.Data.Object)
	})

	t.Run("unmarshals null as nil object", func(t *testing.T) {
		var got treeDTO
		require.NoError(t, json.Unmarshal([]byte(`{"name":"test","data":null}`), &got))
		assert.Nil(t, got.Data.Object)
	})

	t.Run("rejects non-object values", func(t *testing.T) {
		for _, data := range []string{`"{}"`, `[]`, `42`} {
			var got treeDTO
			assert.Error(t, json.Unmarshal([]byte(`{"data":`+data+`}`), &got), data)
		}
	})

	t.Run("round-trips", func(t *testing.T) {
		in := treeDTO{Name: "test", Data: *From(map[string]any{"key": "value", "list": []any{"a", map[string]any{"b": true}}})}
		b, err := json.Marshal(in)
		require.NoError(t, err)
		var out treeDTO
		require.NoError(t, json.Unmarshal(b, &out))
		assert.Equal(t, in, out)
	})
}

func TestUnstructuredYAML(t *testing.T) {
	t.Run("marshals as a tree", func(t *testing.T) {
		got, err := yaml.Marshal(treeDTO{Name: "test", Data: *From(map[string]any{"key": "value"})})
		require.NoError(t, err)
		assert.Equal(t, "name: test\ndata:\n    key: value\n", string(got))
	})

	t.Run("marshals through a pointer field", func(t *testing.T) {
		got, err := yaml.Marshal(treeDTOPtr{Name: "test", Data: From(map[string]any{"key": "value"})})
		require.NoError(t, err)
		assert.Equal(t, "name: test\ndata:\n    key: value\n", string(got))
	})

	t.Run("omitempty skips nil and empty objects", func(t *testing.T) {
		for _, data := range []Unstructured{{}, *New()} {
			got, err := yaml.Marshal(treeDTOOmit{Name: "test", Data: data})
			require.NoError(t, err)
			assert.Equal(t, "name: test\n", string(got))
		}
	})

	t.Run("unmarshals a tree", func(t *testing.T) {
		var got treeDTO
		require.NoError(t, yaml.Unmarshal([]byte("name: test\ndata:\n  key: value\n  obj:\n    a: [1]\n"), &got))
		assert.Equal(t, map[string]any{"key": "value", "obj": map[string]any{"a": []any{1}}}, got.Data.Object)
	})

	t.Run("unmarshals JSON-flavoured YAML", func(t *testing.T) {
		var got treeDTO
		require.NoError(t, yaml.Unmarshal([]byte(`{"name":"test","data":{"key":"value"}}`), &got))
		assert.Equal(t, map[string]any{"key": "value"}, got.Data.Object)
	})

	t.Run("rejects non-object values", func(t *testing.T) {
		var got treeDTO
		assert.Error(t, yaml.Unmarshal([]byte("data: [a, b]\n"), &got))
	})

	t.Run("round-trips", func(t *testing.T) {
		in := treeDTO{Name: "test", Data: *From(map[string]any{"key": "value", "list": []any{"a", map[string]any{"b": true}}})}
		b, err := yaml.Marshal(in)
		require.NoError(t, err)
		var out treeDTO
		require.NoError(t, yaml.Unmarshal(b, &out))
		assert.Equal(t, in, out)
	})
}

func TestUnstructuredAccessors(t *testing.T) {
	u := New().
		Put("s", "str").
		Put("b", true).
		Put("l", []any{"x"})

	assert.Equal(t, "str", u.GetString("s"))
	assert.Equal(t, "", u.GetString("b"), "wrong type yields zero value")
	assert.True(t, u.GetBool("b"))
	assert.False(t, u.GetBool("missing"))
	assert.Equal(t, []any{"x"}, u.GetSlice("l"))
	assert.Nil(t, u.GetSlice("s"))
	assert.Nil(t, u.Get("missing"))

	u.Remove("s")
	assert.Nil(t, u.Get("s"))

	var zero Unstructured
	zero.Put("k", "v")
	assert.Equal(t, "v", zero.GetString("k"), "Put initialises a nil object")
}
