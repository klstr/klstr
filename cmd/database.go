package cmd

import (
	klstr "github.com/klstr/klstr/pkg"
	"github.com/spf13/cobra"
)

func NewDatabaseCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "database",
		Short: "db",
		Long:  "Manage db",
	}
	cmd.AddCommand(newDBCloneCommand())
	return cmd
}

func newDBCloneCommand() *cobra.Command {
	var (
		dbname  string
		dbtype  string
		dbiname string
	)
	cmd := &cobra.Command{
		Use:   "clone",
		Short: "Clone an existing database",
		Long:  "Clone an existing mysql or postgres database from one namespace to another",
		Run: func(cmd *cobra.Command, args []string) {
			err := klstr.CloneDB(&klstr.DatabaseConfig{
				Name:    dbname,
				DBType:  dbtype,
				DBIName: dbiname,
			}, kubeConfig)
			if err != nil {
				panic(err)
			}
		},
	}
	cmd.Flags().StringVar(&dbname, "name", "", "--name=dbname1")
	cmd.Flags().StringVar(&dbtype, "type", "pg", "--type=pg/mysql")
	cmd.Flags().StringVar(&dbiname, "instance-name", "", "--instance-name=db1")
	return cmd
}
