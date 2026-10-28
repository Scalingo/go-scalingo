package scalingo

import (
	"context"
	"net/url"

	"github.com/Scalingo/go-utils/errors/v3"
)

type OrganizationsService interface {
	OrganizationsList(ctx context.Context) ([]Organization, error)
}

var _ OrganizationsService = (*Client)(nil)

type Organization struct {
	ID   string `json:"id"`
	Slug string `json:"slug"`
	Role string `json:"role"`
}

type OrganizationsRes struct {
	Organizations []Organization `json:"organizations"`
}

func (c *Client) OrganizationsList(ctx context.Context) ([]Organization, error) {
	var res OrganizationsRes
	// The API only returns the role of the current user when asked to.
	params := url.Values{"with_role": {"true"}}
	err := c.AuthAPI().ResourceList(ctx, "organizations", params, &res)
	if err != nil {
		return nil, errors.Wrap(ctx, err, "list organizations")
	}
	return res.Organizations, nil
}
