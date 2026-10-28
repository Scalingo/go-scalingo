package scalingo

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOrganizationsClient_OrganizationsList(t *testing.T) {
	tests := map[string]struct {
		responseStatus int
		responseBody   string
		expected       []Organization
		expectedError  string
	}{
		"it should decode the organizations with their role": {
			responseStatus: http.StatusOK,
			responseBody: `{"organizations": [
				{"id": "org-1", "name": "Acme", "slug": "acme", "role": "admin"},
				{"id": "org-2", "name": "Globex", "slug": "globex", "role": "member"}
			]}`,
			expected: []Organization{
				{ID: "org-1", Slug: "acme", Role: "admin"},
				{ID: "org-2", Slug: "globex", Role: "member"},
			},
		},
		"it should return no organization when the user has none": {
			responseStatus: http.StatusOK,
			responseBody:   `{"organizations": []}`,
			expected:       []Organization{},
		},
		"it should return an error when the user is not authenticated": {
			responseStatus: http.StatusUnauthorized,
			responseBody:   `{"error": "unauthorized"}`,
			expectedError:  "list organizations",
		},
		"it should return an error when the API fails": {
			responseStatus: http.StatusInternalServerError,
			responseBody:   `INVALID`,
			expectedError:  "list organizations",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := t.Context()

			authServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodGet, r.Method)
				assert.Equal(t, "/v1/organizations", r.URL.Path)
				assert.Equal(t, "true", r.URL.Query().Get("with_role"))
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			defer authServer.Close()

			client, err := New(ctx, ClientConfig{
				AuthEndpoint:         authServer.URL,
				StaticTokenGenerator: NewStaticTokenGenerator("static-token"),
			})
			require.NoError(t, err)

			organizations, err := client.OrganizationsList(ctx)
			if test.expectedError != "" {
				require.ErrorContains(t, err, test.expectedError)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expected, organizations)
		})
	}
}
