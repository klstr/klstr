package manifests

import (
	log "github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
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
	sset, err := si.Create(getStatefulSetSpec())
	if err != nil {
		log.Errorf("unable to create oklog deployment %v", err)
		return err
	}
	log.Infof("Created oklog statefulset %+v", sset)
	return nil
}

const OkLogImage = "oklog/oklog:v0.3.2"

func getStatefulSetSpec() *appsv1.StatefulSet {
	var replicaCount int32 = 3
	var storageClassName = "ssd"
	quantity, err := resource.ParseQuantity("10Gi")
	if err != nil {
		log.Errorf("Error parsing resource quantity")
		return nil
	}
	return &appsv1.StatefulSet{
		ObjectMeta: getMeta(),
		Spec: appsv1.StatefulSetSpec{
			Replicas: &replicaCount,
			Selector: &metav1.LabelSelector{MatchLabels: map[string]string{"app": "oklog"}},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: getMeta(),
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						corev1.Container{
							Name:            "oklog",
							Image:           OkLogImage,
							ImagePullPolicy: corev1.PullAlways,
							Env: []corev1.EnvVar{
								corev1.EnvVar{
									Name: "POD_IP",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											FieldPath: "status.podIP",
										},
									},
								},
								corev1.EnvVar{
									Name: "POD_NAMESPACE",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											FieldPath: "metadata.namespace",
										},
									},
								},
							},
							Ports: []corev1.ContainerPort{
								corev1.ContainerPort{
									Name:          "api",
									ContainerPort: 7650,
								},
								corev1.ContainerPort{
									Name:          "ingest-fast",
									ContainerPort: 7651,
								},
								corev1.ContainerPort{
									Name:          "ingest-durable",
									ContainerPort: 7652,
								},
								corev1.ContainerPort{
									Name:          "ingest-bulk",
									ContainerPort: 7653,
								},
								corev1.ContainerPort{
									Name:          "cluster",
									ContainerPort: 7659,
								},
							},
							Args: []string{
								"ingeststore",
								"--debug",
								"--api=tcp://0.0.0.0:7650",
								"--ingest.fast=tcp://0.0.0.0:7651",
								"--ingest.durable=tcp://0.0.0.0:7652",
								"--ingest.bulk=tcp://0.0.0.0:7653",
								"--cluster=tcp://$(POD_IP):7659",
								"--peer=oklog-0.oklog",
								"--peer=oklog-1.oklog",
								"--peer=oklog-2.oklog",
							},
							VolumeMounts: []corev1.VolumeMount{
								corev1.VolumeMount{
									Name:      "oklog",
									MountPath: "/data",
								},
							},
						},
					},
				},
			},
			VolumeClaimTemplates: []corev1.PersistentVolumeClaim{
				corev1.PersistentVolumeClaim{
					ObjectMeta: getMeta(),
					Spec: corev1.PersistentVolumeClaimSpec{
						AccessModes: []corev1.PersistentVolumeAccessMode{
							corev1.ReadWriteOnce,
						},
						StorageClassName: &storageClassName,
						Resources: corev1.ResourceRequirements{
							Requests: corev1.ResourceList{
								corev1.ResourceStorage: quantity,
							},
						},
					},
				},
			},
		},
	}
}

func createService(si typedcorev1.ServiceInterface) error {
	svc, err := si.Create(getServiceSpec())
	if err != nil {
		log.Errorf("unable to create oklog service %s", err)
		return err
	}
	log.Infof("Created service %+v", svc)
	return nil
}

func getServiceSpec() *corev1.Service {
	return &corev1.Service{
		ObjectMeta: getMeta(),
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{"app": "oklog"},
			Ports: []corev1.ServicePort{
				corev1.ServicePort{
					Name:       "api-default",
					Port:       7650,
					TargetPort: intstr.FromInt(7650),
					Protocol:   corev1.ProtocolTCP,
				},
				corev1.ServicePort{
					Name:       "ingest-fast",
					Port:       7651,
					TargetPort: intstr.FromInt(7651),
					Protocol:   corev1.ProtocolTCP,
				},
				corev1.ServicePort{
					Name:       "ingest-durable",
					Port:       7652,
					TargetPort: intstr.FromInt(7652),
					Protocol:   corev1.ProtocolTCP,
				},
				corev1.ServicePort{
					Name:       "ingest-bulk",
					Port:       7653,
					TargetPort: intstr.FromInt(7653),
					Protocol:   corev1.ProtocolTCP,
				},
				corev1.ServicePort{
					Name:       "cluster",
					Port:       7659,
					TargetPort: intstr.FromInt(7659),
					Protocol:   corev1.ProtocolTCP,
				},
			},
			ClusterIP: "None",
		},
	}
}

func getMeta() metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name:   "oklog",
		Labels: map[string]string{"app": "oklog"},
	}
}
