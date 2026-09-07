package server

import (
	"testing"

	"github.com/gravitee-io-labs/gravitee-automation-sdks/am/pkg/sdk/reporter"
)

func testReporter() Reporter {
	return Reporter{Key: "test", Name: new("Test reporter")}
}

func testReporterSDK() reporter.Reporter {
	return reporter.Reporter{Key: "test", Name: new("Test reporter")}
}

func TestListReporters(t *testing.T) {
	am, srv := createAMServer(t)
	am.Reporters.Put(testReporter())

	assertListEqual(t, reportersURL(srv), []Reporter{testReporter()})
}

func TestGetReporter(t *testing.T) {
	am, srv := createAMServer(t)
	am.Reporters.Put(testReporter())

	assertGetEqual(t, reportersURL(srv), "test", testReporter())
}

func TestGetReporter404(t *testing.T) {
	_, srv := createAMServer(t)

	assertGet404(t, reportersURL(srv), "test", "Reporter")
}

func TestPutGetReporter(t *testing.T) {
	_, srv := createAMServer(t)
	body := testReporter()

	assertPutEqual(t, reportersURL(srv), body)
	assertGetEqual(t, reportersURL(srv), "test", body)
}

func TestDeleteReporter(t *testing.T) {
	am, srv := createAMServer(t)
	am.Reporters.Put(testReporter())

	assertDeleteGone(t, reportersURL(srv), "test", "Reporter")
}

func TestListReportersSDK(t *testing.T) {
	am, srv := createAMServer(t)
	am.Reporters.Put(testReporter())
	client := newAMClient(t, srv)

	res, err := client.Reporters.ListReportersWithResponse(t.Context(), defaultDomainKey)
	assertSDKOK(t, res, err, []reporter.Reporter{testReporterSDK()})
}

func TestGetReporterSDK(t *testing.T) {
	am, srv := createAMServer(t)
	am.Reporters.Put(testReporter())
	client := newAMClient(t, srv)

	res, err := client.Reporters.GetReporterWithResponse(t.Context(), defaultDomainKey, "test")
	assertSDKOK(t, res, err, testReporterSDK())
}

func TestGetReporter404SDK(t *testing.T) {
	_, srv := createAMServer(t)
	client := newAMClient(t, srv)

	res, err := client.Reporters.GetReporterWithResponse(t.Context(), defaultDomainKey, "test")
	assertSDK404(t, res, err)
}

func TestPutGetReporterSDK(t *testing.T) {
	_, srv := createAMServer(t)
	client := newAMClient(t, srv)
	body := testReporterSDK()

	put, err := client.Reporters.UpsertReporterWithResponse(t.Context(), defaultDomainKey, body)
	assertSDKOK(t, put, err, body)

	get, err := client.Reporters.GetReporterWithResponse(t.Context(), defaultDomainKey, "test")
	assertSDKOK(t, get, err, body)
}

func TestDeleteReporterSDK(t *testing.T) {
	am, srv := createAMServer(t)
	am.Reporters.Put(testReporter())
	client := newAMClient(t, srv)

	del, err := client.Reporters.DeleteReporterWithResponse(t.Context(), defaultDomainKey, "test")
	assertSDKNoContent(t, del, err)

	get, err := client.Reporters.GetReporterWithResponse(t.Context(), defaultDomainKey, "test")
	assertSDK404(t, get, err)
}
