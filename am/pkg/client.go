package pkg

import (
	"context"
	"net/http"
	"time"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/gravitee-io-labs/gravitee-automation-sdks/am/pkg/sdk/domain"
	"github.com/gravitee-io-labs/gravitee-automation-sdks/common/pkg/auth"
	"github.com/gravitee-io-labs/gravitee-automation-sdks/common/pkg/errors"
)

type AMClient struct {
	orgID string
	envID string
	domain.ClientWithResponsesInterface
}

func NewClient(apiContext auth.APIContext, timeoutMs int) (*AMClient, error) {

	requestEditor, err := apiContext.RequestEditor()
	if err != nil {
		return nil, err
	}

	client, err := domain.NewClientWithResponses(
		apiContext.BaseURL,
		domain.WithRequestEditorFn(requestEditor),
		domain.WithHTTPClient(&http.Client{
			Timeout: time.Duration(timeoutMs) * time.Millisecond,
		}))

	if err != nil {
		return nil, err
	}

	return &AMClient{orgID: apiContext.OrgID, envID: apiContext.EnvID, ClientWithResponsesInterface: client}, nil
}

func (c *AMClient) GetDomain(ctx context.Context, domainKey string) (*domain.AutomationDomain, error) {
	if resp, err := c.AutomationGetDomainWithResponse(ctx, c.orgID, c.envID, domainKey); err != nil {
		return nil, errors.NewClientError(err)
	} else {
		return respond(resp.JSON200, resp, resp.Body, resp.JSON403, resp.JSON404)
	}
}

func (c *AMClient) UpsertDomain(ctx context.Context, domainBody domain.AutomationDomain) (*domain.AutomationDomain, error) {
	if resp, err := c.AutomationCreateOrUpdateDomainWithResponse(ctx, c.orgID, c.envID, domainBody); err != nil {
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

func respond[T any](entity *T, statusCoder openapi3filter.StatusCoder, body []byte, errs ...*domain.Error) (*T, error) {

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
