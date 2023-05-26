package scalingo

import (
	"fmt"
	"time"
)

type MaintenanceWindow struct {
	WeekdayUTC      int `json:"weekday_utc"`
	StartingHourUTC int `json:"starting_hour_utc"`
	DurationInHour  int `json:"duration_in_hour"`
}

func (m MaintenanceWindow) String() string {
	return m.formatMaintenanceWindowWithTimezone(time.Local)
}

func (m MaintenanceWindow) formatMaintenanceWindowWithTimezone(location *time.Location) string {
	weekdayUTC := time.Weekday(m.WeekdayUTC)
	weekdayLocal, startingHourLocal := convertUTCDayAndHourToTimezone(
		weekdayUTC,
		m.StartingHourUTC,
		location,
	)

	return fmt.Sprintf("%ss at %02d:00 (%02d hours)",
		weekdayLocal.String(),
		startingHourLocal,
		m.DurationInHour,
	)
}

func convertUTCDayAndHourToTimezone(weekdayUTC time.Weekday, hourUTC int, location *time.Location) (time.Weekday, int) {
	newTimezoneDate := beginningOfWeek(time.Now().UTC())

	newTimezoneDate = newTimezoneDate.AddDate(0, 0, int(weekdayUTC)-1)
	newTimezoneDate = newTimezoneDate.Add(time.Duration(hourUTC) * time.Hour)

	newTimezoneDate = newTimezoneDate.In(location)

	return newTimezoneDate.Weekday(), newTimezoneDate.Hour()
}

func beginningOfWeek(t time.Time) time.Time {
	t = beginningOfDay(t)
	weekday := int(t.Weekday())

	weekStartDayInt := int(time.Monday)

	if weekday < weekStartDayInt {
		weekday = weekday + 7 - weekStartDayInt
	} else {
		weekday = weekday - weekStartDayInt
	}
	return t.AddDate(0, 0, -weekday)
}

func beginningOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}
