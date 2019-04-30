package services

type DeployOutput struct {
	DeploymentID      string `json:"deploymentId"`
	TaskDefinitionArn string `json:"taskDefinitionArn"`
}

type ListDeploymentsOutput struct {
	DeploymentIDs []string `json:"deploymentIds"`
}

type GenericOutput struct {
	Status        string `json:"status"`
	StatusMessage string `json:"statusMessage"`
}

type ContinueDeploymentOutput struct {
}

type ContinueLatestOutput struct {
	DeploymentID string `json:"deploymentId"`
}

type RollbackLatestOutput struct {
	DeploymentID string `json:"deploymentId"`
}

type WaitForStateOutput struct {
	Waited string `json:"waited"`
	State  string `json:"state"`
}
