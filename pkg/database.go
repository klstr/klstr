package klstr

import (
	"errors"
	"io/ioutil"

	"github.com/klstr/klstr/pkg/command_jobs"
	"github.com/klstr/klstr/pkg/util"
	log "github.com/sirupsen/logrus"
	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/client-go/kubernetes"
	typedbatchv1 "k8s.io/client-go/kubernetes/typed/batch/v1"
	"k8s.io/client-go/tools/clientcmd"
)

type DatabaseConfig struct {
	Name    string
	DBType  string
	DBIName string
}

type DatabaseJob struct {
	cs *kubernetes.Clientset
	dc *DatabaseConfig
}

func CloneDB(dc *DatabaseConfig, kubeconfig string) error {
	if kubeconfig == "" {
		return errors.New("Kubeconfig is empty")
	}
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return err
	}
	cs, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}
	dj := DatabaseJob{
		cs: cs,
		dc: dc,
	}
	err = dj.CreateCloneDBJob()
	if err != nil {
		return err
	}
	return nil
}

func (dj *DatabaseJob) CreateCloneDBJob() error {
	ji := dj.cs.BatchV1().Jobs("klstr")
	jobobj, err := getJobFromFile(ji, dj.dc)
	if err != nil {
		return err
	}
	job, err := ji.Create(jobobj)
	if err != nil {
		log.Errorf("unable to create db clone job %v", err)
		return err
	}
	log.Infof("Created db clone job %+v", job)
	return nil
}

func getJobFromFile(ji typedbatchv1.JobInterface, dc *DatabaseConfig) (*batchv1.Job, error) {
	data, err := ioutil.ReadFile("k8s/jobs/clone_db.yaml")
	if err != nil {
		return nil, err
	}
	schemaDecoder := util.NewSchemaDecoder(data)
	object, err := schemaDecoder.Decode()
	if err != nil {
		return nil, err
	}
	job := object.(*batchv1.Job)
	buildJobCommand(job, dc)
	return job, nil
}

func buildJobCommand(object *batchv1.Job, dc *DatabaseConfig) {
	cj, err := command_jobs.CreateCommandJob(dc.DBType, command_jobs.CommandJobOptions{
		DBName:  dc.Name,
		DBIName: dc.DBIName,
	})
	if err != nil {
	}
	cj.BuildCommand(object)
}
