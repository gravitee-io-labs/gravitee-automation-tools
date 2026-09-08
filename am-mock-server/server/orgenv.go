package server

func (r ListDomainsRequestObject) GetOrgId() string  { return r.OrgId }
func (r ListDomainsRequestObject) GetEnvId() string   { return r.EnvId }
func (r UpsertDomainRequestObject) GetOrgId() string  { return r.OrgId }
func (r UpsertDomainRequestObject) GetEnvId() string   { return r.EnvId }
func (r DeleteDomainRequestObject) GetOrgId() string  { return r.OrgId }
func (r DeleteDomainRequestObject) GetEnvId() string   { return r.EnvId }
func (r GetDomainRequestObject) GetOrgId() string     { return r.OrgId }
func (r GetDomainRequestObject) GetEnvId() string      { return r.EnvId }

func (r ListCertificatesRequestObject) GetOrgId() string  { return r.OrgId }
func (r ListCertificatesRequestObject) GetEnvId() string   { return r.EnvId }
func (r UpsertCertificateRequestObject) GetOrgId() string  { return r.OrgId }
func (r UpsertCertificateRequestObject) GetEnvId() string   { return r.EnvId }
func (r DeleteCertificateRequestObject) GetOrgId() string  { return r.OrgId }
func (r DeleteCertificateRequestObject) GetEnvId() string   { return r.EnvId }
func (r GetCertificateRequestObject) GetOrgId() string     { return r.OrgId }
func (r GetCertificateRequestObject) GetEnvId() string      { return r.EnvId }

func (r ListIdentityProvidersRequestObject) GetOrgId() string  { return r.OrgId }
func (r ListIdentityProvidersRequestObject) GetEnvId() string   { return r.EnvId }
func (r UpsertIdentityProviderRequestObject) GetOrgId() string  { return r.OrgId }
func (r UpsertIdentityProviderRequestObject) GetEnvId() string   { return r.EnvId }
func (r DeleteIdentityProviderRequestObject) GetOrgId() string  { return r.OrgId }
func (r DeleteIdentityProviderRequestObject) GetEnvId() string   { return r.EnvId }
func (r GetIdentityProviderRequestObject) GetOrgId() string     { return r.OrgId }
func (r GetIdentityProviderRequestObject) GetEnvId() string      { return r.EnvId }

func (r ListReportersRequestObject) GetOrgId() string  { return r.OrgId }
func (r ListReportersRequestObject) GetEnvId() string   { return r.EnvId }
func (r UpsertReporterRequestObject) GetOrgId() string  { return r.OrgId }
func (r UpsertReporterRequestObject) GetEnvId() string   { return r.EnvId }
func (r DeleteReporterRequestObject) GetOrgId() string  { return r.OrgId }
func (r DeleteReporterRequestObject) GetEnvId() string   { return r.EnvId }
func (r GetReporterRequestObject) GetOrgId() string     { return r.OrgId }
func (r GetReporterRequestObject) GetEnvId() string      { return r.EnvId }
