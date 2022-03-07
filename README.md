# ECS-GO

Tool for deploying new version of ecs service with usage od codedeploy.

Tool updates docker image in task definition, it searches given cluster for service, and then updates task definition associated with it, at the end it triggers deployment through code deploy.

## Supported commands

```bash
deploy, d                             Deploys new version of app, takes newest task definition and updates docker image by creating new revision
scale, s                              Changes number of tasks
continue-deployment, cd               Allows active deployment to continue deployment
force-continue-deployment, cd         Forces active deployment to continue deployment (kills replacement task before its time elapsed
list-deployments, ld                  Deploys new version of app
rollback-deployment, rd               Rollbacks active deployment
wait-for-state, cd                    waits until given deployment reaches given state
continue-latest-deployment, cd        Trigger continue on latest deployment for given app and group
force-continue-latest-deployment, cd  Trigger force continue on latest deployment for given app and group
rollback-latest-deployment, cd        Trigger rollback on latest deployment for given app and group
wait-for-latest, cd                   waits until given deployment reaches given state
get-live-variant, glv                 Gets live variant of service
tag-image, ti                         Tags ECR image
help, h                               Shows a list of commands or help for one command
```

## Sample call

For deployment:

```bash
ecs-go deploy --clusterName cluster_name \
    --serviceName service_name \
    --codedeployGroup codedeploy_group \
    --codedeployApp codedeploy_group \
    --image nginx:latest
```

## AWS Assume role

If your AWS configuration (`~/.aws`) requires to assume role then you can try calling calling this like below:

```bash
export AWS_SDK_LOAD_CONFIG=true
export AWS_PROFILE=aws_profile_with_role_to_assume
ecs-go deploy .....
```

or

```bash
export AWS_SESSION_TOKEN=token
export AWS_SECRET_ACCESS_KEY=access_secret_key
export AWS_ACCESS_KEY_ID=aws_key
ecs-go deploy .....
```
