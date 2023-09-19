package scalingo

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCronTasksClient_CronTasksGet(t *testing.T) {
	ctx := context.Background()
	const appName = "my-app"

	tests := []struct {
		action           string
		testedClientCall func(c CronTasksService) error
		expectedEndpoint string
		expectedMethod   string
		response         interface{}
		responseStatus   int
	}{
		{
			action: "get",
			testedClientCall: func(c CronTasksService) error {
				_, err := c.CronTasksGet(ctx, appName)
				return err
			},
			expectedEndpoint: "/v1/apps/my-app/cron_tasks",
			expectedMethod:   "GET",
			response:         CronTasks{},
			responseStatus:   200,
		},
	}

	for _, test := range tests {
		for msg, run := range map[string]struct {
			invalidResponse bool
		}{
			"it should fail if it fails to " + test.action + " the subresource": {
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
						_, err := w.Write([]byte("INVALID"))
						require.NoError(t, err)
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

				c, err := New(ctx, ClientConfig{
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
