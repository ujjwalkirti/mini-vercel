package buildlog

type Event struct {
	ProjectID    string `json:"project_id"`
	DeploymentID string `json:"deployment_id"`
	Log          string `json:"log"`
}
