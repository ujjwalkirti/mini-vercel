package ecs

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

type EnvVar struct {
	Name  string
	Value string
}

type Service struct {
	client         *ecs.Client
	cluster        string
	taskDef        string
	subnets        []string
	securityGrp    string
	assignPublicIP types.AssignPublicIp
	launchType     types.LaunchType
	count          int32
	imageName      string
}

func New(cfg aws.Config, cluster, taskDef string, subnets []string, securityGrp, assignPublicIP, launchType, imageName string, count int) *Service {
	// Parse AssignPublicIp
	publicIP := types.AssignPublicIpEnabled
	if assignPublicIP == "DISABLED" {
		publicIP = types.AssignPublicIpDisabled
	}

	// Parse LaunchType
	launch := types.LaunchTypeFargate
	if launchType == "EC2" {
		launch = types.LaunchTypeEc2
	}

	return &Service{
		client:         ecs.NewFromConfig(cfg),
		cluster:        cluster,
		taskDef:        taskDef,
		subnets:        subnets,
		securityGrp:    securityGrp,
		assignPublicIP: publicIP,
		launchType:     launch,
		count:          int32(count),
		imageName:      imageName,
	}
}

// RunTask triggers an ECS task with the provided environment variables
func (s *Service) RunTask(ctx context.Context, envVars []EnvVar) (*string, error) {
	// Convert EnvVar to ECS KeyValuePair
	var ecsEnvVars []types.KeyValuePair
	for _, env := range envVars {
		ecsEnvVars = append(ecsEnvVars, types.KeyValuePair{
			Name:  aws.String(env.Name),
			Value: aws.String(env.Value),
		})
	}

	// Run the ECS task
	input := &ecs.RunTaskInput{
		Cluster:        aws.String(s.cluster),
		TaskDefinition: aws.String(s.taskDef),
		LaunchType:     s.launchType,
		Count:          aws.Int32(s.count),
		NetworkConfiguration: &types.NetworkConfiguration{
			AwsvpcConfiguration: &types.AwsVpcConfiguration{
				Subnets:        s.subnets,
				SecurityGroups: []string{s.securityGrp},
				AssignPublicIp: s.assignPublicIP,
			},
		},
		Overrides: &types.TaskOverride{
			ContainerOverrides: []types.ContainerOverride{
				{
					Name:        aws.String(s.imageName),
					Environment: ecsEnvVars,
				},
			},
		},
	}

	result, err := s.client.RunTask(ctx, input)
	if err != nil {
		return nil, err
	}

	// Return the task ARN
	if len(result.Tasks) > 0 {
		return result.Tasks[0].TaskArn, nil
	}

	return nil, nil
}
