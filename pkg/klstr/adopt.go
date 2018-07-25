package klstr

import (
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type AdoptOptions struct {
	KubeConfig  string
	SkipLogging bool
	SkipMetrics bool
}

func AdoptCluster(ao *AdoptOptions) {
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
	pods, err := clientSet.CoreV1().Pods("kube-system").List(metav1.ListOptions{})
	for _, p := range pods.Items {
		log.Infof("Found pod %s", p.Name)
	}
}
