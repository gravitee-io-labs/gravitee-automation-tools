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

type stringDTO struct {
	Name string      `json:"name" yaml:"name"`
	Data Stringified `json:"data" yaml:"data"`
}

type stringDTOPtr struct {
	Name string       `json:"name" yaml:"name"`
	Data *Stringified `json:"data,omitempty" yaml:"data,omitempty"`
}

type stringDTOOmit struct {
	Name string      `json:"name" yaml:"name"`
	Data Stringified `json:"data,omitzero" yaml:"data,omitempty"`
}

func TestStringifiedJSON(t *testing.T) {
	t.Run("marshals as an escaped JSON string", func(t *testing.T) {
		got, err := json.Marshal(stringDTO{Name: "test", Data: *StringifiedFrom(map[string]any{"key": "value"})})
		require.NoError(t, err)
		assert.Equal(t, `{"name":"test","data":"{\"key\":\"value\"}"}`, string(got))
	})

	t.Run("marshals through a pointer field", func(t *testing.T) {
		got, err := json.Marshal(stringDTOPtr{Name: "test", Data: StringifiedFrom(map[string]any{"key": "value"})})
		require.NoError(t, err)
		assert.Equal(t, `{"name":"test","data":"{\"key\":\"value\"}"}`, string(got))
	})

	t.Run("marshals nested values", func(t *testing.T) {
		data := StringifiedFrom(map[string]any{"nested": map[string]any{"list": []any{1, "two"}}})
		got, err := json.Marshal(stringDTO{Name: "test", Data: *data})
		require.NoError(t, err)
		assert.Equal(t, `{"name":"test","data":"{\"nested\":{\"list\":[1,\"two\"]}}"}`, string(got))
	})

	t.Run("marshals nil object as null and empty object as {}", func(t *testing.T) {
		got, err := json.Marshal(stringDTO{Name: "test"})
		require.NoError(t, err)
		assert.Equal(t, `{"name":"test","data":null}`, string(got))

		got, err = json.Marshal(stringDTO{Name: "test", Data: *StringifiedFrom(map[string]any{})})
		require.NoError(t, err)
		assert.Equal(t, `{"name":"test","data":"{}"}`, string(got))
	})

	t.Run("omitzero skips nil and empty objects", func(t *testing.T) {
		for _, data := range []Stringified{{}, *StringifiedFrom(map[string]any{})} {
			got, err := json.Marshal(stringDTOOmit{Name: "test", Data: data})
			require.NoError(t, err)
			assert.Equal(t, `{"name":"test"}`, string(got))
		}
	})

	t.Run("unmarshals an escaped JSON string", func(t *testing.T) {
		var got stringDTO
		require.NoError(t, json.Unmarshal([]byte(`{"name":"test","data":"{\"key\":\"value\",\"obj\":{\"a\":[1]}}"}`), &got))
		assert.Equal(t, "test", got.Name)
		assert.Equal(t, map[string]any{"key": "value", "obj": map[string]any{"a": []any{float64(1)}}}, got.Data.Object)
	})

	t.Run("unmarshals null and empty string as nil object", func(t *testing.T) {
		for _, data := range []string{`null`, `""`} {
			var got stringDTO
			require.NoError(t, json.Unmarshal([]byte(`{"data":`+data+`}`), &got), data)
			assert.Nil(t, got.Data.Object, data)
		}
	})

	t.Run("unmarshals a JSON tree", func(t *testing.T) {
		var got stringDTO
		require.NoError(t, json.Unmarshal([]byte(`{"data": {"key":"value"}}`), &got))
		assert.Equal(t, map[string]any{"key": "value"}, got.Data.Object)
	})

	t.Run("rejects a string that is not a JSON object", func(t *testing.T) {
		for _, data := range []string{`"not json"`, `"[1,2]"`, `"42"`} {
			var got stringDTO
			assert.Error(t, json.Unmarshal([]byte(`{"data":`+data+`}`), &got), data)
		}
	})

	t.Run("round-trips", func(t *testing.T) {
		in := stringDTO{Name: "test", Data: *StringifiedFrom(map[string]any{"key": "va\"lue", "list": []any{"a", map[string]any{"b": true}}})}
		b, err := json.Marshal(in)
		require.NoError(t, err)
		var out stringDTO
		require.NoError(t, json.Unmarshal(b, &out))
		assert.Equal(t, in, out)
	})
}

func TestStringifiedYAML(t *testing.T) {
	t.Run("marshals as a JSON string", func(t *testing.T) {
		got, err := yaml.Marshal(stringDTO{Name: "test", Data: *StringifiedFrom(map[string]any{"key": "value"})})
		require.NoError(t, err)
		assert.Equal(t, "name: test\ndata: '{\"key\":\"value\"}'\n", string(got))
	})

	t.Run("marshals through a pointer field", func(t *testing.T) {
		got, err := yaml.Marshal(stringDTOPtr{Name: "test", Data: StringifiedFrom(map[string]any{"key": "value"})})
		require.NoError(t, err)
		assert.Equal(t, "name: test\ndata: '{\"key\":\"value\"}'\n", string(got))
	})

	t.Run("marshals nil object as null", func(t *testing.T) {
		got, err := yaml.Marshal(stringDTO{Name: "test"})
		require.NoError(t, err)
		assert.Equal(t, "name: test\ndata: null\n", string(got))
	})

	t.Run("omitempty skips nil and empty objects", func(t *testing.T) {
		for _, data := range []Stringified{{}, *StringifiedFrom(map[string]any{})} {
			got, err := yaml.Marshal(stringDTOOmit{Name: "test", Data: data})
			require.NoError(t, err)
			assert.Equal(t, "name: test\n", string(got))
		}
	})

	t.Run("unmarshals quoted and block strings", func(t *testing.T) {
		for _, doc := range []string{
			"data: '{\"key\":\"value\"}'\n",
			"data: |\n  {\"key\": \"value\"}\n",
		} {
			var got stringDTO
			require.NoError(t, yaml.Unmarshal([]byte(doc), &got), doc)
			assert.Equal(t, map[string]any{"key": "value"}, got.Data.Object, doc)
		}
	})

	t.Run("unmarshals a YAML tree", func(t *testing.T) {
		var got stringDTO
		require.NoError(t, yaml.Unmarshal([]byte("data:\n  key: value\n"), &got))
		assert.Equal(t, map[string]any{"key": "value"}, got.Data.Object)
	})

	t.Run("rejects a YAML sequence", func(t *testing.T) {
		var got stringDTO
		assert.Error(t, yaml.Unmarshal([]byte("data: [a, b]\n"), &got))
	})

	t.Run("round-trips", func(t *testing.T) {
		in := stringDTO{Name: "test", Data: *StringifiedFrom(map[string]any{"key": "va'lue", "list": []any{"a", map[string]any{"b": true}}})}
		b, err := yaml.Marshal(in)
		require.NoError(t, err)
		var out stringDTO
		require.NoError(t, yaml.Unmarshal(b, &out))
		assert.Equal(t, in, out)
	})
}

func TestStringifiedMapsToUnstructured(t *testing.T) {
	var fromAPI stringDTO
	require.NoError(t, json.Unmarshal([]byte(`{"name":"test","data":"{\"key\":\"value\"}"}`), &fromAPI))

	got, err := json.Marshal(treeDTO{Name: fromAPI.Name, Data: fromAPI.Data.Unstructured})
	require.NoError(t, err)
	assert.Equal(t, `{"name":"test","data":{"key":"value"}}`, string(got))
}

func TestUnstructuredMapsToStringifiedViaJSON(t *testing.T) {
	b, err := json.Marshal(treeDTO{Name: "test", Data: *From(map[string]any{"key": "value"})})
	require.NoError(t, err)

	var toAPI stringDTO
	require.NoError(t, json.Unmarshal(b, &toAPI))

	got, err := json.Marshal(toAPI)
	require.NoError(t, err)
	assert.Equal(t, `{"name":"test","data":"{\"key\":\"value\"}"}`, string(got))
}
