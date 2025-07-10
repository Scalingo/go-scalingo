package scalingo

import (
	"encoding/json"
	"fmt"
	"testing"
)

var eventsSpecializeCases = map[string]struct {
	Event *Event

	// Expected
	DetailedEventName   string
	DetailedEventString string
}{
	"test event specialization": {
		Event: &Event{
			User: EventUser{
				Username: "user1",
			},
			Type:        EventRestart,
			RawTypeData: json.RawMessage([]byte(`{"scope": ["web"]}`)),
		},
		DetailedEventName:   "*scalingo.EventRestartType",
		DetailedEventString: "containers [web] began to restart",
	},
	"test edit app event without with null force_https": {
		Event: &Event{
			Type:        EventEditApp,
			RawTypeData: json.RawMessage([]byte(`{"force_https": null}`)),
		},
		DetailedEventName:   "*scalingo.EventEditAppType",
		DetailedEventString: "application settings have been updated",
	},
	"test edit app event, force https enabled": {
		Event: &Event{
			Type:        EventEditApp,
			RawTypeData: json.RawMessage([]byte(`{"force_https": true}`)),
		},
		DetailedEventName:   "*scalingo.EventEditAppType",
		DetailedEventString: "application settings have been updated, Force HTTPS has been enabled",
	},
	"test edit app event, force https disabled": {
		Event: &Event{
			Type:        EventEditApp,
			RawTypeData: json.RawMessage([]byte(`{"force_https": false}`)),
		},
		DetailedEventName:   "*scalingo.EventEditAppType",
		DetailedEventString: "application settings have been updated, Force HTTPS has been disabled",
	},
	"test app run event for a command run by an operator": {
		Event: &Event{
			Type:        EventRun,
			RawTypeData: json.RawMessage([]byte(`{"command": ""}`)),
		},
		DetailedEventName:   "*scalingo.EventRunType",
		DetailedEventString: "one-off container for maintenance/support purposes",
	},
	"test app run event for a command run by a collaborator": {
		Event: &Event{
			Type:        EventRun,
			RawTypeData: json.RawMessage([]byte(`{"command": "bash"}`)),
		},
		DetailedEventName:   "*scalingo.EventRunType",
		DetailedEventString: "one-off container with command 'bash'",
	},
	"test collaborator role change to limited collaborator": {
		Event: &Event{
			User: EventUser{
				Username: "user1",
			},
			Type:        EventChangeCollaboratorRole,
			RawTypeData: json.RawMessage([]byte(`{"collaborator": {"username": "user1","email":"user1@scalingo.com","is_limited":true}}`)),
		},
		DetailedEventName:   "*scalingo.EventChangeCollaboratorRoleType",
		DetailedEventString: "'user1' (user1@scalingo.com) is now a Limited collaborator",
	},
	"test collaborator role change to collaborator": {
		Event: &Event{
			User: EventUser{
				Username: "user1",
			},
			Type:        EventChangeCollaboratorRole,
			RawTypeData: json.RawMessage([]byte(`{"collaborator": {"username": "user1","email":"user1@scalingo.com","is_limited":false}}`)),
		},
		DetailedEventName:   "*scalingo.EventChangeCollaboratorRoleType",
		DetailedEventString: "'user1' (user1@scalingo.com) is now a Collaborator",
	},
	"test collaborator role change on pending collaborator": {
		Event: &Event{
			User:        EventUser{},
			Type:        EventChangeCollaboratorRole,
			RawTypeData: json.RawMessage([]byte(`{"collaborator": {"email":"user1@scalingo.com","is_limited":false}}`)),
		},
		DetailedEventName:   "*scalingo.EventChangeCollaboratorRoleType",
		DetailedEventString: "user1@scalingo.com is now a Collaborator",
	},
	"test update app project event": {
		Event: &Event{
			User: EventUser{
				Username: "user1",
			},
			Type:        EventUpdateAppProject,
			RawTypeData: json.RawMessage([]byte(`{"old_project_name": "old-project", "new_project_name": "new-project", "old_project_id": "proj-123", "new_project_id": "proj-456"}`)),
		},
		DetailedEventName:   "*scalingo.EventUpdateAppProjectType",
		DetailedEventString: "the application has been moved from project 'old-project' (proj-123) to 'new-project' (proj-456)",
	},
	"test edit project values": {
		Event: &Event{
			User:        EventUser{},
			Type:        EventEditProject,
			RawTypeData: json.RawMessage([]byte(`{"updated_values": [{"name":"name","value":"new-project-name","old_value":"initial-name"},{"name":"default","value":"true","old_value":"false"}]}`)),
		},
		DetailedEventName:   "*scalingo.EventEditProjectType",
		DetailedEventString: "project settings have been updated: name modified from 'initial-name' to 'new-project-name', default modified from 'false' to 'true'",
	},
}

func TestEvent_Specialize(t *testing.T) {
	for msg, c := range eventsSpecializeCases {
		t.Run(msg, func(t *testing.T) {
			dev := c.Event.Specialize()

			tdev := fmt.Sprintf("%T", dev)
			if tdev != c.DetailedEventName {
				t.Errorf("Expecting event of type %v, got %v", c.DetailedEventName, tdev)
			}

			if dev.String() != c.DetailedEventString {
				t.Errorf("Expecting event string\n===\n%s\n=== got\n%s\n===", c.DetailedEventString, dev.String())
			}
		})
	}
}
