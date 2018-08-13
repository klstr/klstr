package manifests

import (
	"io/ioutil"

	prometheusop "github.com/coreos/prometheus-operator/pkg/client/monitoring"
	prometheusopv1 "github.com/coreos/prometheus-operator/pkg/client/monitoring/v1"
	k8sutil "github.com/coreos/prometheus-operator/pkg/k8sutil"
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
	ps *prometheusop.Clientset
}

func NewPrometheusOperatorInstaller(
	cs *kubernetes.Clientset,
	ps *prometheusop.Clientset,
) *PrometheusOperatorInstaller {
	return &PrometheusOperatorInstaller{
		cs: cs,
		ps: ps,
	}
}

func (pi *PrometheusOperatorInstaller) InstallService() error {
	err := ensurePrometheusOperator(pi.cs)
	if err != nil {
		return err
	}
	err = ensurePrometheusRbac(pi.cs)
	if err != nil {
		return err
	}
	err = ensurePrometheusPersisted(pi.ps)
	if err != nil {
		return err
	}
	return ensurePrometheusService(pi.cs)
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
			log.Errorf("unable to create prometheus operator %s", err)
			return err
		}
	}
	return nil
}

func ensurePrometheusRbac(cs *kubernetes.Clientset) error {
	crbi := cs.RbacV1beta1().ClusterRoleBindings()
	crbList, err := crbi.List(metav1.ListOptions{
		FieldSelector: "metadata.name=prometheus-operator",
	})
	if err != nil {
		log.Errorf("unable to list any cluster role bindings")
		return err
	}
	if len(crbList.Items) > 0 {
		log.Infof("Found cluster role binding %+v", crbList.Items[0])
	} else {
		log.Infof("creating prometheus rbac")
		err := createPrometheusRbac(cs)
		if err != nil {
			log.Errorf("unable to create prometheus rbac %s", err)
			return err
		}
	}
	return nil
}

func ensurePrometheusPersisted(ps *prometheusop.Clientset) error {
	log.Info("Waiting for Prometheus CRD to be ready...")
	k8sutil.WaitForCRDReady(
		ps.MonitoringV1().Prometheuses("default").List,
	)
	log.Info("Wait over, resuming")
	pi := ps.MonitoringV1().Prometheuses("default")
	pList, err := pi.List(metav1.ListOptions{
		FieldSelector: "metadata.name=prometheus2",
	})
	if err != nil {
		log.Errorf("unable to list any prometheuses")
		return err
	}
	plItems := pList.(*prometheusopv1.PrometheusList).Items
	if len(plItems) > 0 {
		log.Infof("Found prometheus %+v", plItems[0])
	} else {
		log.Info("creating prometheus")
		err := createPrometheusPersisted(ps)
		if err != nil {
			log.Errorf("unable to create prometheus persisted %s", err)
		}
	}
	return nil
}

func ensurePrometheusService(cs *kubernetes.Clientset) error {
	si := cs.CoreV1().Services("default")
	serviceList, err := si.List(metav1.ListOptions{
		FieldSelector: "metadata.name=prometheus",
	})
	if err != nil {
		log.Errorf("unable to list any services")
		return err
	}
	if len(serviceList.Items) > 0 {
		log.Info("Found service %+v", serviceList.Items[0])
	} else {
		log.Info("creating prometheus service")
		err := createPrometheusService(cs)
		if err != nil {
			log.Errorf("unable to create prometheus service %s", err)
		}
	}
	return nil
}

func createPrometheusOperator(cs *kubernetes.Clientset) error {
	objects, err := getPrometheusOperatorSpecFromFile()
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

func createPrometheusRbac(cs *kubernetes.Clientset) error {
	objects, err := getPrometheusRbacSpecFromFile()
	if err != nil {
		return err
	}
	for _, object := range objects {
		err := createObject(cs, object)
		if err != nil {
			log.Errorf("unable to create object: %+v", object)
		}
	}
	return nil
}

func createPrometheusPersisted(ps *prometheusop.Clientset) error {
	object, err := getPrometheusPersistedSpecFromFile()
	if err != nil {
		return err
	}
	return createPrometheusObject(ps, object)
}

func createPrometheusService(cs *kubernetes.Clientset) error {
	svcObj, err := getPrometheusServiceSpecFromFile()
	if err != nil {
		return err
	}
	return createObject(cs, svcObj)
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

func createPrometheusObject(ps *prometheusop.Clientset, object runtime.Object) error {
	var kobj runtime.Object
	var err error
	switch o := object.(type) {
	case *prometheusopv1.Prometheus:
		kobj, err = ps.MonitoringV1().Prometheuses("default").Create(o)
	case *prometheusopv1.ServiceMonitor:
		kobj, err = ps.MonitoringV1().ServiceMonitors("default").Create(o)
	case *prometheusopv1.Alertmanager:
		kobj, err = ps.MonitoringV1().Alertmanagers("default").Create(o)
	}
	if err != nil {
		log.Errorf("unable to create prometheus object: %v", err)
		return err
	}
	log.Infof("Created prometheus %+v", kobj)
	return nil
}

func getPrometheusOperatorSpecFromFile() ([]runtime.Object, error) {
	data, err := ioutil.ReadFile("k8s/monitoring/prometheus-operator.yaml")
	if err != nil {
		return nil, err
	}
	schemaDecoder := util.NewSchemaDecoder(data)
	return schemaDecoder.MultiDecode()
}

func getPrometheusRbacSpecFromFile() ([]runtime.Object, error) {
	data, err := ioutil.ReadFile("k8s/monitoring/prometheus-rbac.yaml")
	if err != nil {
		return nil, err
	}
	schemaDecoder := util.NewSchemaDecoder(data)
	return schemaDecoder.MultiDecode()
}

func getPrometheusPersistedSpecFromFile() (runtime.Object, error) {
	data, err := ioutil.ReadFile("k8s/monitoring/prometheus-persisted.yaml")
	if err != nil {
		return nil, err
	}
	schemaDecoder := util.NewSchemaDecoder(data)
	object := &prometheusopv1.Prometheus{}
	return schemaDecoder.Decode(object)
}

func getPrometheusServiceSpecFromFile() (*corev1.Service, error) {
	data, err := ioutil.ReadFile("k8s/monitoring/prometheus-service.yaml")
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
