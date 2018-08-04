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

func AdoptCluster(ao *AdoptOptions) {
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
	err = manifests.GetGrafanaDeployment(clientSet, "aws")
	if err != nil {
		panic(err)
	}
}
