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
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDryRunPut_OffPersists(t *testing.T) {
	am, srv := createAMServer(t)
	body := Domain{Key: "test", Name: "Test domain"}

	resp := httpPut(t, domainsURL(srv)+"?dryRun=true", body)
	defer resp.Body.Close()
	failOnNotOK(t, resp)

	got, exists := defaultTenant(am).Domains.Get("test")
	require.True(t, exists)
	assert.Equal(t, body, got)
}

func TestDryRunPut_DoesNotPersist(t *testing.T) {
	am, srv := createAMServerWithDryRunReject(t, true)
	body := Domain{Key: "test", Name: "Test domain"}

	resp := httpPut(t, domainsURL(srv)+"?dryRun=true", body)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, []DryRunError{{Severity: new(SeverityError), Message: new(dryRunMessage)}}, decodeToSliceOf[DryRunError](t, resp))
	_, exists := defaultTenant(am).Domains.Get("test")
	assert.False(t, exists)
	assertGet404(t, domainsURL(srv), "test", "Domain")
}

func TestDryRunPut_FalseStillPersists(t *testing.T) {
	am, srv := createAMServerWithDryRunReject(t, true)
	body := Domain{Key: "test", Name: "Test domain"}

	resp := httpPut(t, domainsURL(srv)+"?dryRun=false", body)
	defer resp.Body.Close()
	failOnNotOK(t, resp)

	got, exists := defaultTenant(am).Domains.Get("test")
	require.True(t, exists)
	assert.Equal(t, body, got)
}

func TestDryRunPut_LeavesExistingUnchanged(t *testing.T) {
	am, srv := createAMServerWithDryRunReject(t, true)
	existing := Domain{Key: "test", Name: "Original"}
	defaultTenant(am).Domains.Put(existing)

	resp := httpPut(t, domainsURL(srv)+"?dryRun=true", Domain{Key: "test", Name: "Changed"})
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, []DryRunError{{Severity: new(SeverityError), Message: new(dryRunMessage)}}, decodeToSliceOf[DryRunError](t, resp))

	got, exists := defaultTenant(am).Domains.Get("test")
	require.True(t, exists)
	assert.Equal(t, existing, got)
}

func TestDryRunPut_NestedDoesNotPersist(t *testing.T) {
	am, srv := createAMServerWithDryRunReject(t, true)
	defaultTenant(am).Domains.Put(Domain{Key: defaultDomainKey, Name: "Test"})

	resp := httpPut(t, certificatesURL(srv)+"?dryRun=true", testCertificate())
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, []DryRunError{{Severity: new(SeverityError), Message: new(dryRunMessage)}}, decodeToSliceOf[DryRunError](t, resp))
	assert.Empty(t, defaultTenant(am).Certificates.GetAll())
}
