package scalingo

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLogsClient_LogsURL(t *testing.T) {
	ctx := t.Context()
	const appName = "my-app"

	t.Run("it parses logs URL response", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodGet, r.Method)
			assert.Equal(t, "/v1/apps/"+appName+"/logs", r.URL.Path)
			err := json.NewEncoder(w).Encode(LogsURLRes{
				LogsURL: "https://logs-osc-fr1.scalingo.com?token=test-token&timestamp=1700000000",
				App: &App{
					ID:   "app-id",
					Name: appName,
				},
			})
			assert.NoError(t, err)
		}))
		defer ts.Close()

		c, err := New(ctx, ClientConfig{
			APIEndpoint: ts.URL,
			APIToken:    "test",
		})
		require.NoError(t, err)

		c.authClient = MockAuth(ctrl)

		logsURLRes, err := c.LogsURL(ctx, appName)
		require.NoError(t, err)
		require.NotNil(t, logsURLRes)
		assert.Equal(t, "https://logs-osc-fr1.scalingo.com?token=test-token&timestamp=1700000000", logsURLRes.LogsURL)
		require.NotNil(t, logsURLRes.App)
		assert.Equal(t, appName, logsURLRes.App.Name)
	})

	t.Run("it fails if logs URL response is invalid JSON", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := w.Write([]byte("INVALID_JSON"))
			assert.NoError(t, err)
		}))
		defer ts.Close()

		c, err := New(ctx, ClientConfig{
			APIEndpoint: ts.URL,
			APIToken:    "test",
		})
		require.NoError(t, err)

		c.authClient = MockAuth(ctrl)

		_, err = c.LogsURL(ctx, appName)
		require.Error(t, err)
	})
}
