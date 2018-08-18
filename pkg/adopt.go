package klstr

import (
	"os"

	prometheusop "github.com/coreos/prometheus-operator/pkg/client/monitoring"
	prometheusopv1 "github.com/coreos/prometheus-operator/pkg/client/monitoring/v1"
	"github.com/klstr/klstr/pkg/manifests"
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
	ao         AdoptOptions
	clientSet  *kubernetes.Clientset
	pclientSet *prometheusop.Clientset
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
	pclientSet, err := prometheusop.NewForConfig(
		&prometheusopv1.DefaultCrdKinds,
		"monitoring.coreos.com",
		config,
	)
	if err != nil {
		log.Errorf("Unable to create prometheus operator client from config - %s", err.Error())
		panic(err)
	}
	adopter := &Adopter{
		ao:         ao,
		clientSet:  clientSet,
		pclientSet: pclientSet,
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
	err = manifests.NewPrometheusOperatorInstaller(a.clientSet, a.pclientSet).InstallService()
	if err != nil {
		log.Errorf("Unable to install prometheus operator - %s", err)
		panic(err)
	}
}
