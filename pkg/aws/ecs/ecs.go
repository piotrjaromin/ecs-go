package ecs

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"

	"fmt"
)

type ECS interface {
	DescribeServices() ([]*ecs.Service, error)
	GetService(clusterName, serviceName *string) (*ecs.Service, error)
	GetTaskDefinition(taskDefArn *string) (*ecs.TaskDefinition, error)
	GetLatestTaskDefinition(taskDefArn *string) (*ecs.TaskDefinition, error)
	UpdateService(toUpdate *ecs.Service) (*ecs.Service, error)
	UpdateTaskDefinitions(taskDef *ecs.TaskDefinition, image *string, imageIndex int) (*ecs.TaskDefinition, error)
}

type ECSImpl struct {
	svc *ecs.ECS
}

func (e ECSImpl) GetService(clusterName, serviceName *string) (*ecs.Service, error) {

	input := &ecs.DescribeServicesInput{
		Services: []*string{serviceName},
		Cluster:  clusterName,
	}

	svcList, err := e.svc.DescribeServices(input)
	if err != nil {
		return nil, err
	}

	if len(svcList.Services) == 0 {
		return nil, fmt.Errorf("Unable to find %s service in %s cluster", *clusterName, *serviceName)
	}

	return svcList.Services[0], nil
}

func (e ECSImpl) UpdateService(toUpdate *ecs.Service) (*ecs.Service, error) {

	input := &ecs.UpdateServiceInput{
		Cluster:                 toUpdate.ClusterArn,
		DeploymentConfiguration: toUpdate.DeploymentConfiguration,
		DesiredCount:            toUpdate.DesiredCount,
		Service:                 toUpdate.ServiceName,
	}

	updateOutput, err := e.svc.UpdateService(input)
	if err != nil {
		return nil, err
	}

	return updateOutput.Service, nil
}

func (e ECSImpl) UpdateTaskDefinitions(taskDef *ecs.TaskDefinition, image *string, imageIndex int) (*ecs.TaskDefinition, error) {

	if len(taskDef.ContainerDefinitions) == 0 {
		return nil, fmt.Errorf("No task definitions defined")
	}

	taskDef.ContainerDefinitions[0].Image = image

	input := &ecs.RegisterTaskDefinitionInput{
		ContainerDefinitions:    taskDef.ContainerDefinitions,
		Volumes:                 taskDef.Volumes,
		Cpu:                     taskDef.Cpu,
		Family:                  taskDef.Family,
		IpcMode:                 taskDef.IpcMode,
		Memory:                  taskDef.Memory,
		NetworkMode:             taskDef.NetworkMode,
		PidMode:                 taskDef.PidMode,
		PlacementConstraints:    taskDef.PlacementConstraints,
		RequiresCompatibilities: taskDef.RequiresCompatibilities,
		TaskRoleArn:             taskDef.TaskRoleArn,
		ExecutionRoleArn:        taskDef.ExecutionRoleArn,
		ProxyConfiguration:      taskDef.ProxyConfiguration,
	}

	output, err := e.svc.RegisterTaskDefinition(input)
	if err != nil {
		return nil, err
	}

	return output.TaskDefinition, nil
}

func (e ECSImpl) GetTaskDefinition(taskDefArn *string) (*ecs.TaskDefinition, error) {
	input := ecs.DescribeTaskDefinitionInput{
		TaskDefinition: taskDefArn,
	}

	out, err := e.svc.DescribeTaskDefinition(&input)
	if err != nil {
		return nil, err
	}

	return out.TaskDefinition, nil
}

func (e ECSImpl) GetLatestTaskDefinition(currentTaskDefArn *string) (*ecs.TaskDefinition, error) {
	taskDefFamily := GetFamilyFromTaskDefArn(*currentTaskDefArn)
	listTaskInput := &ecs.ListTaskDefinitionsInput{
		FamilyPrefix: taskDefFamily,
		Sort:         aws.String("DESC"), // newest first
		Status:       aws.String("ACTIVE"),
	}

	taskDefList, err := e.svc.ListTaskDefinitions(listTaskInput)
	if err != nil {
		return nil, err
	}

	if len(taskDefList.TaskDefinitionArns) == 0 {
		return nil, fmt.Errorf("There was no active revision for task family: %s", *taskDefFamily)
	}

	taskDefArn := taskDefList.TaskDefinitionArns[0]
	return e.GetTaskDefinition(taskDefArn)
}

func (e ECSImpl) DescribeServices() ([]*ecs.Service, error) {
	clsArns := make([]*string, 0, 10)
	err := e.svc.ListClustersPages(&ecs.ListClustersInput{}, func(page *ecs.ListClustersOutput, lastPage bool) bool {
		clsArns = append(clsArns, page.ClusterArns...)
		return true
	})
	if err != nil {
		return nil, err
	}

	allSvc := make([]*ecs.Service, 0, 10)
	for _, arn := range clsArns {
		svcArns := make([]*string, 0, 10)
		err := e.svc.ListServicesPages(&ecs.ListServicesInput{Cluster: arn}, func(page *ecs.ListServicesOutput, lastPage bool) bool {
			svcArns = append(svcArns, page.ServiceArns...)
			return true
		})
		if err != nil {
			return nil, err
		}
		// describe operation takes max 10 service arns
		for i := 0; i < len(svcArns); i += 10 {
			end := i + 10
			if end > len(svcArns) {
				end = len(svcArns)
			}
			output, err := e.svc.DescribeServices(&ecs.DescribeServicesInput{
				Cluster:  arn,
				Services: svcArns[i:end],
			})
			if err != nil {
				return nil, err
			}
			allSvc = append(allSvc, output.Services...)
		}
	}

	return allSvc, nil
}

func GetFamilyFromTaskDefArn(arn string) *string {
	// currentTaskDefArn has format arn:aws:ecs:eu-central-1:ACCOUNT_NUMBER:task-definition/SERVICE_NAME:REVISION
	// we want to exctract SERVICE_NAME
	taskDefFamilyWithRev := strings.Split(arn, ":")[5]
	taskDefFamily := strings.Split(taskDefFamilyWithRev, "/")[1]
	return &taskDefFamily
}

// func (e ECSImpl) CreateTaskDefinition() (*ecs.TaskDefinition, error) {
// 	input := ecs.RegisterTaskDefinitionInput{}

// 	out, err := e.svc.RegisterTaskDefinition(input)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return out.TaskDefinition, nil
// }

func NewEcs(sess *session.Session) ECS {

	return &ECSImpl{
		svc: ecs.New(sess),
	}
}
