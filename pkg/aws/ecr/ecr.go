package ecr

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
)

type ECR interface {
	TagImage(repositoryName, currentTag, newTag *string) error
}

type ECRImpl struct {
	svc *ecr.ECR
}

func (e ECRImpl) TagImage(repositoryName, currentTag, newTag *string) error {
	manifest, err := e.getImageManifest(repositoryName, currentTag)
	if err != nil {
		return err
	}
	err = e.putImage(repositoryName, manifest, newTag)
	if err != nil {
		return err
	}
	return nil
}

func (e ECRImpl) getImageManifest(repositoryName, imageTag *string) (*string, error) {
	batchGetInput := &ecr.BatchGetImageInput{
		RepositoryName: repositoryName,
		ImageIds: []*ecr.ImageIdentifier{
			&ecr.ImageIdentifier{
				ImageTag: imageTag,
			},
		},
	}
	batchGetOutput, err := e.svc.BatchGetImage(batchGetInput)
	if err != nil {
		return nil, err
	}
	if len(batchGetOutput.Images) == 0 {
		return nil, fmt.Errorf("No images in repository %s with tag: %s", *repositoryName, *imageTag)
	}
	return batchGetOutput.Images[0].ImageManifest, nil
}

func (e ECRImpl) putImage(repositoryName, manifest, imageTag *string) error {
	putInput := &ecr.PutImageInput{
		RepositoryName: repositoryName,
		ImageManifest:  manifest,
		ImageTag:       imageTag,
	}
	_, err := e.svc.PutImage(putInput)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			if awsErr.Code() == ecr.ErrCodeImageAlreadyExistsException {
				return nil
			}
		}
		return err
	}
	return nil
}

func NewEcr(sess *session.Session) ECR {
	return &ECRImpl{
		svc: ecr.New(sess),
	}
}
