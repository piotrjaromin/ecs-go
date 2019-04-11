package codedeploy

import (
	"crypto/sha256"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/codedeploy"
	"github.com/aws/aws-sdk-go/service/codedeploy/codedeployiface"
	"time"
)

var appSec = "AppSpecContent"

var deployFail = "DEPLOYMENT_FAILURE"
var deployStopOnAlarm = "DEPLOYMENT_STOP_ON_ALARM"
var deployStopOnReq = "DEPLOYMENT_STOP_ON_REQUEST"

var events = []*string{&deployFail, &deployStopOnAlarm, &deployStopOnReq}

var inProgress = "InProgress"
var ready = "Ready"
var created = "Created"
var queued = "Queued"

var includeStatuses = []*string{&inProgress, &ready, &created, &queued}

var readyWait = "READY_WAIT"

type CodeDeploy interface {
	ContinueDeployment(deploymentID *string) (*codedeploy.ContinueDeploymentOutput, error)
	RollbackDeployment(deploymentID *string) (*codedeploy.StopDeploymentOutput, error)
	ListDeployments(codedeployApp, codedeployGroup *string) ([]*string, error)
	CreateDeployment(codedeployApp, codedeployGroup, taskDefinitionArn, containerName, containerPort *string) (*string, error)
}

type CodeDeployImpl struct {
	svc codedeployiface.CodeDeployAPI
}

func (d CodeDeployImpl) ContinueDeployment(deploymentID *string) (*codedeploy.ContinueDeploymentOutput, error) {

	input := &codedeploy.ContinueDeploymentInput{
		DeploymentId:       deploymentID,
		DeploymentWaitType: &readyWait,
	}

	return d.svc.ContinueDeployment(input)
}

func (d CodeDeployImpl) RollbackDeployment(deploymentID *string) (*codedeploy.StopDeploymentOutput, error) {

	autoRollbackEnabled := true

	input := &codedeploy.StopDeploymentInput{
		AutoRollbackEnabled: &autoRollbackEnabled,
		DeploymentId:        deploymentID,
	}

	return d.svc.StopDeployment(input)
}

func (d CodeDeployImpl) ListDeployments(codedeployApp, codedeployGroup *string) ([]*string, error) {

	input := &codedeploy.ListDeploymentsInput{
		ApplicationName:     codedeployApp,
		DeploymentGroupName: codedeployGroup,
		IncludeOnlyStatuses: includeStatuses,
	}

	output, err := d.svc.ListDeployments(input)
	if err != nil {
		return []*string{}, err
	}
	return output.Deployments, nil
}

func (d CodeDeployImpl) CreateDeployment(codedeployApp, codedeployGroup, taskDefinitionArn, containerName, containerPort *string) (*string, error) {

	desc := fmt.Sprintf("Handled by ecs-go at %v", time.Now())

	enabled := true

	appSecContent := appSpec(taskDefinitionArn, containerName, containerPort)

	h := sha256.New()
	h.Write([]byte(appSecContent))

	appSecSha256 := fmt.Sprintf("%x", h.Sum(nil))

	input := &codedeploy.CreateDeploymentInput{
		ApplicationName:     codedeployApp,
		DeploymentGroupName: codedeployGroup,
		Description:         &desc,
		AutoRollbackConfiguration: &codedeploy.AutoRollbackConfiguration{
			Enabled: &enabled,
			Events:  events,
		},
		Revision: &codedeploy.RevisionLocation{
			AppSpecContent: &codedeploy.AppSpecContent{
				Content: &appSecContent,
				Sha256:  &appSecSha256,
			},
			RevisionType: &appSec,
		},
	}
	output, err := d.svc.CreateDeployment(input)

	if err != nil {
		return nil, err
	}

	return output.DeploymentId, nil
}

func appSpec(taskDefinitionWithRevisionArn, containerName, containerPort *string) string {
	return fmt.Sprintf(`{
		"version": 1,
		"Resources": [
			{
				"TargetService": {
					"Type": "AWS::ECS::Service",
					"Properties": {
						"TaskDefinition": "%s",
						"LoadBalancerInfo": {
							"ContainerName": "%s",
							"ContainerPort": %s
						}
					}
				}
			}
		]
	}`, *taskDefinitionWithRevisionArn, *containerName, *containerPort)
}

func NewCodeDeploy(sess *session.Session) CodeDeploy {
	return &CodeDeployImpl{
		svc: codedeploy.New(sess),
	}
}
