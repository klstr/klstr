package manifests

import (
	"io/ioutil"

	"github.com/klstr/klstr/pkg/klstr/util"
	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	extnv1beta1 "k8s.io/api/extensions/v1beta1"
	rbacv1beta1 "k8s.io/api/rbac/v1beta1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
)

type PrometheusOperatorInstaller struct {
	cs *kubernetes.Clientset
}

func NewPrometheusOperatorInstaller(cs *kubernetes.Clientset) *PrometheusOperatorInstaller {
	return &PrometheusOperatorInstaller{cs: cs}
}

func (pi *PrometheusOperatorInstaller) InstallService() error {
	err := ensurePrometheusOperator(pi.cs)
	if err != nil {
		return err
	}
	return nil
}

func ensurePrometheusOperator(cs *kubernetes.Clientset) error {
	di := cs.AppsV1().Deployments("default")
	deploymentList, err := di.List(metav1.ListOptions{
		LabelSelector: "k8s-app=prometheus-operator",
	})
	if err != nil {
		log.Errorf("unable to list any deployments")
		return err
	}
	if len(deploymentList.Items) > 0 {
		log.Infof("Found deployment %+v", deploymentList.Items[0])
	} else {
		log.Infof("creating prometheus operator")
		err = createPrometheusOperator(cs)
		if err != nil {
			log.Errorf("unable to create prometheus operator %v", err)
			return err
		}
	}
	return nil
}

func createPrometheusOperator(cs *kubernetes.Clientset) error {
	objects, err := getSpecFromFile()
	if err != nil {
		return err
	}
	for _, object := range objects {
		err = createObject(cs, object)
		if err != nil {
			log.Errorf("unable to create object: %+v", object)
		}
	}
	return nil
}

func createObject(cs *kubernetes.Clientset, object runtime.Object) error {
	var kobj runtime.Object
	var err error
	switch o := object.(type) {
	case *extnv1beta1.Deployment:
		kobj, err = cs.ExtensionsV1beta1().Deployments("default").Create(o)
	case *corev1.Service:
		kobj, err = cs.CoreV1().Services("default").Create(o)
	case *rbacv1beta1.ClusterRoleBinding:
		kobj, err = cs.RbacV1beta1().ClusterRoleBindings().Create(o)
	case *rbacv1beta1.ClusterRole:
		kobj, err = cs.RbacV1beta1().ClusterRoles().Create(o)
	case *corev1.ServiceAccount:
		kobj, err = cs.CoreV1().ServiceAccounts("default").Create(o)
	}
	if err != nil {
		log.Errorf("unable to create prometheus operator %v", err)
		return err
	}
	log.Infof("Created prometheus operator %+v", kobj)
	return nil
}

func getSpecFromFile() ([]runtime.Object, error) {
	data, err := ioutil.ReadFile("k8s/monitoring/prometheus-operator.yaml")
	if err != nil {
		return nil, err
	}
	schemaDecoder := util.NewSchemaDecoder(data)
	return schemaDecoder.MultiDecode()
}
