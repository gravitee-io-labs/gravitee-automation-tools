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

// Package apicontext holds connection settings for an Automation API client: URL, org/env, and auth.
package apicontext

import (
	"context"
	"net/http"

	"github.com/gravitee-io-labs/gravitee-automation-sdks/common/pkg/errors"
	"github.com/oapi-codegen/oapi-codegen/v2/pkg/securityprovider"
)

const defaultOrgEnv = "DEFAULT"

// APIContext is the input to NewClient. Empty OrgID and EnvID become "DEFAULT" at request time.
// Auth must be exactly one of bearer or basic; that is enforced by AuthInterceptor, not here.
type APIContext struct {
	BaseURL string
	OrgID   string
	EnvID   string
	Auth    Auth
}

// Auth is the credentials for one APIContext. Set BearerToken or BasicAuth, not both. Pointers are optional slots.
type Auth struct {
	// BearerToken, when non-nil, selects bearer auth. The pointed string is the token value.
	BearerToken *string
	// BasicAuth, when non-nil, selects HTTP Basic. Username and Password are sent as-is.
	BasicAuth *BasicAuth
}

func (c APIContext) isBearerAuth() bool {
	return c.Auth.BearerToken != nil
}

func (c APIContext) isBasicAuth() bool {
	return c.Auth.BasicAuth != nil
}

// BasicAuth is a username/password pair. Both fields are required when this struct is used.
type BasicAuth struct {
	Username string
	Password string
}

// AuthInterceptor returns a request editor that sets Authorization. It errors if both or neither auth mode is set.
func (c APIContext) AuthInterceptor() (func(ctx context.Context, req *http.Request) error, error) {
	if c.isBearerAuth() && c.isBasicAuth() {
		return nil, errors.NewClientError(errors.ManyAuthProvided)
	} else if c.isBearerAuth() {
		bearerAuth, err := securityprovider.NewSecurityProviderBearerToken(*c.Auth.BearerToken)
		if err != nil {
			return nil, errors.NewClientError(err)
		}
		return bearerAuth.Intercept, nil
	} else if c.isBasicAuth() {
		basicAuth, err := securityprovider.NewSecurityProviderBasicAuth(
			c.Auth.BasicAuth.Username,
			c.Auth.BasicAuth.Password)
		if err != nil {
			return nil, errors.NewClientError(err)
		}
		return basicAuth.Intercept, nil
	}
	return nil, errors.NewClientError(errors.NoAuthProvided)
}

// GetOrgIdOrDefault returns OrgID, or "DEFAULT" when OrgID is empty.
func (c APIContext) GetOrgIdOrDefault() string {
	if c.OrgID != "" {
		return c.OrgID
	}
	return defaultOrgEnv
}

// GetEnvIdOrDefault returns EnvID, or "DEFAULT" when EnvID is empty.
func (c APIContext) GetEnvIdOrDefault() string {
	if c.EnvID != "" {
		return c.EnvID
	}
	return defaultOrgEnv
}
