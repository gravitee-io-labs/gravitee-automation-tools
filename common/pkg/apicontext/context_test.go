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
	"encoding/base64"
	"net/http"
	"testing"

	"github.com/gravitee-io-labs/gravitee-automation-sdks/common/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestAPIContext_RequestEditor(t *testing.T) {
	type fields struct {
		BaseURL string
		Auth    Auth
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr error
		auth    string
	}{
		{
			name: "Basic Auth",
			fields: fields{
				Auth: Auth{
					BasicAuth: &BasicAuth{Username: "admin", Password: "admin"},
				},
			},
			wantErr: nil,
			auth:    "Basic " + base64.URLEncoding.EncodeToString([]byte("admin:admin")),
		}, {
			name: "Bearer Auth",
			fields: fields{
				Auth: Auth{
					BearerToken: new("123456789"),
				},
			},
			wantErr: nil,
			auth:    "Bearer 123456789",
		}, {
			name: "Bearer+Basic Auth",
			fields: fields{
				Auth: Auth{
					BearerToken: new("123456789"),
					BasicAuth:   &BasicAuth{Username: "admin", Password: "admin"},
				},
			},
			wantErr: errors.ManyAuthProvided,
		}, {
			name:    "No Auth",
			wantErr: errors.NoAuthProvided,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := APIContext{
				BaseURL: tt.fields.BaseURL,
				Auth:    tt.fields.Auth,
			}
			interceptor, err := c.AuthInterceptor()
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}
			req := http.Request{Header: http.Header{}}
			assert.NoError(t, interceptor(context.Background(), &req))
			assert.Equal(t, tt.auth, req.Header.Get("Authorization"))
		})
	}
}

func TestAPIContext_isBasicAuth(t *testing.T) {
	assert.True(t, APIContext{
		Auth: Auth{
			BasicAuth: &BasicAuth{Username: "admin", Password: "admin"},
		},
	}.isBasicAuth())
	assert.True(t, APIContext{
		Auth: Auth{
			BasicAuth:   &BasicAuth{Username: "admin", Password: "admin"},
			BearerToken: new("123456789"),
		},
	}.isBasicAuth())
	assert.False(t, APIContext{Auth: Auth{BearerToken: new("foo")}}.isBasicAuth())
	assert.False(t, APIContext{}.isBasicAuth())
}

func TestAPIContext_isBearerAuth(t *testing.T) {
	assert.True(t, APIContext{
		Auth: Auth{
			BearerToken: new("123456789"),
		},
	}.isBearerAuth())
	assert.True(t, APIContext{
		Auth: Auth{
			BasicAuth:   &BasicAuth{Username: "admin", Password: "admin"},
			BearerToken: new("123456789"),
		},
	}.isBasicAuth())
	assert.False(t, APIContext{Auth: Auth{BasicAuth: &BasicAuth{
		Username: "admin",
		Password: "admin",
	}}}.isBearerAuth())
	assert.False(t, APIContext{}.isBearerAuth())
}
