package scalingo

import (
	"context"
	stderrors "errors"

	httpclient "github.com/Scalingo/go-scalingo/v8/http"
	"github.com/Scalingo/go-utils/errors/v2"
)

const firewallRulesResource = "firewall_rules"

var ErrFirewallRuleNotFound = stderrors.New("firewall rule not found")

type FirewallRuleType string

const (
	FirewallRuleTypeManagedRange FirewallRuleType = "managed_range"
	FirewallRuleTypeCustomRange  FirewallRuleType = "custom_range"
)

type FirewallRulesService interface {
	FirewallRulesCreate(ctx context.Context, database DatabaseNG, params FirewallRuleCreateParams) (FirewallRule, error)
	FirewallRulesList(ctx context.Context, database DatabaseNG) ([]FirewallRule, error)
	FirewallRulesDestroy(ctx context.Context, database DatabaseNG, firewallRuleID string) error
	FirewallRulesGetManagedRanges(ctx context.Context, database DatabaseNG) ([]FirewallManagedRange, error)
}

type FirewallManagedRange struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type FirewallRuleCreateParams struct {
	Type    FirewallRuleType `json:"type"`
	CIDR    string           `json:"cidr,omitempty"`
	Label   string           `json:"label,omitempty"`
	RangeID string           `json:"range_id,omitempty"`
}

type FirewallRule struct {
	ID         string           `json:"id"`
	Type       FirewallRuleType `json:"type"`
	CIDR       string           `json:"cidr"`
	Label      string           `json:"label"`
	RangeID    string           `json:"range_id"`
	DatabaseID string           `json:"database_id"`
}

type FirewallManagedRangesResponse struct {
	ManagedRanges []FirewallManagedRange `json:"ranges"`
}

type FirewallRulesResponse struct {
	FirewallRules []FirewallRule `json:"rules"`
}

var _ FirewallRulesService = (*PreviewClient)(nil)

func (c *PreviewClient) FirewallRulesCreate(ctx context.Context, database DatabaseNG, params FirewallRuleCreateParams) (FirewallRule, error) {
	var res FirewallRule

	addonID, err := c.getAddonID(ctx, database)
	if err != nil {
		return res, errors.Wrap(ctx, err, "get addon ID")
	}

	err = c.parent.DBAPI(database.App.ID, addonID).SubresourceAdd(ctx, "databases", addonID, firewallRulesResource, params, &res)
	if err != nil {
		return res, errors.Wrap(ctx, err, "create firewall rule")
	}
	return res, nil
}

func (c *PreviewClient) FirewallRulesList(ctx context.Context, database DatabaseNG) ([]FirewallRule, error) {
	var res FirewallRulesResponse

	addonID, err := c.getAddonID(ctx, database)
	if err != nil {
		return nil, errors.Wrap(ctx, err, "get addon ID")
	}

	err = c.parent.DBAPI(database.App.ID, addonID).SubresourceList(ctx, "databases", addonID, firewallRulesResource, nil, &res)
	if err != nil {
		return nil, errors.Wrap(ctx, err, "list firewall rules")
	}

	return res.FirewallRules, nil
}

func (c *PreviewClient) FirewallRulesDestroy(ctx context.Context, database DatabaseNG, firewallRuleID string) error {
	addonID, err := c.getAddonID(ctx, database)
	if err != nil {
		return errors.Wrap(ctx, err, "get addon ID")
	}

	err = c.parent.DBAPI(database.App.ID, addonID).SubresourceDelete(ctx, "databases", addonID, firewallRulesResource, firewallRuleID)
	if err != nil {
		return errors.Wrap(ctx, err, "destroy firewall rule")
	}
	return nil
}

func (c *PreviewClient) FirewallRulesGetManagedRanges(ctx context.Context, database DatabaseNG) ([]FirewallManagedRange, error) {
	var res FirewallManagedRangesResponse

	addonID, err := c.getAddonID(ctx, database)
	if err != nil {
		return nil, errors.Wrap(ctx, err, "get addon ID")
	}

	req := &httpclient.APIRequest{
		Method:   "GET",
		Endpoint: "/firewall/managed_ranges",
	}

	err = c.parent.DBAPI(database.App.ID, addonID).DoRequest(ctx, req, &res)
	if err != nil {
		return nil, errors.Wrap(ctx, err, "get managed ranges")
	}

	return res.ManagedRanges, nil
}

func (c *PreviewClient) getAddonID(ctx context.Context, database DatabaseNG) (string, error) {
	var res []*Addon

	res, err := c.parent.AddonsList(ctx, database.ID)
	if err != nil {
		return "", errors.Wrap(ctx, err, "list addons")
	}

	// There is only one addon per app for databases next gen
	return res[0].ID, nil
}
