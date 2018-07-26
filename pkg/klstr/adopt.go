package klstr

import (
	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
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
	selector := labels.NewSelector()
	req, err := labels.NewRequirement("app", selection.Equals, []string{"helm"})
	if err != nil {
		log.Errorf("Unable to create filter for helm - %s", err.Error())
		panic(err)
	}
	selector = selector.Add(*req)
	log.Info(selector.String())
	pods, err := clientSet.CoreV1().Pods("kube-system").List(metav1.ListOptions{
		LabelSelector: selector.String(),
	})

	pod, err := clientSet.CoreV1().Pods("default").Create(&corev1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind: "Pod",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   "klstr",
			Labels: map[string]string{"app": "klstr"},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				corev1.Container{Name: "klstr", Image: "quay.io/klstr/klstr:latest", Args: []string{"loop"}},
			},
		},
	})
	if err != nil {
		log.Infof("Unable to create pod - %v", err)
	} else {
		log.Infof("Pod created %s", pod.GetName())
	}

	for _, p := range pods.Items {
		log.Infof("Found pod %s", p.Name)
	}
}
