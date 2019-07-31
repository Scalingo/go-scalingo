package scalingo

import (
	"gopkg.in/errgo.v1"

	"github.com/Scalingo/go-scalingo/http"
)

type SCMType string

// Type of SCM integrations
const (
	SCMGithubType           SCMType = "github"             // GitHub
	SCMGithubEnterpriseType SCMType = "github-enterprise"  // GitHub Enterprise (private instance)
	SCMGitlabType           SCMType = "gitlab"             // GitLab.com
	SCMGitlabSelfHostedType SCMType = "gitlab-self-hosted" // GitLab self-hosted (private instance)
)

func (t SCMType) Str() string {
	return string(t)
}

type SCMIntegrationsService interface {
	SCMIntegrationsList() ([]SCMIntegration, error)
	SCMIntegrationsCreate(scmType SCMType, url string, accessToken string) (*SCMIntegration, error)
	SCMIntegrationsDestroy(id string) error
}

var _ SCMIntegrationsService = (*Client)(nil)

type SCMIntegration struct {
	ID          string  `json:"id,omitempty"`
	SCMType     SCMType `json:"scm_type"`
	Url         string  `json:"url"`
	AccessToken string  `json:"access_token"`
	Uid         string  `json:"uid,omitempty"`
	Username    string  `json:"username,omitempty"`
	Email       string  `json:"email,omitempty"`
	AvatarUrl   string  `json:"avatar_url,omitempty"`
	ProfileUrl  string  `json:"profile_url,omitempty"`
}

type SCMIntegrationRes struct {
	SCMIntegration SCMIntegration `json:"scm_integration"`
}

type SCMIntegrationsRes struct {
	SCMIntegrations []SCMIntegration `json:"scm_integrations"`
}

func (c *Client) SCMIntegrationsList() ([]SCMIntegration, error) {
	var res SCMIntegrationsRes

	err := c.AuthAPI().ResourceList("scm_integrations", nil, &res)
	if err != nil {
		return nil, errgo.Mask(err)
	}
	return res.SCMIntegrations, nil
}

func (c *Client) SCMIntegrationsShow(id string) (*SCMIntegration, error) {
	var res SCMIntegrationRes

	err := c.AuthAPI().ResourceGet("scm_integrations", id, nil, &res)
	if err != nil {
		return nil, errgo.Mask(err)
	}
	return &res.SCMIntegration, nil
}

func (c *Client) SCMIntegrationsCreate(scmType SCMType, url string, accessToken string) (*SCMIntegration, error) {
	payload := SCMIntegrationRes{SCMIntegration{
		SCMType:     scmType,
		Url:         url,
		AccessToken: accessToken,
	}}
	var res SCMIntegrationRes

	err := c.AuthAPI().ResourceAdd("scm_integrations", payload, &res)
	if err != nil {
		return nil, errgo.Mask(err)
	}

	return &res.SCMIntegration, nil
}

func (c *Client) SCMIntegrationsDestroy(id string) error {
	err := c.AuthAPI().ResourceDelete("scm_integrations", id)
	if err != nil {
		return errgo.Mask(err)
	}
	return nil
}

func (c *Client) SCMIntegrationsImportKeys(id string) ([]Key, error) {
	var res KeysRes

	var err = c.AuthAPI().DoRequest(&http.APIRequest{
		Method:   "POST",
		Endpoint: "/scm_integrations/" + id + "/import_keys",
		Params:   nil,
		Expected: http.Statuses{201},
	}, &res)
	if err != nil {
		return nil, errgo.Mask(err)
	}
	return res.Keys, nil
}
