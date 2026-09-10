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

// Package auth is file-based authentication and route authorization for the mock server.
package auth

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config is the YAML file shape. Load once, then pass to NewRegistry; do not mutate after that.
type Config struct {
	Users       map[string]UserConfig `yaml:"users"`
	Permissions []RoutePermission     `yaml:"permissions"`
}

// UserConfig is one user in the YAML. Empty Permissions means all permissions (AllPermissions on User).
// Token indexes bearer auth; Password indexes basic auth by the map key as username.
type UserConfig struct {
	Password    string   `yaml:"password"`
	Token       string   `yaml:"token"`
	Permissions []string `yaml:"permissions"`
}

// RoutePermission binds HTTP methods on Path to permission names. Empty method fields are omitted from the registry.
// Path is appended to the registry basePath (OpenAPI path template, not a concrete URL).
type RoutePermission struct {
	Path   string `yaml:"path"`
	Get    string `yaml:"get"`
	Put    string `yaml:"put"`
	Delete string `yaml:"delete"`
}

// LoadConfig reads and unmarshals a YAML auth file. The file is not watched; reload by calling again.
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading auth config: %w", err)
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing auth config: %w", err)
	}
	return &cfg, nil
}

// User is a resolved principal after authentication. Name and Password come from YAML.
// AllPermissions is true when YAML Permissions was empty; otherwise HasPermission uses the listed names.
type User struct {
	Name           string
	Password       string
	permissions    map[string]bool
	AllPermissions bool
}

// HasPermission reports whether perm is granted. AllPermissions short-circuits to true.
func (u *User) HasPermission(perm string) bool {
	if u.AllPermissions {
		return true
	}
	return u.permissions[perm]
}

// Registry is an immutable lookup built from Config. Do not share a Config mutation after NewRegistry.
type Registry struct {
	tokenIndex map[string]*User
	basicIndex map[string]*User
	routePerms map[string]string
}

// NewRegistry indexes users and route permissions. basePath is prepended to each RoutePermission.Path.
func NewRegistry(cfg Config, basePath string) *Registry {
	reg := &Registry{
		tokenIndex: make(map[string]*User),
		basicIndex: make(map[string]*User),
		routePerms: make(map[string]string),
	}

	for name, uc := range cfg.Users {
		u := &User{
			Name:           name,
			Password:       uc.Password,
			permissions:    make(map[string]bool),
			AllPermissions: len(uc.Permissions) == 0,
		}
		for _, p := range uc.Permissions {
			u.permissions[p] = true
		}
		if uc.Token != "" {
			reg.tokenIndex[uc.Token] = u
		}
		if uc.Password != "" {
			reg.basicIndex[name] = u
		}
	}

	for _, rp := range cfg.Permissions {
		path := basePath + rp.Path
		if rp.Get != "" {
			reg.routePerms[routeKey(http.MethodGet, path)] = rp.Get
		}
		if rp.Put != "" {
			reg.routePerms[routeKey(http.MethodPut, path)] = rp.Put
		}
		if rp.Delete != "" {
			reg.routePerms[routeKey(http.MethodDelete, path)] = rp.Delete
		}
	}

	return reg
}

// AuthenticateBearer returns the user for token. ok is false when the token is unknown.
func (reg *Registry) AuthenticateBearer(token string) (*User, bool) {
	u, ok := reg.tokenIndex[token]
	return u, ok
}

// AuthenticateBasic returns the user for username/password. ok is false on unknown user or wrong password.
func (reg *Registry) AuthenticateBasic(username, password string) (*User, bool) {
	u, ok := reg.basicIndex[username]
	if !ok || u.Password != password {
		return nil, false
	}
	return u, true
}

// RequiredPermission returns the permission name for method + routePattern. ok is false when the route is unconfigured (allow).
func (reg *Registry) RequiredPermission(method, routePattern string) (string, bool) {
	perm, ok := reg.routePerms[routeKey(method, routePattern)]
	return perm, ok
}

func routeKey(method, path string) string {
	return strings.ToUpper(method) + " " + path
}
