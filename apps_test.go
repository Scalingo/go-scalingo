package scalingo

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAppsClient_Update(t *testing.T) {
	ctx := context.Background()
	appName := "my-app"

	runs := map[string]struct {
		testedClientCall func(c AppsService) error
		expectedEndpoint string
		expectedMethod   string
		expectedParams   string
		response         interface{}
		responseStatus   int
	}{
		"it should enable the app router_logs attribute": {
			testedClientCall: func(c AppsService) error {
				_, err := c.AppsRouterLogs(ctx, appName, true)
				return err
			},
			expectedEndpoint: "/v1/apps/my-app",
			expectedMethod:   "PUT",
			expectedParams:   `{"router_logs":true}`,
			response:         &AppResponse{},
			responseStatus:   200,
		},
		"it should disable the app router_logs attribute": {
			testedClientCall: func(c AppsService) error {
				_, err := c.AppsRouterLogs(ctx, appName, false)
				return err
			},
			expectedEndpoint: "/v1/apps/my-app",
			expectedMethod:   "PUT",
			expectedParams:   `{"router_logs":false}`,
			response:         &AppResponse{},
			responseStatus:   200,
		},

		"it should enable the app force_https attribute": {
			testedClientCall: func(c AppsService) error {
				_, err := c.AppsForceHTTPS(ctx, appName, true)
				return err
			},
			expectedEndpoint: "/v1/apps/my-app",
			expectedMethod:   "PUT",
			expectedParams:   `{"force_https":true}`,
			response:         &AppResponse{},
			responseStatus:   200,
		},
		"it should disable the app force_https attribute": {
			testedClientCall: func(c AppsService) error {
				_, err := c.AppsForceHTTPS(ctx, appName, false)
				return err
			},
			expectedEndpoint: "/v1/apps/my-app",
			expectedMethod:   "PUT",
			expectedParams:   `{"force_https":false}`,
			response:         &AppResponse{},
			responseStatus:   200,
		},

		"it should enable the app sticky_session attribute": {
			testedClientCall: func(c AppsService) error {
				_, err := c.AppsStickySession(ctx, appName, true)
				return err
			},
			expectedEndpoint: "/v1/apps/my-app",
			expectedMethod:   "PUT",
			expectedParams:   `{"sticky_session":true}`,
			response:         &AppResponse{},
			responseStatus:   200,
		},
		"it should disable the app sticky_session attribute": {
			testedClientCall: func(c AppsService) error {
				_, err := c.AppsStickySession(ctx, appName, false)
				return err
			},
			expectedEndpoint: "/v1/apps/my-app",
			expectedMethod:   "PUT",
			expectedParams:   `{"sticky_session":false}`,
			response:         &AppResponse{},
			responseStatus:   200,
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
				require.NoError(t, err)
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
