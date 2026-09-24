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

func testReporter() Reporter {
	return Reporter{Key: "test", Name: new("Test reporter")}.WithDefaults()
}

func testReporterSDK() sdk.Reporter {
	return sdk.Reporter{Key: "test", Name: new("Test reporter")}.WithDefaults()
}

func TestListReporters(t *testing.T) {
	am, srv := createAMServerWithDomain(t)
	defaultTenant(am).Reporters.Put(newChild(testReporter(), defaultDomainKey))

	assertListEqual(t, reportersURL(srv), []Reporter{testReporter()})
}

func TestGetReporter(t *testing.T) {
	am, srv := createAMServerWithDomain(t)
	defaultTenant(am).Reporters.Put(newChild(testReporter(), defaultDomainKey))

	assertGetEqual(t, reportersURL(srv), "test", testReporter())
}

func TestGetReporter404(t *testing.T) {
	_, srv := createAMServerWithDomain(t)

	assertGet404(t, reportersURL(srv), "test", "Reporter")
}

func TestPutGetReporter(t *testing.T) {
	_, srv := createAMServerWithDomain(t)
	body := testReporter()

	assertPutEqual(t, reportersURL(srv), body)
	assertGetEqual(t, reportersURL(srv), "test", body)
}

func TestDeleteReporter(t *testing.T) {
	am, srv := createAMServerWithDomain(t)
	defaultTenant(am).Reporters.Put(newChild(testReporter(), defaultDomainKey))

	assertDeleteGone(t, reportersURL(srv), "test", "Reporter")
}

func TestListReportersSDK(t *testing.T) {
	am, srv := createAMServerWithDomain(t)
	defaultTenant(am).Reporters.Put(newChild(testReporter(), defaultDomainKey))
	client := newAMClient(t, srv)

	res, err := client.ListReportersWithResponse(t.Context(), defaultDomainKey)
	assertSDKOK(t, res, err, []sdk.Reporter{testReporterSDK()})
}

func TestGetReporterSDK(t *testing.T) {
	am, srv := createAMServerWithDomain(t)
	defaultTenant(am).Reporters.Put(newChild(testReporter(), defaultDomainKey))
	client := newAMClient(t, srv)

	res, err := client.GetReporterWithResponse(t.Context(), defaultDomainKey, "test")
	assertSDKOK(t, res, err, testReporterSDK())
}

func TestGetReporter404SDK(t *testing.T) {
	_, srv := createAMServerWithDomain(t)
	client := newAMClient(t, srv)

	res, err := client.GetReporterWithResponse(t.Context(), defaultDomainKey, "test")
	assertSDK404(t, res, err)
}

func TestPutGetReporterSDK(t *testing.T) {
	_, srv := createAMServerWithDomain(t)
	client := newAMClient(t, srv)
	body := testReporterSDK()

	put, err := client.UpsertReporterWithResponse(t.Context(), defaultDomainKey, body)
	assertSDKOK(t, put, err, body)

	get, err := client.GetReporterWithResponse(t.Context(), defaultDomainKey, "test")
	assertSDKOK(t, get, err, body)
}

func TestDeleteReporterSDK(t *testing.T) {
	am, srv := createAMServerWithDomain(t)
	defaultTenant(am).Reporters.Put(newChild(testReporter(), defaultDomainKey))
	client := newAMClient(t, srv)

	del, err := client.DeleteReporterWithResponse(t.Context(), defaultDomainKey, "test")
	assertSDKNoContent(t, del, err)

	get, err := client.GetReporterWithResponse(t.Context(), defaultDomainKey, "test")
	assertSDK404(t, get, err)
}
