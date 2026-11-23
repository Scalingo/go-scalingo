package scalingo

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPreviewClient_DatabasesList(t *testing.T) {
	t.Run("it should list the personal databases without organization_id", func(t *testing.T) {
		client := newDatabasesListTestClient(t, "", http.StatusOK, `{"databases": []}`)

		databases, err := client.DatabasesList(t.Context())
		require.NoError(t, err)
		assert.Equal(t, []DatabaseNG{}, databases)
	})
}

func TestPreviewClient_DatabasesListByOrganization(t *testing.T) {
	const organizationID = "org-11111111-2222-3333-4444-555555555555"
	const organizationQuery = "organization_id=" + organizationID

	t.Run("it should return an empty list when the organization has no database", func(t *testing.T) {
		client := newDatabasesListTestClient(t, organizationQuery, http.StatusOK, `{"databases": []}`)

		databases, err := client.DatabasesListByOrganization(t.Context(), organizationID)
		require.NoError(t, err)
		assert.Equal(t, []DatabaseNG{}, databases)
	})

	t.Run("it should return an error when the API refuses the organization", func(t *testing.T) {
		client := newDatabasesListTestClient(t, organizationQuery, http.StatusForbidden, `{"error": "You are not allowed to access this organization"}`)

		_, err := client.DatabasesListByOrganization(t.Context(), organizationID)
		require.ErrorContains(t, err, "list databases")
	})

	t.Run("it should hydrate the databases of the organization", func(t *testing.T) {
		const (
			databaseID  = "db-app-id"
			appName     = "acme-postgresql"
			projectID   = "pr-acme-uuid"
			projectName = "acme-project"
			addonID     = "ad-postgresql-id"
			addonToken  = "addon-token"
		)

		apiServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var body string
			switch r.Method + " " + r.URL.Path {
			case http.MethodGet + " /v1/databases":
				assert.Equal(t, organizationQuery, r.URL.RawQuery)
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
	})
}

// newDatabasesListTestClient returns a preview client whose API only answers once to GET /v1/databases with the expected raw query.
func newDatabasesListTestClient(t *testing.T, expectedQuery string, responseStatus int, responseBody string) *PreviewClient {
	t.Helper()

	callCount := 0
	apiServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/v1/databases", r.URL.Path)
		assert.Equal(t, expectedQuery, r.URL.RawQuery)
		w.WriteHeader(responseStatus)
		_, err := w.Write([]byte(responseBody))
		assert.NoError(t, err)
	}))
	t.Cleanup(func() {
		apiServer.Close()
		assert.Equal(t, 1, callCount)
	})

	client, err := New(t.Context(), ClientConfig{
		APIEndpoint:          apiServer.URL,
		StaticTokenGenerator: NewStaticTokenGenerator("static-token"),
	})
	require.NoError(t, err)

	return client.Preview()
}
