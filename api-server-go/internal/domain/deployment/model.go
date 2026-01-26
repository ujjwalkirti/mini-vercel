package deployment

import "time"

type Status string

const (
	NotStarted Status = "NOT_STARTED"
	Queued     Status = "QUEUED"
	InProgress Status = "IN_PROGRESS"
	Ready      Status = "READY"
	Fail       Status = "FAIL"
)

type Deployment struct {
	ID        string    `json:"id"`
	ProjectID string    `json:"projectId"`
	Status    Status    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
