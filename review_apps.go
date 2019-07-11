package scalingo

import (
	"time"
)

type ReviewApp struct {
	ID                  string      `json:"id"`
	RepoLinkID          string      `json:"repo_link_id"`
	AppID               string      `json:"app_id"`
	AppName             string      `json:"app_name"`
	ParentAppID         string      `json:"parent_app_id"`
	ParentAppName       string      `json:"parent_app_name"`
	CreatedAt           time.Time   `json:"created_at"`
	StaleDeletionDate   time.Time   `json:"stale_deletion_date"`
	OnCloseDeletionDate time.Time   `json:"on_close_deletion_date"`
	LastDeployment      *Deployment `json:"last_deployment"`
}
