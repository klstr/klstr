package cmd

import (
	"github.com/klstr/klstr/pkg/controller"
	"github.com/spf13/cobra"
)

func NewControllerCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "controller",
		Short: "launch as a controller",
		Long:  "launch as a controller with in cluster config",
		Run: func(cmd *cobra.Command, args []string) {
			err := controller.SetupController()
			if err != nil {
				panic(err)
			}
		},
	}
}
