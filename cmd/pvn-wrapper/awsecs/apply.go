package awsecs

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/prodvana/pvn-wrapper/cmdutil"
	"github.com/spf13/cobra"
)

var applyFlags = struct {
	taskDefinitionFile   string
	ecsClusterName       string
	ecsServiceName       string
	pvnServiceId         string
	pvnServiceVersion    string
	networkConfiguration string
	desiredCount         int
}{}

type registerTaskDefinitionOutput struct {
	TaskDefinition struct {
		TaskDefinitionArn string `json:"taskDefinitionArn"`
	} `json:"taskDefinition"`
}

func registerTaskDefinitionIfNeeded(taskDefPath string) (string, error) {
	// TODO(naphat) skip registering if task definition already exists
	taskDefPath, err := filepath.Abs(taskDefPath)
	if err != nil {
		return "", errors.Wrap(err, "failed to make abs path")
	}
	registerCmd := exec.Command(
		awsPath,
		"ecs",
		"register-task-definition",
		"--cli-input-json",
		fmt.Sprintf("file://%s", taskDefPath),
	)
	output, err := registerCmd.Output()
	if err != nil {
		return "", err
	}
	var registerOutput registerTaskDefinitionOutput
	if err := json.Unmarshal(output, &registerOutput); err != nil {
		return "", errors.Wrap(err, "failed to unmarshal register-task-definition output")
	}
	taskArn := registerOutput.TaskDefinition.TaskDefinitionArn
	if taskArn == "" {
		return "", errors.Errorf("got empty task definition arn. Register output: %s", string(output))
	}
	return taskArn, nil
}

type describeServicesOutput struct {
	Services []struct {
		Status string `json:"status"`
	} `json:"services"`
	Failures []struct {
		Reason string `json:"reason"`
	} `json:"failures"`
}

func describeService(clusterName, serviceName string) (*describeServicesOutput, error) {
	describeCmd := exec.Command(
		awsPath,
		"ecs",
		"describe-services",
		"--cluster",
		clusterName,
		"--services",
		serviceName,
	)
	output, err := cmdutil.RunCmdOutput(describeCmd)
	if err != nil {
		return nil, err
	}
	var describeOutput describeServicesOutput
	if err := json.Unmarshal(output, &describeOutput); err != nil {
		return nil, errors.Wrap(err, "failed to unmarsal describe-services output")
	}
	return &describeOutput, nil
}

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Create or update an ECS service",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO(naphat) patch task definition file with our own tags
		taskArn, err := registerTaskDefinitionIfNeeded(applyFlags.taskDefinitionFile)
		if err != nil {
			return err
		}
		serviceOutput, err := describeService(applyFlags.ecsClusterName, applyFlags.ecsServiceName)
		if err != nil {
			return err
		}
		if len(serviceOutput.Failures) > 0 {
			if serviceOutput.Failures[0].Reason != "MISSING" {
				return errors.Errorf("unexpected failure reason: %s", serviceOutput.Failures[0].Reason)
			}
		} else {
			if len(serviceOutput.Services) != 1 {
				return errors.Errorf("unexpected number of services: %d", len(serviceOutput.Services))
			}
		}
		commonArgs := []string{
			"--task-definition",
			taskArn,
			"--desired-count",
			fmt.Sprintf("%d", applyFlags.desiredCount),
			"--network-configuration",
			applyFlags.networkConfiguration,
		}
		if serviceOutput.Services[0].Status == "INACTIVE" || serviceOutput.Services[0].Status == "MISSING" {
			// create service
			createCmd := exec.Command(awsPath, append([]string{
				"ecs",
				"create-service",
				"--service-name",
				applyFlags.ecsServiceName,
				"--launch-type=FARGATE",
			}, commonArgs...)...)
			_, err := cmdutil.RunCmdOutput(createCmd)
			if err != nil {
				return err
			}
		} else {
			// update service
			updateCmd := exec.Command(awsPath, append([]string{
				"ecs",
				"update-service",
				"--service",
				applyFlags.ecsServiceName,
			}, commonArgs...)...)
			_, err := cmdutil.RunCmdOutput(updateCmd)
			if err != nil {
				return nil
			}
		}
		waitCmd := exec.Command(awsPath, "ecs", "wait", "services-stable", "--services", applyFlags.ecsServiceName)
		_, err = cmdutil.RunCmdOutput(waitCmd)
		return err
	},
}

func init() {
	RootCmd.AddCommand(applyCmd)

	applyCmd.Flags().StringVar(&applyFlags.taskDefinitionFile, "task-definition-file", "", "Path to task definition file")
	cmdutil.Must(applyCmd.MarkFlagRequired("task-definition-file"))
	applyCmd.Flags().StringVar(&applyFlags.ecsServiceName, "ecs-service-name", "", "Name of ECS service")
	cmdutil.Must(applyCmd.MarkFlagRequired("ecs-service-name"))
	applyCmd.Flags().StringVar(&applyFlags.ecsClusterName, "ecs-cluster-name", "", "Name of ECS cluster")
	cmdutil.Must(applyCmd.MarkFlagRequired("ecs-cluster-name"))
	applyCmd.Flags().StringVar(&applyFlags.pvnServiceId, "pvn-service-id", "", "Prodvana Service ID")
	cmdutil.Must(applyCmd.MarkFlagRequired("pvn-service-id"))
	applyCmd.Flags().StringVar(&applyFlags.pvnServiceVersion, "pvn-service-version", "", "Prodvana Service Version")
	cmdutil.Must(applyCmd.MarkFlagRequired("pvn-service-version"))
	applyCmd.Flags().StringVar(&applyFlags.networkConfiguration, "network-configuration", "", "awsvpc network configuration")
	cmdutil.Must(applyCmd.MarkFlagRequired("network-configuration"))
	applyCmd.Flags().IntVar(&applyFlags.desiredCount, "desired-count", 0, "Number of instances desired")
	cmdutil.Must(applyCmd.MarkFlagRequired("desired-count"))
}
