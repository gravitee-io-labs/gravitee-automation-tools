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

	"github.com/gravitee-io-labs/gravitee-automation-tools/am-sdk/pkg/sdk/certificate"
)

func testCertificate() Certificate {
	return Certificate{Key: "test", Name: new("Test certificate")}
}

func testCertificateSDK() certificate.Certificate {
	return certificate.Certificate{Key: "test", Name: new("Test certificate")}
}

func TestListCertificates(t *testing.T) {
	am, srv := createAMServerWithDomain(t)
	defaultTenant(am).Certificates.Put(newChild(testCertificate(), defaultDomainKey))

	assertListEqual(t, certificatesURL(srv), []Certificate{testCertificate()})
}

func TestGetCertificate(t *testing.T) {
	am, srv := createAMServerWithDomain(t)
	defaultTenant(am).Certificates.Put(newChild(testCertificate(), defaultDomainKey))

	assertGetEqual(t, certificatesURL(srv), "test", testCertificate())
}

func TestGetCertificate404(t *testing.T) {
	_, srv := createAMServerWithDomain(t)

	assertGet404(t, certificatesURL(srv), "test", "Certificate")
}

func TestPutGetCertificate(t *testing.T) {
	_, srv := createAMServerWithDomain(t)
	body := testCertificate()

	assertPutEqual(t, certificatesURL(srv), body)
	assertGetEqual(t, certificatesURL(srv), "test", body)
}

func TestDeleteCertificate(t *testing.T) {
	am, srv := createAMServerWithDomain(t)
	defaultTenant(am).Certificates.Put(newChild(testCertificate(), defaultDomainKey))

	assertDeleteGone(t, certificatesURL(srv), "test", "Certificate")
}

func TestListCertificatesSDK(t *testing.T) {
	am, srv := createAMServerWithDomain(t)
	defaultTenant(am).Certificates.Put(newChild(testCertificate(), defaultDomainKey))
	client := newAMClient(t, srv)

	res, err := client.Certificates.ListCertificatesWithResponse(t.Context(), defaultDomainKey)
	assertSDKOK(t, res, err, []certificate.Certificate{testCertificateSDK()})
}

func TestGetCertificateSDK(t *testing.T) {
	am, srv := createAMServerWithDomain(t)
	defaultTenant(am).Certificates.Put(newChild(testCertificate(), defaultDomainKey))
	client := newAMClient(t, srv)

	res, err := client.Certificates.GetCertificateWithResponse(t.Context(), defaultDomainKey, "test")
	assertSDKOK(t, res, err, testCertificateSDK())
}

func TestGetCertificate404SDK(t *testing.T) {
	_, srv := createAMServerWithDomain(t)
	client := newAMClient(t, srv)

	res, err := client.Certificates.GetCertificateWithResponse(t.Context(), defaultDomainKey, "test")
	assertSDK404(t, res, err)
}

func TestPutGetCertificateSDK(t *testing.T) {
	_, srv := createAMServerWithDomain(t)
	client := newAMClient(t, srv)
	body := testCertificateSDK()

	put, err := client.Certificates.UpsertCertificateWithResponse(t.Context(), defaultDomainKey, body)
	assertSDKOK(t, put, err, body)

	get, err := client.Certificates.GetCertificateWithResponse(t.Context(), defaultDomainKey, "test")
	assertSDKOK(t, get, err, body)
}

func TestDeleteCertificateSDK(t *testing.T) {
	am, srv := createAMServerWithDomain(t)
	defaultTenant(am).Certificates.Put(newChild(testCertificate(), defaultDomainKey))
	client := newAMClient(t, srv)

	del, err := client.Certificates.DeleteCertificateWithResponse(t.Context(), defaultDomainKey, "test")
	assertSDKNoContent(t, del, err)

	get, err := client.Certificates.GetCertificateWithResponse(t.Context(), defaultDomainKey, "test")
	assertSDK404(t, get, err)
}
