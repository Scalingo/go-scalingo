package scalingo

import (
	"encoding/json"
	"reflect"
	"testing"
)

var eventsSpecializeCases = []struct {
	Event *Event

	// Expected
	DetailedEventName   string
	DetailedEventString string
}{{
	Event: &Event{
		User: EventUser{
			Username: "user1",
		},
		Type:        EventRestart,
		RawTypeData: json.RawMessage([]byte(`{"scope": "web"}`)),
	},
	DetailedEventName:   "EventRestartType",
	DetailedEventString: "user1 restarted 'web' containers",
}}

func TestEvent_Specialize(t *testing.T) {
	for _, c := range eventsSpecializeCases {
		dev := c.Event.Specialize()

		tdev := reflect.ValueOf(dev).Elem().Type().Name()
		if tdev != c.DetailedEventName {
			t.Errorf("Expecting event of type %v, got %v", c.DetailedEventName, tdev)
		}

		if dev.String() != c.DetailedEventString {
			t.Errorf("Expecting event string\n===\n%s\n=== got\n%s\n===", c.DetailedEventString, dev.String())
		}
	}
}
