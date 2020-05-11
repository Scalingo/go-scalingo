package scalingo

import (
	httpclient "github.com/Scalingo/go-scalingo/http"
	"gopkg.in/errgo.v1"
)

type LogDrainsService interface {
	LogDrainsList(app string) ([]LogDrain, error)
	LogDrainAdd(app string, params LogDrainAddParams) (*LogDrainRes, error)
}

var _ LogDrainsService = (*Client)(nil)

type LogDrain struct {
	AppID string `json:"app_id"`
	URL   string `json:"url"`
	Type  string `json:"type"`
	Token string `json:"token"`
	Param string `json:"param"`
}

type logDrainReq struct {
	Drain LogDrain `json:"drain"`
}

type LogDrainRes struct {
	Drain LogDrain `json:"drain"`
	Error string   `json:"error"`
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
	URL   string `json:"url"`
	Type  string `json:"type"`
	Token string `json:"token"`
	Param string `json:"param"`
}

func (c *Client) LogDrainAdd(app string, params LogDrainAddParams) (*LogDrainRes, error) {
	var logDrainRes LogDrainRes
	payload := logDrainReq{
		Drain: LogDrain{
			URL:   params.URL,
			Type:  params.Type,
			Token: params.Token,
			Param: params.Param,
		},
	}

	req := &httpclient.APIRequest{
		Method:   "POST",
		Endpoint: "/apps/" + app + "/log_drains",
		Expected: httpclient.Statuses{201, 422},
		Params:   payload,
	}

	err := c.ScalingoAPI().DoRequest(req, &logDrainRes)
	if err != nil {
		return nil, errgo.Notef(err, "fail to add drain")
	}

	if logDrainRes.Error != "" {
		return nil, errgo.Notef(err, logDrainRes.Error)
	}

	return &logDrainRes, nil
}
