package command_jobs

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	batchv1 "k8s.io/api/batch/v1"
)

type CommandJobOptions struct {
	DBName  string
	DBIName string
}

type CommandJob interface {
	BuildCommand(object *batchv1.Job)
}

type CommandJobFactory func(options CommandJobOptions) CommandJob

var commandJobFactories = make(map[string]CommandJobFactory)

func RegisterCommandJobFactory(
	name string,
	commandJobFactory CommandJobFactory,
) {
	if commandJobFactory == nil {
		log.Errorf("Command Job Factory %s does not exist", name)
		return
	}
	_, registered := commandJobFactories[name]
	if registered {
		log.Errorf("Command Job Factory %s already registered. Ignoring.", name)
		return
	}
	commandJobFactories[name] = commandJobFactory
}

func init() {
	RegisterCommandJobFactory("pg", NewPGCommandJob)
	RegisterCommandJobFactory("mysql", NewMySQLCommandJob)
}

func CreateCommandJob(dbType string, options CommandJobOptions) (CommandJob, error) {
	commandJob, ok := commandJobFactories[dbType]
	if !ok {
		availableCommandJobs := make([]string, len(commandJobFactories))
		for cj := range commandJobFactories {
			availableCommandJobs = append(availableCommandJobs, cj)
		}
		return nil, fmt.Errorf("Invalid DB Type: %s. Myst be one of: %s", dbType, strings.Join(availableCommandJobs, ", "))
	}
	return commandJob(options), nil
}
