package pkg

import (
	"context"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gravitee-io-labs/gravitee-automation-sdks/common/pkg/apicontext"
	"github.com/gravitee-io-labs/gravitee-automation-sdks/common/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestNewClient_NoAuth(t *testing.T) {
	var ac apicontext.APIContext
	_, err := NewClient(ac, 0)
	assert.Error(t, err, errors.NoAuthProvided)
}

func TestNewClient_InvalidAuth(t *testing.T) {
	ac := apicontext.APIContext{
		Auth: apicontext.Auth{
			BasicAuth:   &apicontext.BasicAuth{},
			BearerToken: new(""),
		},
	}
	_, err := NewClient(ac, 0)
	assert.Error(t, err, errors.ManyAuthProvided)
}

func TestNewClient_InvalidURL(t *testing.T) {
	ac := apicontext.APIContext{
		BaseURL: "::",
	}
	_, err := NewClient(ac, 0)
	assert.Error(t, err)
}

func TestNewClient_Valid(t *testing.T) {
	ac := apicontext.APIContext{
		BaseURL: "http://localhost/automation/",
		Auth:    apicontext.Auth{BearerToken: new("123")},
	}
	client, err := NewClient(ac, 0)
	assert.NoError(t, err)
	assert.NotNil(t, client.Certificates)
	assert.NotNil(t, client.Domains)
	assert.NotNil(t, client.IdentityProviders)
	assert.NotNil(t, client.Reporters)
}

func TestNewClient_Call(t *testing.T) {

	tests := []struct {
		name               string
		apiContextSupplier func(string) apicontext.APIContext
		expectedAuth       string
		expectedUri        string
	}{
		{
			name: "default with token",
			apiContextSupplier: func(serverURL string) apicontext.APIContext {
				return apicontext.APIContext{
					BaseURL: serverURL + "/automation",
					Auth: apicontext.Auth{
						BearerToken: new("123"),
					},
				}
			},
			expectedAuth: "Bearer 123",
			expectedUri:  "/automation/organizations/DEFAULT/environments/DEFAULT/domains",
		}, {
			name: "default with trailing slash with token",
			apiContextSupplier: func(serverURL string) apicontext.APIContext {
				return apicontext.APIContext{
					BaseURL: serverURL + "/automation/",
					Auth: apicontext.Auth{
						BearerToken: new("123"),
					},
				}
			},
			expectedAuth: "Bearer 123",
			expectedUri:  "/automation/organizations/DEFAULT/environments/DEFAULT/domains",
		}, {
			name: "given org and env with basic",
			apiContextSupplier: func(serverURL string) apicontext.APIContext {
				return apicontext.APIContext{
					BaseURL: serverURL + "/automation",
					OrgID:   "foo",
					EnvID:   "bar",
					Auth: apicontext.Auth{
						BasicAuth: &apicontext.BasicAuth{
							Username: "admin",
							Password: "admin",
						},
					},
				}
			},
			expectedAuth: "Basic " + base64.URLEncoding.EncodeToString([]byte("admin:admin")),
			expectedUri:  "/automation/organizations/foo/environments/bar/domains",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var givenAuth string
			var givenUri string
			s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				givenUri = r.RequestURI
				givenAuth = r.Header.Get("Authorization")
				_, err := w.Write([]byte("[]"))
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
			}))
			defer s.Close()
			client, err := NewClient(tt.apiContextSupplier(s.URL), 0)
			assert.NoError(t, err)
			r, err := client.Domains.AutomationListDomainsWithResponse(context.Background())
			assert.NoError(t, err)
			assert.NotNil(t, r)
			assert.Equal(t, r.StatusCode(), http.StatusOK)
			assert.Equal(t, tt.expectedAuth, givenAuth)
			assert.Equal(t, tt.expectedUri, givenUri)
		})
	}
}
