package codedeploy

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/codedeploy"
	"github.com/aws/aws-sdk-go/service/codedeploy/codedeployiface"

	"github.com/stretchr/testify/assert"
	"testing"
)

var testDeploy1 = "firstDeployID"
var testDeploy2 = "secondDeployID"
var codedeployApp = "codedeployApp"
var codedeployGroup = "codedeployGroup"

type mockedCodedeploy struct {
	codedeployiface.CodeDeployAPI
	ListDeployOut         codedeploy.ListDeploymentsOutput
	ContinueDeploymentOut codedeploy.ContinueDeploymentOutput
	StopDeploymentOut     codedeploy.StopDeploymentOutput
	CreateDeploymentOut   codedeploy.CreateDeploymentOutput
	GetDeploymentOut      codedeploy.GetDeploymentOutput
}

func (m mockedCodedeploy) ListDeployments(input *codedeploy.ListDeploymentsInput) (*codedeploy.ListDeploymentsOutput, error) {
	return &m.ListDeployOut, nil
}

func (m mockedCodedeploy) StopDeployment(deploymentID *codedeploy.StopDeploymentInput) (*codedeploy.StopDeploymentOutput, error) {
	return &m.StopDeploymentOut, nil
}

func (m mockedCodedeploy) ContinueDeployment(deploymentID *codedeploy.ContinueDeploymentInput) (*codedeploy.ContinueDeploymentOutput, error) {
	return &m.ContinueDeploymentOut, nil
}

func (m mockedCodedeploy) CreateDeployment(*codedeploy.CreateDeploymentInput) (*codedeploy.CreateDeploymentOutput, error) {
	return &m.CreateDeploymentOut, nil
}

func (m mockedCodedeploy) GetDeployment(*codedeploy.GetDeploymentInput) (*codedeploy.GetDeploymentOutput, error) {
	return &m.GetDeploymentOut, nil
}

var testCodeSvc = CodeDeployImpl{
	svc: mockedCodedeploy{
		ListDeployOut: codedeploy.ListDeploymentsOutput{
			Deployments: []*string{&testDeploy1, &testDeploy2},
		},
		ContinueDeploymentOut: codedeploy.ContinueDeploymentOutput{},
		StopDeploymentOut:     codedeploy.StopDeploymentOutput{},
		CreateDeploymentOut: codedeploy.CreateDeploymentOutput{
			DeploymentId: aws.String("created-deployment"),
		},
		GetDeploymentOut: codedeploy.GetDeploymentOutput{
			DeploymentInfo: &codedeploy.DeploymentInfo{
				Status: aws.String("Ready"),
			},
		},
	},
}

func TestListDeployments(t *testing.T) {

	deployments, err := testCodeSvc.ListDeployments(&codedeployApp, &codedeployGroup)

	assert.Nil(t, err, "Error should not exist")
	assert.Equal(t, deployments, []*string{&testDeploy1, &testDeploy2})
}

func TestContinueDeployment(t *testing.T) {
	_, err := testCodeSvc.ContinueDeployment(&testDeploy1)

	assert.Nil(t, err, "Error should not exist")
}

func TestRollbackDeployment(t *testing.T) {
	_, err := testCodeSvc.RollbackDeployment(&testDeploy1)

	assert.Nil(t, err, "Error should not exist")
}

func TestCreateDeployment(t *testing.T) {
	taskDefinitionArn := "taskDefinitionArn"
	containerName := "test-name"
	containerPort := "3000"

	id, err := testCodeSvc.CreateDeployment(&codedeployApp, &codedeployGroup, &taskDefinitionArn, &containerName, &containerPort)

	assert.Nil(t, err, "Error should not exist")
	assert.Equal(t, *id, "created-deployment")
}

func TestWaitForSate(t *testing.T) {

	state := "Ready"
	err := testCodeSvc.WaitForSate(&testDeploy1, &state, 5)

	assert.Nil(t, err, "Error should not exist")
}
