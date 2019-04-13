package ecs

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"

	"fmt"
)

type ECS interface {
	GetService(clusterName, serviceName *string) (*ecs.Service, error)
	GetTaskDefinition(taskDefArn *string) (*ecs.TaskDefinition, error)
	UpdateTaskDefinitions(taskDef *ecs.TaskDefinition, image *string) (*ecs.TaskDefinition, error)
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

func (e ECSImpl) UpdateTaskDefinitions(taskDef *ecs.TaskDefinition, image *string) (*ecs.TaskDefinition, error) {

	if len(taskDef.ContainerDefinitions) == 0 {
		return nil, fmt.Errorf("No task definitions defnied")
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
