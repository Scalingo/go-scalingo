package scalingo

import (
	"time"

	"github.com/Scalingo/go-scalingo/http"
)

type SCMRepoLinkService interface {
	SCMRepoLinkShow(app string) (*SCMRepoLink, error)
	SCMRepoLinkCreate(app string, params SCMRepoLinkParams) (*SCMRepoLink, error)
	SCMRepoLinkUpdate(app string, params SCMRepoLinkParams) (*SCMRepoLink, error)
	SCMRepoLinkDelete(app string) error

	SCMRepoLinkManualDeploy(app, branch string) error
	SCMRepoLinkManualReviewApp(app, pullRequestId string) error
	SCMRepoLinkDeployments(app string) ([]*Deployment, error)
	SCMRepoLinkReviewApps(app string) ([]*ReviewApp, error)
}

type SCMRepoLinkParams struct {
	Source                   *string `json:"source,omitempty"`
	Branch                   *string `json:"branch,omitempty"`
	AuthIntegrationUUID      *string `json:"auth_integration_uuid,omitempty"`
	ScmIntegrationUUID       *string `json:"scm_integration_uuid,omitempty"`
	AutoDeployEnabled        *bool   `json:"auto_deploy_enabled,omitempty"`
	DeployReviewAppsEnabled  *bool   `json:"deploy_review_apps_enabled,omitempty"`
	DestroyOnCloseEnabled    *bool   `json:"delete_on_close_enabled,omitempty"`
	HoursBeforeDeleteOnClose *uint   `json:"hours_before_delete_on_close,omitempty"`
	DestroyStaleEnabled      *bool   `json:"delete_stale_enabled,omitempty"`
	HoursBeforeDeleteStale   *uint   `json:"hours_before_delete_stale,omitempty"`
}

type SCMRepoLink struct {
	ID                       string            `json:"id"`
	AppID                    string            `json:"app_id"`
	Linker                   SCMRepoLinkLinker `json:"linker"`
	Owner                    string            `json:"owner"`
	Repo                     string            `json:"repo"`
	Branch                   string            `json:"branch"`
	CreatedAt                time.Time         `json:"created_at"`
	UpdatedAt                time.Time         `json:"updated_at"`
	AutoDeployEnabled        bool              `json:"auto_deploy_enabled"`
	ScmIntegrationUUID       string            `json:"scm_integration_uuid"`
	AuthIntegrationID        string            `json:"auth_integration_id"`
	DeployReviewAppsEnabled  bool              `json:"deploy_review_apps_enabled"`
	DeleteOnCloseEnabled     bool              `json:"delete_on_close_enabled"`
	DeleteStaleEnabled       bool              `json:"delete_stale_enabled"`
	HoursBeforeDeleteOnClose uint              `json:"hours_before_delete_on_close"`
	HoursBeforeDeleteStale   uint              `json:"hours_before_delete_stale"`
	LastAutoDeployAt         time.Time         `json:"last_auto_deploy_at"`
}

type SCMRepoLinkLinker struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	ID       string `json:"id"`
}

type ScmRepoLinkResponse struct {
	SCMRepoLink *SCMRepoLink `json:"scm_repo_link"`
}

type SCMRepoLinkDeploymentsResponse struct {
	Deployments []*Deployment `json:"deployments"`
}

type SCMRepoLinkReviewAppsResponse struct {
	ReviewApps []*ReviewApp `json:"review_apps"`
}

var _ SCMRepoLinkService = (*Client)(nil)

func (c *Client) SCMRepoLinkShow(app string) (*SCMRepoLink, error) {
	var res ScmRepoLinkResponse
	err := c.ScalingoAPI().DoRequest(&http.APIRequest{
		Method:   "GET",
		Endpoint: "/apps/" + app + "/scm_repo_link",
		Expected: http.Statuses{200},
	}, &res)
	if err != nil {
		return nil, err
	}
	return res.SCMRepoLink, nil
}

func (c *Client) SCMRepoLinkCreate(app string, params SCMRepoLinkParams) (*SCMRepoLink, error) {
	var res ScmRepoLinkResponse
	err := c.ScalingoAPI().DoRequest(&http.APIRequest{
		Method:   "POST",
		Endpoint: "/apps/" + app + "/scm_repo_link",
		Expected: http.Statuses{201},
		Params:   map[string]SCMRepoLinkParams{"scm_repo_link": params},
	}, &res)
	if err != nil {
		return nil, err
	}

	return res.SCMRepoLink, nil
}

func (c *Client) SCMRepoLinkUpdate(app string, params SCMRepoLinkParams) (*SCMRepoLink, error) {
	var res ScmRepoLinkResponse
	err := c.ScalingoAPI().DoRequest(&http.APIRequest{
		Method:   "UPDATE",
		Endpoint: "/apps/" + app + "/scm_repo_link",
		Expected: http.Statuses{200},
		Params:   map[string]SCMRepoLinkParams{"scm_repo_link": params},
	}, &res)
	if err != nil {
		return nil, err
	}

	return res.SCMRepoLink, nil
}

func (c *Client) SCMRepoLinkDelete(app string) error {
	_, err := c.ScalingoAPI().Do(&http.APIRequest{
		Method:   "DELETE",
		Endpoint: "/apps/" + app + "/scm_repo_link",
		Expected: http.Statuses{204},
	})
	return err
}

func (c *Client) SCMRepoLinkManualDeploy(app, branch string) error {
	_, err := c.ScalingoAPI().Do(&http.APIRequest{
		Method:   "POST",
		Endpoint: "/apps/" + app + "/scm_repo_link/manual_deploy",
		Expected: http.Statuses{200},
		Params:   map[string]string{"branch": branch},
	})
	return err
}

func (c *Client) SCMRepoLinkManualReviewApp(app, pullRequestId string) error {
	_, err := c.ScalingoAPI().Do(&http.APIRequest{
		Method:   "POST",
		Endpoint: "/apps/" + app + "/scm_repo_link/manual_review_app",
		Expected: http.Statuses{200},
		Params:   map[string]string{"pull_request_id": pullRequestId},
	})
	return err
}

func (c *Client) SCMRepoLinkDeployments(app string) ([]*Deployment, error) {
	var res SCMRepoLinkDeploymentsResponse

	err := c.ScalingoAPI().DoRequest(&http.APIRequest{
		Method:   "GET",
		Endpoint: "/apps/" + app + "/scm_repo_link/deployments",
		Expected: http.Statuses{200},
	}, &res)
	if err != nil {
		return nil, err
	}
	return res.Deployments, nil
}

func (c *Client) SCMRepoLinkReviewApps(app string) ([]*ReviewApp, error) {
	var res SCMRepoLinkReviewAppsResponse

	err := c.ScalingoAPI().DoRequest(&http.APIRequest{
		Method:   "GET",
		Endpoint: "/apps/" + app + "/scm_repo_link/review_apps",
		Expected: http.Statuses{200},
	}, &res)
	if err != nil {
		return nil, err
	}
	return res.ReviewApps, nil
}
