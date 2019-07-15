package scalingo

import (
	"time"

	"github.com/Scalingo/go-scalingo/http"
)

type ScmRepoLinkService interface {
	ScmRepoLinkShow(app string) (*ScmRepoLink, error)
	ScmRepoLinkAdd(app string, params ScmRepoLinkParams) (*ScmRepoLink, error)
	ScmRepoLinkUpdate(app, id string, params ScmRepoLinkParams) (*ScmRepoLink, error)
	ScmRepoLinkDelete(app, id string) error

	ScmRepoLinkManualDeploy(app, id, branch string) error
	ScmRepoLinkManualReviewApp(app, id, pullRequestId string) error
	ScmRepoLinkDeployments(app, id string) ([]*Deployment, error)
	ScmRepoLinkReviewApps(app, id string) ([]*ReviewApp, error)
}

type ScmRepoLinkParams struct {
	Source                   *string `json:"source,omitempty"`
	Branch                   *string `json:"branch,omitempty"`
	AuthIntegrationID        *string `json:"auth_integration_id,omitempty"`
	AutoDeployEnabled        *bool   `json:"auto_deploy_enabled,omitempty"`
	DeployReviewAppsEnabled  *bool   `json:"deploy_review_apps_enabled,omitempty"`
	DestroyOnCloseEnabled    *bool   `json:"delete_on_close_enabled,omitempty"`
	HoursBeforeDeleteOnClose *uint   `json:"hours_before_delete_on_close,omitempty"`
	DestroyStaleEnabled      *bool   `json:"delete_stale_enabled,omitempty"`
	HoursBeforeDeleteStale   *uint   `json:"hours_before_delete_stale,omitempty"`
}

type ScmRepoLink struct {
	ID                       string            `json:"id"`
	AppID                    string            `json:"app_id"`
	Linker                   ScmRepoLinkLinker `json:"linker"`
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
	DeleteOnStaleEnabled     bool              `json:"delete_stale_enabled"`
	HoursBeforeDeleteOnClose uint              `json:"hours_before_delete_on_close"`
	HoursBeforeDeleteStale   uint              `json:"hours_before_delete_stale"`
	LastAutoDeployAt         time.Time         `json:"last_auto_deploy_at"`
}

type ScmRepoLinkLinker struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	ID       string `json:"id"`
}

type ScmRepoLinkResponse struct {
	ScmRepoLink *ScmRepoLink `json:"scm_repo_link"`
}

type ScmRepoLinksResponse struct {
	ScmRepoLinks []*ScmRepoLink `json:"scm_repo_links"`
}

type ScmRepoLinkDeploymentsResponse struct {
	Deployments []*Deployment `json:"deployments"`
}

type ScmRepoLinkReviewAppsResponse struct {
	ReviewApps []*ReviewApp `json:"review_apps"`
}

var _ ScmRepoLinkService = (*Client)(nil)

func (c *Client) ScmRepoLinkShow(app string) (*ScmRepoLink, error) {
	var res ScmRepoLinksResponse
	err := c.ScalingoAPI().SubresourceList("apps", app, "scm_repo_links", nil, &res)
	if err != nil {
		return nil, err
	}
	if len(res.ScmRepoLinks) == 0 {
		return nil, nil
	}
	return res.ScmRepoLinks[0], nil
}

func (c *Client) ScmRepoLinkAdd(app string, params ScmRepoLinkParams) (*ScmRepoLink, error) {
	payload := map[string]ScmRepoLinkParams{
		"scm_repo_link": params,
	}

	var res ScmRepoLinkResponse
	err := c.ScalingoAPI().SubresourceAdd("apps", app, "scm_repo_links", payload, &res)
	if err != nil {
		return nil, err
	}
	return res.ScmRepoLink, nil
}

func (c *Client) ScmRepoLinkUpdate(app, id string, params ScmRepoLinkParams) (*ScmRepoLink, error) {
	payload := map[string]ScmRepoLinkParams{
		"scm_repo_link": params,
	}

	var res ScmRepoLinkResponse
	err := c.ScalingoAPI().SubresourceUpdate("apps", app, "scm_repo_links", id, payload, &res)
	if err != nil {
		return nil, err
	}
	return res.ScmRepoLink, nil
}

func (c *Client) ScmRepoLinkDelete(app, id string) error {
	return c.ScalingoAPI().SubresourceDelete("apps", app, "scm_repo_links", id)
}

func (c *Client) ScmRepoLinkManualDeploy(app, id, branch string) error {
	_, err := c.ScalingoAPI().Do(&http.APIRequest{
		Method:   "POST",
		Endpoint: "/apps/" + app + "/scm_repo_links/" + id + "/manual_deploy",
		Expected: http.Statuses{200},
		Params:   map[string]string{"branch": branch},
	})
	return err
}

func (c *Client) ScmRepoLinkManualReviewApp(app, id, pullRequestId string) error {
	_, err := c.ScalingoAPI().Do(&http.APIRequest{
		Method:   "POST",
		Endpoint: "/apps/" + app + "/scm_repo_links/" + id + "/manual_review_app",
		Expected: http.Statuses{200},
		Params:   map[string]string{"pull_request_id": pullRequestId},
	})
	return err
}

func (c *Client) ScmRepoLinkDeployments(app, id string) ([]*Deployment, error) {
	var res ScmRepoLinkDeploymentsResponse

	err := c.ScalingoAPI().DoRequest(&http.APIRequest{
		Method:   "GET",
		Endpoint: "/apps/" + app + "/scm_repo_links/" + id + "/deployments",
		Expected: http.Statuses{200},
	}, &res)
	if err != nil {
		return nil, err
	}
	return res.Deployments, nil
}

func (c *Client) ScmRepoLinkReviewApps(app, id string) ([]*ReviewApp, error) {
	var res ScmRepoLinkReviewAppsResponse

	err := c.ScalingoAPI().DoRequest(&http.APIRequest{
		Method:   "GET",
		Endpoint: "/apps/" + app + "/scm_repo_links/" + id + "/review_apps",
		Expected: http.Statuses{200},
	}, &res)
	if err != nil {
		return nil, err
	}
	return res.ReviewApps, nil
}
