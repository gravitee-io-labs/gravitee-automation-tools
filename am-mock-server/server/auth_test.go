package server

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gravitee-io-labs/gravitee-automation-sdks/common/pkg/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testAuthConfig() auth.Config {
	return auth.Config{
		Users: map[string]auth.UserConfig{
			"admin": {
				Password:    "admin",
				Token:       "admin-token",
				Permissions: nil,
			},
			"reader": {
				Token:       "reader-token",
				Permissions: []string{"DOMAIN_READ"},
			},
		},
		Permissions: []auth.RoutePermission{
			{
				Path: "/organizations/{orgId}/environments/{envId}/domains",
				Get:  "DOMAIN_READ",
				Put:  "DOMAIN_UPDATE",
			},
			{
				Path:   "/organizations/{orgId}/environments/{envId}/domains/{domainKey}",
				Get:    "DOMAIN_READ",
				Delete: "DOMAIN_DELETE",
			},
		},
	}
}

func createAMServerWithAuth(t *testing.T) (*MockAM, *httptest.Server) {
	t.Helper()
	cfg := testAuthConfig()
	reg := auth.NewRegistry(cfg, BasePath)
	am := NewMockAM(t.Context())
	srv := httptest.NewServer(NewWithPath(am, BasePath, reg))
	t.Cleanup(srv.Close)
	return am, srv
}

func TestAuth_BearerOK(t *testing.T) {
	am, srv := createAMServerWithAuth(t)
	defaultTenant(am).Domains.Put(Domain{Key: "test", Name: "Test domain"})

	req, _ := http.NewRequest(http.MethodGet, domainsURL(srv), nil)
	req.Header.Set("Authorization", "Bearer admin-token")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestAuth_BasicOK(t *testing.T) {
	am, srv := createAMServerWithAuth(t)
	defaultTenant(am).Domains.Put(Domain{Key: "test", Name: "Test domain"})

	req, _ := http.NewRequest(http.MethodGet, domainsURL(srv), nil)
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("admin:admin")))

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestAuth_NoHeader_401(t *testing.T) {
	_, srv := createAMServerWithAuth(t)

	resp, err := http.Get(domainsURL(srv))
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestAuth_BadToken_401(t *testing.T) {
	_, srv := createAMServerWithAuth(t)

	req, _ := http.NewRequest(http.MethodGet, domainsURL(srv), nil)
	req.Header.Set("Authorization", "Bearer wrong-token")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestAuth_WrongPassword_401(t *testing.T) {
	_, srv := createAMServerWithAuth(t)

	req, _ := http.NewRequest(http.MethodGet, domainsURL(srv), nil)
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("admin:wrong")))

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestAuth_ReaderCanRead_200(t *testing.T) {
	am, srv := createAMServerWithAuth(t)
	defaultTenant(am).Domains.Put(Domain{Key: "test", Name: "Test domain"})

	req, _ := http.NewRequest(http.MethodGet, domainsURL(srv), nil)
	req.Header.Set("Authorization", "Bearer reader-token")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestAuth_ReaderCannotWrite_403(t *testing.T) {
	_, srv := createAMServerWithAuth(t)

	body := encode(t, Domain{Key: "test", Name: "Test domain"})
	req, _ := http.NewRequest(http.MethodPut, domainsURL(srv), body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer reader-token")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusForbidden, resp.StatusCode)
}

func TestAuth_AdminCanWrite_200(t *testing.T) {
	_, srv := createAMServerWithAuth(t)

	body := encode(t, Domain{Key: "test", Name: "Test domain"})
	req, _ := http.NewRequest(http.MethodPut, domainsURL(srv), body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer admin-token")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestAuth_NoAuthConfig_OpenAccess(t *testing.T) {
	am, srv := createAMServer(t)
	defaultTenant(am).Domains.Put(Domain{Key: "test", Name: "Test domain"})

	resp, err := http.Get(domainsURL(srv))
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
