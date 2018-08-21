package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var kubeConfig string

var RootCmd = &cobra.Command{
	Use:   "klstr",
	Short: "klstr - friendly neighborhood kubernetes helper",
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&kubeConfig, "kubeconfig", "", "kubeconfig to use for interacting with klstr")

	RootCmd.AddCommand(NewAdoptCommand())
	RootCmd.AddCommand(NewUsersCommand())
	RootCmd.AddCommand(NewCreateCommand())
	RootCmd.AddCommand(NewDeleteCommand())
	RootCmd.AddCommand(NewDBInstancesCommand())
	RootCmd.AddCommand(NewControllerCommand())
	RootCmd.AddCommand(NewDatabaseCommand())
}

func initConfig() {
}
