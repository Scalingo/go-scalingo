package scalingo

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestPreviewClient_DatabaseNetPeerings(t *testing.T) {
	ctx := t.Context()

	const databaseID = "db-id"
	const netPeeringID = "np-id"

	netPeering := DatabaseNetPeering{
		ID:                       netPeeringID,
		DatabaseID:               databaseID,
		Status:                   "active",
		OutscaleNetPeeringID:     "vpcx-1234",
		OutscaleSourceNetID:      "vpc-1234",
		OutscaleSourceNetIPRange: "10.0.0.0/16",
		OutscaleSourceAccountID:  "123456789012",
	}

	runs := map[string]struct {
		testedClientCall func(c *PreviewClient) error
		expectedMethod   string
		expectedPath     string
		expectedBody     string
		responseStatus   int
		responseBody     any
	}{
		"create database net peering": {
			testedClientCall: func(c *PreviewClient) error {
				res, err := c.DatabaseNetPeeringCreate(ctx, databaseID, DatabaseNetPeeringCreateParams{
					OutscaleNetPeeringID: "vpcx-1234",
				})
				if err != nil {
					return err
				}
				assert.Equal(t, netPeeringID, res.ID)
				return nil
			},
			expectedMethod: "POST",
			expectedPath:   "/v1/databases/db-id/net_peerings",
			expectedBody:   `{"outscale_net_peering_id":"vpcx-1234"}`,
			responseStatus: 201,
			responseBody: DatabaseNetPeeringResponse{
				NetPeering: netPeering,
			},
		},
		"list database net peerings": {
			testedClientCall: func(c *PreviewClient) error {
				res, err := c.DatabaseNetPeeringsList(ctx, databaseID)
				if err != nil {
					return err
				}
				require.Len(t, res, 1)
				assert.Equal(t, netPeeringID, res[0].ID)
				return nil
			},
			expectedMethod: "GET",
			expectedPath:   "/v1/databases/db-id/net_peerings",
			responseStatus: 200,
			responseBody: DatabaseNetPeeringsResponse{
				NetPeerings: []DatabaseNetPeering{netPeering},
			},
		},
		"show database net peering": {
			testedClientCall: func(c *PreviewClient) error {
				res, err := c.DatabaseNetPeeringShow(ctx, databaseID, netPeeringID)
				if err != nil {
					return err
				}
				assert.Equal(t, netPeeringID, res.ID)
				return nil
			},
			expectedMethod: "GET",
			expectedPath:   "/v1/databases/db-id/net_peerings/np-id",
			responseStatus: 200,
			responseBody: DatabaseNetPeeringResponse{
				NetPeering: netPeering,
			},
		},
		"destroy database net peering": {
			testedClientCall: func(c *PreviewClient) error {
				return c.DatabaseNetPeeringDestroy(ctx, databaseID, netPeeringID)
			},
			expectedMethod: "DELETE",
			expectedPath:   "/v1/databases/db-id/net_peerings/np-id",
			responseStatus: 204,
		},
	}

	for msg, run := range runs {
		t.Run(msg, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, run.expectedMethod, r.Method)
				assert.Equal(t, run.expectedPath, r.URL.Path)

				if run.expectedBody != "" {
					buf := new(bytes.Buffer)
					_, err := buf.ReadFrom(r.Body)
					if !assert.NoError(t, err) {
						return
					}
					assert.Equal(t, run.expectedBody, buf.String())
				}

				w.WriteHeader(run.responseStatus)
				if run.responseBody != nil {
					err := json.NewEncoder(w).Encode(run.responseBody)
					assert.NoError(t, err)
				}
			}))
			defer ts.Close()

			client, err := New(ctx, ClientConfig{
				APIEndpoint: ts.URL,
				APIToken:    "test",
			})
			require.NoError(t, err)
			client.authClient = MockAuth(ctrl)

			err = run.testedClientCall(client.Preview())
			require.NoError(t, err)
		})
	}
}

func TestPreviewClient_DatabaseNetworkConfigurationShow(t *testing.T) {
	ctx := t.Context()

	const databaseID = "db-id"

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/v1/databases/db-id/network_configuration", r.URL.Path)
		w.WriteHeader(200)
		err := json.NewEncoder(w).Encode(DatabaseNetworkConfigurationResponse{
			NetworkConfiguration: DatabaseNetworkConfiguration{
				OutscaleAccountID: "oac-123",
				OutscaleNetID:     "vpc-123",
				IPRange:           "10.0.0.0/24",
			},
		})
		assert.NoError(t, err)
	}))
	defer ts.Close()

	client, err := New(ctx, ClientConfig{
		APIEndpoint: ts.URL,
		APIToken:    "test",
	})
	require.NoError(t, err)
	client.authClient = MockAuth(ctrl)

	res, err := client.Preview().DatabaseNetworkConfigurationShow(ctx, databaseID)
	require.NoError(t, err)
	assert.Equal(t, "oac-123", res.OutscaleAccountID)
	assert.Equal(t, "vpc-123", res.OutscaleNetID)
	assert.Equal(t, "10.0.0.0/24", res.IPRange)
}
