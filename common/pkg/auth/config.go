package auth

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Users       map[string]UserConfig `yaml:"users"`
	Permissions []RoutePermission     `yaml:"permissions"`
}

type UserConfig struct {
	Password    string   `yaml:"password"`
	Token       string   `yaml:"token"`
	Permissions []string `yaml:"permissions"`
}

type RoutePermission struct {
	Path   string `yaml:"path"`
	Get    string `yaml:"get"`
	Put    string `yaml:"put"`
	Delete string `yaml:"delete"`
}

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

type User struct {
	Name           string
	Password       string
	permissions    map[string]bool
	AllPermissions bool
}

func (u *User) HasPermission(perm string) bool {
	if u.AllPermissions {
		return true
	}
	return u.permissions[perm]
}

type Registry struct {
	tokenIndex map[string]*User
	basicIndex map[string]*User
	routePerms map[string]string
}

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

func (reg *Registry) AuthenticateBearer(token string) (*User, bool) {
	u, ok := reg.tokenIndex[token]
	return u, ok
}

func (reg *Registry) AuthenticateBasic(username, password string) (*User, bool) {
	u, ok := reg.basicIndex[username]
	if !ok || u.Password != password {
		return nil, false
	}
	return u, true
}

func (reg *Registry) RequiredPermission(method, routePattern string) (string, bool) {
	perm, ok := reg.routePerms[routeKey(method, routePattern)]
	return perm, ok
}

func routeKey(method, path string) string {
	return strings.ToUpper(method) + " " + path
}
