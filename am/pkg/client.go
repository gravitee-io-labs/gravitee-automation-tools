package pkg

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gravitee-io-labs/gravitee-automation-sdks/am/pkg/sdk/certificate"
	"github.com/gravitee-io-labs/gravitee-automation-sdks/am/pkg/sdk/domain"
	"github.com/gravitee-io-labs/gravitee-automation-sdks/am/pkg/sdk/identityprovider"
	"github.com/gravitee-io-labs/gravitee-automation-sdks/am/pkg/sdk/reporter"
	"github.com/gravitee-io-labs/gravitee-automation-sdks/common/pkg/apicontext"
	"github.com/gravitee-io-labs/gravitee-automation-sdks/common/pkg/errors"
)

const defaultOrgEnv = "DEFAULT"

type AMClient struct {
	Domains           domain.ClientWithResponsesInterface
	Certificates      certificate.ClientWithResponsesInterface
	IdentityProviders identityprovider.ClientWithResponsesInterface
	Reporters         reporter.ClientWithResponsesInterface
}

type requestEditor = func(context.Context, *http.Request) error

func NewClient(ac apicontext.APIContext, timeoutMs int) (*AMClient, error) {

	baseUrl, err := url.Parse(strings.TrimRight(ac.BaseURL, "/"))
	if err != nil {
		return nil, errors.NewClientError(err)
	}

	editor, err := ac.AuthInterceptor()
	if err != nil {
		return nil, err
	}

	server, err := domain.NewScopedServerURL(
		domain.ScopedServerURLBaseUrlVariable(baseUrl.String()),
		domain.ScopedServerURLEnvIdVariable(orDefault(ac.EnvID)),
		domain.ScopedServerURLOrgIdVariable(orDefault(ac.OrgID)),
	)
	if err != nil {
		return nil, err
	}

	httpClient := &http.Client{Timeout: time.Duration(timeoutMs) * time.Millisecond}

	domains, err := newDomainClient(server, httpClient, editor)
	if err != nil {
		return nil, err
	}
	certificates, err := newCertificateClient(server, httpClient, editor)
	if err != nil {
		return nil, err
	}
	identities, err := newIdentityProviderClient(server, httpClient, editor)
	if err != nil {
		return nil, err
	}
	reporters, err := newReporterClient(server, httpClient, editor)
	if err != nil {
		return nil, err
	}

	return &AMClient{
		Domains:           domains,
		Certificates:      certificates,
		IdentityProviders: identities,
		Reporters:         reporters,
	}, nil
}

func orDefault(id string) string {
	if id == "" {
		return defaultOrgEnv
	}
	return id
}

func newDomainClient(server string, httpClient *http.Client, editor requestEditor) (domain.ClientWithResponsesInterface, error) {
	return domain.NewClientWithResponses(
		server,
		domain.WithHTTPClient(httpClient),
		domain.WithRequestEditorFn(editor),
	)
}

func newCertificateClient(server string, httpClient *http.Client, editor requestEditor) (certificate.ClientWithResponsesInterface, error) {
	return certificate.NewClientWithResponses(
		server,
		certificate.WithHTTPClient(httpClient),
		certificate.WithRequestEditorFn(editor),
	)
}

func newIdentityProviderClient(server string, httpClient *http.Client, editor requestEditor) (identityprovider.ClientWithResponsesInterface, error) {
	return identityprovider.NewClientWithResponses(
		server,
		identityprovider.WithHTTPClient(httpClient),
		identityprovider.WithRequestEditorFn(editor),
	)
}

func newReporterClient(server string, httpClient *http.Client, editor requestEditor) (reporter.ClientWithResponsesInterface, error) {
	return reporter.NewClientWithResponses(
		server,
		reporter.WithHTTPClient(httpClient),
		reporter.WithRequestEditorFn(editor),
	)
}
