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

package apicontext

import (
	"context"
	"net/http"

	"github.com/gravitee-io-labs/gravitee-automation-sdks/common/pkg/errors"
	"github.com/oapi-codegen/oapi-codegen/v2/pkg/securityprovider"
)

const defaultOrgEnv = "DEFAULT"

type APIContext struct {
	BaseURL string
	OrgID   string
	EnvID   string
	Auth    Auth
}

type Auth struct {
	// The bearer token used to authenticate against the API instance
	BearerToken *string
	// The Basic credentials used to authenticate against the API instance.
	BasicAuth *BasicAuth
}

func (c APIContext) isBearerAuth() bool {
	return c.Auth.BearerToken != nil
}

func (c APIContext) isBasicAuth() bool {
	return c.Auth.BasicAuth != nil
}

type BasicAuth struct {
	Username string
	Password string
}

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

func (c APIContext) GetOrgIDOrDefault() string {
	if c.OrgID != "" {
		return c.OrgID
	}
	return defaultOrgEnv
}

func (c APIContext) GetEnvIDOrDefault() string {
	if c.EnvID != "" {
		return c.EnvID
	}
	return defaultOrgEnv
}
