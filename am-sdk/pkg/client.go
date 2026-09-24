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
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gravitee-io-labs/gravitee-automation-tools/am-sdk/v2/pkg/sdk"
	"github.com/gravitee-io-labs/gravitee-automation-tools/common/pkg/apicontext"
	"github.com/gravitee-io-labs/gravitee-automation-tools/common/pkg/errors"
)

// AMClient is the generated Automation API client, sharing one base URL, org/env, auth, and HTTP timeout.
type AMClient struct {
	sdk.ClientWithResponsesInterface
}

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

	server, err := sdk.NewScopedServerURL(
		sdk.ScopedServerURLBaseUrlVariable(baseUrl.String()),
		sdk.ScopedServerURLEnvIdVariable(ac.GetEnvIdOrDefault()),
		sdk.ScopedServerURLOrgIdVariable(ac.GetOrgIdOrDefault()),
	)
	if err != nil {
		return nil, err
	}

	client, err := sdk.NewClientWithResponses(
		server,
		sdk.WithHTTPClient(&http.Client{Timeout: time.Duration(timeoutMs) * time.Millisecond}),
		sdk.WithRequestEditorFn(editor),
	)
	if err != nil {
		return nil, err
	}

	return &AMClient{client}, nil
}
