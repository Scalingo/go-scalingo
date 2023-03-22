package scalingo

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Scalingo/go-scalingo/v6/http/httpmock"
)

func TestStackIsDeprecated(t *testing.T) {

	isDeprecatedCases := map[string]struct {
		date       time.Time
		deprecated bool
	}{
		"test isDeprecated is false when deprecation date is null":          {time.Time{}, false},
		"test isDeprecated is true when deprecation date is today's date":   {time.Now(), true},
		"test isDeprecated is true when deprecation date is in the past":    {time.Now().AddDate(0, 0, -1), true},
		"test isDeprecated is false when deprecation date is in the future": {time.Now().AddDate(0, 0, 1), false},
	}

	for message, test := range isDeprecatedCases {
		t.Run(message, func(t *testing.T) {
			stack := Stack{DeprecatedAt: DeprecationDate{test.date}}
			assert.Equal(t, test.deprecated, stack.IsDeprecated())
		})
	}
}

func TestStacksList(t *testing.T) {
	ctx := context.Background()

	stacksListCases := map[string]struct {
		json          string
		expectedStack Stack
	}{
		"test StacksList with stacks with deprecated_at being null": {
			json:          `{"stacks": [{"deprecated_at": null}]}`,
			expectedStack: Stack{DeprecatedAt: DeprecationDate{time.Time{}}},
		},
		"test StacksList with stacks with deprecated_at is empty": {
			json:          `{"stacks": [{}]}`,
			expectedStack: Stack{DeprecatedAt: DeprecationDate{time.Time{}}},
		},
		"test StacksList with stacks with deprecated_at being set": {
			json:          `{"stacks": [{"deprecated_at": "2022-08-31"}]}`,
			expectedStack: Stack{DeprecatedAt: DeprecationDate{time.Date(2022, 8, 31, 0, 0, 0, 0, time.UTC)}},
		},
	}

	for message, test := range stacksListCases {
		t.Run(message, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			client, err := New(ctx, ClientConfig{})
			require.NoError(t, err)
			apiMock := httpmock.NewMockClient(ctrl)
			client.apiClient = apiMock

			apiMock.EXPECT().DoRequest(gomock.Any(), gomock.Any(), gomock.Any()).Do(func(_ context.Context, _, res interface{}) {
				err := json.Unmarshal([]byte(test.json), &res)
				require.NoError(t, err)
			}).Return(nil)

			stacks, _ := client.StacksList(ctx)
			assert.Equal(t, test.expectedStack.DeprecatedAt, stacks[0].DeprecatedAt)
		})
	}
}
