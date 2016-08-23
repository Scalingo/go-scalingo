package scalingo

import "time"

type Container struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	DeletedAt time.Time `json:"deleted_at"`
	Command   string    `json:"command"`
	Type      string    `json:"type"`
	TypeIndex int       `json:"type_index"`
	State     string    `json:"state"`
	Size      string    `json:"size"`
	App       *App      `json:"app"`
}
