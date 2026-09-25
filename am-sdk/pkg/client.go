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

// AMClient is the generated Automation API client, sharing one base URL, org/env, auth, and HTTP client.
type AMClient struct {
	sdk.ClientWithResponsesInterface
	server string
	editor sdk.RequestEditorFn
}

// NewClient builds an AMClient from ac. Trailing slashes are stripped from BaseURL.
// OrgID and EnvID are baked into the server URL (empty becomes "DEFAULT").
// Auth must be exactly one of bearer or basic.
// Without timeoutMs the HTTP client has no timeout. Otherwise timeoutMs[0] is the timeout
// in milliseconds (0 means no timeout); use WithHTTPClient for any other HTTP setting.
func NewClient(ac apicontext.APIContext, timeoutMs ...int) (*AMClient, error) {
	if len(timeoutMs) > 0 {
		client, err := NewClient(ac)
		if err != nil {
			return nil, err
		}
		return client.WithHTTPClient(&http.Client{Timeout: time.Duration(timeoutMs[0]) * time.Millisecond})
	}

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

	return (&AMClient{server: server, editor: editor}).WithHTTPClient(&http.Client{})
}

// WithHTTPClient returns a new AMClient with the same server URL and auth, sending requests through httpClient.
// The receiver is left unchanged.
func (c *AMClient) WithHTTPClient(httpClient *http.Client) (*AMClient, error) {
	client, err := sdk.NewClientWithResponses(
		c.server,
		sdk.WithHTTPClient(httpClient),
		sdk.WithRequestEditorFn(c.editor),
	)
	if err != nil {
		return nil, err
	}
	return &AMClient{ClientWithResponsesInterface: client, server: c.server, editor: c.editor}, nil
}
