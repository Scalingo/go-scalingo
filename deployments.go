package scalingo

import (
	"net/http"

	"gopkg.in/errgo.v1"
)

type Deployment struct {
	ID        string          `json:"id"`
	AppID     string          `json:"app_id"`
	CreatedAt string          `json:"created_at"`
	Status    string          `json:"status"`
	GitRef    string          `json:"git_ref"`
	Image     string          `json:"image"`
	Registry  string          `json:"registry"`
	Duration  int             `json:"duration"`
	User      User            `json:"pusher"`
	Links     DeploymentLinks `json:"links"`
}

type DeploymentLinks struct {
	Output string `json:"output"`
}

func (c *Client) DeploymentList(app string) ([]*Deployment, error) {
	req := &APIRequest{
		Client:   c,
		Endpoint: "/apps/" + app + "/deployments",
	}

	res, err := req.Do()

	if err != nil {
		return []*Deployment{}, errgo.Mask(err, errgo.Any)
	}

	defer res.Body.Close()

	deployMap := map[string][]*Deployment{}
	err = ParseJSON(res, &deployMap)

	if err != nil {
		return []*Deployment{}, errgo.Mask(err, errgo.Any)
	}

	return deployMap["deployments"], nil
}

func (c *Client) DeploymentLogs(app string, deploy string) (*http.Response, error) {
	req := &APIRequest{
		Client:   c,
		Endpoint: "/apps/" + app + "/deployments/" + deploy,
	}

	return req.Do()
}
