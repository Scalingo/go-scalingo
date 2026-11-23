package scalingo

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestAppsClient_Update(t *testing.T) {
	ctx := t.Context()
	const appName = "my-app"

	runs := map[string]struct {
		testedClientCall func(c AppsService) error
		expectedEndpoint string
		expectedMethod   string
		expectedParams   string
		response         any
		responseStatus   int
	}{
		"it should enable the app router_logs attribute": {
			testedClientCall: func(c AppsService) error {
				_, err := c.AppsRouterLogs(ctx, appName, true)
				return err
			},
			expectedEndpoint: "/v1/apps/my-app",
			expectedMethod:   http.MethodPut,
			expectedParams:   `{"router_logs":true}`,
			response:         &AppResponse{},
			responseStatus:   http.StatusOK,
		},
		"it should disable the app router_logs attribute": {
			testedClientCall: func(c AppsService) error {
				_, err := c.AppsRouterLogs(ctx, appName, false)
				return err
			},
			expectedEndpoint: "/v1/apps/my-app",
			expectedMethod:   http.MethodPut,
			expectedParams:   `{"router_logs":false}`,
			response:         &AppResponse{},
			responseStatus:   http.StatusOK,
		},

		"it should enable the app force_https attribute": {
			testedClientCall: func(c AppsService) error {
				_, err := c.AppsForceHTTPS(ctx, appName, true)
				return err
			},
			expectedEndpoint: "/v1/apps/my-app",
			expectedMethod:   http.MethodPut,
			expectedParams:   `{"force_https":true}`,
			response:         &AppResponse{},
			responseStatus:   http.StatusOK,
		},
		"it should disable the app force_https attribute": {
			testedClientCall: func(c AppsService) error {
				_, err := c.AppsForceHTTPS(ctx, appName, false)
				return err
			},
			expectedEndpoint: "/v1/apps/my-app",
			expectedMethod:   http.MethodPut,
			expectedParams:   `{"force_https":false}`,
			response:         &AppResponse{},
			responseStatus:   http.StatusOK,
		},

		"it should enable the app sticky_session attribute": {
			testedClientCall: func(c AppsService) error {
				_, err := c.AppsStickySession(ctx, appName, true)
				return err
			},
			expectedEndpoint: "/v1/apps/my-app",
			expectedMethod:   http.MethodPut,
			expectedParams:   `{"sticky_session":true}`,
			response:         &AppResponse{},
			responseStatus:   http.StatusOK,
		},
		"it should disable the app sticky_session attribute": {
			testedClientCall: func(c AppsService) error {
				_, err := c.AppsStickySession(ctx, appName, false)
				return err
			},
			expectedEndpoint: "/v1/apps/my-app",
			expectedMethod:   http.MethodPut,
			expectedParams:   `{"sticky_session":false}`,
			response:         &AppResponse{},
			responseStatus:   http.StatusOK,
		},

		"it should disable multipart upload on the app router": {
			testedClientCall: func(c AppsService) error {
				_, err := c.AppsRouterDisableMultipartUpload(ctx, appName)
				return err
			},
			expectedEndpoint: "/v1/apps/my-app",
			expectedMethod:   http.MethodPut,
			expectedParams:   `{"router_disable_multipart_upload":true}`,
			response:         &AppResponse{},
			responseStatus:   http.StatusOK,
		},
		"it should enable multipart upload on the app router": {
			testedClientCall: func(c AppsService) error {
				_, err := c.AppsRouterEnableMultipartUpload(ctx, appName)
				return err
			},
			expectedEndpoint: "/v1/apps/my-app",
			expectedMethod:   http.MethodPut,
			expectedParams:   `{"router_disable_multipart_upload":false}`,
			response:         &AppResponse{},
			responseStatus:   http.StatusOK,
		},
	}

	for msg, run := range runs {
		t.Run(msg, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, run.expectedMethod, r.Method)
				assert.Equal(t, run.expectedEndpoint, r.URL.Path)
				buf := new(bytes.Buffer)
				_, err := buf.ReadFrom(r.Body)
				if !assert.NoError(t, err) {
					return
				}
				assert.Equal(t, run.expectedParams, buf.String())

				if run.responseStatus != 0 {
					w.WriteHeader(run.responseStatus)
				}
				if run.response != nil {
					err := json.NewEncoder(w).Encode(&run.response)
					assert.NoError(t, err)
				}
			}))
			defer ts.Close()

			c, err := New(ctx, ClientConfig{
				APIEndpoint: ts.URL,
				APIToken:    "test",
			})
			require.NoError(t, err)

			c.authClient = MockAuth(ctrl)

			err = run.testedClientCall(c)
			require.NoError(t, err)
		})
	}
}

func TestAppsFirewallRulesList(t *testing.T) {
	ctx := t.Context()
	client, cleanup := newTestClient(t, appsTestRequest{
		expectedEndpoint: "/v1/apps/my-app/firewall_rules",
		expectedMethod:   http.MethodGet,
		response: &AppFirewallRulesResponse{Rules: []AppFirewallRule{{
			ID:    "rule-id",
			AppID: "app-id",
			CIDR:  "192.0.2.0/24",
			Label: "office",
		}}},
		responseStatus: http.StatusOK,
	})
	defer cleanup()

	rules, err := client.AppsFirewallRulesList(ctx, "my-app")
	require.NoError(t, err)
	require.Len(t, rules, 1)
	assert.Equal(t, "app-id", rules[0].AppID)
}

func TestAppsFirewallRuleShow(t *testing.T) {
	ctx := t.Context()
	client, cleanup := newTestClient(t, appsTestRequest{
		expectedEndpoint: "/v1/apps/my-app/firewall_rules/rule-id",
		expectedMethod:   http.MethodGet,
		response: &AppFirewallRuleResponse{Rule: &AppFirewallRule{
			ID:    "rule-id",
			AppID: "app-id",
			CIDR:  "192.0.2.0/24",
			Label: "office",
		}},
		responseStatus: http.StatusOK,
	})
	defer cleanup()

	rule, err := client.AppsFirewallRuleShow(ctx, "my-app", "rule-id")
	require.NoError(t, err)
	assert.Equal(t, "app-id", rule.AppID)
}

func TestAppsFirewallRuleCreate(t *testing.T) {
	ctx := t.Context()
	client, cleanup := newTestClient(t, appsTestRequest{
		expectedEndpoint: "/v1/apps/my-app/firewall_rules",
		expectedMethod:   http.MethodPost,
		expectedParams:   `{"firewall_rule":{"cidr":"192.0.2.0/24","label":"office"}}`,
		response: &AppFirewallRuleResponse{Rule: &AppFirewallRule{
			ID:    "rule-id",
			AppID: "app-id",
			CIDR:  "192.0.2.0/24",
			Label: "office",
		}},
		responseStatus: http.StatusCreated,
	})
	defer cleanup()

	rule, err := client.AppsFirewallRuleCreate(ctx, "my-app", AppFirewallRuleParams{
		CIDR:  "192.0.2.0/24",
		Label: "office",
	})
	require.NoError(t, err)
	assert.Equal(t, "app-id", rule.AppID)
}

func TestAppsFirewallRuleDelete(t *testing.T) {
	ctx := t.Context()
	client, cleanup := newTestClient(t, appsTestRequest{
		expectedEndpoint: "/v1/apps/my-app/firewall_rules/rule-id",
		expectedMethod:   http.MethodDelete,
		responseStatus:   http.StatusNoContent,
	})
	defer cleanup()

	err := client.AppsFirewallRuleDelete(ctx, "my-app", "rule-id")
	require.NoError(t, err)
}

func TestAppsClient_List(t *testing.T) {
	const organizationID = "org-11111111-2222-3333-4444-555555555555"
	const organizationQuery = "organization_id=" + organizationID

	tests := map[string]struct {
		list           func(ctx context.Context, c *Client) ([]*App, error)
		responseStatus int
		responseBody   string
		expectedQuery  string
		expected       []*App
		expectedError  string
	}{
		"it should not send any organization_id when listing the personal apps": {
			list: func(ctx context.Context, c *Client) ([]*App, error) {
				return c.AppsList(ctx)
			},
			responseStatus: http.StatusOK,
			responseBody:   `{"apps": [{"id": "app-personal", "name": "personal-app", "project": {"id": "pr-1", "name": "default"}}]}`,
			expectedQuery:  "",
			expected:       []*App{{ID: "app-personal", Name: "personal-app", Project: appProject{ID: "pr-1", Name: "default"}}},
		},
		"it should send the organization_id when listing the apps of an organization": {
			list: func(ctx context.Context, c *Client) ([]*App, error) {
				return c.AppsListByOrganization(ctx, organizationID)
			},
			responseStatus: http.StatusOK,
			responseBody:   `{"apps": [{"id": "app-2", "name": "acme-front", "project": {"id": "pr-2", "name": "front"}}]}`,
			expectedQuery:  organizationQuery,
			expected:       []*App{{ID: "app-2", Name: "acme-front", Project: appProject{ID: "pr-2", Name: "front"}}},
		},
		"it should return an empty list when the organization has no app": {
			list: func(ctx context.Context, c *Client) ([]*App, error) {
				return c.AppsListByOrganization(ctx, organizationID)
			},
			responseStatus: http.StatusOK,
			responseBody:   `{"apps": []}`,
			expectedQuery:  organizationQuery,
			expected:       []*App{},
		},
		"it should return an error when the API refuses to list the apps of the organization": {
			list: func(ctx context.Context, c *Client) ([]*App, error) {
				return c.AppsListByOrganization(ctx, organizationID)
			},
			responseStatus: http.StatusForbidden,
			responseBody:   `{"error": "You are not allowed to access this organization"}`,
			expectedQuery:  organizationQuery,
			expectedError:  "list apps",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			apiServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodGet, r.Method)
				assert.Equal(t, "/v1/apps", r.URL.Path)
				assert.Equal(t, test.expectedQuery, r.URL.RawQuery)
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			defer apiServer.Close()

			client, err := New(t.Context(), ClientConfig{
				APIEndpoint:          apiServer.URL,
				StaticTokenGenerator: NewStaticTokenGenerator("static-token"),
			})
			require.NoError(t, err)

			apps, err := test.list(t.Context(), client)
			if test.expectedError != "" {
				require.ErrorContains(t, err, test.expectedError)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expected, apps)
		})
	}
}

type appsTestRequest struct {
	expectedEndpoint string
	expectedMethod   string
	expectedParams   string
	response         any
	responseStatus   int
}
