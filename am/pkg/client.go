package pkg

import (
	"context"
	"net/http"
	"time"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/gravitee-io-labs/gravitee-automation-sdks/am/pkg/sdk"
	"github.com/gravitee-io-labs/gravitee-automation-sdks/common/pkg/auth"
	"github.com/gravitee-io-labs/gravitee-automation-sdks/common/pkg/errors"
)

type AMClient struct {
	orgID string
	envID string
	sdk.ClientWithResponsesInterface
}

func NewClient(apiContext auth.APIContext, timeoutMs int) (*AMClient, error) {

	authOption, err := apiContext.AuthClientOption()
	if err != nil {
		return nil, err
	}

	client, err := sdk.NewClientWithResponses(
		apiContext.BaseURL,
		authOption,
		sdk.WithHTTPClient(&http.Client{
			Timeout: time.Duration(timeoutMs) * time.Millisecond,
		}))

	if err != nil {
		return nil, err
	}

	return &AMClient{orgID: apiContext.OrgID, envID: apiContext.EnvID, ClientWithResponsesInterface: client}, nil
}

func (c *AMClient) GetDomain(ctx context.Context, domainKey string) (*sdk.AutomationDomain, error) {
	if resp, err := c.AutomationGetDomainWithResponse(ctx, c.orgID, c.envID, domainKey); err != nil {
		return nil, errors.NewClientError(err)
	} else {
		return respond(resp.JSON200, resp, resp.Body, resp.JSON403, resp.JSON404)
	}
}

func (c *AMClient) UpsertDomain(ctx context.Context, domain sdk.AutomationDomain) (*sdk.AutomationDomain, error) {
	if resp, err := c.AutomationCreateOrUpdateDomainWithResponse(ctx, c.orgID, c.envID, domain); err != nil {
		return nil, errors.NewClientError(err)
	} else {
		return respond(resp.JSON200, resp, resp.Body, resp.JSON400, resp.JSON403)
	}
}

func (c *AMClient) DeleteDomain(ctx context.Context, domainKey string) error {
	if resp, err := c.AutomationDeleteDomainWithResponse(ctx, c.orgID, c.orgID, domainKey); err != nil {
		return errors.NewClientError(err)
	} else {
		_, err := respond(empty(), resp, resp.Body, resp.JSON403, resp.JSONDefault)
		return err
	}
}

func respond[T any](entity *T, statusCoder openapi3filter.StatusCoder, body []byte, errs ...*sdk.Error) (*T, error) {

	// r is never nil according to generated code
	switch {
	case statusCoder.StatusCode() < 300:
		return entity, nil
	case len(errs) > 0:
		return nil, errors.HttpError{Status: statusCoder.StatusCode(), Body: deref(errs[0].Message)}
	default:
		return nil, errors.HttpError{Status: statusCoder.StatusCode(), Body: string(body)}
	}
}

func deref(message *string) string {
	if message == nil {
		return ""
	}
	return *message
}

func empty() *struct{} {
	return nil
}
