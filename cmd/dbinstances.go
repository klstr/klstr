package cmd

import (
	klstr "github.com/klstr/klstr/pkg"
	"github.com/spf13/cobra"
)

func NewDBInstancesCommand() *cobra.Command {
	dbiCommand := &cobra.Command{
		Use:   "dbinstances",
		Short: "dbi",
		Long:  "manage db instances",
	}
	dbiCommand.AddCommand(newDBIRegisterCommand())
	return dbiCommand
}

var dbiname string
var dbtype string
var host string
var port int
var username string
var password string

func newDBIRegisterCommand() *cobra.Command {
	dbiRegisterCmd := &cobra.Command{
		Use:   "register",
		Short: "Register a database instance",
		Long:  "Register a mysql or postgres instance often with admin credentials",
		Run: func(cmd *cobra.Command, args []string) {
			err := klstr.RegisterDBInstance(&klstr.DBInstanceRegistration{
				Name:     dbiname,
				Host:     host,
				Port:     port,
				DBType:   dbtype,
				Username: username,
				Password: password,
			}, kubeConfig)
			if err != nil {
				panic(err)
			}
		},
	}
	dbiRegisterCmd.Flags().StringVar(&dbiname, "name", "", "--name=stolon")
	dbiRegisterCmd.Flags().StringVar(&dbtype, "type", "pg", "--type=pg/mysql")
	dbiRegisterCmd.Flags().IntVar(&port, "port", 5432, "--port=5432")
	dbiRegisterCmd.Flags().StringVar(&host, "host", "postgres", "--host=postgres")
	dbiRegisterCmd.Flags().StringVar(&username, "username", "postgres", "--username=postgres")
	dbiRegisterCmd.Flags().StringVar(&password, "password", "password1", "--password=password1")
	return dbiRegisterCmd
}
