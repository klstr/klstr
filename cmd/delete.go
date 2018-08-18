package cmd

import (
	klstr "github.com/klstr/klstr/pkg"
	"github.com/spf13/cobra"
)

func NewDeleteCommand() *cobra.Command {
	var tag string
	var domain string
	var cloudProvider string

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a kubernetes cluster",
		Long:  "Delete kubernetes cluster and related networking ",
		Run: func(cmd *cobra.Command, args []string) {
			deleter := klstr.NewDeleter(klstr.ClusterOptions{
				Tag:           tag,
				Domain:        domain,
				CloudProvider: cloudProvider,
			})
			deleter.DeleteCluster()
		},
	}

	cmd.Flags().StringVar(&tag, "tag", "awsklstr", "Name of the cluster")
	cmd.Flags().StringVar(&domain, "domain", "dev.klstr.io", "Base Domain of the cluster used for ingress hosts")
	cmd.Flags().StringVar(&cloudProvider, "cloud-provider", "aws", "Cloud Provider used for provisioning")

	return cmd
}
