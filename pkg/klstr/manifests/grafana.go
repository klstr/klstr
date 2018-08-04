package manifests

import (
	"io/ioutil"

	"github.com/klstr/klstr/pkg/klstr/util"
	log "github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	typedappsv1 "k8s.io/client-go/kubernetes/typed/apps/v1"
)

type GrafanaInstaller struct {
	cs *kubernetes.Clientset
}

func NewGrafanaInstaller(cs *kubernetes.Clientset) *GrafanaInstaller {
	return &GrafanaInstaller{cs: cs}
}

func (gi *GrafanaInstaller) InstallService() error {
	err := ensureGrafanaDeployment(gi.cs)
	if err != nil {
		return err
	}
	return ensureGrafanaService(gi.cs)
}

const GrafanaImage = "grafana/grafana:5.2.2"

func ensureGrafanaDeployment(cs *kubernetes.Clientset) error {
	di := cs.AppsV1().Deployments("default")
	deploymentList, err := di.List(metav1.ListOptions{LabelSelector: "app=grafana"})
	if err != nil {
		log.Errorf("unable to list any deployments")
		return err
	}
	if len(deploymentList.Items) > 0 {
		log.Infof("Found deployment %+v", deploymentList.Items[0])
	} else {
		log.Infof("creating deployment")
		err = createGrafanaDeployment(di)
		if err != nil {
			log.Errorf("unable to create deployment %v", err)
			return err
		}
	}
	return nil
}
func ensureGrafanaService(cs *kubernetes.Clientset) error {
	si := cs.CoreV1().Services("default")
	s, err := si.Get("grafana", metav1.GetOptions{})
	if err == nil {
		log.Infof("Found grafana service %+v", s)
		return nil
	}
	sobj, err := getGrafanaServiceSpecFromFile()
	if err != nil {
		log.Info("unable to decode from service file ", err)
		return err
	}
	s, err = si.Create(sobj)
	if err != nil {
		log.Errorf("unable to create service %s", err)
		return err
	}
	log.Infof("Created service %+v", s)
	return nil
}

func createGrafanaDeployment(di typedappsv1.DeploymentInterface) error {
	depObj, err := getGrafanaDeplomentSpecFromFile()
	if err != nil {
		return err
	}
	dep, err := di.Create(depObj)
	if err != nil {
		return err
	}
	log.Infof("Created a deployment ", dep)
	return nil
}

func getGrafanaServiceSpecFromFile() (*corev1.Service, error) {
	data, err := ioutil.ReadFile("k8s/monitoring/grafana-service.yaml")
	if err != nil {
		return nil, err
	}
	schemaDecoder := util.NewSchemaDecoder(data)
	object := &corev1.Service{}
	err = schemaDecoder.Decode(object)
	if err != nil {
		return nil, err
	}
	return object, nil
}

func getGrafanaDeplomentSpecFromFile() (*appsv1.Deployment, error) {
	data, err := ioutil.ReadFile("k8s/monitoring/grafana-deployment.yaml")
	if err != nil {
		return nil, err
	}
	schemaDecoder := util.NewSchemaDecoder(data)
	object := &appsv1.Deployment{}
	err = schemaDecoder.Decode(object)
	if err != nil {
		return nil, err
	}
	return object, nil
}

func getGrafanaMeta() metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name:   "grafana",
		Labels: map[string]string{"app": "grafana"},
	}
}
