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

package server

import (
	"testing"

	"github.com/gravitee-io-labs/gravitee-automation-tools/am-sdk/v2/pkg/sdk"
)

func TestListDomains(t *testing.T) {
	am, srv := createAMServer(t)
	defaultTenant(am).Domains.Put(Domain{Key: "test", Name: "Test domain"})

	assertListEqual(t, domainsURL(srv), []Domain{{Key: "test", Name: "Test domain"}})
}

func TestGetDomain(t *testing.T) {
	am, srv := createAMServer(t)
	defaultTenant(am).Domains.Put(Domain{Key: "test", Name: "Test domain"})

	assertGetEqual(t, domainsURL(srv), "test", Domain{Key: "test", Name: "Test domain"})
}

func TestGetDomain404(t *testing.T) {
	_, srv := createAMServer(t)

	assertGet404(t, domainsURL(srv), "test", "Domain")
}

func TestPutGetDomain(t *testing.T) {
	_, srv := createAMServer(t)
	body := Domain{Key: "test", Name: "Test domain"}.WithDefaults()

	assertPutEqual(t, domainsURL(srv), body)
	assertGetEqual(t, domainsURL(srv), "test", body)
}

func TestDeleteDomain(t *testing.T) {
	am, srv := createAMServer(t)
	defaultTenant(am).Domains.Put(Domain{Key: "test", Name: "Test domain"})

	assertDeleteGone(t, domainsURL(srv), "test", "Domain")
}

func TestListDomainsSDK(t *testing.T) {
	am, srv := createAMServer(t)
	defaultTenant(am).Domains.Put(Domain{Key: "test", Name: "Test domain"})
	client := newAMClient(t, srv)

	res, err := client.ListDomainsWithResponse(t.Context())
	assertSDKOK(t, res, err, []sdk.Domain{{Key: "test", Name: "Test domain"}})
}

func TestGetDomainSDK(t *testing.T) {
	am, srv := createAMServer(t)
	defaultTenant(am).Domains.Put(Domain{Key: "test", Name: "Test domain"})
	client := newAMClient(t, srv)

	res, err := client.GetDomainWithResponse(t.Context(), "test")
	assertSDKOK(t, res, err, sdk.Domain{Key: "test", Name: "Test domain"})
}

func TestGetDomain404SDK(t *testing.T) {
	_, srv := createAMServer(t)
	client := newAMClient(t, srv)

	res, err := client.GetDomainWithResponse(t.Context(), "test")
	assertSDK404(t, res, err)
}

func TestPutGetDomainSDK(t *testing.T) {
	_, srv := createAMServer(t)
	client := newAMClient(t, srv)
	body := sdk.Domain{Key: "test", Name: "Test domain"}.WithDefaults()

	put, err := client.UpsertDomainWithResponse(t.Context(), nil, body)
	assertSDKOK(t, put, err, body)

	get, err := client.GetDomainWithResponse(t.Context(), "test")
	assertSDKOK(t, get, err, body)
}

func TestDeleteDomainSDK(t *testing.T) {
	am, srv := createAMServer(t)
	defaultTenant(am).Domains.Put(Domain{Key: "test", Name: "Test domain"})
	client := newAMClient(t, srv)

	del, err := client.DeleteDomainWithResponse(t.Context(), "test")
	assertSDKNoContent(t, del, err)

	get, err := client.GetDomainWithResponse(t.Context(), "test")
	assertSDK404(t, get, err)
}
