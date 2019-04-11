package codedeploy

import (
	"github.com/aws/aws-sdk-go/service/codedeploy"
	"github.com/aws/aws-sdk-go/service/codedeploy/codedeployiface"

	"github.com/stretchr/testify/assert"
	"testing"
)

var testDeploy1 = "firstDeployID"
var testDeploy2 = "secondDeployID"

type mockedCodedeploy struct {
	codedeployiface.CodeDeployAPI
	ListDeployOut         codedeploy.ListDeploymentsOutput
	ContinueDeploymentOut codedeploy.ContinueDeploymentOutput
	StopDeploymentOut     codedeploy.StopDeploymentOutput
}

func (m mockedCodedeploy) ListDeployments(input *codedeploy.ListDeploymentsInput) (*codedeploy.ListDeploymentsOutput, error) {
	return &m.ListDeployOut, nil
}

func (m mockedCodedeploy) ContinueDeployment(deploymentID *codedeploy.ContinueDeploymentInput) (*codedeploy.ContinueDeploymentOutput, error) {
	return &m.ContinueDeploymentOut, nil
}

func (m mockedCodedeploy) StopDeployment(deploymentID *codedeploy.StopDeploymentInput) (*codedeploy.StopDeploymentOutput, error) {
	return &m.StopDeploymentOut, nil
}

var testCodeSvc = CodeDeployImpl{
	svc: mockedCodedeploy{
		ListDeployOut: codedeploy.ListDeploymentsOutput{
			Deployments: []*string{&testDeploy1, &testDeploy2},
		},
		ContinueDeploymentOut: codedeploy.ContinueDeploymentOutput{},
		StopDeploymentOut:     codedeploy.StopDeploymentOutput{},
	},
}

func TestListDeployments(t *testing.T) {
	codedeployApp := "codedeployApp"
	codedeployGroup := "codedeployGroup"

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
