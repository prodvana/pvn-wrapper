package main

import (
	"fmt"
	"log"

	"github.com/prodvana/pvn-wrapper/cmd/pvn-wrapper/awsecs"
	"github.com/prodvana/pvn-wrapper/cmd/pvn-wrapper/fly"
	"github.com/prodvana/pvn-wrapper/cmd/pvn-wrapper/googlecloudrun"
	"github.com/prodvana/pvn-wrapper/cmd/pvn-wrapper/pulumi"
	"github.com/prodvana/pvn-wrapper/cmd/pvn-wrapper/terraform"
	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

var rootCmd = &cobra.Command{
	Use:              "pvn-wrapper",
	Short:            "pvn-wrapper is used to facilitate executions of jobs in Prodvana.",
	Long:             `pvn-wrapper is used to facilitate executions of jobs in Prodvana.`,
	TraverseChildren: true,
}

func init() {
	rootCmd.AddCommand(awsecs.RootCmd)
	rootCmd.AddCommand(terraform.RootCmd)
	rootCmd.AddCommand(pulumi.RootCmd)
	rootCmd.AddCommand(googlecloudrun.RootCmd)
	rootCmd.AddCommand(fly.RootCmd)
	rootCmd.Version = version
	rootCmd.SetVersionTemplate(fmt.Sprintf("{{ .Version }} (%s %s)\n", commit, date))
}

func main() {
	log.SetFlags(log.Lmicroseconds | log.Lshortfile)
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
