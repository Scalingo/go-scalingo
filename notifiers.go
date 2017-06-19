package scalingo

import errgo "gopkg.in/errgo.v1"

type NotifierParam struct {
	Notifier interface{} `json:"notifier"`
}

type NotifierRes struct {
	Notifier DetailedNotifier `json:"notifier"`
}

type NotifiersRes struct {
	Notifiers []*Notifier `json:"notifiers"`
}

func (c *Client) NotifiersList(app string) (Notifiers, error) {
	var notifiersRes NotifiersRes
	err := c.subresourceList(app, "notifiers", nil, &notifiersRes)
	if err != nil {
		return nil, errgo.Mask(err)
	}
	var notifiers Notifiers
	for _, not := range notifiersRes.Notifiers {
		notifiers = append(notifiers, not.Specialize())
	}
	return notifiers, nil
}

func (c *Client) NotifierProvision(app, notifierType string, params map[string]interface{}) (NotifierRes, error) {
	var notifierRes NotifierRes

	notifier := NewNotifier(notifierType, params)
	notifierParams := &NotifierParam{Notifier: notifier}

	err := c.subresourceAdd(app, "notifiers", notifierParams, &notifierRes)
	if err != nil {
		return NotifierRes{}, errgo.Mask(err, errgo.Any)
	}
	return notifierRes, nil
}
