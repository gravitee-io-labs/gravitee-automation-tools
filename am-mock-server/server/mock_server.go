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
	"sync"

	"github.com/gravitee-io-labs/gravitee-automation-sdks/common/pkg/store"
)

// MockAM is an in-memory StrictServer. Tenants are created on first request, keyed by org+env. Data is process-local.
type MockAM struct {
	tenants map[string]*tenant
	mutex   sync.Mutex
}

type Child[T store.Identifiable] struct {
	ParentIdentity string
	Self           T
}

func (m Child[T]) Identity() string {
	return m.Self.Identity()
}

type tenant struct {
	Domains           *store.Store[Domain]
	Certificates      *store.Store[Child[Certificate]]
	IdentityProviders *store.Store[Child[IdentityProvider]]
	Reporters         *store.Store[Child[Reporter]]
}

func newTenant() *tenant {
	return &tenant{
		Domains:           store.NewStore[Domain](),
		Certificates:      store.NewStore[Child[Certificate]](),
		IdentityProviders: store.NewStore[Child[IdentityProvider]](),
		Reporters:         store.NewStore[Child[Reporter]](),
	}
}

func (m *MockAM) getTenant(aware store.OrgEnvAware) *tenant {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	tenantKey := aware.GetOrgId() + "-" + aware.GetEnvId()
	et, ok := m.tenants[tenantKey]
	if !ok {
		nt := newTenant()
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

// NewMockAM returns an empty mock.
func NewMockAM() *MockAM {
	return &MockAM{
		tenants: make(map[string]*tenant),
		mutex:   sync.Mutex{},
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
	tenant := m.getTenant(req)
	m.mutex.Lock()
	defer m.mutex.Unlock()
	tenant.Domains.DeleteByKey(req.DomainKey)
	deleteChildren(tenant.Certificates, req.DomainKey)
	deleteChildren(tenant.IdentityProviders, req.DomainKey)
	deleteChildren(tenant.Reporters, req.DomainKey)
	return DeleteDomain204Response{}, nil
}

func deleteChildren[T store.Identifiable](s *store.Store[Child[T]], parentIdentity string) {
	all := s.GetAll()
	for _, child := range all {
		if child.ParentIdentity == parentIdentity {
			s.DeleteByKey(child.Self.Identity())
		}
	}
}

func (m *MockAM) GetDomain(_ context.Context, req GetDomainRequestObject) (GetDomainResponseObject, error) {
	domain, ok := m.getTenant(req).Domains.Get(req.DomainKey)
	if ok {
		return GetDomain200JSONResponse(domain), nil
	}
	return GetDomain404JSONResponse(notFoundError("Domain", req.DomainKey)), nil
}

func (m *MockAM) ListCertificates(_ context.Context, req ListCertificatesRequestObject) (ListCertificatesResponseObject, error) {
	return ListCertificates200JSONResponse(asSelves(ofParent(m.getTenant(req).Certificates.GetAll(), req.DomainKey))), nil
}

func (m *MockAM) UpsertCertificate(_ context.Context, req UpsertCertificateRequestObject) (UpsertCertificateResponseObject, error) {
	if req.Body != nil {
		body := *req.Body
		m.getTenant(req).Certificates.Put(newChild(body, req.DomainKey))
		return UpsertCertificate200JSONResponse(body), nil
	}
	return UpsertCertificate400JSONResponse(emptyBodyError()), nil
}

func (m *MockAM) DeleteCertificate(_ context.Context, req DeleteCertificateRequestObject) (DeleteCertificateResponseObject, error) {
	deleteChild(m.getTenant(req).Certificates, req.DomainKey, req.CertKey)
	return DeleteCertificate204Response{}, nil
}

func (m *MockAM) GetCertificate(_ context.Context, req GetCertificateRequestObject) (GetCertificateResponseObject, error) {
	cert, ok := isChildOf(m.getTenant(req).Certificates, req.DomainKey, req.CertKey)
	if ok {
		return GetCertificate200JSONResponse(cert.Self), nil
	}
	return GetCertificate404JSONResponse(notFoundError("Certificate", req.CertKey)), nil
}

func (m *MockAM) ListIdentityProviders(_ context.Context, req ListIdentityProvidersRequestObject) (ListIdentityProvidersResponseObject, error) {
	return ListIdentityProviders200JSONResponse(asSelves(ofParent(m.getTenant(req).IdentityProviders.GetAll(), req.DomainKey))), nil
}

func (m *MockAM) UpsertIdentityProvider(_ context.Context, req UpsertIdentityProviderRequestObject) (UpsertIdentityProviderResponseObject, error) {
	if req.Body != nil {
		body := *req.Body
		m.getTenant(req).IdentityProviders.Put(newChild(body, req.DomainKey))
		return UpsertIdentityProvider200JSONResponse(body), nil
	}
	return UpsertIdentityProvider400JSONResponse(emptyBodyError()), nil
}

func (m *MockAM) DeleteIdentityProvider(_ context.Context, req DeleteIdentityProviderRequestObject) (DeleteIdentityProviderResponseObject, error) {
	deleteChild(m.getTenant(req).IdentityProviders, req.DomainKey, req.IdentityKey)
	return DeleteIdentityProvider204Response{}, nil
}

func (m *MockAM) GetIdentityProvider(_ context.Context, req GetIdentityProviderRequestObject) (GetIdentityProviderResponseObject, error) {
	idp, ok := isChildOf(m.getTenant(req).IdentityProviders, req.DomainKey, req.IdentityKey)
	if ok {
		return GetIdentityProvider200JSONResponse(idp.Self), nil
	}
	return GetIdentityProvider404JSONResponse(notFoundError("IdentityProvider", req.IdentityKey)), nil
}

func (m *MockAM) ListReporters(_ context.Context, req ListReportersRequestObject) (ListReportersResponseObject, error) {
	return ListReporters200JSONResponse(asSelves(ofParent(m.getTenant(req).Reporters.GetAll(), req.DomainKey))), nil
}

func (m *MockAM) UpsertReporter(_ context.Context, req UpsertReporterRequestObject) (UpsertReporterResponseObject, error) {
	if req.Body != nil {
		body := *req.Body
		m.getTenant(req).Reporters.Put(newChild(body, req.DomainKey))
		return UpsertReporter200JSONResponse(body), nil
	}
	return UpsertReporter400JSONResponse(emptyBodyError()), nil
}

func (m *MockAM) DeleteReporter(_ context.Context, req DeleteReporterRequestObject) (DeleteReporterResponseObject, error) {
	deleteChild(m.getTenant(req).Reporters, req.DomainKey, req.ReporterKey)
	return DeleteReporter204Response{}, nil
}

func (m *MockAM) GetReporter(_ context.Context, req GetReporterRequestObject) (GetReporterResponseObject, error) {
	reporter, ok := isChildOf(m.getTenant(req).Reporters, req.DomainKey, req.ReporterKey)
	if ok {
		return GetReporter200JSONResponse(reporter.Self), nil
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

func asSelves[T store.Identifiable](children []Child[T]) []T {
	mapped := make([]T, len(children))
	for i, child := range children {
		mapped[i] = child.Self
	}
	return mapped
}

func ofParent[T store.Identifiable](children []Child[T], parentIdentity string) []Child[T] {
	matched := make([]Child[T], 0)
	for _, child := range children {
		if child.ParentIdentity == parentIdentity {
			matched = append(matched, child)
		}
	}
	return matched
}

func isChildOf[T store.Identifiable](s *store.Store[Child[T]], parentIdentity, id string) (Child[T], bool) {
	child, ok := s.Get(id)
	if !ok || child.ParentIdentity != parentIdentity {
		return Child[T]{}, false
	}
	return child, true
}

func deleteChild[T store.Identifiable](s *store.Store[Child[T]], parentIdentity, id string) {
	if _, ok := isChildOf(s, parentIdentity, id); ok {
		s.DeleteByKey(id)
	}
}

func newChild[T store.Identifiable](identifiable T, parentIdentity string) Child[T] {
	return Child[T]{
		ParentIdentity: parentIdentity,
		Self:           identifiable,
	}
}
