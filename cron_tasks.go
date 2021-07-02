package scalingo

import (
	"gopkg.in/errgo.v1"
)

type CronTasksService interface {
	CronTasksAdd(app string, params CronTasks) (CronTasks, error)
}

var _ CronTasksService = (*Client)(nil)

type Job struct {
	Command string `json:"command"`
	Size    string `json:"size,omitempty"`
}

type CronTasks struct {
	Jobs []Job `json:"jobs"`
}

func (c *Client) CronTasksAdd(app string, params CronTasks) (CronTasks, error) {
	resp := CronTasks{}
	err := c.ScalingoAPI().SubresourceAdd("apps", app, "cron_tasks", params, &resp)
	if err != nil {
		return CronTasks{}, errgo.Notef(err, "fail to add cron tasks")
	}
	return resp, nil
}
