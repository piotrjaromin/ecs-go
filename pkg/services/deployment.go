package services

import (
	"github.com/piotrjaromin/ecs-go/pkg/aws"
	"github.com/piotrjaromin/ecs-go/pkg/aws/codedeploy"
	"github.com/piotrjaromin/ecs-go/pkg/aws/ecs"

	"fmt"
	"strconv"
	"strings"
)

type Deployment interface {
	Deploy(clusterName, serviceName, image, codedeployApp, codedeployGroup *string) (*DeployOutput, error)
	ContinueDeployment(deploymentId *string) error
	RollbackDeployment(deploymentId *string) error
	ListDeployments(codedeployApp, codedeployGroup *string) (*ListDeploymentsOutput, error)
}

type DeploymentImpl struct {
	ecs        ecs.ECS
	codedeploy codedeploy.CodeDeploy
}

func (d DeploymentImpl) ContinueDeployment(deploymentId *string) error {
	return d.codedeploy.ContinueDeployment(deploymentId)
}

func (d DeploymentImpl) RollbackDeployment(deploymentId *string) error {
	return d.codedeploy.RollbackDeployment(deploymentId)
}

func (d DeploymentImpl) ListDeployments(codedeployApp, codedeployGroup *string) (*ListDeploymentsOutput, error) {
	deploymentIDs, err := d.codedeploy.ListDeployments(codedeployApp, codedeployGroup)
	if err != nil {
		return nil, err
	}

	ids := make([]string, len(deploymentIDs))
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
