package scalingo

import "gopkg.in/errgo.v1"

type Notification struct {
	ID              string         `json:"id"`
	Type            string         `json:"type"`
	WebHookURL      string         `json:"webhook_url"`
}

type NotificationsRes struct {
	Notifications []*Notification `json:"notifications"`
}

type NotificationRes struct {
	Notification     Notification    `json:"notification"`
	Message          string   `json:"message,omitempty"`
	Variables        []string `json:"variables,omitempty"`
}

func (c *Client) NotificationsList(app string) ([]*Notification, error) {
	var notificationsRes NotificationsRes
	err := c.subresourceList(app, "notifications", nil, &notificationsRes)
	if err != nil {
		return nil, errgo.Mask(err, errgo.Any)
	}
	return notificationsRes.Notifications, nil
}

func (c *Client) NotificationProvision(app, webHookURL string) (NotificationRes, error) {
	var notificationRes NotificationRes
	err := c.subresourceAdd(app, "notifications", NotificationRes{Notification: Notification{WebHookURL: webHookURL}}, &notificationsRes)
	if err != nil {
		return NotificationRes{}, errgo.Mask(err, errgo.Any)
	}
	return notificationRes, nil
}
