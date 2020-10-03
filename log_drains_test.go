package scalingo

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLogDrainsClient(t *testing.T) {
	appName := "my-app"
	logDrainID := "my-id"
	logDrainURL := "tcp+tls://localhost:8080"
	addonID := "addon_uuid"

	tests := []struct {
		action           string
		testedClientCall func(c LogDrainsService) error
		expectedEndpoint string
		expectedMethod   string
		response         interface{}
		responseStatus   int
		noBody           bool
	}{
		{
			action: "list",
			testedClientCall: func(c LogDrainsService) error {
				_, err := c.LogDrainsList(appName)
				return err
			},
			expectedEndpoint: "/v1/apps/my-app/log_drains",
			expectedMethod:   "GET",
			response: LogDrainsRes{
				[]LogDrain{
					{
						AppID: logDrainID,
						URL:   logDrainURL,
					},
				}},
		},
		{
			action: "addon list",
			testedClientCall: func(c LogDrainsService) error {
				_, err := c.LogDrainsAddonList(appName, addonID)
				return err
			},
			expectedEndpoint: "/v1/apps/my-app/addons/" + addonID + "/log_drains",
			expectedMethod:   "GET",
			response: LogDrainsRes{
				[]LogDrain{
					{
						AppID: addonID,
						URL:   logDrainURL,
					},
				}},
		},
		{
			action: "add",
			testedClientCall: func(c LogDrainsService) error {
				_, err := c.LogDrainAdd(appName, LogDrainAddParams{
					Type: "syslog",
					Host: "localhost",
					Port: "8080",
				})
				return err
			},
			expectedEndpoint: "/v1/apps/my-app/log_drains",
			expectedMethod:   "POST",
			response: LogDrain{
				AppID: logDrainID,
				URL:   logDrainURL,
			},
			responseStatus: 201,
		},
		{
			action: "addon add",
			testedClientCall: func(c LogDrainsService) error {
				_, err := c.LogDrainAddonAdd(appName, addonID, LogDrainAddParams{
					Type: "syslog",
					Host: "localhost",
					Port: "8080",
				})
				return err
			},
			expectedEndpoint: "/v1/apps/my-app/addons/" + addonID + "/log_drains",
			expectedMethod:   "POST",
			response: LogDrain{
				AppID: logDrainID,
				URL:   logDrainURL,
			},
			responseStatus: 201,
		},
		{
			action: "remove",
			testedClientCall: func(c LogDrainsService) error {
				err := c.LogDrainRemove(appName, logDrainURL)
				return err
			},
			expectedEndpoint: "/v1/apps/my-app/log_drains",
			expectedMethod:   "DELETE",
			response:         nil,
			responseStatus:   204,
		},
		{
			action: "addon remove",
			testedClientCall: func(c LogDrainsService) error {
				err := c.LogDrainAddonRemove(appName, addonID, logDrainURL)
				return err
			},
			expectedEndpoint: "/v1/apps/my-app/addons/" + addonID + "/log_drains",
			expectedMethod:   "DELETE",
			response:         nil,
			responseStatus:   204,
		},
	}

	for _, test := range tests {
		for msg, run := range map[string]struct {
			invalidResponse bool
		}{
			"it should fail if it fails to " + test.action + "the subresource": {
				invalidResponse: true,
			},
			"it should succeed if it succeeds to " + test.action + " the subresource": {
				invalidResponse: false,
			},
		} {
			t.Run(msg, func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()
				handler := func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, test.expectedMethod, r.Method)
					assert.Equal(t, test.expectedEndpoint, r.URL.Path)
					if run.invalidResponse {
						w.WriteHeader(500)
						w.Write([]byte("INVALID"))
					} else {
						if test.responseStatus != 0 {
							w.WriteHeader(test.responseStatus)
						}
						if test.response != nil {
							err := json.NewEncoder(w).Encode(&test.response)
							assert.NoError(t, err)
						}
					}
				}
				ts := httptest.NewServer(http.HandlerFunc(handler))
				defer ts.Close()

				c, err := New(ClientConfig{
					APIEndpoint: ts.URL,
					APIToken:    "test",
				})
				require.NoError(t, err)

				c.authClient = MockAuth(ctrl)

				err = test.testedClientCall(c)
				if run.invalidResponse {
					require.Error(t, err)
				} else {
					require.NoError(t, err)
				}
			})
		}
	}
}
