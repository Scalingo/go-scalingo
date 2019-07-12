package scalingo

import (
	"time"

	errgo "gopkg.in/errgo.v1"
)

const (
	RegionMigrationStatusScheduled      RegionMigrationStatus = "scheduled"
	RegionMigrationStatusPreflightError RegionMigrationStatus = "preflight-error"
	RegionMigrationStatusRunning        RegionMigrationStatus = "running"
	RegionMigrationStatusError          RegionMigrationStatus = "error"
	RegionMigrationStatusDone           RegionMigrationStatus = "done"

	StepStatusRunning StepStatus = "running"
	StepStatusDone    StepStatus = "done"
	StepStatusError   StepStatus = "error"
)

type RegionMigrationService interface {
	CreateRegionMigration(appID string) (RegionMigration, error)
	ShowRegionMigration(appID, migrationID string) (RegionMigration, error)
	ListRegionMigrations(appID string) ([]RegionMigration, error)
}

type RegionMigration struct {
	ID         string
	AppName    string
	AppID      string
	NewAppID   string
	Status     RegionMigrationStatus
	StartedAt  time.Time
	FinishedAt time.Time
	Steps      Steps
}

type StepStatus string
type RegionMigrationStatus string
type Steps []Step

type Step struct {
	ID     string
	Name   string
	Status StepStatus
	Logs   string
}

func (c *Client) CreateRegionMigration(appID string) (RegionMigration, error) {
	var migration RegionMigration

	err := c.ScalingoAPI().SubresourceAdd("apps", appID, "region_migrations", nil, &migration)
	if err != nil {
		return migration, errgo.Notef(err, "fail to create migration")
	}

	return migration, nil
}

func (c *Client) ShowRegionMigration(appID, migrationID string) (RegionMigration, error) {
	var migration RegionMigration

	err := c.ScalingoAPI().SubresourceGet("apps", appID, "region_migrations", migrationID, nil, &migration)
	if err != nil {
		return migration, errgo.Notef(err, "fail to get migration")
	}

	return migration, nil
}

func (c *Client) ListRegionMigrations(appID string) ([]RegionMigration, error) {
	var migrations []RegionMigration

	err := c.ScalingoAPI().SubresourceList("apps", appID, "region_migrations", nil, &migrations)
	if err != nil {
		return migrations, errgo.Notef(err, "fail to list migrations")
	}

	return migrations, nil
}
