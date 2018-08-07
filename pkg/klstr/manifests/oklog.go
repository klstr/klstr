package manifests

import (
	"fmt"
	"io/ioutil"

	"github.com/klstr/klstr/pkg/klstr/util"
	log "github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	typedappsv1 "k8s.io/client-go/kubernetes/typed/apps/v1"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

type OkLogInstaller struct {
	cs *kubernetes.Clientset
}

func NewOkLogInstaller(cs *kubernetes.Clientset) *OkLogInstaller {
	return &OkLogInstaller{cs: cs}
}

func (oi *OkLogInstaller) InstallService() error {
	err := ensureStatefulSet(oi.cs)
	if err != nil {
		return err
	}
	return ensureService(oi.cs)
}

func ensureStatefulSet(cs *kubernetes.Clientset) error {
	si := cs.AppsV1().StatefulSets("default")
	statefulList, err := si.List(metav1.ListOptions{LabelSelector: "app=oklog"})
	if err != nil {
		log.Errorf("unable to list any statefulset")
		return err
	}
	if len(statefulList.Items) > 0 {
		log.Infof("Found oklog statefulset: %+v", statefulList.Items[0])
	} else {
		log.Infof("creating oklog statefulset")
		err = createStatefulSet(si)
		if err != nil {
			log.Errorf("unable to create stateful set %v", err)
			return err
		}
	}
	return nil
}

func ensureService(cs *kubernetes.Clientset) error {
	si := cs.CoreV1().Services("default")
	s, err := si.Get("oklog", metav1.GetOptions{})
	if err == nil {
		log.Infof("Found oklog service %+v", s)
		return nil
	}
	err = createService(si)
	log.Info("Created oklog service %+v", s)
	return nil
}

func createStatefulSet(si typedappsv1.StatefulSetInterface) error {
	ssObj, err := getStatefulSetSpecFromFile()
	if err != nil {
		return err
	}
	sset, err := si.Create(ssObj)
	if err != nil {
		log.Errorf("unable to create oklog deployment %v", err)
		return err
	}
	log.Infof("Created oklog statefulset %+v", sset)
	return nil
}

const OkLogImage = "oklog/oklog:v0.3.2"

func getStatefulSetSpecFromFile() (*appsv1.StatefulSet, error) {
	data, err := ioutil.ReadFile("k8s/logging/oklog-ss.yaml")
	if err != nil {
		return nil, err
	}
	schemaDecoder := util.NewSchemaDecoder(data)
	object, err := schemaDecoder.Decode()
	if err != nil {
		return nil, err
	}
	sobj := object.(*appsv1.StatefulSet)
	buildOkLogArgs(sobj)
	return sobj, nil
}

func buildOkLogArgs(object *appsv1.StatefulSet) {
	prefix := object.ObjectMeta.Name
	args := []string{
		"ingeststore",
		"--debug",
		"--api=tcp://0.0.0.0:7650",
		"--ingest.fast=tcp://0.0.0.0:7651",
		"--ingest.durable=tcp://0.0.0.0:7652",
		"--ingest.bulk=tcp://0.0.0.0:7653",
		"--cluster=tcp://$(POD_IP):7659",
	}
	for i := 0; i < int(*object.Spec.Replicas); i++ {
		args = append(args, fmt.Sprintf("--peer=%s-%d", prefix, i))
	}
	object.Spec.Template.Spec.Containers[0].Args = args
}

func createService(si typedcorev1.ServiceInterface) error {
	svcObj, err := getServiceSpecFromFile()
	if err != nil {
		return err
	}
	svc, err := si.Create(svcObj)
	if err != nil {
		log.Errorf("unable to create oklog service %s", err)
		return err
	}
	log.Infof("Created service %+v", svc)
	return nil
}

func getServiceSpecFromFile() (*corev1.Service, error) {
	data, err := ioutil.ReadFile("k8s/logging/oklog-service.yaml")
	if err != nil {
		return nil, err
	}
	schemaDecoder := util.NewSchemaDecoder(data)
	object, err := schemaDecoder.Decode()
	if err != nil {
		return nil, err
	}
	return object.(*corev1.Service), nil
}

func getMeta() metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name:   "oklog",
		Labels: map[string]string{"app": "oklog"},
	}
}
