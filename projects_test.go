package scalingo

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProjectsClient_List(t *testing.T) {
	const organizationID = "org-11111111-2222-3333-4444-555555555555"

	tests := map[string]struct {
		list           func(ctx context.Context, c *Client) ([]Project, error)
		responseStatus int
		responseBody   string
		expectedQuery  url.Values
		expected       []Project
		expectedError  string
	}{
		"it should not send any organization_id when listing the personal projects": {
			list: func(ctx context.Context, c *Client) ([]Project, error) {
				return c.ProjectsList(ctx)
			},
			responseStatus: http.StatusOK,
			responseBody:   `{"projects": [{"id": "pr-1", "name": "personal", "default": true, "organization_id": ""}]}`,
			expectedQuery:  url.Values{},
			expected:       []Project{{ID: "pr-1", Name: "personal", Default: true}},
		},
		"it should send the organization_id and decode the organization of the projects": {
			list: func(ctx context.Context, c *Client) ([]Project, error) {
				return c.ProjectsListByOrganization(ctx, organizationID)
			},
			responseStatus: http.StatusOK,
			responseBody: `{"projects": [
				{"id": "pr-2", "name": "acme-app", "default": false, "organization_id": "` + organizationID + `"},
				{"id": "pr-3", "name": "acme-db", "default": false, "organization_id": "` + organizationID + `"}
			]}`,
			expectedQuery: url.Values{"organization_id": {organizationID}},
			expected: []Project{
				{ID: "pr-2", Name: "acme-app", OrganizationID: organizationID},
				{ID: "pr-3", Name: "acme-db", OrganizationID: organizationID},
			},
		},
		"it should return an empty list when the organization has no project": {
			list: func(ctx context.Context, c *Client) ([]Project, error) {
				return c.ProjectsListByOrganization(ctx, organizationID)
			},
			responseStatus: http.StatusOK,
			responseBody:   `{"projects": []}`,
			expectedQuery:  url.Values{"organization_id": {organizationID}},
			expected:       []Project{},
		},
		"it should return an error when the API refuses the organization": {
			list: func(ctx context.Context, c *Client) ([]Project, error) {
				return c.ProjectsListByOrganization(ctx, organizationID)
			},
			responseStatus: http.StatusNotFound,
			responseBody:   `{"error": "not found"}`,
			expectedQuery:  url.Values{"organization_id": {organizationID}},
			expectedError:  "list projects",
		},
		"it should return an error when the API fails": {
			list: func(ctx context.Context, c *Client) ([]Project, error) {
				return c.ProjectsList(ctx)
			},
			responseStatus: http.StatusInternalServerError,
			responseBody:   `INVALID`,
			expectedQuery:  url.Values{},
			expectedError:  "list projects",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			apiServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodGet, r.Method)
				assert.Equal(t, "/v1/projects", r.URL.Path)
				assert.Equal(t, test.expectedQuery, r.URL.Query())
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

			projects, err := test.list(t.Context(), client)
			if test.expectedError != "" {
				require.ErrorContains(t, err, test.expectedError)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expected, projects)
		})
	}
}

func TestProjectsClient_Add(t *testing.T) {
	const organizationID = "org-11111111-2222-4333-8444-555555555555"

	tests := map[string]struct {
		params       ProjectAddParams
		expectedBody string
		responseBody string
		expected     Project
	}{
		"it should send the organization_id when creating a project in an organization": {
			params:       ProjectAddParams{Name: "my-project", Default: false, OrganizationID: organizationID},
			expectedBody: `{"project":{"name":"my-project","default":false,"organization_id":"` + organizationID + `"}}`,
			responseBody: `{"project": {"id": "prj-1", "name": "my-project", "default": false, "organization_id": "` + organizationID + `"}}`,
			expected:     Project{ID: "prj-1", Name: "my-project", OrganizationID: organizationID},
		},
		"it should not send any organization_id when creating a personal project": {
			params:       ProjectAddParams{Name: "my-project", Default: true},
			expectedBody: `{"project":{"name":"my-project","default":true}}`,
			responseBody: `{"project": {"id": "prj-2", "name": "my-project", "default": true, "organization_id": ""}}`,
			expected:     Project{ID: "prj-2", Name: "my-project", Default: true},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			apiServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodPost, r.Method)
				assert.Equal(t, "/v1/projects", r.URL.Path)
				assert.Empty(t, r.URL.Query())

				body, err := io.ReadAll(r.Body)
				if !assert.NoError(t, err) {
					return
				}
				assert.JSONEq(t, test.expectedBody, string(body))

				w.WriteHeader(http.StatusCreated)
				_, err = w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			defer apiServer.Close()

			client, err := New(t.Context(), ClientConfig{
				APIEndpoint:          apiServer.URL,
				StaticTokenGenerator: NewStaticTokenGenerator("static-token"),
			})
			require.NoError(t, err)

			project, err := client.ProjectAdd(t.Context(), test.params)
			require.NoError(t, err)
			assert.Equal(t, test.expected, project)
		})
	}
}
