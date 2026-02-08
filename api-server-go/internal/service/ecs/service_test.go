package ecs

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

// MockECSClient is a mock implementation of the ECS client
type MockECSClient struct {
	RunTaskFunc func(ctx context.Context, params *ecs.RunTaskInput, optFns ...func(*ecs.Options)) (*ecs.RunTaskOutput, error)
}

func (m *MockECSClient) RunTask(ctx context.Context, params *ecs.RunTaskInput, optFns ...func(*ecs.Options)) (*ecs.RunTaskOutput, error) {
	if m.RunTaskFunc != nil {
		return m.RunTaskFunc(ctx, params, optFns...)
	}
	return nil, errors.New("RunTaskFunc not set")
}

// TestNew tests service initialization
func TestNew(t *testing.T) {
	t.Run("should create service with all parameters", func(t *testing.T) {
		// Arrange
		cfg := aws.Config{}
		cluster := "test-cluster"
		taskDef := "test-task-def"
		subnets := []string{"subnet-1", "subnet-2"}
		securityGrp := "sg-123"
		assignPublicIP := "ENABLED"
		launchType := "FARGATE"
		imageName := "test-image"
		count := 1

		// Act
		service := New(cfg, cluster, taskDef, subnets, securityGrp, assignPublicIP, launchType, imageName, count)

		// Assert
		if service == nil {
			t.Fatal("expected service to not be nil")
		}
		if service.cluster != cluster {
			t.Errorf("expected cluster %s, got %s", cluster, service.cluster)
		}
		if service.taskDef != taskDef {
			t.Errorf("expected taskDef %s, got %s", taskDef, service.taskDef)
		}
		if len(service.subnets) != len(subnets) {
			t.Errorf("expected %d subnets, got %d", len(subnets), len(service.subnets))
		}
		if service.securityGrp != securityGrp {
			t.Errorf("expected securityGrp %s, got %s", securityGrp, service.securityGrp)
		}
		if service.assignPublicIP != types.AssignPublicIpEnabled {
			t.Errorf("expected assignPublicIP to be ENABLED")
		}
		if service.launchType != types.LaunchTypeFargate {
			t.Errorf("expected launchType to be FARGATE")
		}
		if service.count != 1 {
			t.Errorf("expected count 1, got %d", service.count)
		}
	})

	t.Run("should handle DISABLED assignPublicIP", func(t *testing.T) {
		// Arrange
		cfg := aws.Config{}
		assignPublicIP := "DISABLED"

		// Act
		service := New(cfg, "cluster", "taskDef", []string{}, "sg", assignPublicIP, "FARGATE", "image", 1)

		// Assert
		if service.assignPublicIP != types.AssignPublicIpDisabled {
			t.Errorf("expected assignPublicIP to be DISABLED")
		}
	})

	t.Run("should handle EC2 launch type", func(t *testing.T) {
		// Arrange
		cfg := aws.Config{}
		launchType := "EC2"

		// Act
		service := New(cfg, "cluster", "taskDef", []string{}, "sg", "ENABLED", launchType, "image", 1)

		// Assert
		if service.launchType != types.LaunchTypeEc2 {
			t.Errorf("expected launchType to be EC2")
		}
	})

	t.Run("should handle empty subnets", func(t *testing.T) {
		// Arrange
		cfg := aws.Config{}
		subnets := []string{}

		// Act
		service := New(cfg, "cluster", "taskDef", subnets, "sg", "ENABLED", "FARGATE", "image", 1)

		// Assert
		if service == nil {
			t.Fatal("expected service to not be nil")
		}
		if len(service.subnets) != 0 {
			t.Errorf("expected 0 subnets, got %d", len(service.subnets))
		}
	})

	t.Run("should handle missing optional parameters", func(t *testing.T) {
		// Arrange
		cfg := aws.Config{}

		// Act
		service := New(cfg, "cluster", "taskDef", []string{}, "", "", "", "", 0)

		// Assert
		if service == nil {
			t.Fatal("expected service to not be nil")
		}
		if service.count != 0 {
			t.Errorf("expected count 0, got %d", service.count)
		}
	})
}

// TestRunTask tests running an ECS task
func TestRunTask(t *testing.T) {
	t.Run("should run task with correct configuration", func(t *testing.T) {
		// Arrange
		taskARN := "arn:aws:ecs:us-east-1:123456789012:task/test-task"
		mockClient := &MockECSClient{
			RunTaskFunc: func(ctx context.Context, params *ecs.RunTaskInput, optFns ...func(*ecs.Options)) (*ecs.RunTaskOutput, error) {
				// Verify parameters
				if *params.Cluster != "test-cluster" {
					t.Errorf("expected cluster test-cluster, got %s", *params.Cluster)
				}
				if *params.TaskDefinition != "test-task-def" {
					t.Errorf("expected taskDef test-task-def, got %s", *params.TaskDefinition)
				}

				return &ecs.RunTaskOutput{
					Tasks: []types.Task{
						{TaskArn: aws.String(taskARN)},
					},
				}, nil
			},
		}

		service := &Service{
			client:         mockClient,
			cluster:        "test-cluster",
			taskDef:        "test-task-def",
			subnets:        []string{"subnet-1"},
			securityGrp:    "sg-123",
			assignPublicIP: types.AssignPublicIpEnabled,
			launchType:     types.LaunchTypeFargate,
			count:          1,
			imageName:      "test-image",
		}

		envVars := []EnvVar{
			{Name: "PROJECT_ID", Value: "project-123"},
			{Name: "GIT_URL", Value: "https://github.com/user/repo"},
		}

		// Act
		result, err := service.RunTask(context.Background(), envVars)

		// Assert
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if result == nil {
			t.Fatal("expected result to not be nil")
		}
		if *result != taskARN {
			t.Errorf("expected taskARN %s, got %s", taskARN, *result)
		}
	})

	t.Run("should include all environment variables", func(t *testing.T) {
		// Arrange
		mockClient := &MockECSClient{
			RunTaskFunc: func(ctx context.Context, params *ecs.RunTaskInput, optFns ...func(*ecs.Options)) (*ecs.RunTaskOutput, error) {
				// Verify environment variables
				if len(params.Overrides.ContainerOverrides) == 0 {
					t.Error("expected container overrides")
					return nil, errors.New("no container overrides")
				}

				envVars := params.Overrides.ContainerOverrides[0].Environment
				if len(envVars) != 3 {
					t.Errorf("expected 3 env vars, got %d", len(envVars))
				}

				return &ecs.RunTaskOutput{
					Tasks: []types.Task{
						{TaskArn: aws.String("arn:test")},
					},
				}, nil
			},
		}

		service := &Service{
			client:      mockClient,
			cluster:     "test-cluster",
			taskDef:     "test-task-def",
			subnets:     []string{"subnet-1"},
			securityGrp: "sg-123",
			launchType:  types.LaunchTypeFargate,
			count:       1,
			imageName:   "test-image",
		}

		envVars := []EnvVar{
			{Name: "VAR1", Value: "value1"},
			{Name: "VAR2", Value: "value2"},
			{Name: "VAR3", Value: "value3"},
		}

		// Act
		_, err := service.RunTask(context.Background(), envVars)

		// Assert
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("should use correct cluster name", func(t *testing.T) {
		// Arrange
		expectedCluster := "production-cluster"
		mockClient := &MockECSClient{
			RunTaskFunc: func(ctx context.Context, params *ecs.RunTaskInput, optFns ...func(*ecs.Options)) (*ecs.RunTaskOutput, error) {
				if *params.Cluster != expectedCluster {
					t.Errorf("expected cluster %s, got %s", expectedCluster, *params.Cluster)
				}
				return &ecs.RunTaskOutput{
					Tasks: []types.Task{{TaskArn: aws.String("arn:test")}},
				}, nil
			},
		}

		service := &Service{
			client:      mockClient,
			cluster:     expectedCluster,
			taskDef:     "test-task-def",
			subnets:     []string{},
			securityGrp: "sg-123",
			launchType:  types.LaunchTypeFargate,
			count:       1,
			imageName:   "test-image",
		}

		// Act
		_, err := service.RunTask(context.Background(), []EnvVar{})

		// Assert
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("should use correct task definition", func(t *testing.T) {
		// Arrange
		expectedTaskDef := "my-task-definition:5"
		mockClient := &MockECSClient{
			RunTaskFunc: func(ctx context.Context, params *ecs.RunTaskInput, optFns ...func(*ecs.Options)) (*ecs.RunTaskOutput, error) {
				if *params.TaskDefinition != expectedTaskDef {
					t.Errorf("expected taskDef %s, got %s", expectedTaskDef, *params.TaskDefinition)
				}
				return &ecs.RunTaskOutput{
					Tasks: []types.Task{{TaskArn: aws.String("arn:test")}},
				}, nil
			},
		}

		service := &Service{
			client:      mockClient,
			cluster:     "test-cluster",
			taskDef:     expectedTaskDef,
			subnets:     []string{},
			securityGrp: "sg-123",
			launchType:  types.LaunchTypeFargate,
			count:       1,
			imageName:   "test-image",
		}

		// Act
		_, err := service.RunTask(context.Background(), []EnvVar{})

		// Assert
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("should configure network with correct subnets", func(t *testing.T) {
		// Arrange
		expectedSubnets := []string{"subnet-1", "subnet-2", "subnet-3"}
		mockClient := &MockECSClient{
			RunTaskFunc: func(ctx context.Context, params *ecs.RunTaskInput, optFns ...func(*ecs.Options)) (*ecs.RunTaskOutput, error) {
				subnets := params.NetworkConfiguration.AwsvpcConfiguration.Subnets
				if len(subnets) != len(expectedSubnets) {
					t.Errorf("expected %d subnets, got %d", len(expectedSubnets), len(subnets))
				}
				return &ecs.RunTaskOutput{
					Tasks: []types.Task{{TaskArn: aws.String("arn:test")}},
				}, nil
			},
		}

		service := &Service{
			client:      mockClient,
			cluster:     "test-cluster",
			taskDef:     "test-task-def",
			subnets:     expectedSubnets,
			securityGrp: "sg-123",
			launchType:  types.LaunchTypeFargate,
			count:       1,
			imageName:   "test-image",
		}

		// Act
		_, err := service.RunTask(context.Background(), []EnvVar{})

		// Assert
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("should configure security group", func(t *testing.T) {
		// Arrange
		expectedSG := "sg-abcdef123"
		mockClient := &MockECSClient{
			RunTaskFunc: func(ctx context.Context, params *ecs.RunTaskInput, optFns ...func(*ecs.Options)) (*ecs.RunTaskOutput, error) {
				sgs := params.NetworkConfiguration.AwsvpcConfiguration.SecurityGroups
				if len(sgs) != 1 || sgs[0] != expectedSG {
					t.Errorf("expected security group %s, got %v", expectedSG, sgs)
				}
				return &ecs.RunTaskOutput{
					Tasks: []types.Task{{TaskArn: aws.String("arn:test")}},
				}, nil
			},
		}

		service := &Service{
			client:      mockClient,
			cluster:     "test-cluster",
			taskDef:     "test-task-def",
			subnets:     []string{"subnet-1"},
			securityGrp: expectedSG,
			launchType:  types.LaunchTypeFargate,
			count:       1,
			imageName:   "test-image",
		}

		// Act
		_, err := service.RunTask(context.Background(), []EnvVar{})

		// Assert
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("should set task count correctly", func(t *testing.T) {
		// Arrange
		expectedCount := int32(3)
		mockClient := &MockECSClient{
			RunTaskFunc: func(ctx context.Context, params *ecs.RunTaskInput, optFns ...func(*ecs.Options)) (*ecs.RunTaskOutput, error) {
				if *params.Count != expectedCount {
					t.Errorf("expected count %d, got %d", expectedCount, *params.Count)
				}
				return &ecs.RunTaskOutput{
					Tasks: []types.Task{{TaskArn: aws.String("arn:test")}},
				}, nil
			},
		}

		service := &Service{
			client:      mockClient,
			cluster:     "test-cluster",
			taskDef:     "test-task-def",
			subnets:     []string{"subnet-1"},
			securityGrp: "sg-123",
			launchType:  types.LaunchTypeFargate,
			count:       expectedCount,
			imageName:   "test-image",
		}

		// Act
		_, err := service.RunTask(context.Background(), []EnvVar{})

		// Assert
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("should return error when AWS API fails", func(t *testing.T) {
		// Arrange
		expectedError := errors.New("ECS API error")
		mockClient := &MockECSClient{
			RunTaskFunc: func(ctx context.Context, params *ecs.RunTaskInput, optFns ...func(*ecs.Options)) (*ecs.RunTaskOutput, error) {
				return nil, expectedError
			},
		}

		service := &Service{
			client:      mockClient,
			cluster:     "test-cluster",
			taskDef:     "test-task-def",
			subnets:     []string{"subnet-1"},
			securityGrp: "sg-123",
			launchType:  types.LaunchTypeFargate,
			count:       1,
			imageName:   "test-image",
		}

		// Act
		_, err := service.RunTask(context.Background(), []EnvVar{})

		// Assert
		if err == nil {
			t.Error("expected error but got nil")
		}
		if err.Error() != expectedError.Error() {
			t.Errorf("expected error %v, got %v", expectedError, err)
		}
	})

	t.Run("should handle context cancellation", func(t *testing.T) {
		// Arrange
		mockClient := &MockECSClient{
			RunTaskFunc: func(ctx context.Context, params *ecs.RunTaskInput, optFns ...func(*ecs.Options)) (*ecs.RunTaskOutput, error) {
				return nil, context.Canceled
			},
		}

		service := &Service{
			client:      mockClient,
			cluster:     "test-cluster",
			taskDef:     "test-task-def",
			subnets:     []string{"subnet-1"},
			securityGrp: "sg-123",
			launchType:  types.LaunchTypeFargate,
			count:       1,
			imageName:   "test-image",
		}

		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		// Act
		_, err := service.RunTask(ctx, []EnvVar{})

		// Assert
		if err == nil {
			t.Error("expected error but got nil")
		}
		if !errors.Is(err, context.Canceled) {
			t.Errorf("expected context.Canceled error, got %v", err)
		}
	})

	t.Run("should handle empty environment variables", func(t *testing.T) {
		// Arrange
		mockClient := &MockECSClient{
			RunTaskFunc: func(ctx context.Context, params *ecs.RunTaskInput, optFns ...func(*ecs.Options)) (*ecs.RunTaskOutput, error) {
				envVars := params.Overrides.ContainerOverrides[0].Environment
				if len(envVars) != 0 {
					t.Errorf("expected 0 env vars, got %d", len(envVars))
				}
				return &ecs.RunTaskOutput{
					Tasks: []types.Task{{TaskArn: aws.String("arn:test")}},
				}, nil
			},
		}

		service := &Service{
			client:      mockClient,
			cluster:     "test-cluster",
			taskDef:     "test-task-def",
			subnets:     []string{"subnet-1"},
			securityGrp: "sg-123",
			launchType:  types.LaunchTypeFargate,
			count:       1,
			imageName:   "test-image",
		}

		// Act
		_, err := service.RunTask(context.Background(), []EnvVar{})

		// Assert
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("should return nil when no tasks are created", func(t *testing.T) {
		// Arrange
		mockClient := &MockECSClient{
			RunTaskFunc: func(ctx context.Context, params *ecs.RunTaskInput, optFns ...func(*ecs.Options)) (*ecs.RunTaskOutput, error) {
				return &ecs.RunTaskOutput{
					Tasks: []types.Task{}, // Empty tasks
				}, nil
			},
		}

		service := &Service{
			client:      mockClient,
			cluster:     "test-cluster",
			taskDef:     "test-task-def",
			subnets:     []string{"subnet-1"},
			securityGrp: "sg-123",
			launchType:  types.LaunchTypeFargate,
			count:       1,
			imageName:   "test-image",
		}

		// Act
		result, err := service.RunTask(context.Background(), []EnvVar{})

		// Assert
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if result != nil {
			t.Error("expected nil result when no tasks created")
		}
	})
}

// TestEnvVar tests the EnvVar struct
func TestEnvVar(t *testing.T) {
	t.Run("should have name and value fields", func(t *testing.T) {
		// Arrange & Act
		env := EnvVar{
			Name:  "TEST_VAR",
			Value: "test_value",
		}

		// Assert
		if env.Name != "TEST_VAR" {
			t.Errorf("expected Name TEST_VAR, got %s", env.Name)
		}
		if env.Value != "test_value" {
			t.Errorf("expected Value test_value, got %s", env.Value)
		}
	})
}
