package codedeploy

import (
	"crypto/sha256"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/codedeploy"
	"github.com/aws/aws-sdk-go/service/codedeploy/codedeployiface"
	"time"
)

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
		DeploymentWaitType: aws.String("READY_WAIT"),
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
		IncludeOnlyStatuses: []*string{aws.String("InProgress"), aws.String("Ready"), aws.String("Created"), aws.String("Queued")},
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
			Events:  []*string{aws.String("DEPLOYMENT_FAILURE"), aws.String("DEPLOYMENT_STOP_ON_ALARM"), aws.String("DEPLOYMENT_STOP_ON_REQUEST")},
		},
		Revision: &codedeploy.RevisionLocation{
			AppSpecContent: &codedeploy.AppSpecContent{
				Content: &appSecContent,
				Sha256:  &appSecSha256,
			},
			RevisionType: aws.String("AppSpecContent"),
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
