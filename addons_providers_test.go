package scalingo

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddonProvidersClient(t *testing.T) {
	ctx := t.Context()

	tests := []struct {
		action           string
		testedClientCall func(c AddonProvidersService) error
		expectedEndpoint string
		expectedMethod   string
		expectedQuery    string
		response         interface{}
		responseStatus   int
		noBody           bool
	}{
		{
			action: "list",
			testedClientCall: func(c AddonProvidersService) error {
				_, err := c.AddonProvidersList(ctx)
				return err
			},
			expectedEndpoint: "/v1/addon_providers",
			expectedMethod:   "GET",
			expectedQuery:    "",
			response:         AddonProvidersListResponse{},
		},
		{
			action: "plans list",
			testedClientCall: func(c AddonProvidersService) error {
				_, err := c.AddonProviderPlansList(ctx, "the-addon", AddonProviderPlansListOpts{ShowAll: false})
				return err
			},
			expectedEndpoint: "/v1/addon_providers/the-addon/plans",
			expectedMethod:   "GET",
			expectedQuery:    "show_all=false",
			response:         AddonProviderPlansListResponse{},
		},
		{
			action: "all plans list",
			testedClientCall: func(c AddonProvidersService) error {
				_, err := c.AddonProviderPlansList(ctx, "the-addon", AddonProviderPlansListOpts{ShowAll: true})
				return err
			},
			expectedEndpoint: "/v1/addon_providers/the-addon/plans",
			expectedMethod:   "GET",
			expectedQuery:    "show_all=true",
			response:         AddonProviderPlansListResponse{},
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
					assert.Equal(t, test.expectedQuery, r.URL.RawQuery)
					if run.invalidResponse {
						w.WriteHeader(500)
						_, err := w.Write([]byte("INVALID"))
						assert.NoError(t, err)
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
