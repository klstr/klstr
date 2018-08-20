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
		fromdbname string
		todbname   string
		dbtype     string
		dbiname    string
	)
	cmd := &cobra.Command{
		Use:   "clone",
		Short: "Clone an existing database",
		Long:  "Clone an existing mysql or postgres database from one namespace to another",
		Run: func(cmd *cobra.Command, args []string) {
			err := klstr.CloneDB(&klstr.DatabaseConfig{
				FromDBName: fromdbname,
				ToDBName:   todbname,
				DBType:     dbtype,
				DBIName:    dbiname,
			}, kubeConfig)
			if err != nil {
				panic(err)
			}
		},
	}
	cmd.Flags().StringVar(&fromdbname, "from-db", "", "--from-db=dbname1")
	cmd.Flags().StringVar(&todbname, "to-db", "", "--to-db=dbname2")
	cmd.Flags().StringVar(&dbtype, "type", "pg", "--type=pg/mysql")
	cmd.Flags().StringVar(&dbiname, "instance-name", "", "--instance-name=db1")
	return cmd
}
