package scalingo

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDomainsClient_DomainCanonical(t *testing.T) {
	appName := "my-app"
	domainID := "domain-id"

	runs := map[string]struct {
		testedClientCall func(c DomainsService) error
		expectedEndpoint string
		expectedMethod   string
		expectedParams   string
		responseStatus   int
		mockDomainsList  func(t *testing.T, w http.ResponseWriter, r *http.Request)
		expectedError    string
	}{
		"it should set the domain as canonical": {
			testedClientCall: func(c DomainsService) error {
				_, err := c.DomainSetCanonical(appName, domainID)
				return err
			},
			expectedEndpoint: "/v1/apps/my-app/domains/domain-id",
			expectedMethod:   "PATCH",
			expectedParams:   `"canonical":true`,
			responseStatus:   200,
		},
		"it should unset the domain as canonical": {
			mockDomainsList: func(t *testing.T, w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
				err := json.NewEncoder(w).Encode(DomainsRes{Domains: []Domain{
					{
						ID:        "domain-id",
						Canonical: true,
					},
				}})
				assert.NoError(t, err)
			},
			testedClientCall: func(c DomainsService) error {
				_, err := c.DomainUnsetCanonical(appName)
				return err
			},
			expectedEndpoint: "/v1/apps/my-app/domains/domain-id",
			expectedMethod:   "PATCH",
			expectedParams:   `"canonical":false`,
			responseStatus:   200,
		},
		"it should return an error if there is no canonical domain": {
			mockDomainsList: func(t *testing.T, w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
				err := json.NewEncoder(w).Encode(DomainsRes{Domains: []Domain{}})
				assert.NoError(t, err)
			},
			testedClientCall: func(c DomainsService) error {
				_, err := c.DomainUnsetCanonical(appName)
				return err
			},
			expectedError: "no canonical domain configured",
		},
	}

	for msg, run := range runs {
		t.Run(msg, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// If request domains list
				if r.Method == "GET" && r.URL.Path == "/v1/apps/my-app/domains" {
					require.NotNil(t, run.mockDomainsList)
					run.mockDomainsList(t, w, r)
					return
				}
				assert.Equal(t, run.expectedMethod, r.Method)
				assert.Equal(t, run.expectedEndpoint, r.URL.Path)
				buf := new(bytes.Buffer)
				buf.ReadFrom(r.Body)
				assert.Contains(t, buf.String(), run.expectedParams)

				if run.responseStatus != 0 {
					w.WriteHeader(run.responseStatus)
				}
				err := json.NewEncoder(w).Encode(DomainRes{})
				assert.NoError(t, err)
			}))
			defer ts.Close()

			c, err := New(ClientConfig{
				APIEndpoint: ts.URL,
				APIToken:    "test",
			})
			c.authClient = MockAuth(ctrl)
			require.NoError(t, err)

			err = run.testedClientCall(c)
			if run.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), run.expectedError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
