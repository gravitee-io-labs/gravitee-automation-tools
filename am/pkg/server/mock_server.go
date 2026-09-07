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

func (m *MockAM) ListDomains(context.Context, ListDomainsRequestObject) (ListDomainsResponseObject, error) {
	return ListDomains200JSONResponse(m.Domains.GetAll()), nil
}

func (m *MockAM) UpsertDomain(_ context.Context, req UpsertDomainRequestObject) (UpsertDomainResponseObject, error) {
	if req.Body != nil {
		body := *req.Body
		m.Domains.Put(body)
		return UpsertDomain200JSONResponse(body), nil
	}
	return UpsertDomain400JSONResponse(emptyBodyError()), nil
}

func (m *MockAM) DeleteDomain(_ context.Context, req DeleteDomainRequestObject) (DeleteDomainResponseObject, error) {
	m.Domains.DeleteByKey(req.DomainKey)
	return DeleteDomain204Response{}, nil
}

func (m *MockAM) GetDomain(_ context.Context, req GetDomainRequestObject) (GetDomainResponseObject, error) {
	domain, ok := m.Domains.Get(req.DomainKey)
	if ok {
		return GetDomain200JSONResponse(domain), nil
	}
	return GetDomain404JSONResponse(notFoundError("Domain", req.DomainKey)), nil
}

func (m *MockAM) ListCertificates(context.Context, ListCertificatesRequestObject) (ListCertificatesResponseObject, error) {
	return ListCertificates200JSONResponse(m.Certificates.GetAll()), nil
}

func (m *MockAM) UpsertCertificate(_ context.Context, req UpsertCertificateRequestObject) (UpsertCertificateResponseObject, error) {
	if req.Body != nil {
		body := *req.Body
		m.Certificates.Put(body)
		return UpsertCertificate200JSONResponse(body), nil
	}
	return UpsertCertificate400JSONResponse(emptyBodyError()), nil
}

func (m *MockAM) DeleteCertificate(_ context.Context, req DeleteCertificateRequestObject) (DeleteCertificateResponseObject, error) {
	m.Certificates.DeleteByKey(req.CertKey)
	return DeleteCertificate204Response{}, nil
}

func (m *MockAM) GetCertificate(_ context.Context, req GetCertificateRequestObject) (GetCertificateResponseObject, error) {
	cert, ok := m.Certificates.Get(req.CertKey)
	if ok {
		return GetCertificate200JSONResponse(cert), nil
	}
	return GetCertificate404JSONResponse(notFoundError("Certificate", req.CertKey)), nil
}

func (m *MockAM) ListIdentityProviders(context.Context, ListIdentityProvidersRequestObject) (ListIdentityProvidersResponseObject, error) {
	return ListIdentityProviders200JSONResponse(m.IdentityProviders.GetAll()), nil
}

func (m *MockAM) UpsertIdentityProvider(_ context.Context, req UpsertIdentityProviderRequestObject) (UpsertIdentityProviderResponseObject, error) {
	if req.Body != nil {
		body := *req.Body
		m.IdentityProviders.Put(body)
		return UpsertIdentityProvider200JSONResponse(body), nil
	}
	return UpsertIdentityProvider400JSONResponse(emptyBodyError()), nil
}

func (m *MockAM) DeleteIdentityProvider(_ context.Context, req DeleteIdentityProviderRequestObject) (DeleteIdentityProviderResponseObject, error) {
	m.IdentityProviders.DeleteByKey(req.IdentityKey)
	return DeleteIdentityProvider204Response{}, nil
}

func (m *MockAM) GetIdentityProvider(_ context.Context, req GetIdentityProviderRequestObject) (GetIdentityProviderResponseObject, error) {
	idp, ok := m.IdentityProviders.Get(req.IdentityKey)
	if ok {
		return GetIdentityProvider200JSONResponse(idp), nil
	}
	return GetIdentityProvider404JSONResponse(notFoundError("IdentityProvider", req.IdentityKey)), nil
}

func (m *MockAM) ListReporters(context.Context, ListReportersRequestObject) (ListReportersResponseObject, error) {
	return ListReporters200JSONResponse(m.Reporters.GetAll()), nil
}

func (m *MockAM) UpsertReporter(_ context.Context, req UpsertReporterRequestObject) (UpsertReporterResponseObject, error) {
	if req.Body != nil {
		body := *req.Body
		m.Reporters.Put(body)
		return UpsertReporter200JSONResponse(body), nil
	}
	return UpsertReporter400JSONResponse(emptyBodyError()), nil
}

func (m *MockAM) DeleteReporter(_ context.Context, req DeleteReporterRequestObject) (DeleteReporterResponseObject, error) {
	m.Reporters.DeleteByKey(req.ReporterKey)
	return DeleteReporter204Response{}, nil
}

func (m *MockAM) GetReporter(_ context.Context, req GetReporterRequestObject) (GetReporterResponseObject, error) {
	reporter, ok := m.Reporters.Get(req.ReporterKey)
	if ok {
		return GetReporter200JSONResponse(reporter), nil
	}
	return GetReporter404JSONResponse(notFoundError("Reporter", req.ReporterKey)), nil
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
