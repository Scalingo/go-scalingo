package scalingo

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	t.Run("static token generator should be used if present", func(t *testing.T) {
		ctx := context.Background()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			require.NotEmpty(t, auth)
			split := strings.Split(auth, " ")
			require.Len(t, split, 2)
			assert.Equal(t, "static-token", split[1])
			w.WriteHeader(200)
			_, err := w.Write([]byte(`{"apps": []}`))
			require.NoError(t, err)
		}))
		defer server.Close()

		client, err := New(ctx, ClientConfig{
			APIEndpoint:          server.URL,
			StaticTokenGenerator: NewStaticTokenGenerator("static-token"),
		})
		require.NoError(t, err)

		_, err = client.AppsList(ctx)
		require.NoError(t, err)
	})

	t.Run("it should exchange the API token for a JWT", func(t *testing.T) {
		ctx := context.Background()
		claims := &jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
		}
		jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
		jwt, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
		require.NoError(t, err)

		authserver := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.Contains(r.URL.Path, "exchange") {
				_, password, ok := r.BasicAuth()
				require.True(t, ok)
				assert.Equal(t, "api-token", password)
				w.WriteHeader(200)
				_, err := w.Write([]byte(fmt.Sprintf(`{"token": "%v"}`, jwt)))
				require.NoError(t, err)
			}
			if strings.Contains(r.URL.Path, "self") {
				auth := r.Header.Get("Authorization")
				require.NotEmpty(t, auth)
				split := strings.Split(auth, " ")
				require.Len(t, split, 2)
				assert.Equal(t, jwt, split[1])
				w.WriteHeader(200)
				_, err := w.Write([]byte(`{"user": {}}`))
				require.NoError(t, err)
			}
			assert.Equal(t, "bar", r.Header.Get("Foo"))
		}))
		defer authserver.Close()

		client, err := New(ctx, ClientConfig{
			AuthEndpoint: authserver.URL,
			APIToken:     "api-token",
			ExtraHeaders: http.Header{"Foo": {"bar"}},
		})
		require.NoError(t, err)

		_, err = client.Self(ctx)
		require.NoError(t, err)
	})
}
