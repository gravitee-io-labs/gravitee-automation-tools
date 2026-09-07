package server

import (
	"testing"

	"github.com/gravitee-io-labs/gravitee-automation-sdks/am/pkg/sdk/domain"
)

func TestListDomains(t *testing.T) {
	am, srv := createAMServer(t)
	am.Domains.Put(Domain{Key: "test", Name: "Test domain"})

	assertListEqual(t, domainsURL(srv), []Domain{{Key: "test", Name: "Test domain"}})
}

func TestGetDomain(t *testing.T) {
	am, srv := createAMServer(t)
	am.Domains.Put(Domain{Key: "test", Name: "Test domain"})

	assertGetEqual(t, domainsURL(srv), "test", Domain{Key: "test", Name: "Test domain"})
}

func TestGetDomain404(t *testing.T) {
	_, srv := createAMServer(t)

	assertGet404(t, domainsURL(srv), "test", "Domain")
}

func TestPutGetDomain(t *testing.T) {
	_, srv := createAMServer(t)
	body := Domain{Key: "test", Name: "Test domain"}

	assertPutEqual(t, domainsURL(srv), body)
	assertGetEqual(t, domainsURL(srv), "test", body)
}

func TestDeleteDomain(t *testing.T) {
	am, srv := createAMServer(t)
	am.Domains.Put(Domain{Key: "test", Name: "Test domain"})

	assertDeleteGone(t, domainsURL(srv), "test", "Domain")
}

func TestListDomainsSDK(t *testing.T) {
	am, srv := createAMServer(t)
	am.Domains.Put(Domain{Key: "test", Name: "Test domain"})
	client := newAMClient(t, srv)

	res, err := client.Domains.AutomationListDomainsWithResponse(t.Context())
	assertSDKOK(t, res, err, []domain.Domain{{Key: "test", Name: "Test domain"}})
}

func TestGetDomainSDK(t *testing.T) {
	am, srv := createAMServer(t)
	am.Domains.Put(Domain{Key: "test", Name: "Test domain"})
	client := newAMClient(t, srv)

	res, err := client.Domains.AutomationGetDomainWithResponse(t.Context(), "test")
	assertSDKOK(t, res, err, domain.Domain{Key: "test", Name: "Test domain"})
}

func TestGetDomain404SDK(t *testing.T) {
	_, srv := createAMServer(t)
	client := newAMClient(t, srv)

	res, err := client.Domains.AutomationGetDomainWithResponse(t.Context(), "test")
	assertSDK404(t, res, err)
}

func TestPutGetDomainSDK(t *testing.T) {
	_, srv := createAMServer(t)
	client := newAMClient(t, srv)
	body := domain.Domain{Key: "test", Name: "Test domain"}

	put, err := client.Domains.AutomationCreateOrUpdateDomainWithResponse(t.Context(), body)
	assertSDKOK(t, put, err, body)

	get, err := client.Domains.AutomationGetDomainWithResponse(t.Context(), "test")
	assertSDKOK(t, get, err, body)
}

func TestDeleteDomainSDK(t *testing.T) {
	am, srv := createAMServer(t)
	am.Domains.Put(Domain{Key: "test", Name: "Test domain"})
	client := newAMClient(t, srv)

	del, err := client.Domains.AutomationDeleteDomainWithResponse(t.Context(), "test")
	assertSDKNoContent(t, del, err)

	get, err := client.Domains.AutomationGetDomainWithResponse(t.Context(), "test")
	assertSDK404(t, get, err)
}
