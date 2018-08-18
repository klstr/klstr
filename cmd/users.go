package cmd

import (
	klstr "github.com/klstr/klstr/pkg"
	"github.com/spf13/cobra"
)

func NewUsersCommand() *cobra.Command {
	usersCmd := &cobra.Command{
		Use: "users",
	}
	usersCmd.AddCommand(newCreateUsersCommand())
	return usersCmd
}

var name string

func newCreateUsersCommand() *cobra.Command {
	createUserCmd := &cobra.Command{
		Use:   "create",
		Long:  "Create users in the current cluster",
		Short: "create users",
		Run: func(cmd *cobra.Command, args []string) {
			err := klstr.NewUser(name, kubeConfig)
			if err != nil {
				panic(err)
			}
		},
	}
	createUserCmd.Flags().StringVar(&name, "name", "", "set the username")
	return createUserCmd
}
