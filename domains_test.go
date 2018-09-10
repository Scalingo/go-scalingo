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

func TestDomainsClient_DomainCanonical(t *testing.T) {
	appName := "my-app"
	domainID := "domain-id"

	runs := map[string]struct {
		testedClientCall func(c DomainsService) error
		expectedEndpoint string
		expectedMethod   string
		expectedParams   string
		responseStatus   int
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
			testedClientCall: func(c DomainsService) error {
				_, err := c.DomainUnsetCanonical(appName, domainID)
				return err
			},
			expectedEndpoint: "/v1/apps/my-app/domains/domain-id",
			expectedMethod:   "PATCH",
			expectedParams:   `"canonical":false`,
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
				assert.Contains(t, buf.String(), run.expectedParams)

				if run.responseStatus != 0 {
					w.WriteHeader(run.responseStatus)
				}
				err := json.NewEncoder(w).Encode(DomainRes{})
				assert.NoError(t, err)
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
