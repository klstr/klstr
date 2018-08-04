package klstr

import (
	"os"

	"github.com/klstr/klstr/pkg/klstr/manifests"
	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type AdoptOptions struct {
	KubeConfig  string
	SkipLogging bool
	SkipMetrics bool
}
type Adopter struct {
	ao        AdoptOptions
	clientSet *kubernetes.Clientset
}

type ServiceInstaller interface {
	InstallService() error
}

func NewAdopter(ao AdoptOptions) *Adopter {
	if ao.KubeConfig == "" {
		ao.KubeConfig = os.Getenv("KUBECONFIG")
	}
	config, err := clientcmd.BuildConfigFromFlags("", ao.KubeConfig)
	if err != nil {
		log.Errorf("Unable to setup client config - %s", err.Error())
		panic(err)
	}
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Errorf("Unable to create client from config - %s", err.Error())
		panic(err)
	}
	adopter := &Adopter{
		ao:        ao,
		clientSet: clientSet,
	}
	return adopter
}

func (a *Adopter) AdoptCluster() {
	err := manifests.NewGrafanaInstaller(a.clientSet).InstallService()
	if err != nil {
		log.Errorf("Unable to install grafana - %s", err)
		return
	}
	err = manifests.NewOkLogInstaller(a.clientSet).InstallService()
	if err != nil {
		log.Errorf("Unable to install oklog - %s", err)
		panic(err)
	}
}
