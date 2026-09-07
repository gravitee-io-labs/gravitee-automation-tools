package server

import (
	"testing"

	"github.com/gravitee-io-labs/gravitee-automation-sdks/am/pkg/sdk/certificate"
)

func testCertificate() Certificate {
	return Certificate{Key: "test", Name: new("Test certificate")}
}

func testCertificateSDK() certificate.Certificate {
	return certificate.Certificate{Key: "test", Name: new("Test certificate")}
}

func TestListCertificates(t *testing.T) {
	am, srv := createAMServer(t)
	am.Certificates.Put(testCertificate())

	assertListEqual(t, certificatesURL(srv), []Certificate{testCertificate()})
}

func TestGetCertificate(t *testing.T) {
	am, srv := createAMServer(t)
	am.Certificates.Put(testCertificate())

	assertGetEqual(t, certificatesURL(srv), "test", testCertificate())
}

func TestGetCertificate404(t *testing.T) {
	_, srv := createAMServer(t)

	assertGet404(t, certificatesURL(srv), "test", "Certificate")
}

func TestPutGetCertificate(t *testing.T) {
	_, srv := createAMServer(t)
	body := testCertificate()

	assertPutEqual(t, certificatesURL(srv), body)
	assertGetEqual(t, certificatesURL(srv), "test", body)
}

func TestDeleteCertificate(t *testing.T) {
	am, srv := createAMServer(t)
	am.Certificates.Put(testCertificate())

	assertDeleteGone(t, certificatesURL(srv), "test", "Certificate")
}

func TestListCertificatesSDK(t *testing.T) {
	am, srv := createAMServer(t)
	am.Certificates.Put(testCertificate())
	client := newAMClient(t, srv)

	res, err := client.Certificates.ListCertificatesWithResponse(t.Context(), defaultDomainKey)
	assertSDKOK(t, res, err, []certificate.Certificate{testCertificateSDK()})
}

func TestGetCertificateSDK(t *testing.T) {
	am, srv := createAMServer(t)
	am.Certificates.Put(testCertificate())
	client := newAMClient(t, srv)

	res, err := client.Certificates.GetCertificateWithResponse(t.Context(), defaultDomainKey, "test")
	assertSDKOK(t, res, err, testCertificateSDK())
}

func TestGetCertificate404SDK(t *testing.T) {
	_, srv := createAMServer(t)
	client := newAMClient(t, srv)

	res, err := client.Certificates.GetCertificateWithResponse(t.Context(), defaultDomainKey, "test")
	assertSDK404(t, res, err)
}

func TestPutGetCertificateSDK(t *testing.T) {
	_, srv := createAMServer(t)
	client := newAMClient(t, srv)
	body := testCertificateSDK()

	put, err := client.Certificates.UpsertCertificateWithResponse(t.Context(), defaultDomainKey, body)
	assertSDKOK(t, put, err, body)

	get, err := client.Certificates.GetCertificateWithResponse(t.Context(), defaultDomainKey, "test")
	assertSDKOK(t, get, err, body)
}

func TestDeleteCertificateSDK(t *testing.T) {
	am, srv := createAMServer(t)
	am.Certificates.Put(testCertificate())
	client := newAMClient(t, srv)

	del, err := client.Certificates.DeleteCertificateWithResponse(t.Context(), defaultDomainKey, "test")
	assertSDKNoContent(t, del, err)

	get, err := client.Certificates.GetCertificateWithResponse(t.Context(), defaultDomainKey, "test")
	assertSDK404(t, get, err)
}
