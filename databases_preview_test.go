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

func TestPreviewClient_DatabasesListByOrganization_Hydration(t *testing.T) {
	const (
		organizationID = "org-11111111-2222-3333-4444-555555555555"
		databaseID     = "db-app-id"
		appName        = "acme-postgresql"
		projectID      = "pr-acme-uuid"
		projectName    = "acme-project"
		addonID        = "ad-postgresql-id"
		addonToken     = "addon-token"
	)

	apiServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body string
		switch r.Method + " " + r.URL.Path {
		case http.MethodGet + " /v1/databases":
			assert.Equal(t, "organization_id="+organizationID, r.URL.RawQuery)
			body = `{"databases": [{"id": "` + databaseID + `", "name": "` + appName + `", "project_id": "` + projectID + `",
				"organization_id": "` + organizationID + `", "technology": "postgresql", "plan": "postgresql-dr-starter-4096"}]}`
		case http.MethodGet + " /v1/apps/" + appName + "/addons":
			body = `{"addons": [{"id": "` + addonID + `"}]}`
		case http.MethodGet + " /v1/apps/" + appName:
			body = `{"app": {"id": "` + databaseID + `", "name": "` + appName + `", "project": {"id": "` + projectID + `", "name": "` + projectName + `"}}}`
		case http.MethodPost + " /v1/apps/" + appName + "/addons/" + addonID + "/token":
			body = `{"addon": {"token": "` + addonToken + `"}}`
		case http.MethodGet + " /api/databases/" + addonID:
			assert.Equal(t, "Bearer "+addonToken, r.Header.Get("Authorization"))
			body = `{"database": {"id": "` + addonID + `", "status": "running"}}`
		default:
			assert.Failf(t, "unexpected request", "%s %s", r.Method, r.URL.String())
			w.WriteHeader(http.StatusNotFound)
			return
		}
		_, err := w.Write([]byte(body))
		assert.NoError(t, err)
	}))
	defer apiServer.Close()

	client, err := New(t.Context(), ClientConfig{
		APIEndpoint:          apiServer.URL,
		DatabaseAPIEndpoint:  apiServer.URL,
		StaticTokenGenerator: NewStaticTokenGenerator("static-token"),
	})
	require.NoError(t, err)

	databases, err := client.Preview().DatabasesListByOrganization(t.Context(), organizationID)
	require.NoError(t, err)

	require.Len(t, databases, 1)
	database := databases[0]
	assert.Equal(t, databaseID, database.ID)
	assert.Equal(t, appName, database.Name)
	assert.Equal(t, projectID, database.ProjectID)
	assert.Equal(t, projectName, database.App.Project.Name)
	assert.Equal(t, DatabaseStatusRunning, database.Database.Status)
}
