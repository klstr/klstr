package cmd

import (
	"github.com/spf13/cobra"
)

func NewDatabasesCommand() *cobra.Command {
	dbsCommand := &cobra.Command {
		Use: "databases",
	}
	dbsCommand.AddCommand(newDatabasesCreateCommand())
	return dbsCommand
}

var dbtype string
var instance string
var namespace string
var name string

func newDatabasesCreateCommand() *cobra.Command {
	createCommand := &cobra.Command {
		Use: "create",
		Short: "create a database",
		Long: "Create a database given an instance and in a given namespace",
		Run: func(cmd *cobra.Command, args []string) {
			
		}
	}
	createCommand.Flags().StringVar(&name,"name","","set the database name")
	createCommand.Flags().StringVar(&instance,"instance","","instance name to use")
	createCommand.Flags().StringVar(&namespace,"namespace","","namespace with within which we should create the secret")
	createCommand.Flags().StringVar(&dbtype,"type","","database type") // should be eliminated when we turn it into a controller
	return createCommand
}