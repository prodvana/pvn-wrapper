package googlecloudrun

import (
	"github.com/prodvana/pvn-wrapper/cmdutil"
	"github.com/spf13/cobra"
)

var commonFlags = struct {
	gcpProject        string
	region            string
	specFile          string
	pvnServiceId      string
	pvnServiceVersion string
}{}

func registerCommonFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&commonFlags.gcpProject, "gcp-project", "", "GCP project")
	cmdutil.Must(cmd.MarkFlagRequired("gcp-project"))
	cmd.Flags().StringVar(&commonFlags.region, "region", "", "GCP region")
	cmdutil.Must(cmd.MarkFlagRequired("region"))
	cmd.Flags().StringVar(&commonFlags.specFile, "spec-file", "", "Path to service spec file")
	cmdutil.Must(cmd.MarkFlagRequired("spec-file"))
	cmd.Flags().StringVar(&commonFlags.pvnServiceId, "pvn-service-id", "", "Prodvana Service ID")
	cmdutil.Must(cmd.MarkFlagRequired("pvn-service-id"))
	cmd.Flags().StringVar(&commonFlags.pvnServiceVersion, "pvn-service-version", "", "Prodvana Service Version")
	cmdutil.Must(cmd.MarkFlagRequired("pvn-service-version"))
}
