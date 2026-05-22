package scalingo

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
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
	"test delete project": {
		Event: &Event{
			User:        EventUser{},
			Type:        EventDeleteProject,
			ProjectName: "proj-123",
			RawTypeData: json.RawMessage([]byte(`{"project_id": "p-123456"}`)),
		},
		DetailedEventName:   "*scalingo.EventDeleteProjectType",
		DetailedEventString: "The project 'proj-123' has been deleted",
	},
	"test backup succeeded with automated method": {
		Event: &Event{
			Type:        EventDatabaseBackupSucceeded,
			RawTypeData: json.RawMessage([]byte(`{"addon_uuid":"addon-123","addon_name":"my-addon","resource_id":"db-123","backup_id":"backup-123","backup_method":"automated"}`)),
		},
		DetailedEventName:   "*scalingo.EventDatabaseBackupSucceededType",
		DetailedEventString: "Backup backup-123 for addon 'my-addon' (db-123) succeeded",
	},
	"test backup succeeded with manual method": {
		Event: &Event{
			Type:        EventDatabaseBackupSucceeded,
			RawTypeData: json.RawMessage([]byte(`{"addon_uuid":"addon-123","addon_name":"my-addon","resource_id":"db-123","backup_id":"backup-123","backup_method":"manual"}`)),
		},
		DetailedEventName:   "*scalingo.EventDatabaseBackupSucceededType",
		DetailedEventString: "Manual backup backup-123 for addon 'my-addon' (db-123) succeeded",
	},
	"test backup succeeded without method": {
		Event: &Event{
			Type:        EventDatabaseBackupSucceeded,
			RawTypeData: json.RawMessage([]byte(`{"addon_uuid":"addon-123","addon_name":"my-addon","resource_id":"db-123","backup_id":"backup-123"}`)),
		},
		DetailedEventName:   "*scalingo.EventDatabaseBackupSucceededType",
		DetailedEventString: "Backup backup-123 for addon 'my-addon' (db-123) succeeded",
	},
	"test backup failed with periodic method": {
		Event: &Event{
			Type:        EventDatabaseBackupFailed,
			RawTypeData: json.RawMessage([]byte(`{"addon_uuid":"addon-123","addon_name":"my-addon","resource_id":"db-123","backup_id":"backup-123","backup_method":"periodic"}`)),
		},
		DetailedEventName:   "*scalingo.EventDatabaseBackupFailedType",
		DetailedEventString: "Periodic backup backup-123 for addon 'my-addon' (db-123) failed",
	},
	"test backup failed with manual method": {
		Event: &Event{
			Type:        EventDatabaseBackupFailed,
			RawTypeData: json.RawMessage([]byte(`{"addon_uuid":"addon-123","addon_name":"my-addon","resource_id":"db-123","backup_id":"backup-123","backup_method":"manual"}`)),
		},
		DetailedEventName:   "*scalingo.EventDatabaseBackupFailedType",
		DetailedEventString: "Manual backup backup-123 for addon 'my-addon' (db-123) failed",
	},
	"test continuous backup healthy": {
		Event: &Event{
			Type:        EventDatabaseContinuousBackupHealthy,
			RawTypeData: json.RawMessage([]byte(`{"addon_name":"PostgreSQL","resource_id":"db-123","addon_uuid":"ad-b48f347e-2e3d-4b48-8dee-40eed0539ee1","status":"healthy","error":null,"recoverable":true,"checked_at":"2026-05-28T08:59:00.011Z","unrecoverable_duration_seconds":0}`)),
		},
		DetailedEventName:   "*scalingo.EventDatabaseContinuousBackupHealthyType",
		DetailedEventString: "Point-in-time recovery for database 'db-123' is healthy",
	},
	"test continuous backup delayed": {
		Event: &Event{
			Type:        EventDatabaseContinuousBackupDelayed,
			RawTypeData: json.RawMessage([]byte(`{"addon_name":"PostgreSQL","resource_id":"db-123","addon_uuid":"ad-b48f347e-2e3d-4b48-8dee-40eed0539ee1","status":"pgbackrest_error","error":"latest restorable point is more than 12 hours behind","recoverable":true,"checked_at":"2026-05-28T08:56:00.009Z","unrecoverable_duration_seconds":180}`)),
		},
		DetailedEventName:   "*scalingo.EventDatabaseContinuousBackupDelayedType",
		DetailedEventString: "Point-in-time recovery for database 'db-123' is delayed (status: pgBackRest error, error: latest restorable point is more than 12 hours behind)",
	},
	"test continuous backup stale": {
		Event: &Event{
			Type:        EventDatabaseContinuousBackupStale,
			RawTypeData: json.RawMessage([]byte(`{"addon_name":"PostgreSQL","resource_id":"db-123","addon_uuid":"ad-b48f347e-2e3d-4b48-8dee-40eed0539ee1","status":"pgbackrest_error","error":"repository metadata is unavailable","recoverable":false,"checked_at":"2026-05-28T08:57:00.009Z","unrecoverable_duration_seconds":240}`)),
		},
		DetailedEventName:   "*scalingo.EventDatabaseContinuousBackupStaleType",
		DetailedEventString: "Point-in-time recovery for database 'db-123' is stale (status: pgBackRest error, error: repository metadata is unavailable)",
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

func TestEventDatabaseContinuousBackupWho(t *testing.T) {
	testCases := map[string]struct {
		eventType   EventTypeName
		rawTypeData string
		user        EventUser
		expectedWho string
	}{
		"healthy falls back to event user": {
			eventType:   EventDatabaseContinuousBackupHealthy,
			rawTypeData: `{"addon_name":"PostgreSQL","resource_id":"db-123","addon_uuid":"ad-b48f347e-2e3d-4b48-8dee-40eed0539ee1","status":"healthy","error":null,"recoverable":true,"checked_at":"2026-05-28T08:59:00.011Z","unrecoverable_duration_seconds":0}`,
			user: EventUser{
				Username: "pitr-monitor",
				Email:    "pitr-monitor@scalingo.test",
			},
			expectedWho: "pitr-monitor (pitr-monitor@scalingo.test)",
		},
		"delayed falls back to event user": {
			eventType:   EventDatabaseContinuousBackupDelayed,
			rawTypeData: `{"addon_name":"PostgreSQL","resource_id":"db-123","addon_uuid":"ad-b48f347e-2e3d-4b48-8dee-40eed0539ee1","status":"pgbackrest_error","error":"pgbackrest info returned no stanza: no valid backups","recoverable":false,"checked_at":"2026-05-28T08:56:00.009Z","unrecoverable_duration_seconds":180}`,
			user: EventUser{
				Username: "backup-service",
				Email:    "backup-service@scalingo.test",
			},
			expectedWho: "backup-service (backup-service@scalingo.test)",
		},
	}

	for msg, tc := range testCases {
		t.Run(msg, func(t *testing.T) {
			event := &Event{
				User:        tc.user,
				Type:        tc.eventType,
				RawTypeData: json.RawMessage([]byte(tc.rawTypeData)),
			}

			detailedEvent := event.Specialize()

			if got := detailedEvent.Who(); got != tc.expectedWho {
				t.Errorf("expected who %q, got %q", tc.expectedWho, got)
			}
		})
	}
}

func TestEventDatabaseContinuousBackupTypeData(t *testing.T) {
	event := &Event{
		Type:        EventDatabaseContinuousBackupStale,
		RawTypeData: json.RawMessage([]byte(`{"addon_name":"PostgreSQL","resource_id":"db-123","addon_uuid":"ad-b48f347e-2e3d-4b48-8dee-40eed0539ee1","status":"pgbackrest_error","error":"pgbackrest info returned no stanza: no valid backups","recoverable":false,"checked_at":"2026-05-28T08:57:00.009Z","unrecoverable_duration_seconds":240}`)),
	}

	detailedEvent, ok := event.Specialize().(*EventDatabaseContinuousBackupStaleType)
	if !ok {
		t.Errorf("expected *EventDatabaseContinuousBackupStaleType, got %T", event.Specialize())
	}

	expectedCheckedAt := time.Date(2026, time.May, 28, 8, 57, 0, 9000000, time.UTC)
	if detailedEvent.TypeData.AddonName != "PostgreSQL" {
		t.Errorf("expected addon_name to be parsed, got %q", detailedEvent.TypeData.AddonName)
	}
	if detailedEvent.TypeData.ResourceID != "db-123" {
		t.Errorf("expected resource_id to be parsed, got %q", detailedEvent.TypeData.ResourceID)
	}
	if detailedEvent.TypeData.AddonUUID != "ad-b48f347e-2e3d-4b48-8dee-40eed0539ee1" {
		t.Errorf("expected addon_uuid to be parsed, got %q", detailedEvent.TypeData.AddonUUID)
	}
	if detailedEvent.TypeData.Status != "pgbackrest_error" {
		t.Errorf("expected status to be parsed, got %q", detailedEvent.TypeData.Status)
	}
	if detailedEvent.TypeData.Error != "pgbackrest info returned no stanza: no valid backups" {
		t.Errorf("expected error to be parsed, got %q", detailedEvent.TypeData.Error)
	}
	if detailedEvent.TypeData.Recoverable {
		t.Errorf("expected recoverable to be false")
	}
	if !detailedEvent.TypeData.CheckedAt.Equal(expectedCheckedAt) {
		t.Errorf("expected checked_at to be parsed, got %s", detailedEvent.TypeData.CheckedAt)
	}
	if detailedEvent.TypeData.UnrecoverableDurationSeconds != 240 {
		t.Errorf("expected unrecoverable_duration_seconds to be parsed, got %d", detailedEvent.TypeData.UnrecoverableDurationSeconds)
	}
}
