package scalingo

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTokensClient_TokenDelete(t *testing.T) {
	tests := map[string]struct {
		responseStatus int
		expectedError  string
	}{
		"it should delete the token": {
			responseStatus: http.StatusNoContent,
		},
		"it should return an error if the token does not exist": {
			responseStatus: http.StatusNotFound,
			expectedError:  "delete token",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := t.Context()

			authServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodDelete, r.Method)
				assert.Equal(t, "/v1/tokens/token-id", r.URL.Path)
				w.WriteHeader(test.responseStatus)
			}))
			defer authServer.Close()

			client, err := New(ctx, ClientConfig{
				AuthEndpoint:         authServer.URL,
				StaticTokenGenerator: NewStaticTokenGenerator("static-token"),
			})
			require.NoError(t, err)

			err = client.TokenDelete(ctx, "token-id")
			if test.expectedError != "" {
				require.ErrorContains(t, err, test.expectedError)
				return
			}
			require.NoError(t, err)
		})
	}
}
