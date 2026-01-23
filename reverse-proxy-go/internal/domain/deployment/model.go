package deployment

import "time"

type Status string

const (
	StatusNotStarted Status = "NOT_STARTED"
	StatusQueued     Status = "QUEUED"
	StatusInProgress Status = "IN_PROGRESS"
	StatusReady      Status = "READY"
	StatusFail       Status = "FAIL"
)

type Deployment struct {
	ID        string    `json:"id"`
	ProjectID string    `json:"project_id"`
	Status    Status    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
