package server

import (
	"testing"

	"github.com/gravitee-io-labs/gravitee-automation-sdks/am/pkg/sdk/identityprovider"
)

func testIdentityProvider() IdentityProvider {
	return IdentityProvider{Key: "test", Name: new("Test identity provider")}
}

func testIdentityProviderSDK() identityprovider.IdentityProvider {
	return identityprovider.IdentityProvider{Key: "test", Name: new("Test identity provider")}
}

func TestListIdentityProviders(t *testing.T) {
	am, srv := createAMServer(t)
	defaultTenant(am).IdentityProviders.Put(testIdentityProvider())

	assertListEqual(t, identitiesURL(srv), []IdentityProvider{testIdentityProvider()})
}

func TestGetIdentityProvider(t *testing.T) {
	am, srv := createAMServer(t)
	defaultTenant(am).IdentityProviders.Put(testIdentityProvider())

	assertGetEqual(t, identitiesURL(srv), "test", testIdentityProvider())
}

func TestGetIdentityProvider404(t *testing.T) {
	_, srv := createAMServer(t)

	assertGet404(t, identitiesURL(srv), "test", "IdentityProvider")
}

func TestPutGetIdentityProvider(t *testing.T) {
	_, srv := createAMServer(t)
	body := testIdentityProvider()

	assertPutEqual(t, identitiesURL(srv), body)
	assertGetEqual(t, identitiesURL(srv), "test", body)
}

func TestDeleteIdentityProvider(t *testing.T) {
	am, srv := createAMServer(t)
	defaultTenant(am).IdentityProviders.Put(testIdentityProvider())

	assertDeleteGone(t, identitiesURL(srv), "test", "IdentityProvider")
}

func TestListIdentityProvidersSDK(t *testing.T) {
	am, srv := createAMServer(t)
	defaultTenant(am).IdentityProviders.Put(testIdentityProvider())
	client := newAMClient(t, srv)

	res, err := client.IdentityProviders.ListIdentityProvidersWithResponse(t.Context(), defaultDomainKey)
	assertSDKOK(t, res, err, []identityprovider.IdentityProvider{testIdentityProviderSDK()})
}

func TestGetIdentityProviderSDK(t *testing.T) {
	am, srv := createAMServer(t)
	defaultTenant(am).IdentityProviders.Put(testIdentityProvider())
	client := newAMClient(t, srv)

	res, err := client.IdentityProviders.GetIdentityProviderWithResponse(t.Context(), defaultDomainKey, "test")
	assertSDKOK(t, res, err, testIdentityProviderSDK())
}

func TestGetIdentityProvider404SDK(t *testing.T) {
	_, srv := createAMServer(t)
	client := newAMClient(t, srv)

	res, err := client.IdentityProviders.GetIdentityProviderWithResponse(t.Context(), defaultDomainKey, "test")
	assertSDK404(t, res, err)
}

func TestPutGetIdentityProviderSDK(t *testing.T) {
	_, srv := createAMServer(t)
	client := newAMClient(t, srv)
	body := testIdentityProviderSDK()

	put, err := client.IdentityProviders.UpsertIdentityProviderWithResponse(t.Context(), defaultDomainKey, body)
	assertSDKOK(t, put, err, body)

	get, err := client.IdentityProviders.GetIdentityProviderWithResponse(t.Context(), defaultDomainKey, "test")
	assertSDKOK(t, get, err, body)
}

func TestDeleteIdentityProviderSDK(t *testing.T) {
	am, srv := createAMServer(t)
	defaultTenant(am).IdentityProviders.Put(testIdentityProvider())
	client := newAMClient(t, srv)

	del, err := client.IdentityProviders.DeleteIdentityProviderWithResponse(t.Context(), defaultDomainKey, "test")
	assertSDKNoContent(t, del, err)

	get, err := client.IdentityProviders.GetIdentityProviderWithResponse(t.Context(), defaultDomainKey, "test")
	assertSDK404(t, get, err)
}
