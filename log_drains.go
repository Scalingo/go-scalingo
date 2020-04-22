package scalingo

import (
	"gopkg.in/errgo.v1"
)

type LogDrainsService interface {
	LogDrainsList(app string) ([]LogDrain, error)
	LogDrainAdd(app string, params LogDrainAddParams) (*LogDrain, error)
	LogDrainRemove(app string, id string) error
}

var _ LogDrainsService = (*Client)(nil)

type LogDrain struct {
	ID        string `json:"id"`
	Type      string `json:"drain_type"`
	Status    string `json:"status"`
	TargetURL string `json:"targetURL"`
}

type LogDrainsRes struct {
	LogDrains []LogDrain `json:"log_drains"`
}

type LogDrainRes struct {
	LogDrain LogDrain `json:"log_drain"`
}

func (c *Client) LogDrainsList(app string) ([]LogDrain, error) {
	var logDrainsRes LogDrainsRes
	err := c.ScalingoAPI().SubresourceList("apps", app, "log_drains", nil, &logDrainsRes)
	if err != nil {
		return nil, errgo.Mask(err)
	}
	return logDrainsRes.LogDrains, nil
}

type LogDrainAddParams struct {
	Type      string `json:"drain_type"`
	Status    string `json:"status"`
	TargetURL string `json:"targetURL"`
}

func (c *Client) LogDrainAdd(app string, params LogDrainAddParams) (*LogDrain, error) {
	var logDrainRes LogDrainRes
	err := c.ScalingoAPI().SubresourceAdd("apps", app, "log_drains", LogDrainRes{
		LogDrain: LogDrain{
			Type:      params.Type,
			TargetURL: params.TargetURL,
			Status:    params.Status,
		},
	}, &logDrainRes)
	if err != nil {
		return nil, errgo.Mask(err)
	}
	return &logDrainRes.LogDrain, nil
}

func (c *Client) LogDrainInfo(app, id string) (*LogDrain, error) {
	var logDrainRes LogDrainRes
	err := c.ScalingoAPI().SubresourceGet("apps", app, "log_drains", id, nil, &logDrainRes)
	if err != nil {
		return nil, errgo.Mask(err)
	}
	return &logDrainRes.LogDrain, nil
}

func (c *Client) LogDrainRemove(app, id string) error {
	err := c.ScalingoAPI().SubresourceDelete("apps", app, "log_drains", id)
	if err != nil {
		return errgo.Mask(err)
	}
	return nil
}
