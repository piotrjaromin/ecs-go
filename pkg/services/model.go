package services

type DeployOutput struct {
	DeploymentID      string `json:"deploymentId"`
	TaskDefinitionArn string `json:"taskDefinitionArn"`
}

type ListDeploymentsOutput struct {
	DeploymentIDs []string `json:"deploymentIds"`
}

type RollbackDeploymentOutput struct {
	Status        string `json:"status"`
	StatusMessage string `json:"statusMessage"`
}

type ContinueDeploymentOutput struct {
}
