package cmd

import (
	"github.com/klstr/klstr/pkg/klstr"
	"github.com/spf13/cobra"
)

var skipLogging bool
var skipMetrics bool

func NewAdoptCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "adopt",
		Short: "Adopt a new kubernetes cluster",
		Long:  "Adopts a new kubernetes cluster by installing klstr components",
		Run: func(cmd *cobra.Command, args []string) {
			klstr.AdoptCluster(&klstr.AdoptOptions{
				KubeConfig:  kubeConfig,
				SkipLogging: skipLogging,
				SkipMetrics: skipMetrics,
			})
		},
	}
	cmd.Flags().BoolVar(&skipLogging, "skip-logging", false, "Do not install elastic search logging stack")
	cmd.Flags().BoolVar(&skipMetrics, "skip-metrics", false, "Do not install prometheus and grafana")
	return cmd
}
