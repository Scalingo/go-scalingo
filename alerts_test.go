package scalingo

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestAlertsClient(t *testing.T) {
	appName := "my-app"
	alertID := "my-id"

	tests := []struct {
		action           string
		testedClientCall func(c *AlertsClient) error
		expectMock       func(mock *MockSubresourceService) *gomock.Call
	}{
		{
			action: "list",
			testedClientCall: func(c *AlertsClient) error {
				_, err := c.AlertsList(appName)
				return err
			},
			expectMock: func(mock *MockSubresourceService) *gomock.Call {
				return mock.EXPECT().subresourceList(appName, "alerts", nil, gomock.Any())
			},
		},
		{
			action: "add",
			testedClientCall: func(c *AlertsClient) error {
				_, err := c.AlertAdd(appName, AlertAddParams{})
				return err
			},
			expectMock: func(mock *MockSubresourceService) *gomock.Call {
				return mock.EXPECT().subresourceAdd(appName, "alerts", gomock.Any(), gomock.Any())
			},
		},
		{
			action: "show",
			testedClientCall: func(c *AlertsClient) error {
				_, err := c.AlertShow(appName, alertID)
				return err
			},
			expectMock: func(mock *MockSubresourceService) *gomock.Call {
				return mock.EXPECT().subresourceGet(appName, "alerts", alertID, nil, gomock.Any())
			},
		},
		{
			action: "update",
			testedClientCall: func(c *AlertsClient) error {
				_, err := c.AlertUpdate(appName, alertID, AlertUpdateParams{})
				return err
			},
			expectMock: func(mock *MockSubresourceService) *gomock.Call {
				return mock.EXPECT().subresourceUpdate(appName, "alerts", alertID, gomock.Any(), gomock.Any())
			},
		},
		{
			action: "remove",
			testedClientCall: func(c *AlertsClient) error {
				return c.AlertRemove(appName, alertID)
			},
			expectMock: func(mock *MockSubresourceService) *gomock.Call {
				return mock.EXPECT().subresourceDelete(appName, "alerts", alertID)
			},
		},
	}

	for _, test := range tests {
		for msg, run := range map[string]struct {
			expectMockReturnCall func(call *gomock.Call)
			expectedError        string
		}{
			"it should fail if it fails to " + test.action + "the subresource": {
				expectMockReturnCall: func(call *gomock.Call) {
					call.Return(errors.New("error " + test.action))
				},
				expectedError: "error " + test.action,
			},
			"it should succeed if it succeeds to " + test.action + " the subresource": {
				expectMockReturnCall: func(call *gomock.Call) {
					call.Return(nil)
				},
			},
		} {
			t.Run(msg, func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()
				mock := NewMockSubresourceService(ctrl)
				c := &AlertsClient{
					SubresourceService: mock,
				}
				run.expectMockReturnCall(test.expectMock(mock))

				err := test.testedClientCall(c)
				if run.expectedError != "" {
					require.EqualError(t, err, run.expectedError)
				} else {
					require.NoError(t, err)
				}
			})
		}
	}
}
