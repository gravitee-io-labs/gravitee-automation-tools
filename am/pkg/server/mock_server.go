package server

import (
	"context"
	"net/http"

	"github.com/gravitee-io-labs/gravitee-automation-sdks/am/pkg/server/store"
)

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
	return AutomationCreateOrUpdateDomain400JSONResponse(emptyBodyError()), nil
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
	return AutomationGetDomain404JSONResponse(notFoundError("Domain", req.DomainKey)), nil
}

func (m *MockAM) AutomationListCertificates(context.Context, AutomationListCertificatesRequestObject) (AutomationListCertificatesResponseObject, error) {
	return AutomationListCertificates200JSONResponse(m.Certificates.GetAll()), nil
}

func (m *MockAM) AutomationCreateOrUpdateCertificate(_ context.Context, req AutomationCreateOrUpdateCertificateRequestObject) (AutomationCreateOrUpdateCertificateResponseObject, error) {
	if req.Body != nil {
		body := *req.Body
		m.Certificates.Put(body)
		return AutomationCreateOrUpdateCertificate200JSONResponse(body), nil
	}
	return AutomationCreateOrUpdateCertificate400JSONResponse(emptyBodyError()), nil
}

func (m *MockAM) AutomationDeleteCertificate(_ context.Context, req AutomationDeleteCertificateRequestObject) (AutomationDeleteCertificateResponseObject, error) {
	m.Certificates.DeleteByKey(req.CertKey)
	return AutomationDeleteCertificate204Response{}, nil
}

func (m *MockAM) AutomationGetCertificate(_ context.Context, req AutomationGetCertificateRequestObject) (AutomationGetCertificateResponseObject, error) {
	cert, ok := m.Certificates.Get(req.CertKey)
	if ok {
		return AutomationGetCertificate200JSONResponse(cert), nil
	}
	return AutomationGetCertificate404JSONResponse(notFoundError("Certificate", req.CertKey)), nil
}

func (m *MockAM) AutomationListIdentityProviders(context.Context, AutomationListIdentityProvidersRequestObject) (AutomationListIdentityProvidersResponseObject, error) {
	return AutomationListIdentityProviders200JSONResponse(m.IdentityProviders.GetAll()), nil
}

func (m *MockAM) AutomationCreateOrUpdateIdentityProvider(_ context.Context, req AutomationCreateOrUpdateIdentityProviderRequestObject) (AutomationCreateOrUpdateIdentityProviderResponseObject, error) {
	if req.Body != nil {
		body := *req.Body
		m.IdentityProviders.Put(body)
		return AutomationCreateOrUpdateIdentityProvider200JSONResponse(body), nil
	}
	return AutomationCreateOrUpdateIdentityProvider400JSONResponse(emptyBodyError()), nil
}

func (m *MockAM) AutomationDeleteIdentityProvider(_ context.Context, req AutomationDeleteIdentityProviderRequestObject) (AutomationDeleteIdentityProviderResponseObject, error) {
	m.IdentityProviders.DeleteByKey(req.IdentityKey)
	return AutomationDeleteIdentityProvider204Response{}, nil
}

func (m *MockAM) AutomationGetIdentityProvider(_ context.Context, req AutomationGetIdentityProviderRequestObject) (AutomationGetIdentityProviderResponseObject, error) {
	idp, ok := m.IdentityProviders.Get(req.IdentityKey)
	if ok {
		return AutomationGetIdentityProvider200JSONResponse(idp), nil
	}
	return AutomationGetIdentityProvider404JSONResponse(notFoundError("IdentityProvider", req.IdentityKey)), nil
}

func (m *MockAM) AutomationListReporters(context.Context, AutomationListReportersRequestObject) (AutomationListReportersResponseObject, error) {
	return AutomationListReporters200JSONResponse(m.Reporters.GetAll()), nil
}

func (m *MockAM) AutomationCreateOrUpdateReporter(_ context.Context, req AutomationCreateOrUpdateReporterRequestObject) (AutomationCreateOrUpdateReporterResponseObject, error) {
	if req.Body != nil {
		body := *req.Body
		m.Reporters.Put(body)
		return AutomationCreateOrUpdateReporter200JSONResponse(body), nil
	}
	return AutomationCreateOrUpdateReporter400JSONResponse(emptyBodyError()), nil
}

func (m *MockAM) AutomationDeleteReporter(_ context.Context, req AutomationDeleteReporterRequestObject) (AutomationDeleteReporterResponseObject, error) {
	m.Reporters.DeleteByKey(req.ReporterKey)
	return AutomationDeleteReporter204Response{}, nil
}

func (m *MockAM) AutomationGetReporter(_ context.Context, req AutomationGetReporterRequestObject) (AutomationGetReporterResponseObject, error) {
	reporter, ok := m.Reporters.Get(req.ReporterKey)
	if ok {
		return AutomationGetReporter200JSONResponse(reporter), nil
	}
	return AutomationGetReporter404JSONResponse(notFoundError("Reporter", req.ReporterKey)), nil
}

func emptyBodyError() Error {
	return Error{
		HttpStatus: new(int32(http.StatusBadRequest)),
		Message:    new("Empty body"),
	}
}

func notFoundError(kind, key string) Error {
	return Error{
		HttpStatus: new(int32(http.StatusNotFound)),
		Message:    new(kind + " [" + key + "] not found"),
	}
}
