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

package pkg

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gravitee-io/gravitee-automation-sdk/common/pkg/apicontext"
	"github.com/gravitee-io/gravitee-automation-sdk/common/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClient_Valid(t *testing.T) {
	ac := apicontext.APIContext{
		BaseURL: "http://localhost/automation/",
		Auth:    apicontext.Auth{BearerToken: new("123")},
	}
	client, err := NewClient(ac)
	assert.NoError(t, err)
	assert.NotNil(t, client.ClientWithResponsesInterface)
}

func TestNewClient_WithoutTimeout(t *testing.T) {
	s := slowServer(t, 50*time.Millisecond)
	client, err := NewClient(bearerContext(s.URL))
	require.NoError(t, err)

	_, err = client.ListDomainsWithResponse(context.Background())
	assert.NoError(t, err)
}

func TestNewClient_WithTimeout(t *testing.T) {
	s := slowServer(t, 200*time.Millisecond)
	client, err := NewClient(bearerContext(s.URL), 10)
	require.NoError(t, err)

	_, err = client.ListDomainsWithResponse(context.Background())
	assert.ErrorContains(t, err, "Client.Timeout exceeded")
}

func TestNewClient_WithTimeoutAndInvalidAuth(t *testing.T) {
	_, err := NewClient(apicontext.APIContext{BaseURL: "http://localhost"}, 10)
	assert.ErrorIs(t, err, errors.NoAuthProvided)
}

func TestNewClient_CallsScopedURLWithAuth(t *testing.T) {
	var givenAuth, givenURI string
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		givenAuth = r.Header.Get("Authorization")
		givenURI = r.RequestURI
		_, _ = w.Write([]byte("[]"))
	}))
	defer s.Close()
	ac := bearerContext(s.URL)
	ac.OrgID, ac.EnvID = "foo", "bar"
	client, err := NewClient(ac)
	require.NoError(t, err)

	_, err = client.ListDomainsWithResponse(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "Bearer 123", givenAuth)
	assert.Equal(t, "/automation/organizations/foo/environments/bar/domains", givenURI)
}

func TestWithHTTPClient_KeepsWrapperType(t *testing.T) {
	s := slowServer(t, 200*time.Millisecond)
	client, err := NewClient(bearerContext(s.URL))
	require.NoError(t, err)

	var withTimeout *AMClient
	withTimeout, err = client.WithHTTPClient(&http.Client{Timeout: 10 * time.Millisecond})
	require.NoError(t, err)

	_, err = withTimeout.ListDomainsWithResponse(context.Background())
	assert.ErrorContains(t, err, "Client.Timeout exceeded")
}

func bearerContext(serverURL string) apicontext.APIContext {
	return apicontext.APIContext{
		BaseURL: serverURL + "/automation",
		Auth:    apicontext.Auth{BearerToken: new("123")},
	}
}

func slowServer(t *testing.T, delay time.Duration) *httptest.Server {
	t.Helper()
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		time.Sleep(delay)
		_, _ = w.Write([]byte("[]"))
	}))
	t.Cleanup(s.Close)
	return s
}
