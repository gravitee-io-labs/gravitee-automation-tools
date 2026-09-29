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
	"time"

	"github.com/gravitee-io/gravitee-automation-sdk/am-sdk/v2/pkg/sdk"
	"github.com/gravitee-io/gravitee-automation-sdk/common/pkg/apicontext"
)

// AMClient is the generated Automation API client, sharing one base URL, org/env, auth, and HTTP client.
//
// Deprecated: use sdk.AMClient instead.
type AMClient struct {
	*sdk.AMClient
}

// NewClient builds an AMClient from ac. Trailing slashes are stripped from BaseURL.
// OrgID and EnvID are baked into the server URL (empty becomes "DEFAULT").
// Auth must be exactly one of bearer or basic.
// Without timeoutMs the HTTP client has no timeout. Otherwise timeoutMs[0] is the timeout
// in milliseconds (0 means no timeout); use WithHTTPClient for any other HTTP setting.
//
// Deprecated: use sdk.NewAMClient instead.
func NewClient(ac apicontext.APIContext, timeoutMs ...int) (*AMClient, error) {
	client, err := sdk.NewAMClient(ac)
	if err != nil {
		return nil, err
	}
	if len(timeoutMs) == 0 {
		return &AMClient{AMClient: client}, nil
	}
	return (&AMClient{AMClient: client}).WithHTTPClient(&http.Client{Timeout: time.Duration(timeoutMs[0]) * time.Millisecond})
}

// WithHTTPClient returns a new AMClient with the same server URL and auth, sending requests through httpClient.
// The receiver is left unchanged.
func (c *AMClient) WithHTTPClient(httpClient *http.Client) (*AMClient, error) {
	client, err := c.AMClient.WithHTTPClient(httpClient)
	if err != nil {
		return nil, err
	}
	return &AMClient{AMClient: client}, nil
}
