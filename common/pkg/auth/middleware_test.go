package auth

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testRegistry(t *testing.T) *Registry {
	t.Helper()
	cfg, err := LoadConfig(writeConfigFile(t, testYAML))
	require.NoError(t, err)
	return NewRegistry(*cfg, "/automation")
}

var okHandler = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
})

func TestAuthnMiddleware_BearerOK(t *testing.T) {
	reg := testRegistry(t)
	handler := AuthnMiddleware(reg)(okHandler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer 0123456789")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestAuthnMiddleware_BearerUnknown(t *testing.T) {
	reg := testRegistry(t)
	handler := AuthnMiddleware(reg)(okHandler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer bad-token")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestAuthnMiddleware_BasicOK(t *testing.T) {
	reg := testRegistry(t)
	handler := AuthnMiddleware(reg)(okHandler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("admin:admin")))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestAuthnMiddleware_BasicWrongPassword(t *testing.T) {
	reg := testRegistry(t)
	handler := AuthnMiddleware(reg)(okHandler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("admin:wrong")))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestAuthnMiddleware_NoHeader(t *testing.T) {
	reg := testRegistry(t)
	handler := AuthnMiddleware(reg)(okHandler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestAuthnMiddleware_SetsUserInContext(t *testing.T) {
	reg := testRegistry(t)
	var captured *User
	inner := http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		captured = UserFromContext(r.Context())
	})
	handler := AuthnMiddleware(reg)(inner)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer 0123456789")
	handler.ServeHTTP(httptest.NewRecorder(), req)

	require.NotNil(t, captured)
	assert.Equal(t, "admin", captured.Name)
}

func TestAuthzMiddleware_AllPermissions(t *testing.T) {
	reg := testRegistry(t)
	extract := func(_ *http.Request) RouteInfo {
		return RouteInfo{Method: "GET", RoutePattern: "/automation/organizations/{orgId}/environments/{envId}/domains"}
	}
	handler := AuthzMiddleware(reg, extract)(okHandler)

	admin := &User{Name: "admin", AllPermissions: true}
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req = req.WithContext(withUser(req.Context(), admin))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestAuthzMiddleware_HasPermission(t *testing.T) {
	reg := testRegistry(t)
	extract := func(_ *http.Request) RouteInfo {
		return RouteInfo{Method: "GET", RoutePattern: "/automation/organizations/{orgId}/environments/{envId}/domains"}
	}
	handler := AuthzMiddleware(reg, extract)(okHandler)

	reader := &User{Name: "reader", permissions: map[string]bool{"DOMAIN_READ": true}}
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req = req.WithContext(withUser(req.Context(), reader))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestAuthzMiddleware_MissingPermission(t *testing.T) {
	reg := testRegistry(t)
	extract := func(_ *http.Request) RouteInfo {
		return RouteInfo{Method: "PUT", RoutePattern: "/automation/organizations/{orgId}/environments/{envId}/domains"}
	}
	handler := AuthzMiddleware(reg, extract)(okHandler)

	reader := &User{Name: "reader", permissions: map[string]bool{"DOMAIN_READ": true}}
	req := httptest.NewRequest(http.MethodPut, "/", nil)
	req = req.WithContext(withUser(req.Context(), reader))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusForbidden, rec.Code)

	var body errorBody
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&body))
	assert.Contains(t, body.Message, "DOMAIN_UPDATE")
}

func TestAuthzMiddleware_NoPermissionConfig(t *testing.T) {
	reg := testRegistry(t)
	extract := func(_ *http.Request) RouteInfo {
		return RouteInfo{Method: "POST", RoutePattern: "/automation/something/unconfigured"}
	}
	handler := AuthzMiddleware(reg, extract)(okHandler)

	reader := &User{Name: "reader", permissions: map[string]bool{"DOMAIN_READ": true}}
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	req = req.WithContext(withUser(req.Context(), reader))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestAuthzMiddleware_NoUser(t *testing.T) {
	reg := testRegistry(t)
	extract := func(_ *http.Request) RouteInfo {
		return RouteInfo{Method: "GET", RoutePattern: "/automation/organizations/{orgId}/environments/{envId}/domains"}
	}
	handler := AuthzMiddleware(reg, extract)(okHandler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}
