package scalingo

type BackupsService interface {
	BackupList(app string, addonID string) error
}

func (c *Client) BackupList(app string, addonID string) error {
	return nil
}
