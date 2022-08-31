package scalingo

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/Scalingo/go-scalingo/v5/http/httpmock"
)

var eventsListCases = map[string]struct {
	PaginationOpts PaginationOpts
	App            string
	Body           string
	Code           int
	EventsCount    int
	Error          error
}{
	"test get list of events with no event": {
		App:         "app-1",
		Body:        `{"events": [], "meta": {"pagination": {"prev_page": 1, "current_page": 1, "next_page": 1, "total_pages": 1, "total_count": 0}}}`,
		EventsCount: 0,
	},
	"test get list of events with 1 event": {
		App:         "app-1",
		EventsCount: 1,
		Body: `{"events": [{
		"type": "run",
		"type_data": {
			"command": "bundle exec rails console"
		}
	}], "meta": {"pagination": {"prev_page": 1, "current_page": 1, "next_page": 1, "total_pages": 1, "total_count": 0}}}`,
	},
}

func TestEventsList(t *testing.T) {
	ctx := context.Background()

	for msg, c := range eventsListCases {
		t.Run(msg, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			client, err := New(ctx, ClientConfig{})
			require.NoError(t, err)
			apiMock := httpmock.NewMockClient(ctrl)
			client.apiClient = apiMock

			apiMock.EXPECT().SubresourceList(gomock.Any(), "apps", c.App, "events", c.PaginationOpts.ToMap(), gomock.Any()).Do(func(_ context.Context, _, _, _ string, _ interface{}, res interface{}) {
				err := json.Unmarshal([]byte(c.Body), &res)
				require.NoError(t, err)
			}).Return(nil)

			events, _, err := client.EventsList(ctx, c.App, c.PaginationOpts)
			if len(events) != c.EventsCount {
				t.Errorf("expected %d event, got %v", c.EventsCount, len(events))
			}
			if err != c.Error {
				t.Errorf("expected '%v' error got %v", c.Error, err)
			}
		})
	}
}
