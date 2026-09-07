package server

import (
	"context"
	"errors"
	"net/http"

	"github.com/gravitee-io-labs/gravitee-automation-sdks/am/pkg/server/store"
)

var ErrNotImplemented = errors.New("not implemented")

type MockAM struct {
	Domains           *store.Store[Domain]
	Certificates      *store.Store[Certificate]
	IdentityProviders *store.Store[IdentityProvider]
	Reporters         *store.Store[Reporter]
}

var _ StrictServerInterface = &MockAM{}

func NewMockAM(ctx context.Context) *MockAM {
	return &MockAM{
		Domains:           store.NewStore[Domain](ctx),
		Certificates:      store.NewStore[Certificate](ctx),
		IdentityProviders: store.NewStore[IdentityProvider](ctx),
		Reporters:         store.NewStore[Reporter](ctx),
	}
}

func (m *MockAM) AutomationListDomains(context.Context, AutomationListDomainsRequestObject) (AutomationListDomainsResponseObject, error) {
	return AutomationListDomains200JSONResponse(m.Domains.GetAll()), nil
}

func (m *MockAM) AutomationCreateOrUpdateDomain(_ context.Context, req AutomationCreateOrUpdateDomainRequestObject) (AutomationCreateOrUpdateDomainResponseObject, error) {
	if req.Body != nil {
		body := *req.Body
		m.Domains.Put(body)
		return AutomationCreateOrUpdateDomain200JSONResponse(body), nil
	}
	return AutomationCreateOrUpdateDomain400JSONResponse(Error{
		HttpStatus: new(int32(http.StatusBadRequest)),
		Message:    new("Empty body"),
	}), nil
}

func (m *MockAM) AutomationDeleteDomain(_ context.Context, req AutomationDeleteDomainRequestObject) (AutomationDeleteDomainResponseObject, error) {
	m.Domains.DeleteByKey(req.DomainKey)
	return AutomationDeleteDomain204Response{}, nil
}

func (m *MockAM) AutomationGetDomain(_ context.Context, req AutomationGetDomainRequestObject) (AutomationGetDomainResponseObject, error) {
	domain, ok := m.Domains.Get(req.DomainKey)
	if ok {
		return AutomationGetDomain200JSONResponse(domain), nil
	}
	return AutomationGetDomain404JSONResponse{
		HttpStatus: new(int32(http.StatusNotFound)),
		Message:    new("Domain [" + req.DomainKey + "] not found"),
	}, nil
}

func (MockAM) AutomationListCertificates(context.Context, AutomationListCertificatesRequestObject) (AutomationListCertificatesResponseObject, error) {
	return nil, ErrNotImplemented
}

func (MockAM) AutomationCreateOrUpdateCertificate(context.Context, AutomationCreateOrUpdateCertificateRequestObject) (AutomationCreateOrUpdateCertificateResponseObject, error) {
	return nil, ErrNotImplemented
}

func (MockAM) AutomationDeleteCertificate(context.Context, AutomationDeleteCertificateRequestObject) (AutomationDeleteCertificateResponseObject, error) {
	return nil, ErrNotImplemented
}

func (MockAM) AutomationGetCertificate(context.Context, AutomationGetCertificateRequestObject) (AutomationGetCertificateResponseObject, error) {
	return nil, ErrNotImplemented
}

func (MockAM) AutomationListIdentityProviders(context.Context, AutomationListIdentityProvidersRequestObject) (AutomationListIdentityProvidersResponseObject, error) {
	return nil, ErrNotImplemented
}

func (MockAM) AutomationCreateOrUpdateIdentityProvider(context.Context, AutomationCreateOrUpdateIdentityProviderRequestObject) (AutomationCreateOrUpdateIdentityProviderResponseObject, error) {
	return nil, ErrNotImplemented
}

func (MockAM) AutomationDeleteIdentityProvider(context.Context, AutomationDeleteIdentityProviderRequestObject) (AutomationDeleteIdentityProviderResponseObject, error) {
	return nil, ErrNotImplemented
}

func (MockAM) AutomationGetIdentityProvider(context.Context, AutomationGetIdentityProviderRequestObject) (AutomationGetIdentityProviderResponseObject, error) {
	return nil, ErrNotImplemented
}

func (MockAM) AutomationListReporters(context.Context, AutomationListReportersRequestObject) (AutomationListReportersResponseObject, error) {
	return nil, ErrNotImplemented
}

func (MockAM) AutomationCreateOrUpdateReporter(context.Context, AutomationCreateOrUpdateReporterRequestObject) (AutomationCreateOrUpdateReporterResponseObject, error) {
	return nil, ErrNotImplemented
}

func (MockAM) AutomationDeleteReporter(context.Context, AutomationDeleteReporterRequestObject) (AutomationDeleteReporterResponseObject, error) {
	return nil, ErrNotImplemented
}

func (MockAM) AutomationGetReporter(context.Context, AutomationGetReporterRequestObject) (AutomationGetReporterResponseObject, error) {
	return nil, ErrNotImplemented
}
