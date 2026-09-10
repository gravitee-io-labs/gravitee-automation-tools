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
	"context"
	"net/http"

	"github.com/gravitee-io-labs/gravitee-automation-sdks/common/pkg/store"
)

type MockAM struct {
	ctx     context.Context
	tenants map[string]*tenant
}

type tenant struct {
	Domains           *store.Store[Domain]
	Certificates      *store.Store[Certificate]
	IdentityProviders *store.Store[IdentityProvider]
	Reporters         *store.Store[Reporter]
}

func newTenant(ctx context.Context) *tenant {
	return &tenant{
		Domains:           store.NewStore[Domain](ctx),
		Certificates:      store.NewStore[Certificate](ctx),
		IdentityProviders: store.NewStore[IdentityProvider](ctx),
		Reporters:         store.NewStore[Reporter](ctx),
	}
}

func (m *MockAM) getTenant(aware store.OrgEnvAware) *tenant {
	tenantKey := aware.GetOrgId() + "-" + aware.GetEnvId()
	et, ok := m.tenants[tenantKey]
	if !ok {
		nt := newTenant(m.ctx)
		m.tenants[tenantKey] = nt
		return nt
	}
	return et
}

var _ StrictServerInterface = &MockAM{}

type orgEnv struct {
	org string
	env string
}

func (o orgEnv) GetOrgId() string { return o.org }
func (o orgEnv) GetEnvId() string { return o.env }

func NewMockAM(ctx context.Context) *MockAM {
	return &MockAM{
		ctx:     ctx,
		tenants: make(map[string]*tenant),
	}
}

func (m *MockAM) ListDomains(_ context.Context, r ListDomainsRequestObject) (ListDomainsResponseObject, error) {
	return ListDomains200JSONResponse(m.getTenant(r).Domains.GetAll()), nil
}

func (m *MockAM) UpsertDomain(_ context.Context, req UpsertDomainRequestObject) (UpsertDomainResponseObject, error) {
	if req.Body != nil {
		body := *req.Body
		m.getTenant(req).Domains.Put(body)
		return UpsertDomain200JSONResponse(body), nil
	}
	return UpsertDomain400JSONResponse(emptyBodyError()), nil
}

func (m *MockAM) DeleteDomain(_ context.Context, req DeleteDomainRequestObject) (DeleteDomainResponseObject, error) {
	m.getTenant(req).Domains.DeleteByKey(req.DomainKey)
	return DeleteDomain204Response{}, nil
}

func (m *MockAM) GetDomain(_ context.Context, req GetDomainRequestObject) (GetDomainResponseObject, error) {
	domain, ok := m.getTenant(req).Domains.Get(req.DomainKey)
	if ok {
		return GetDomain200JSONResponse(domain), nil
	}
	return GetDomain404JSONResponse(notFoundError("Domain", req.DomainKey)), nil
}

func (m *MockAM) ListCertificates(_ context.Context, req ListCertificatesRequestObject) (ListCertificatesResponseObject, error) {
	return ListCertificates200JSONResponse(m.getTenant(req).Certificates.GetAll()), nil
}

func (m *MockAM) UpsertCertificate(_ context.Context, req UpsertCertificateRequestObject) (UpsertCertificateResponseObject, error) {
	if req.Body != nil {
		body := *req.Body
		m.getTenant(req).Certificates.Put(body)
		return UpsertCertificate200JSONResponse(body), nil
	}
	return UpsertCertificate400JSONResponse(emptyBodyError()), nil
}

func (m *MockAM) DeleteCertificate(_ context.Context, req DeleteCertificateRequestObject) (DeleteCertificateResponseObject, error) {
	m.getTenant(req).Certificates.DeleteByKey(req.CertKey)
	return DeleteCertificate204Response{}, nil
}

func (m *MockAM) GetCertificate(_ context.Context, req GetCertificateRequestObject) (GetCertificateResponseObject, error) {
	cert, ok := m.getTenant(req).Certificates.Get(req.CertKey)
	if ok {
		return GetCertificate200JSONResponse(cert), nil
	}
	return GetCertificate404JSONResponse(notFoundError("Certificate", req.CertKey)), nil
}

func (m *MockAM) ListIdentityProviders(_ context.Context, req ListIdentityProvidersRequestObject) (ListIdentityProvidersResponseObject, error) {
	return ListIdentityProviders200JSONResponse(m.getTenant(req).IdentityProviders.GetAll()), nil
}

func (m *MockAM) UpsertIdentityProvider(_ context.Context, req UpsertIdentityProviderRequestObject) (UpsertIdentityProviderResponseObject, error) {
	if req.Body != nil {
		body := *req.Body
		m.getTenant(req).IdentityProviders.Put(body)
		return UpsertIdentityProvider200JSONResponse(body), nil
	}
	return UpsertIdentityProvider400JSONResponse(emptyBodyError()), nil
}

func (m *MockAM) DeleteIdentityProvider(_ context.Context, req DeleteIdentityProviderRequestObject) (DeleteIdentityProviderResponseObject, error) {
	m.getTenant(req).IdentityProviders.DeleteByKey(req.IdentityKey)
	return DeleteIdentityProvider204Response{}, nil
}

func (m *MockAM) GetIdentityProvider(_ context.Context, req GetIdentityProviderRequestObject) (GetIdentityProviderResponseObject, error) {
	idp, ok := m.getTenant(req).IdentityProviders.Get(req.IdentityKey)
	if ok {
		return GetIdentityProvider200JSONResponse(idp), nil
	}
	return GetIdentityProvider404JSONResponse(notFoundError("IdentityProvider", req.IdentityKey)), nil
}

func (m *MockAM) ListReporters(_ context.Context, req ListReportersRequestObject) (ListReportersResponseObject, error) {
	return ListReporters200JSONResponse(m.getTenant(req).Reporters.GetAll()), nil
}

func (m *MockAM) UpsertReporter(_ context.Context, req UpsertReporterRequestObject) (UpsertReporterResponseObject, error) {
	if req.Body != nil {
		body := *req.Body
		m.getTenant(req).Reporters.Put(body)
		return UpsertReporter200JSONResponse(body), nil
	}
	return UpsertReporter400JSONResponse(emptyBodyError()), nil
}

func (m *MockAM) DeleteReporter(_ context.Context, req DeleteReporterRequestObject) (DeleteReporterResponseObject, error) {
	m.getTenant(req).Reporters.DeleteByKey(req.ReporterKey)
	return DeleteReporter204Response{}, nil
}

func (m *MockAM) GetReporter(_ context.Context, req GetReporterRequestObject) (GetReporterResponseObject, error) {
	reporter, ok := m.getTenant(req).Reporters.Get(req.ReporterKey)
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
