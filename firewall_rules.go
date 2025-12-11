package scalingo

import (
	"context"
	stderrors "errors"

	"github.com/Scalingo/go-utils/errors/v2"
)

const firewallRulesResource = "firewall_rules"

var ErrFirewallRuleNotFound = stderrors.New("firewall rule not found")

type FirewallRulesService interface {
	FirewallRulesCreate(ctx context.Context, database Database, params FirewallRuleCreateParams) (FirewallRule, error)
	FirewallRulesList(ctx context.Context, database Database) ([]FirewallRule, error)
	FirewallRulesDestroy(ctx context.Context, database Database, firewallRuleID string) error
	FirewallRulesGetManagedRanges(ctx context.Context) ([]string, error)
}

type FirewallRuleCreateParams struct {
	Type    string `json:"type"`
	CIDR    string `json:"cidr"`
	Label   string `json:"label,omitempty"`
	RangeID string `json:"range_id,omitempty"`
}

type FirewallRule struct {
	ID         string `json:"id"`
	Type       string `json:"type"`
	CIDR       string `json:"cidr"`
	Label      string `json:"label"`
	RangeID    string `json:"range_id"`
	DatabaseID string `json:"database_id"`
}

var _ FirewallRulesService = (*PreviewClient)(nil)

func (c *PreviewClient) FirewallRulesCreate(ctx context.Context, database Database, params FirewallRuleCreateParams) (FirewallRule, error) {
	var res FirewallRule

	err := c.parent.DBAPI(database.app_id, database.ID).SubresourceAdd(ctx, databasesResource, database.ID, firewallRulesResource, params, &res)
	if err != nil {
		return res, errors.Wrap(ctx, err, "create firewall rule")
	}
	return res, nil
}

func (c *PreviewClient) FirewallRulesList(ctx context.Context, database Database) ([]FirewallRule, error) {
	var res []FirewallRule

	err := c.parent.DBAPI(app, addonID).SubresourceList(ctx, databasesResource, databaseID, firewallRulesResource, nil, &res)
	if err != nil {
		return nil, errors.Wrap(ctx, err, "list firewall rules")
	}

	return res, nil
}

func (c *PreviewClient) FirewallRulesDestroy(ctx context.Context, database Database, firewallRuleID string) error {
	err := c.parent.DBAPI(app, addonID).SubresourceDelete(ctx, databasesResource, databaseID, firewallRulesResource, firewallRuleID)
	if err != nil {
		return errors.Wrap(ctx, err, "destroy firewall rule")
	}
	return nil
}

func (c *PreviewClient) FirewallRulesGetManagedRanges(ctx context.Context) ([]string, error) {
	var res []string

	err := c.parent.DBAPI(app, addonID).SubresourceList(ctx, databasesResource, "", "managed_ranges", nil, &res)
	if err != nil {
		return nil, errors.Wrap(ctx, err, "get managed ranges")
	}

	return res, nil
}
