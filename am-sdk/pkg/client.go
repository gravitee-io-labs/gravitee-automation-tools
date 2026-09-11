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

// Package pkg is the AM Automation SDK facade. Import it as am.
package pkg

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gravitee-io-labs/gravitee-automation-tools/am-sdk/pkg/sdk/certificate"
	"github.com/gravitee-io-labs/gravitee-automation-tools/am-sdk/pkg/sdk/domain"
	"github.com/gravitee-io-labs/gravitee-automation-tools/am-sdk/pkg/sdk/identityprovider"
	"github.com/gravitee-io-labs/gravitee-automation-tools/am-sdk/pkg/sdk/reporter"
	"github.com/gravitee-io-labs/gravitee-automation-tools/common/pkg/apicontext"
	"github.com/gravitee-io-labs/gravitee-automation-tools/common/pkg/errors"
)

// AMClient groups generated resource clients that share one base URL, org/env, auth, and HTTP timeout.
// Fields are set by NewClient and are safe to read; do not replace them after construction.
type AMClient struct {
	Domains           domain.ClientWithResponsesInterface
	Certificates      certificate.ClientWithResponsesInterface
	IdentityProviders identityprovider.ClientWithResponsesInterface
	Reporters         reporter.ClientWithResponsesInterface
}

type requestEditor = func(context.Context, *http.Request) error

// NewClient builds an AMClient from ac. Trailing slashes are stripped from BaseURL.
// OrgID and EnvID are baked into the server URL (empty becomes "DEFAULT").
// timeoutMs is the HTTP client timeout in milliseconds; 0 means no timeout.
// Auth must be exactly one of bearer or basic.
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
		domain.ScopedServerURLEnvIdVariable(ac.GetEnvIdOrDefault()),
		domain.ScopedServerURLOrgIdVariable(ac.GetOrgIdOrDefault()),
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
