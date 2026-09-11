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

	"github.com/gravitee-io-labs/gravitee-automation-tools/am/pkg/sdk/domain"
	"github.com/gravitee-io-labs/gravitee-automation-tools/common/pkg/response"
	"github.com/stretchr/testify/assert"
)

func TestMultiTenancy_DomainIsolation(t *testing.T) {
	_, srv := createAMServer(t)
	t1 := newTestClient(t, srv, "org-alpha", "env-alpha", "test")
	t2 := newTestClient(t, srv, "org-beta", "env-beta", "test")

	domainT1 := domain.Domain{Key: "shared-key", Name: "Tenant 1 Domain"}
	domainT2 := domain.Domain{Key: "shared-key", Name: "Tenant 2 Domain"}

	t.Run("create in both tenants", func(t *testing.T) {
		put1, err := t1.Domains.UpsertDomainWithResponse(t.Context(), domainT1)
		assertSDKOK(t, put1, err, domainT1)

		put2, err := t2.Domains.UpsertDomainWithResponse(t.Context(), domainT2)
		assertSDKOK(t, put2, err, domainT2)
	})

	t.Run("get returns own tenant data", func(t *testing.T) {
		get1, err := t1.Domains.GetDomainWithResponse(t.Context(), "shared-key")
		assertSDKOK(t, get1, err, domainT1)

		get2, err := t2.Domains.GetDomainWithResponse(t.Context(), "shared-key")
		assertSDKOK(t, get2, err, domainT2)
	})

	t.Run("list returns only own tenant data", func(t *testing.T) {
		list1, err := t1.Domains.ListDomainsWithResponse(t.Context())
		assertSDKOK(t, list1, err, []domain.Domain{domainT1})

		list2, err := t2.Domains.ListDomainsWithResponse(t.Context())
		assertSDKOK(t, list2, err, []domain.Domain{domainT2})
	})

	t.Run("update in tenant 1 does not affect tenant 2", func(t *testing.T) {
		updated := domain.Domain{Key: "shared-key", Name: "Tenant 1 Updated"}
		put, err := t1.Domains.UpsertDomainWithResponse(t.Context(), updated)
		assertSDKOK(t, put, err, updated)

		get1, err := t1.Domains.GetDomainWithResponse(t.Context(), "shared-key")
		assertSDKOK(t, get1, err, updated)

		get2, err := t2.Domains.GetDomainWithResponse(t.Context(), "shared-key")
		assertSDKOK(t, get2, err, domainT2)
	})

	t.Run("delete in tenant 1 does not affect tenant 2", func(t *testing.T) {
		del, err := t1.Domains.DeleteDomainWithResponse(t.Context(), "shared-key")
		assertSDKNoContent(t, del, err)

		get1, err := t1.Domains.GetDomainWithResponse(t.Context(), "shared-key")
		assert.True(t, response.IsNotFound(get1, err))

		get2, err := t2.Domains.GetDomainWithResponse(t.Context(), "shared-key")
		assertSDKOK(t, get2, err, domainT2)
	})
}
