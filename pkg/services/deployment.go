package services

import (
	"github.com/piotrjaromin/ecs-go/pkg/aws"
	"github.com/piotrjaromin/ecs-go/pkg/aws/codedeploy"
	"github.com/piotrjaromin/ecs-go/pkg/aws/ecs"
	"time"

	"fmt"
	"strconv"
	"strings"
)

type Deployment interface {
	Deploy(clusterName, serviceName, image, codedeployApp, codedeployGroup *string) (*DeployOutput, error)
	ContinueDeployment(deploymentId *string) (*ContinueDeploymentOutput, error)
	RollbackDeployment(deploymentId *string) (*RollbackDeploymentOutput, error)
	ListDeployments(codedeployApp, codedeployGroup *string) (*ListDeploymentsOutput, error)
	RollbackLatestDeployment(codedeployApp, codedeployGroup *string) (*RollbackLatestOutput, error)
	ContinueLatestDeployment(codedeployApp, codedeployGroup *string) (*ContinueLatestOutput, error)
	WaitForState(deploymentId, state *string, waitTime int) (*WaitForStateOutput, error)
}

type DeploymentImpl struct {
	ecs        ecs.ECS
	codedeploy codedeploy.CodeDeploy
}

func (d DeploymentImpl) ContinueDeployment(deploymentId *string) (*ContinueDeploymentOutput, error) {
	_, err := d.codedeploy.ContinueDeployment(deploymentId)
	if err != nil {
		return nil, err
	}

	return &ContinueDeploymentOutput{}, nil
}

func (d DeploymentImpl) RollbackDeployment(deploymentId *string) (*RollbackDeploymentOutput, error) {
	output, err := d.codedeploy.RollbackDeployment(deploymentId)
	if err != nil {
		return nil, err
	}

	return &RollbackDeploymentOutput{
		Status:        *output.Status,
		StatusMessage: *output.StatusMessage,
	}, nil
}

func (d DeploymentImpl) ListDeployments(codedeployApp, codedeployGroup *string) (*ListDeploymentsOutput, error) {
	deploymentIDs, err := d.codedeploy.ListDeployments(codedeployApp, codedeployGroup)
	if err != nil {
		return nil, err
	}

	ids := make([]string, 0)
	for _, id := range deploymentIDs {
		ids = append(ids, *id)
	}

	return &ListDeploymentsOutput{
		DeploymentIDs: ids,
	}, nil
}

func (d DeploymentImpl) Deploy(clusterName, serviceName, image, codedeployApp, codedeployGroup *string) (*DeployOutput, error) {
	svc, err := d.ecs.GetService(clusterName, serviceName)
	if err != nil {
		return nil, err
	}

	if len(svc.LoadBalancers) == 0 {
		return nil, fmt.Errorf("Missing load balancers data in service")
	}

	taskDef, err := d.ecs.GetTaskDefinition(svc.TaskDefinition)
	if err != nil {
		return nil, err
	}

	lb := svc.LoadBalancers[0]
	containerName := lb.ContainerName
	containerPort := strconv.FormatInt(*lb.ContainerPort, 10)

	updatedTaskDef, err := d.ecs.UpdateTaskDefinitions(taskDef, image)
	if err != nil {
		return nil, err
	}

	deployment, err := d.codedeploy.CreateDeployment(codedeployApp, codedeployGroup, updatedTaskDef.TaskDefinitionArn, containerName, &containerPort)
	if err != nil {
		return nil, err
	}

	return &DeployOutput{
		DeploymentID:      *deployment,
		TaskDefinitionArn: *updatedTaskDef.TaskDefinitionArn,
	}, nil
}

func getFamilyNameFromArn(taskDefArn string) string {
	familyWithRevision := strings.Split(taskDefArn, "/")[1]
	family := strings.Split(familyWithRevision, ":")[0]
	return family
}

func (d DeploymentImpl) RollbackLatestDeployment(codedeployApp, codedeployGroup *string) (*RollbackLatestOutput, error) {
	deployments, err := d.ListDeployments(codedeployApp, codedeployGroup)
	if err != nil {
		return nil, err
	}

	if len(deployments.DeploymentIDs) == 0 {
		return nil, fmt.Errorf("No running deployments found")
	}

	deploymentID := deployments.DeploymentIDs[0]
	_, err = d.RollbackDeployment(&deploymentID)
	if err != nil {
		return nil, err
	}

	return &RollbackLatestOutput{
		DeploymentID: deploymentID,
	}, nil
}

func (d DeploymentImpl) ContinueLatestDeployment(codedeployApp, codedeployGroup *string) (*ContinueLatestOutput, error) {
	deployments, err := d.ListDeployments(codedeployApp, codedeployGroup)
	if err != nil {
		return nil, err
	}

	if len(deployments.DeploymentIDs) == 0 {
		return nil, fmt.Errorf("No running deployments found")
	}

	deploymentID := deployments.DeploymentIDs[0]
	_, err = d.ContinueDeployment(&deploymentID)
	if err != nil {
		return nil, err
	}

	return &ContinueLatestOutput{
		DeploymentID: deploymentID,
	}, nil
}

func (d DeploymentImpl) WaitForState(deploymentId, state *string, waitTime int) (*WaitForStateOutput, error) {

	startTime := time.Now()

	err := d.codedeploy.WaitForSate(deploymentId, state, waitTime)
	if err != nil {
		return nil, err
	}

	timeTaken := time.Now().Sub(startTime)

	return &WaitForStateOutput{
		State:  *state,
		Waited: fmt.Sprintf("%f", timeTaken.Seconds()),
	}, nil
}

func NewDeployment() (Deployment, error) {
	sess, err := aws.GetSession()
	if err != nil {
		return nil, err
	}

	ecsSvc := ecs.NewEcs(sess)
	codedeploySvc := codedeploy.NewCodeDeploy(sess)
	return &DeploymentImpl{
		ecs:        ecsSvc,
		codedeploy: codedeploySvc,
	}, nil
}
