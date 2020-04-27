package scalingo

import (
	httpclient "github.com/Scalingo/go-scalingo/http"
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

type errorRes struct {
	Error string `json:"error"`
}

type LogDrainRes struct {
	LogDrain LogDrain `json:"drain"`
	Errors   errorRes `json:"errors"`
}

// type LogDrainReq struct {
// 	LogDrain LogDrain `json:"drain"`
// }

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

func (c *Client) LogDrainAdd(app string, params LogDrainAddParams) (*LogDrain, error) {
	var logDrainRes LogDrainRes
	// var logDrainReq LogDrainReq
	payload := LogDrainRes{
		LogDrain: LogDrain{
			URL: params.URL,
		},
	}
	req := &httpclient.APIRequest{
		Method:   "POST",
		Endpoint: "/" + "apps" + "/" + app + "/" + "log_drains",
		Expected: httpclient.Statuses{201, 422},
		Params:   payload,
	}

	err := c.ScalingoAPI().DoRequest(req, &logDrainRes)
	if err != nil {
		return nil, errgo.Mask(err)
	}

	if logDrainRes.Errors.Error != "" {
		return nil, errgo.New(logDrainRes.Errors.Error)
	}
	// fmt.Println("result:", logDrainRes)
	return &logDrainRes.LogDrain, nil
}

func (c *Client) LogDrainInfo(app, id string) (*LogDrain, error) {
	var logDrainRes LogDrain
	err := c.ScalingoAPI().SubresourceGet("apps", app, "log_drains", id, nil, &logDrainRes)
	if err != nil {
		return nil, errgo.Mask(err)
	}
	return &logDrainRes, nil
}

func (c *Client) LogDrainRemove(app, id string) error {
	err := c.ScalingoAPI().SubresourceDelete("apps", app, "log_drains", id)
	if err != nil {
		return errgo.Mask(err)
	}
	return nil
}
