package pkg

import (
	"net/http"
	"strings"
	"time"

	"github.com/gravitee-io-labs/gravitee-automation-sdks/am/pkg/sdk/certificate"
	"github.com/gravitee-io-labs/gravitee-automation-sdks/am/pkg/sdk/domain"
	"github.com/gravitee-io-labs/gravitee-automation-sdks/am/pkg/sdk/identityprovider"
	"github.com/gravitee-io-labs/gravitee-automation-sdks/am/pkg/sdk/reporter"
	"github.com/gravitee-io-labs/gravitee-automation-sdks/common/pkg/auth"
)

type AMClient struct {
	Domains           domain.ClientWithResponsesInterface
	Certificates      certificate.ClientWithResponsesInterface
	IdentityProviders identityprovider.ClientWithResponsesInterface
	Reporters         reporter.ClientWithResponsesInterface
}

func NewClient(apiContext auth.APIContext, timeoutMs int) (*AMClient, error) {
	editor, err := apiContext.RequestEditor()
	if err != nil {
		return nil, err
	}

	server, err := domain.NewScopedServerURL(
		domain.ScopedServerURLBaseUrlVariable(strings.TrimRight(apiContext.BaseURL, "/")),
		domain.ScopedServerURLEnvIdVariable(apiContext.EnvID),
		domain.ScopedServerURLOrgIdVariable(apiContext.OrgID),
	)
	if err != nil {
		return nil, err
	}

	httpClient := &http.Client{Timeout: time.Duration(timeoutMs) * time.Millisecond}

	domains, err := domain.NewClientWithResponses(
		server,
		domain.WithHTTPClient(httpClient),
		domain.WithRequestEditorFn(editor),
	)
	if err != nil {
		return nil, err
	}

	certificates, err := certificate.NewClientWithResponses(
		server,
		certificate.WithHTTPClient(httpClient),
		certificate.WithRequestEditorFn(editor),
	)
	if err != nil {
		return nil, err
	}

	identities, err := identityprovider.NewClientWithResponses(
		server,
		identityprovider.WithHTTPClient(httpClient),
		identityprovider.WithRequestEditorFn(editor),
	)
	if err != nil {
		return nil, err
	}

	reporters, err := reporter.NewClientWithResponses(
		server,
		reporter.WithHTTPClient(httpClient),
		reporter.WithRequestEditorFn(editor),
	)
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
