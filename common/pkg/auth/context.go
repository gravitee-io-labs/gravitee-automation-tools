package auth

import (
	"github.com/gravitee-io-labs/gravitee-automation-sdks/am/pkg/sdk"
	"github.com/oapi-codegen/oapi-codegen/v2/pkg/securityprovider"
)

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

func (c APIContext) AuthClientOption() (sdk.ClientOption, error) {
	if c.isBasicAuth() {
		basicAuth, err := securityprovider.NewSecurityProviderBasicAuth(
			c.Auth.BasicAuth.Username,
			c.Auth.BasicAuth.Password)
		if err != nil {
			return nil, ClientError{err: err}
		}
		return sdk.WithRequestEditorFn(basicAuth.Intercept), nil
	} else if c.isBearerAuth() {
		bearerAuth, err := securityprovider.NewSecurityProviderBearerToken(*c.Auth.BearerToken)
		if err != nil {
			return nil, ClientError{err: err}
		}
		return sdk.WithRequestEditorFn(bearerAuth.Intercept), nil
	}
	return nil, NoAuthProvided
}
