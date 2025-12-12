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
	FirewallRulesCreate(ctx context.Context, AppID string, AddonID string, params FirewallRuleCreateParams) (FirewallRule, error)
	FirewallRulesList(ctx context.Context, AppID string, AddonID string) ([]FirewallRule, error)
	FirewallRulesDestroy(ctx context.Context, AppID string, AddonID string, firewallRuleID string) error
	FirewallRulesGetManagedRanges(ctx context.Context, AppID string, AddonID string) ([]FirewallManagedRange, error)
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

func (c *PreviewClient) FirewallRulesCreate(ctx context.Context, AppID string, AddonID string, params FirewallRuleCreateParams) (FirewallRule, error) {
	var res FirewallRule

	err := c.parent.DBAPI(AppID, AddonID).SubresourceAdd(ctx, "databases", AddonID, firewallRulesResource, params, &res)
	if err != nil {
		return res, errors.Wrap(ctx, err, "create firewall rule")
	}
	return res, nil
}

func (c *PreviewClient) FirewallRulesList(ctx context.Context, AppID string, AddonID string) ([]FirewallRule, error) {
	var res FirewallRulesResponse

	err := c.parent.DBAPI(AppID, AddonID).SubresourceList(ctx, "databases", AddonID, firewallRulesResource, nil, &res)
	if err != nil {
		return nil, errors.Wrap(ctx, err, "list firewall rules")
	}

	return res.FirewallRules, nil
}

func (c *PreviewClient) FirewallRulesDestroy(ctx context.Context, AppID string, AddonID string, firewallRuleID string) error {
	err := c.parent.DBAPI(AppID, AddonID).SubresourceDelete(ctx, "databases", AddonID, firewallRulesResource, firewallRuleID)
	if err != nil {
		return errors.Wrap(ctx, err, "destroy firewall rule")
	}
	return nil
}

func (c *PreviewClient) FirewallRulesGetManagedRanges(ctx context.Context, AppID string, AddonID string) ([]FirewallManagedRange, error) {
	var res FirewallManagedRangesResponse

	req := &httpclient.APIRequest{
		Method:   "GET",
		Endpoint: "/firewall/managed_ranges",
	}

	err := c.parent.DBAPI(AppID, AddonID).DoRequest(ctx, req, &res)
	if err != nil {
		return nil, errors.Wrap(ctx, err, "get managed ranges")
	}

	return res.ManagedRanges, nil
}
