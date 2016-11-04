package scalingo

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/net/websocket"
	"gopkg.in/errgo.v1"
)

type DeploymentStatus string

const (
	StatusSuccess      DeploymentStatus = "success"
	StatusWaitForCI                     = "wait-for-ci"
	StatusBuilding                      = "building"
	StatusStarting                      = "starting"
	StatusPushing                       = "pushing"
	StatusAborted                       = "aborted"
	StatusBuildError                    = "build-error"
	StatusCrashedError                  = "crashed-error"
	StatusTimeoutError                  = "timeout-error"
	StatusHookError                     = "hook-error"
)

type Deployment struct {
	ID             string           `json:"id"`
	AppID          string           `json:"app_id"`
	CreatedAt      *time.Time       `json:"created_at"`
	Status         DeploymentStatus `json:"status"`
	GitRef         string           `json:"git_ref"`
	Image          string           `json:"image"`
	Registry       string           `json:"registry"`
	Duration       int              `json:"duration"`
	PostdeployHook string           `json:"postdeploy_hook"`
	User           *User            `json:"pusher"`
	Links          *DeploymentLinks `json:"links"`
}

func (d *Deployment) IsFinished() bool {
	return d.Status != StatusBuilding && d.Status != StatusStarting && d.Status != StatusPushing
}

type DeploymentList struct {
	Deployments []*Deployment `json:"deployments"`
}

type DeploymentLinks struct {
	Output string `json:"output"`
}

type AuthenticationData struct {
	Token string `json:"token"`
}

type AuthStruct struct {
	Type string             `json:"type"`
	Data AuthenticationData `json:"data"`
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

	var deployments DeploymentList
	err = ParseJSON(res, &deployments)

	if err != nil {
		return []*Deployment{}, errgo.Mask(err, errgo.Any)
	}

	return deployments.Deployments, nil
}

func (c *Client) Deployment(app string, deploy string) (*Deployment, error) {
	req := &APIRequest{
		Client:   c,
		Endpoint: "/apps/" + app + "/deployments/" + deploy,
	}

	res, err := req.Do()

	if err != nil {
		return nil, errgo.Mask(err, errgo.Any)
	}
	defer res.Body.Close()

	var deploymentMap map[string]*Deployment

	err = ParseJSON(res, &deploymentMap)

	if err != nil {
		return nil, errgo.Mask(err, errgo.Any)
	}

	return deploymentMap["deployment"], nil
}

func (c *Client) DeploymentLogs(deployURL string) (*http.Response, error) {
	u, err := url.Parse(deployURL)
	if err != nil {
		return nil, errgo.Mask(err, errgo.Any)
	}
	req := &APIRequest{
		Client:   c,
		Expected: Statuses{200, 404},
		Endpoint: u.Path,
		URL:      u.Scheme + "://" + u.Host,
	}

	return req.Do()
}

func (c *Client) DeploymentStream(deployURL string) (*websocket.Conn, error) {
	authString, err := json.Marshal(&AuthStruct{
		Type: "auth",
		Data: AuthenticationData{
			Token: c.APIToken,
		},
	})
	if err != nil {
		return nil, errgo.Mask(err, errgo.Any)
	}

	conn, err := websocket.Dial(deployURL, "", "http://scalingo-cli.local/"+c.APIVersion)
	if err != nil {
		return nil, errgo.Mask(err, errgo.Any)
	}

	_, err = conn.Write(authString)

	if err != nil {
		return nil, errgo.Mask(err, errgo.Any)
	}
	return conn, nil
}
