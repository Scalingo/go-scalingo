package scalingo

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	t.Run("static token generator should be used if present", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			require.NotEmpty(t, auth)
			split := strings.Split(auth, " ")
			require.Len(t, split, 2)
			assert.Equal(t, "static-token", split[1])
			w.WriteHeader(200)
			w.Write([]byte(`{"apps": []}`))
		}))
		defer server.Close()

		client := NewClient(ClientConfig{
			APIEndpoint:          server.URL,
			StaticTokenGenerator: NewStaticTokenGenerator("static-token"),
		})

		_, err := client.AppsList()
		require.NoError(t, err)
	})

	t.Run("it should exchange the API token for a JWT", func(t *testing.T) {
		authserver := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.Contains(r.URL.Path, "exchange") {
				_, password, ok := r.BasicAuth()
				require.True(t, ok)
				assert.Equal(t, "api-token", password)
				w.WriteHeader(200)
				w.Write([]byte(`{"token": "jwt-token"}`))
			}
			if strings.Contains(r.URL.Path, "self") {
				auth := r.Header.Get("Authorization")
				require.NotEmpty(t, auth)
				split := strings.Split(auth, " ")
				require.Len(t, split, 2)
				assert.Equal(t, "jwt-token", split[1])
				w.WriteHeader(200)
				w.Write([]byte(`{"user": {}}`))
			}
		}))
		defer authserver.Close()

		client := NewClient(ClientConfig{
			AuthEndpoint: authserver.URL,
			APIToken:     "api-token",
		})

		_, err := client.Self()
		require.NoError(t, err)
	})
}
