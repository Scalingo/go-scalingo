package scalingo

import (
	httpclient "github.com/Scalingo/go-scalingo/http"
	"gopkg.in/errgo.v1"
)

type LogDrainsService interface {
	LogDrainsList(app string) ([]LogDrain, error)
	LogDrainAdd(app string, params LogDrainAddParams) (*Client, error)
}

var _ LogDrainsService = (*Client)(nil)

type LogDrain struct {
	AppID string `json:"app_id"`
	URL   string `json:"url"`
}

type logDrainReq struct {
	Drain LogDrain `json:"drain"`
}

type logDrainRes struct {
	Error string `json:"error"`
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
	URL string `json:"url"`
}

func (c *Client) LogDrainAdd(app string, params LogDrainAddParams) (*Client, error) {
	var logDrainRes logDrainRes
	payload := logDrainReq{
		Drain: LogDrain{
			URL: params.URL,
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
		return c, errgo.Notef(err, "fail to add drain")
	}

	if logDrainRes.Error != "" {
		return c, errgo.Notef(err, logDrainRes.Error)
	}

	return c, nil
}
