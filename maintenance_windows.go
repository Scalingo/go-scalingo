package scalingo

import (
	"context"

	"gopkg.in/errgo.v1"
)

type MaintenanceWindow struct {
	WeekdayUTC      int `json:"weekday_utc"`
	StartingHourUTC int `json:"starting_hour_utc"`
	DurationInHour  int `json:"duration_in_hour"`
}

type MaintenanceWindowParams struct {
	WeekdayUTC      *int `json:"weekday_utc,omitempty"`
	StartingHourUTC *int `json:"starting_hour_utc,omitempty"`
}

func (c *Client) DatabaseUpdateMaintenanceWindow(ctx context.Context, app, addonID string, params MaintenanceWindowParams) (Database, error) {
	var dbRes DatabaseRes
	err := c.DBAPI(app, addonID).ResourceUpdate(ctx, "databases", addonID, map[string]interface{}{
		"database": map[string]interface{}{
			"maintenance_window": params,
		},
	}, &dbRes)

	if err != nil {
		return Database{}, errgo.Notef(err, "update database '%v' maintenance window", addonID)
	}
	return dbRes.Database, nil
}
