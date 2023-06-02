package scalingo

import (
	"context"
	"time"

	"gopkg.in/errgo.v1"
)

type Maintenance struct {
	DatabaseID string            `json:"database_id"`
	Status     MaintenanceStatus `json:"status"`
	Type       MaintenanceType   `json:"type"`
	StartedAt  *time.Time        `json:"started_at,omitempty"`
	EndedAt    *time.Time        `json:"ended_at,omitempty"`
}

type MaintenanceStatus string

const (
	MaintenanceStatusScheduled MaintenanceStatus = "scheduled"
	MaintenanceStatusNotified  MaintenanceStatus = "notified"
	MaintenanceStatusQueued    MaintenanceStatus = "queued"
	MaintenanceStatusCancelled MaintenanceStatus = "cancelled"
	MaintenanceStatusRunning   MaintenanceStatus = "running"
	MaintenanceStatusFailed    MaintenanceStatus = "failed"
	MaintenanceStatusDone      MaintenanceStatus = "done"
)

type MaintenanceType string

const (
	MaintenanceTypeNoOp    MaintenanceType = "no-op"
	MaintenanceTypeFailing MaintenanceType = "failing"
)

// MaintenanceListRes is the returned response from ListDatabaseMaintenance
type MaintenanceListRes struct {
	Maintenance []Maintenance  `json:"maintenance"`
	Meta        PaginationMeta `json:"meta"`
}

func (c *Client) DatabaseListMaintenance(ctx context.Context, app, addonID string, opts PaginationOpts) (MaintenanceListRes, error) {
	var maintenanceRes MaintenanceListRes
	err := c.DBAPI(app, addonID).SubresourceList(ctx, "databases", addonID, "maintenance", opts.ToMap(), &maintenanceRes)
	if err != nil {
		return MaintenanceListRes{}, errgo.Notef(err, "list database '%v' maintenance", addonID)
	}
	return maintenanceRes, nil
}