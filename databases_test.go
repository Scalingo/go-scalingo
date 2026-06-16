package scalingo

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestClient_DatabasePostgreSQLResetStats(t *testing.T) {
	ctx := t.Context()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/apps/my-app/addons/addon-id/token":
			assert.Equal(t, http.MethodPost, r.Method)
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(`{"addon":{"token":"addon-token"}}`))
			assert.NoError(t, err)
		case "/api/databases/addon-id/postgresql_reset_stats":
			assert.Equal(t, http.MethodPost, r.Method)
			assert.Equal(t, "Bearer addon-token", r.Header.Get("Authorization"))
			w.WriteHeader(http.StatusOK)
		default:
			assert.Failf(t, "unexpected request", "%s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer ts.Close()

	client, err := New(ctx, ClientConfig{
		APIEndpoint:         ts.URL,
		DatabaseAPIEndpoint: ts.URL,
		APIToken:            "test",
	})
	require.NoError(t, err)
	client.authClient = MockAuth(ctrl)

	err = client.DatabasePostgreSQLResetStats(ctx, "my-app", "addon-id")
	require.NoError(t, err)
}

func TestClient_DatabasePostgreSQLStatStatements(t *testing.T) {
	ctx := t.Context()

	runs := map[string]struct {
		call         func(*Client) ([]DatabasePostgreSQLStatStatement, error)
		expectedPath string
		responseBody any
		expectedRows []DatabasePostgreSQLStatStatement
	}{
		"enable pg_stat_statements": {
			call: func(client *Client) ([]DatabasePostgreSQLStatStatement, error) {
				return nil, client.DatabasePostgreSQLStatStatementsEnable(ctx, "my-app", "addon-id")
			},
			expectedPath: "/api/databases/addon-id/postgresql_stat_statements_enable",
			responseBody: map[string]any{"ok": 1, "result": map[string]any{}},
		},
		"reset pg_stat_statements": {
			call: func(client *Client) ([]DatabasePostgreSQLStatStatement, error) {
				return nil, client.DatabasePostgreSQLStatStatementsReset(ctx, "my-app", "addon-id")
			},
			expectedPath: "/api/databases/addon-id/postgresql_stat_statements_reset",
			responseBody: map[string]any{"ok": 1, "result": map[string]any{}},
		},
		"list pg_stat_statements": {
			call: func(client *Client) ([]DatabasePostgreSQLStatStatement, error) {
				return client.DatabasePostgreSQLStatStatementsList(ctx, "my-app", "addon-id")
			},
			expectedPath: "/api/databases/addon-id/postgresql_stat_statements_list",
			responseBody: map[string]any{
				"ok": 1,
				"result": []map[string]any{
					{
						"user_id":     123,
						"query":       "SELECT 1",
						"calls":       5,
						"rows":        10,
						"total_time":  42.5,
						"min_time":    1.5,
						"max_time":    20.0,
						"mean_time":   8.5,
						"stddev_time": 2.5,
					},
				},
			},
			expectedRows: []DatabasePostgreSQLStatStatement{
				{
					UserID:    123,
					Query:     "SELECT 1",
					Calls:     5,
					Rows:      10,
					TotalTime: 42.5,
					MinTime:   1.5,
					MaxTime:   20.0,
					MeanTime:  8.5,
					StdDev:    2.5,
				},
			},
		},
	}

	for name, run := range runs {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch r.URL.Path {
				case "/v1/apps/my-app/addons/addon-id/token":
					assert.Equal(t, http.MethodPost, r.Method)
					w.WriteHeader(http.StatusOK)
					_, err := w.Write([]byte(`{"addon":{"token":"addon-token"}}`))
					assert.NoError(t, err)
				case run.expectedPath:
					assert.Equal(t, http.MethodPost, r.Method)
					assert.Equal(t, "Bearer addon-token", r.Header.Get("Authorization"))
					w.WriteHeader(http.StatusOK)
					err := json.NewEncoder(w).Encode(run.responseBody)
					assert.NoError(t, err)
				default:
					assert.Failf(t, "unexpected request", "%s %s", r.Method, r.URL.Path)
					w.WriteHeader(http.StatusNotFound)
				}
			}))
			defer ts.Close()

			client, err := New(ctx, ClientConfig{
				APIEndpoint:         ts.URL,
				DatabaseAPIEndpoint: ts.URL,
				APIToken:            "test",
			})
			require.NoError(t, err)
			client.authClient = MockAuth(ctrl)

			rows, err := run.call(client)
			require.NoError(t, err)
			assert.Equal(t, run.expectedRows, rows)
		})
	}
}
