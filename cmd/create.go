package cmd

import (
	klstr "github.com/klstr/klstr/pkg"
	"github.com/spf13/cobra"
)

func NewCreateCommand() *cobra.Command {
	var tag string
	var domain string
	var cloudProvider string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new kubernetes cluster",
		Long:  "Create a new HA cluster and setup kubernetes",
		Run: func(cmd *cobra.Command, args []string) {
			creater := klstr.NewCreater(klstr.ClusterOptions{
				Tag:           tag,
				Domain:        domain,
				CloudProvider: cloudProvider,
			})
			creater.CreateCluster()
		},
	}

	cmd.Flags().StringVar(&tag, "tag", "awsklstr", "Name of the cluster")
	cmd.Flags().StringVar(&domain, "domain", "dev.klstr.io", "Base Domain of the cluster used for ingress hosts")
	cmd.Flags().StringVar(&cloudProvider, "cloud-provider", "aws", "Cloud Provider used for provisioning")

	return cmd
}
