package services

import (
	"time"

	awssdk "github.com/aws/aws-sdk-go/aws"

	"github.com/piotrjaromin/ecs-go/pkg/aws"
	"github.com/piotrjaromin/ecs-go/pkg/aws/codedeploy"
	"github.com/piotrjaromin/ecs-go/pkg/aws/ecr"
	"github.com/piotrjaromin/ecs-go/pkg/aws/ecs"

	"fmt"
	"strconv"
	"strings"
)

const (
	VariantBlue  = "blue"
	VariantGreen = "green"
)

type Deployment interface {
	Deploy(clusterName, serviceName, image *string, imageIndex int, codedeployApp, codedeployGroup *string) (*DeployOutput, error)
	Scale(clusterName, serviceName *string, count uint) (*GenericOutput, error)
	ContinueDeployment(deploymentId *string) (*ContinueDeploymentOutput, error)
	ForceContinueDeployment(deploymentId *string) (*ContinueDeploymentOutput, error)
	RollbackDeployment(deploymentId *string) (*GenericOutput, error)
	ListDeployments(codedeployApp, codedeployGroup *string) (*ListDeploymentsOutput, error)
	RollbackLatestDeployment(codedeployApp, codedeployGroup *string) (*RollbackLatestOutput, error)
	ContinueLatestDeployment(codedeployApp, codedeployGroup *string) (*ContinueLatestOutput, error)
	ForceContinueLatestDeployment(codedeployApp, codedeployGroup *string) (*ContinueLatestOutput, error)
	WaitForState(deploymentId, state *string, waitTime int) (*WaitForStateOutput, error)
	WaitForLatest(codedeployApp, codedeployGroup, state *string, waitTime int) (*WaitForStateOutput, error)
	GetLiveVariant(clusterName, serviceName *string) (*string, error)
	TagImage(repositoryName, currentTag, newTag *string) error
	ListServices() ([]*ListServicesItemOutput, error)
}

type DeploymentImpl struct {
	ecs        ecs.ECS
	ecr        ecr.ECR
	codedeploy codedeploy.CodeDeploy
}

func (d DeploymentImpl) ContinueDeployment(deploymentId *string) (*ContinueDeploymentOutput, error) {
	_, err := d.codedeploy.ContinueDeployment(deploymentId)
	if err != nil {
		return nil, err
	}

	return &ContinueDeploymentOutput{}, nil
}

func (d DeploymentImpl) ForceContinueDeployment(deploymentId *string) (*ContinueDeploymentOutput, error) {
	_, err := d.codedeploy.ForceContinueDeployment(deploymentId)
	if err != nil {
		return nil, err
	}

	return &ContinueDeploymentOutput{}, nil
}

func (d DeploymentImpl) RollbackDeployment(deploymentId *string) (*GenericOutput, error) {
	output, err := d.codedeploy.RollbackDeployment(deploymentId)
	if err != nil {
		return nil, err
	}

	return &GenericOutput{
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

func (d DeploymentImpl) Deploy(clusterName, serviceName, image *string, imageIndex int, codedeployApp, codedeployGroup *string) (*DeployOutput, error) {
	svc, err := d.ecs.GetService(clusterName, serviceName)
	if err != nil {
		return nil, err
	}

	if len(svc.LoadBalancers) == 0 {
		return nil, fmt.Errorf("Missing load balancers data in service")
	}

	taskDef, err := d.ecs.GetLatestTaskDefinition(svc.TaskDefinition)
	if err != nil {
		return nil, err
	}

	lb := svc.LoadBalancers[0]
	containerName := lb.ContainerName
	containerPort := strconv.FormatInt(*lb.ContainerPort, 10)

	updatedTaskDef, err := d.ecs.UpdateTaskDefinitions(taskDef, image, imageIndex)
	if err != nil {
		return nil, err
	}

	deployment, err := d.codedeploy.CreateDeployment(codedeployApp, codedeployGroup, updatedTaskDef.TaskDefinitionArn, containerName, &containerPort)
	if err != nil {
		return nil, err
	}

	return &DeployOutput{
		DeploymentID:            *deployment,
		TaskDefinitionArn:       *updatedTaskDef.TaskDefinitionArn,
		SourceTaskDefinitionArn: fmt.Sprintf("%s:%d", *taskDef.Family, *taskDef.Revision),
	}, nil
}

func (d DeploymentImpl) Scale(clusterName, serviceName *string, count uint) (*GenericOutput, error) {
	svc, err := d.ecs.GetService(clusterName, serviceName)
	if err != nil {
		return nil, err
	}

	svc.SetDesiredCount(int64(count))

	_, err = d.ecs.UpdateService(svc)
	if err != nil {
		return nil, err
	}

	return &GenericOutput{
		Status:        "Accepted",
		StatusMessage: "Scaling instances started",
	}, nil
}

func getFamilyNameFromArn(taskDefArn string) string {
	familyWithRevision := strings.Split(taskDefArn, "/")[1]
	family := strings.Split(familyWithRevision, ":")[0]
	return family
}

func (d DeploymentImpl) RollbackLatestDeployment(codedeployApp, codedeployGroup *string) (*RollbackLatestOutput, error) {
	deploymentID, err := d.getLatestDeployment(codedeployApp, codedeployGroup)
	if err != nil {
		return nil, err
	}

	_, err = d.RollbackDeployment(deploymentID)
	if err != nil {
		return nil, err
	}

	return &RollbackLatestOutput{
		DeploymentID: *deploymentID,
	}, nil
}

func (d DeploymentImpl) ContinueLatestDeployment(codedeployApp, codedeployGroup *string) (*ContinueLatestOutput, error) {
	deploymentID, err := d.getLatestDeployment(codedeployApp, codedeployGroup)
	if err != nil {
		return nil, err
	}

	_, err = d.ContinueDeployment(deploymentID)
	if err != nil {
		return nil, err
	}

	return &ContinueLatestOutput{
		DeploymentID: *deploymentID,
	}, nil
}

func (d DeploymentImpl) ForceContinueLatestDeployment(codedeployApp, codedeployGroup *string) (*ContinueLatestOutput, error) {
	deploymentID, err := d.getLatestDeployment(codedeployApp, codedeployGroup)
	if err != nil {
		return nil, err
	}

	_, err = d.ForceContinueDeployment(deploymentID)
	if err != nil {
		return nil, err
	}

	return &ContinueLatestOutput{
		DeploymentID: *deploymentID,
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

func (d DeploymentImpl) WaitForLatest(codedeployApp, codedeployGroup, state *string, waitTime int) (*WaitForStateOutput, error) {

	startTime := time.Now()
	deploymentID, err := d.getLatestDeployment(codedeployApp, codedeployGroup)
	if err != nil {
		return nil, err
	}

	err = d.codedeploy.WaitForSate(deploymentID, state, waitTime)
	if err != nil {
		return nil, err
	}

	timeTaken := time.Now().Sub(startTime)

	return &WaitForStateOutput{
		State:  *state,
		Waited: fmt.Sprintf("%f", timeTaken.Seconds()),
	}, nil
}

func (d DeploymentImpl) getLatestDeployment(codedeployApp, codedeployGroup *string) (*string, error) {
	deployments, err := d.ListDeployments(codedeployApp, codedeployGroup)
	if err != nil {
		return nil, err
	}

	if len(deployments.DeploymentIDs) == 0 {
		return nil, fmt.Errorf("No running deployments found")
	}

	deploymentID := deployments.DeploymentIDs[0]

	return &deploymentID, nil
}

func (d DeploymentImpl) GetLiveVariant(clusterName, serviceName *string) (*string, error) {
	svc, err := d.ecs.GetService(clusterName, serviceName)
	if err != nil {
		return nil, err
	}
	if len(svc.TaskSets) > 1 {
		return nil, fmt.Errorf("Service is during deployment")
	}
	if len(svc.TaskSets) == 0 {
		return nil, fmt.Errorf("Service has no task set")
	}
	if len(svc.LoadBalancers) == 0 {
		return nil, fmt.Errorf("Missing load balancers data in service")
	}
	tgArn := *svc.LoadBalancers[0].TargetGroupArn
	if strings.Contains(tgArn, VariantBlue) {
		return awssdk.String(VariantBlue), nil
	}
	if strings.Contains(tgArn, VariantGreen) {
		return awssdk.String(VariantGreen), nil
	}
	return nil, fmt.Errorf("Cannot find variant name in target group ARN: %s", tgArn)
}

func (d DeploymentImpl) TagImage(repositoryName, currentTag, newTag *string) error {
	return d.ecr.TagImage(repositoryName, currentTag, newTag)
}

func (d DeploymentImpl) ListServices() ([]*ListServicesItemOutput, error) {
	svcList, err := d.ecs.DescribeServices()
	if err != nil {
		return nil, err
	}
	output := make([]*ListServicesItemOutput, len(svcList))
	for i, svc := range svcList {
		item := &ListServicesItemOutput{
			ServiceName: svc.ServiceName,
			ClusterArn:  svc.ClusterArn,
		}
		if len(svc.LoadBalancers) > 0 {
			item.ContainerPort = svc.LoadBalancers[0].ContainerPort
			item.TargetGroupArn = svc.LoadBalancers[0].TargetGroupArn
		}
		output[i] = item
	}

	return output, nil
}

func NewDeployment() (Deployment, error) {
	sess, err := aws.GetSession()
	if err != nil {
		return nil, err
	}

	ecsSvc := ecs.NewEcs(sess)
	codedeploySvc := codedeploy.NewCodeDeploy(sess)
	ecrSvc := ecr.NewEcr(sess)
	return &DeploymentImpl{
		ecs:        ecsSvc,
		codedeploy: codedeploySvc,
		ecr:        ecrSvc,
	}, nil
}
