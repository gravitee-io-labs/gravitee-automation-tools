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

	"github.com/gravitee-io-labs/gravitee-automation-sdks/am/pkg/sdk/certificate"
	"github.com/gravitee-io-labs/gravitee-automation-sdks/am/pkg/sdk/domain"
	"github.com/gravitee-io-labs/gravitee-automation-sdks/am/pkg/sdk/identityprovider"
	"github.com/gravitee-io-labs/gravitee-automation-sdks/am/pkg/sdk/reporter"
)

func TestNested_IsolatedByDomain(t *testing.T) {
	_, srv := createAMServer(t)
	client := newAMClient(t, srv)

	d1 := domain.Domain{Key: "dom-a", Name: "A"}
	d2 := domain.Domain{Key: "dom-b", Name: "B"}
	put, err := client.Domains.UpsertDomainWithResponse(t.Context(), d1)
	assertSDKOK(t, put, err, d1)
	put, err = client.Domains.UpsertDomainWithResponse(t.Context(), d2)
	assertSDKOK(t, put, err, d2)

	certA := certificate.Certificate{Key: "cert", Name: new("A cert")}
	certB := certificate.Certificate{Key: "cert", Name: new("B cert")}
	up, err := client.Certificates.UpsertCertificateWithResponse(t.Context(), "dom-a", certA)
	assertSDKOK(t, up, err, certA)
	up, err = client.Certificates.UpsertCertificateWithResponse(t.Context(), "dom-b", certB)
	assertSDKOK(t, up, err, certB)

	listA, err := client.Certificates.ListCertificatesWithResponse(t.Context(), "dom-a")
	assertSDKOK(t, listA, err, []certificate.Certificate{certA})
	listB, err := client.Certificates.ListCertificatesWithResponse(t.Context(), "dom-b")
	assertSDKOK(t, listB, err, []certificate.Certificate{certB})

	getA, err := client.Certificates.GetCertificateWithResponse(t.Context(), "dom-a", "cert")
	assertSDKOK(t, getA, err, certA)
	getB, err := client.Certificates.GetCertificateWithResponse(t.Context(), "dom-b", "cert")
	assertSDKOK(t, getB, err, certB)
}

func TestNested_DeleteDomainCascades(t *testing.T) {
	_, srv := createAMServer(t)
	client := newAMClient(t, srv)

	keep := domain.Domain{Key: "keep", Name: "Keep"}
	drop := domain.Domain{Key: "drop", Name: "Drop"}
	put, err := client.Domains.UpsertDomainWithResponse(t.Context(), keep)
	assertSDKOK(t, put, err, keep)
	put, err = client.Domains.UpsertDomainWithResponse(t.Context(), drop)
	assertSDKOK(t, put, err, drop)

	certKeep := certificate.Certificate{Key: "ck", Name: new("keep cert")}
	certDrop := certificate.Certificate{Key: "cd", Name: new("drop cert")}
	idpDrop := identityprovider.IdentityProvider{Key: "id", Name: new("drop idp")}
	repDrop := reporter.Reporter{Key: "rd", Name: new("drop reporter")}

	up, err := client.Certificates.UpsertCertificateWithResponse(t.Context(), "keep", certKeep)
	assertSDKOK(t, up, err, certKeep)
	up, err = client.Certificates.UpsertCertificateWithResponse(t.Context(), "drop", certDrop)
	assertSDKOK(t, up, err, certDrop)
	idp, err := client.IdentityProviders.UpsertIdentityProviderWithResponse(t.Context(), "drop", idpDrop)
	assertSDKOK(t, idp, err, idpDrop)
	rep, err := client.Reporters.UpsertReporterWithResponse(t.Context(), "drop", repDrop)
	assertSDKOK(t, rep, err, repDrop)

	del, err := client.Domains.DeleteDomainWithResponse(t.Context(), "drop")
	assertSDKNoContent(t, del, err)

	gone, err := client.Certificates.GetCertificateWithResponse(t.Context(), "drop", "cd")
	assertSDK404(t, gone, err)
	goneIDP, err := client.IdentityProviders.GetIdentityProviderWithResponse(t.Context(), "drop", "id")
	assertSDK404(t, goneIDP, err)
	goneRep, err := client.Reporters.GetReporterWithResponse(t.Context(), "drop", "rd")
	assertSDK404(t, goneRep, err)

	still, err := client.Certificates.GetCertificateWithResponse(t.Context(), "keep", "ck")
	assertSDKOK(t, still, err, certKeep)
}

func TestNested_MissingDomain404_WithAuth(t *testing.T) {
	_, srv := createAMServerWithAuth(t)
	client := newTestClient(t, srv, "", "", "admin-token")

	list, err := client.Certificates.ListCertificatesWithResponse(t.Context(), "missing")
	assertSDK404(t, list, err)
	get, err := client.Certificates.GetCertificateWithResponse(t.Context(), "missing", "any")
	assertSDK404(t, get, err)
}

func TestNested_DeleteMissing_204(t *testing.T) {
	_, srv := createAMServer(t)
	client := newAMClient(t, srv)

	del, err := client.Domains.DeleteDomainWithResponse(t.Context(), "no-such")
	assertSDKNoContent(t, del, err)
	delCert, err := client.Certificates.DeleteCertificateWithResponse(t.Context(), "no-such", "no-such")
	assertSDK404(t, delCert, err)
}
