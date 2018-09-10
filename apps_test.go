package scalingo

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAppsClient_Update(t *testing.T) {
	appName := "my-app"

	runs := map[string]struct {
		testedClientCall func(c AppsService) error
		expectedEndpoint string
		expectedMethod   string
		expectedParams   string
		response         interface{}
		responseStatus   int
	}{
		"it should enable the app force-https attribute": {
			testedClientCall: func(c AppsService) error {
				_, err := c.AppsForceHTTPS(appName, true)
				return err
			},
			expectedEndpoint: "/v1/apps/my-app",
			expectedMethod:   "PUT",
			expectedParams:   `{"force_https":true}`,
			response:         &AppResponse{},
			responseStatus:   200,
		},
		"it should disable the app force-https attribute": {
			testedClientCall: func(c AppsService) error {
				_, err := c.AppsForceHTTPS(appName, false)
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
				_, err := c.AppsStickySession(appName, true)
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
				_, err := c.AppsStickySession(appName, false)
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
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, run.expectedMethod, r.Method)
				assert.Equal(t, run.expectedEndpoint, r.URL.Path)
				buf := new(bytes.Buffer)
				buf.ReadFrom(r.Body)
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

			c := NewClient(ClientConfig{
				Endpoint:       ts.URL,
				TokenGenerator: NewStaticTokenGenerator("test"),
			})

			err := run.testedClientCall(c)
			require.NoError(t, err)
		})
	}
}
