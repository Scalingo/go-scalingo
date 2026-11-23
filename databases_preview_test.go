package scalingo

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPreviewClient_DatabasesList(t *testing.T) {
	const organizationID = "org-11111111-2222-3333-4444-555555555555"
	const organizationQuery = "organization_id=" + organizationID

	tests := map[string]struct {
		list           func(ctx context.Context, c *PreviewClient) ([]DatabaseNG, error)
		responseStatus int
		responseBody   string
		expectedQuery  string
		expectedError  string
	}{
		"it should not send any organization_id when listing the personal databases": {
			list: func(ctx context.Context, c *PreviewClient) ([]DatabaseNG, error) {
				return c.DatabasesList(ctx)
			},
			responseStatus: http.StatusOK,
			responseBody:   `{"databases": []}`,
			expectedQuery:  "",
		},
		"it should send the organization_id even when the organization has no database": {
			list: func(ctx context.Context, c *PreviewClient) ([]DatabaseNG, error) {
				return c.DatabasesListByOrganization(ctx, organizationID)
			},
			responseStatus: http.StatusOK,
			responseBody:   `{"databases": []}`,
			expectedQuery:  organizationQuery,
		},
		"it should return an error when the API refuses to list the databases of the organization": {
			list: func(ctx context.Context, c *PreviewClient) ([]DatabaseNG, error) {
				return c.DatabasesListByOrganization(ctx, organizationID)
			},
			responseStatus: http.StatusForbidden,
			responseBody:   `{"error": "You are not allowed to access this organization"}`,
			expectedQuery:  organizationQuery,
			expectedError:  "list databases",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			callCount := 0
			apiServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				callCount++
				assert.Equal(t, http.MethodGet, r.Method)
				assert.Equal(t, "/v1/databases", r.URL.Path)
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

			databases, err := test.list(t.Context(), client.Preview())
			assert.Equal(t, 1, callCount)
			if test.expectedError != "" {
				require.ErrorContains(t, err, test.expectedError)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, []DatabaseNG{}, databases)
		})
	}
}
