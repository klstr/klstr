package controller

import (
	"time"

	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func SetupController() error {
	config, err := rest.InClusterConfig()
	if err != nil {
		return err
	}
	cs, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}
	for {
		podList, err := cs.CoreV1().Pods("default").List(metav1.ListOptions{})
		if err != nil {
			return err
		}
		log.Infof("Found %d pods.", len(podList.Items))
		log.Infof("Pods : %v", podList.Items)
		time.Sleep(10 * time.Second)
	}
}
