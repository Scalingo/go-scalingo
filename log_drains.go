package scalingo

import (
	"gopkg.in/errgo.v1"
)

type LogDrainsService interface {
	LogDrainsList(app string) ([]LogDrain, error)
}

var _ LogDrainsService = (*Client)(nil)

type LogDrain struct {
	AppID string `json:"app_id"`
	URL   string `json:"url"`
}

type LogDrainsRes struct {
	LogDrains []LogDrain
}

type LogDrainRes struct {
	LogDrain LogDrain `json:"log_drain"`
}

func (c *Client) LogDrainsList(app string) ([]LogDrain, error) {
	var logDrainsRes []LogDrain
	err := c.ScalingoAPI().SubresourceList("apps", app, "log_drains", nil, &logDrainsRes)
	if err != nil {
		return nil, errgo.Notef(err, "fail to list the log drains")
	}
	return logDrainsRes, nil
}

type LogDrainAddParams struct {
	AppID string `json:"app_id"`
	URL   string `json:"url"`
}

func (c *Client) LogDrainAdd(app string, params LogDrainAddParams) (*LogDrain, error) {
	var logDrainRes LogDrainRes
	err := c.ScalingoAPI().SubresourceAdd("apps", app, "log_drains", LogDrainRes{
		LogDrain: LogDrain{
			AppID: params.AppID,
			URL:   params.URL,
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
