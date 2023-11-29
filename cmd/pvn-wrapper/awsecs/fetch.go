package awsecs

import (
	"github.com/spf13/cobra"
)

func runFetch()

var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Fetch current state of an ECS service",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		serviceOutput, err := describeService(commonFlags.ecsClusterName, commonFlags.ecsServiceName)
		if err != nil {
			return err
		}
		if serviceMissing(serviceOutput) {
		}
		return nil
	},
}

func init() {
	RootCmd.AddCommand(fetchCmd)

	registerCommonFlags(fetchCmd)
}
