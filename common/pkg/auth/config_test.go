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

package auth

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testYAML = `
users:
  admin:
    password: admin
    token: "0123456789"
    permissions: []
  reader:
    token: read-token
    permissions:
      - DOMAIN_READ
permissions:
  - path: /organizations/{orgId}/environments/{envId}/domains
    get: DOMAIN_READ
    put: DOMAIN_UPDATE
  - path: /organizations/{orgId}/environments/{envId}/domains/{domainKey}
    get: DOMAIN_READ
    delete: DOMAIN_DELETE
`

func writeConfigFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "auth.yaml")
	require.NoError(t, os.WriteFile(path, []byte(content), 0o644))
	return path
}

func TestLoadConfig(t *testing.T) {
	cfg, err := LoadConfig(writeConfigFile(t, testYAML))
	require.NoError(t, err)

	assert.Len(t, cfg.Users, 2)
	assert.Equal(t, "admin", cfg.Users["admin"].Password)
	assert.Equal(t, "0123456789", cfg.Users["admin"].Token)
	assert.Empty(t, cfg.Users["admin"].Permissions)
	assert.Equal(t, []string{"DOMAIN_READ"}, cfg.Users["reader"].Permissions)
	assert.Len(t, cfg.Permissions, 2)
}

func TestLoadConfig_MissingFile(t *testing.T) {
	_, err := LoadConfig("/nonexistent/path.yaml")
	assert.Error(t, err)
}

func TestNewRegistry_TokenIndex(t *testing.T) {
	cfg, _ := LoadConfig(writeConfigFile(t, testYAML))
	reg := NewRegistry(*cfg, "/automation")

	u, ok := reg.AuthenticateBearer("0123456789")
	assert.True(t, ok)
	assert.Equal(t, "admin", u.Name)

	u, ok = reg.AuthenticateBearer("read-token")
	assert.True(t, ok)
	assert.Equal(t, "reader", u.Name)

	_, ok = reg.AuthenticateBearer("unknown")
	assert.False(t, ok)
}

func TestNewRegistry_BasicIndex(t *testing.T) {
	cfg, _ := LoadConfig(writeConfigFile(t, testYAML))
	reg := NewRegistry(*cfg, "/automation")

	u, ok := reg.AuthenticateBasic("admin", "admin")
	assert.True(t, ok)
	assert.Equal(t, "admin", u.Name)

	_, ok = reg.AuthenticateBasic("admin", "wrong")
	assert.False(t, ok)

	_, ok = reg.AuthenticateBasic("unknown", "pass")
	assert.False(t, ok)
}

func TestNewRegistry_RoutePerms(t *testing.T) {
	cfg, _ := LoadConfig(writeConfigFile(t, testYAML))
	reg := NewRegistry(*cfg, "/automation")

	perm, ok := reg.RequiredPermission("GET", "/automation/organizations/{orgId}/environments/{envId}/domains")
	assert.True(t, ok)
	assert.Equal(t, "DOMAIN_READ", perm)

	perm, ok = reg.RequiredPermission("PUT", "/automation/organizations/{orgId}/environments/{envId}/domains")
	assert.True(t, ok)
	assert.Equal(t, "DOMAIN_UPDATE", perm)

	perm, ok = reg.RequiredPermission("DELETE", "/automation/organizations/{orgId}/environments/{envId}/domains/{domainKey}")
	assert.True(t, ok)
	assert.Equal(t, "DOMAIN_DELETE", perm)

	_, ok = reg.RequiredPermission("POST", "/automation/organizations/{orgId}/environments/{envId}/domains")
	assert.False(t, ok)
}

func TestUser_HasPermission(t *testing.T) {
	admin := &User{Name: "admin", AllPermissions: true}
	assert.True(t, admin.HasPermission("ANYTHING"))

	reader := &User{Name: "reader", permissions: map[string]bool{"DOMAIN_READ": true}}
	assert.True(t, reader.HasPermission("DOMAIN_READ"))
	assert.False(t, reader.HasPermission("DOMAIN_DELETE"))
}
