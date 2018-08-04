package manifests

import (
	log "github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	typedappsv1 "k8s.io/client-go/kubernetes/typed/apps/v1"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

const GrafanaImage = "grafana/grafana:5.2.2"

func GetGrafanaDeployment(cs *kubernetes.Clientset) error {
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
		err = InstallGrafana(di)
		if err != nil {
			log.Errorf("unable to create deployment %v", err)
			return err
		}
	}
	return nil
}
func GetGrafanaService(cs *kubernetes.Clientset) error {
	si := cs.CoreV1().Services("default")
	return CreateGrafanaService(si)
}

func InstallGrafana(di typedappsv1.DeploymentInterface) error {
	dep, err := di.Create(getGrafanaDeplomentSpec())
	if err != nil {
		return err
	}
	log.Infof("Created a deployment ", dep)
	return nil
}

func CreateGrafanaService(si typedcorev1.ServiceInterface) error {
	s, err := si.Get("grafana", metav1.GetOptions{})
	if err == nil {
		log.Infof("Found grafana service %+v", s)
		return nil
	}
	s, err = si.Create(getGrafanaServiceSpec())
	if err != nil {
		log.Errorf("unable to create service %s", err)
		return err
	}
	log.Infof("Created service %+v", s)
	return nil
}
func getGrafanaServiceSpec() *corev1.Service {
	return &corev1.Service{
		ObjectMeta: getGrafanaMeta(),
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{"app": "grafana"},
			Ports: []corev1.ServicePort{
				corev1.ServicePort{
					TargetPort: intstr.FromInt(3000),
					Name:       "grafana",
					Port:       3000,
				},
			},
		},
	}
}

func getGrafanaDeplomentSpec() *appsv1.Deployment {
	var replicaCount int32 = 1
	deployment := &appsv1.Deployment{
		ObjectMeta: getGrafanaMeta(),
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicaCount,
			Selector: &metav1.LabelSelector{MatchLabels: map[string]string{"app": "grafana"}},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: getGrafanaMeta(),
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						corev1.Container{
							Image:           GrafanaImage,
							Name:            "grafana",
							ImagePullPolicy: corev1.PullAlways,
							Ports:           []corev1.ContainerPort{corev1.ContainerPort{ContainerPort: 3000}},
							Env: []corev1.EnvVar{
								corev1.EnvVar{Name: "GF_AUTH_BASIC_ENABLED", Value: "false"},
								corev1.EnvVar{Name: "GF_AUTH_ANONYMOUS_ENABLED", Value: "true"},
								corev1.EnvVar{Name: "GF_AUTH_ANONYMOUS_ORG_ROLE", Value: "Admin"},
							},
						},
					},
				},
			},
		},
	}
	return deployment
}
func getGrafanaMeta() metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name:   "grafana",
		Labels: map[string]string{"app": "grafana"},
	}
}
