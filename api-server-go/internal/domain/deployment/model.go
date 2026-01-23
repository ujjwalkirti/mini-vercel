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
	ID        string
	ProjectID string
	Status    Status
	CreatedAt time.Time
	UpdatedAt time.Time
}
