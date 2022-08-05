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

func TestDomainsClient_DomainCanonical(t *testing.T) {
	ctx := context.Background()
	appName := "my-app"
	domainID := "domain-id"

	tests := map[string]struct {
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
				_, err := c.DomainSetCanonical(ctx, appName, domainID)
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
				_, err := c.DomainUnsetCanonical(ctx, appName)
				return err
			},
			expectedEndpoint: "/v1/apps/my-app/domains/domain-id",
			expectedMethod:   "PATCH",
			expectedParams:   `"canonical":false`,
			responseStatus:   200,
		},
		"it should unset the domain certificate": {
			testedClientCall: func(c DomainsService) error {
				_, err := c.DomainUnsetCertificate(ctx, appName, domainID)
				return err
			},
			expectedEndpoint: "/v1/apps/my-app/domains/domain-id",
			expectedMethod:   "PATCH",
			expectedParams:   `"tlscert":"","tlskey":""`,
			responseStatus:   200,
		},
		"it should return an error if there is no canonical domain": {
			mockDomainsList: func(t *testing.T, w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
				err := json.NewEncoder(w).Encode(DomainsRes{Domains: []Domain{}})
				assert.NoError(t, err)
			},
			testedClientCall: func(c DomainsService) error {
				_, err := c.DomainUnsetCanonical(ctx, appName)
				return err
			},
			expectedError: "no canonical domain configured",
		},
	}

	for msg, test := range tests {
		t.Run(msg, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// If request domains list
				if r.Method == "GET" && r.URL.Path == "/v1/apps/my-app/domains" {
					require.NotNil(t, test.mockDomainsList)
					test.mockDomainsList(t, w, r)
					return
				}
				assert.Equal(t, test.expectedMethod, r.Method)
				assert.Equal(t, test.expectedEndpoint, r.URL.Path)
				buf := new(bytes.Buffer)
				buf.ReadFrom(r.Body)
				assert.Contains(t, buf.String(), test.expectedParams)

				if test.responseStatus != 0 {
					w.WriteHeader(test.responseStatus)
				}
				err := json.NewEncoder(w).Encode(DomainRes{})
				assert.NoError(t, err)
			}))
			defer ts.Close()

			scalingoClient, err := New(ctx, ClientConfig{
				APIEndpoint: ts.URL,
				APIToken:    "test",
			})
			scalingoClient.authClient = MockAuth(ctrl)
			require.NoError(t, err)

			err = test.testedClientCall(scalingoClient)
			if test.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), test.expectedError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
