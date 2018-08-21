package klstr

import (
	"io/ioutil"

	"github.com/klstr/klstr/pkg/command_jobs"
	"github.com/klstr/klstr/pkg/util"
	log "github.com/sirupsen/logrus"
	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/client-go/kubernetes"
)

type DatabaseConfig struct {
	DBName   string
	ToDBName string
	DBType   string
	DBIName  string
}

type DatabaseJob struct {
	cs *kubernetes.Clientset
	dc *DatabaseConfig
}

func CreateDB(dc *DatabaseConfig, kubeconfig string) error {
	cs, err := util.NewKubeClient(kubeconfig)
	if err != nil {
		return err
	}
	dj := DatabaseJob{
		cs: cs,
		dc: dc,
	}
	err = dj.CreateDBJob()
	if err != nil {
		return err
	}
	return nil
}

func CloneDB(dc *DatabaseConfig, kubeconfig string) error {
	cs, err := util.NewKubeClient(kubeconfig)
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

func (dj *DatabaseJob) CreateDBJob() error {
	ji := dj.cs.BatchV1().Jobs("klstr")
	jobobj, err := getJobFromFile()
	err = buildCreateJobCommand(jobobj, dj.dc)
	if err != nil {
		return err
	}
	job, err := ji.Create(jobobj)
	if err != nil {
		log.Errorf("unable to create db create job %v", err)
		return err
	}
	log.Infof("Created db create job %+v", job)
	return nil
}

func (dj *DatabaseJob) CreateCloneDBJob() error {
	ji := dj.cs.BatchV1().Jobs("klstr")
	jobobj, err := getJobFromFile()
	err = buildCloneJobCommand(jobobj, dj.dc)
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

func getJobFromFile() (*batchv1.Job, error) {
	data, err := ioutil.ReadFile("k8s/jobs/clone_db.yaml")
	if err != nil {
		return nil, err
	}
	schemaDecoder := util.NewSchemaDecoder(data)
	object, err := schemaDecoder.Decode()
	if err != nil {
		return nil, err
	}
	return object.(*batchv1.Job), nil
}

func buildCloneJobCommand(object *batchv1.Job, dc *DatabaseConfig) error {
	cj, err := command_jobs.CreateCommandJob(dc.DBType, command_jobs.CommandJobOptions{
		DBName:   dc.DBName,
		ToDBName: dc.ToDBName,
		DBIName:  dc.DBIName,
	})
	if err != nil {
		log.Errorf("unable to create command job %v", err)
		return err
	}
	cj.BuildCloneCommand(object)
	return nil
}

func buildCreateJobCommand(object *batchv1.Job, dc *DatabaseConfig) error {
	cj, err := command_jobs.CreateCommandJob(dc.DBType, command_jobs.CommandJobOptions{
		DBName:   dc.DBName,
		ToDBName: dc.ToDBName,
		DBIName:  dc.DBIName,
	})
	if err != nil {
		log.Errorf("unable to create command job %v", err)
		return err
	}
	cj.BuildCreateCommand(object)
	return nil
}
