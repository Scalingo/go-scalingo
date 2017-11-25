package scalingo

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"golang.org/x/oauth2"
	"gopkg.in/errgo.v1"
)

const defaultAuthUrl = "https://auth.scalingo.com"

type TokenGenerator interface {
	Token() (string, error)
}

type OAuthTokenGenerator struct {
	source       oauth2.TokenSource
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	Expiry       time.Time `json:"expiry"`
	TokenURL     string    `json:"token_url"`
}

type OAuthApplication struct {
	OAuthID     string `json:"oauth_id"`
	OAuthSecret string `json:"oauth_secret"`
}

type LoginParams struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
	OTP        string `json:"otp"`
}

func (c *Client) IsOTPRequired(err error) bool {
	rerr, ok := err.(*RequestFailedError)
	if !ok {
		return false
	}

	if rerr.Message == "OTP Required" {
		return true
	}
	return false
}

func (c *Client) GetOAuthCredentials(params LoginParams) (*OAuthApplication, error) {
	req := &APIRequest{
		Client:   c,
		NoAuth:   true,
		Method:   "POST",
		URL:      fmt.Sprintf("%s/v1/client/oauth", c.AuthURL()),
		Username: params.Identifier,
		Password: params.Password,
		OTP:      params.OTP,
	}

	resp, err := req.Do()
	if err != nil {
		if c.IsOTPRequired(err) {
			return nil, errgo.New("OTP Required")
		}
		return nil, errgo.Notef(err, "fail to fetch oauth credentials")
	}

	var login OAuthApplication
	err = ParseJSON(resp, &login)
	if err != nil {
		return nil, errgo.Notef(err, "invalid response from authentication service")
	}

	return &login, nil
}

func (c *Client) GetOAuthTokenGenerator(app *OAuthApplication, login LoginParams, scopes []string, redirectURL string) (*OAuthTokenGenerator, error) {
	config := &oauth2.Config{
		ClientID:     app.OAuthID,
		ClientSecret: app.OAuthSecret,
		Scopes:       scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  fmt.Sprintf("%s/oauth/authorize", c.AuthURL()),
			TokenURL: fmt.Sprintf("%s/oauth/token", c.AuthURL()),
		},
		RedirectURL: redirectURL,
	}

	authUrl := config.AuthCodeURL("state", oauth2.AccessTypeOnline)

	// Do not follow redirections
	client := http.Client{
		CheckRedirect: func(r *http.Request, v []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	req, err := http.NewRequest("GET", authUrl, nil)
	if err != nil {
		return nil, errgo.Notef(err, "fail to build auth request")
	}

	req.SetBasicAuth(login.Identifier, login.Password)
	if login.OTP != "" {
		req.Header.Add("X-Authorization-OTP", login.OTP)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, errgo.Notef(err, "fail to make auth request")
	}

	location := resp.Header.Get("Location")
	respUrl, err := url.Parse(location)
	if err != nil {
		return nil, errgo.Notef(err, "fail to parse auth response")
	}
	code := respUrl.Query().Get("code")
	token, err := config.Exchange(context.Background(), code)
	if err != nil {
		return nil, errgo.Notef(err, "fail to exchange key with auth service")
	}

	return &OAuthTokenGenerator{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		Expiry:       token.Expiry,
		TokenURL:     config.Endpoint.TokenURL,
	}, nil
}

func (c *Client) AuthURL() string {
	if os.Getenv("SCALINGO_AUTH_URL") != "" {
		return os.Getenv("SCALINGO_AUTH_URL")
	}
	return defaultAuthUrl
}

func (t *OAuthTokenGenerator) Source() oauth2.TokenSource {
	if t.source != nil {
		return t.source
	}
	configSource := &oauth2.Config{
		Endpoint: oauth2.Endpoint{
			TokenURL: t.TokenURL,
		},
	}

	source := configSource.TokenSource(context.Background(), &oauth2.Token{
		Expiry:       t.Expiry,
		AccessToken:  t.AccessToken,
		RefreshToken: t.RefreshToken,
	})

	t.source = source

	return t.source
}

func (t *OAuthTokenGenerator) Token() (string, error) {
	token, err := t.Source().Token()
	if err != nil {
		return "", errgo.Notef(err, "fail to get oauth2 token")
	}
	return token.AccessToken, nil
}
