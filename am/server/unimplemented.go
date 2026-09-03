package server

import (
	"context"
	"errors"
)

var ErrNotImplemented = errors.New("not implemented")

type StrictUnimplemented struct{}

var _ StrictServerInterface = StrictUnimplemented{}

func (StrictUnimplemented) AutomationListDomains(context.Context, AutomationListDomainsRequestObject) (AutomationListDomainsResponseObject, error) {
	return nil, ErrNotImplemented
}

func (StrictUnimplemented) AutomationCreateOrUpdateDomain(context.Context, AutomationCreateOrUpdateDomainRequestObject) (AutomationCreateOrUpdateDomainResponseObject, error) {
	return nil, ErrNotImplemented
}

func (StrictUnimplemented) AutomationDeleteDomain(context.Context, AutomationDeleteDomainRequestObject) (AutomationDeleteDomainResponseObject, error) {
	return nil, ErrNotImplemented
}

func (StrictUnimplemented) AutomationGetDomain(context.Context, AutomationGetDomainRequestObject) (AutomationGetDomainResponseObject, error) {
	return nil, ErrNotImplemented
}

func (StrictUnimplemented) AutomationListCertificates(context.Context, AutomationListCertificatesRequestObject) (AutomationListCertificatesResponseObject, error) {
	return nil, ErrNotImplemented
}

func (StrictUnimplemented) AutomationCreateOrUpdateCertificate(context.Context, AutomationCreateOrUpdateCertificateRequestObject) (AutomationCreateOrUpdateCertificateResponseObject, error) {
	return nil, ErrNotImplemented
}

func (StrictUnimplemented) AutomationDeleteCertificate(context.Context, AutomationDeleteCertificateRequestObject) (AutomationDeleteCertificateResponseObject, error) {
	return nil, ErrNotImplemented
}

func (StrictUnimplemented) AutomationGetCertificate(context.Context, AutomationGetCertificateRequestObject) (AutomationGetCertificateResponseObject, error) {
	return nil, ErrNotImplemented
}

func (StrictUnimplemented) AutomationListIdentityProviders(context.Context, AutomationListIdentityProvidersRequestObject) (AutomationListIdentityProvidersResponseObject, error) {
	return nil, ErrNotImplemented
}

func (StrictUnimplemented) AutomationCreateOrUpdateIdentityProvider(context.Context, AutomationCreateOrUpdateIdentityProviderRequestObject) (AutomationCreateOrUpdateIdentityProviderResponseObject, error) {
	return nil, ErrNotImplemented
}

func (StrictUnimplemented) AutomationDeleteIdentityProvider(context.Context, AutomationDeleteIdentityProviderRequestObject) (AutomationDeleteIdentityProviderResponseObject, error) {
	return nil, ErrNotImplemented
}

func (StrictUnimplemented) AutomationGetIdentityProvider(context.Context, AutomationGetIdentityProviderRequestObject) (AutomationGetIdentityProviderResponseObject, error) {
	return nil, ErrNotImplemented
}

func (StrictUnimplemented) AutomationListReporters(context.Context, AutomationListReportersRequestObject) (AutomationListReportersResponseObject, error) {
	return nil, ErrNotImplemented
}

func (StrictUnimplemented) AutomationCreateOrUpdateReporter(context.Context, AutomationCreateOrUpdateReporterRequestObject) (AutomationCreateOrUpdateReporterResponseObject, error) {
	return nil, ErrNotImplemented
}

func (StrictUnimplemented) AutomationDeleteReporter(context.Context, AutomationDeleteReporterRequestObject) (AutomationDeleteReporterResponseObject, error) {
	return nil, ErrNotImplemented
}

func (StrictUnimplemented) AutomationGetReporter(context.Context, AutomationGetReporterRequestObject) (AutomationGetReporterResponseObject, error) {
	return nil, ErrNotImplemented
}
