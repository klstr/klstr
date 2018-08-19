package klstr

import (
	"errors"
	"fmt"
	"strconv"

	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type DBInstanceRegistration struct {
	Name     string
	DBType   string
	Host     string
	Port     int
	Username string
	Password string
}

func RegisterDBInstance(dbr *DBInstanceRegistration, kubeconfig string) error {
	if kubeconfig == "" {
		return errors.New("Kubeconfig is empty")
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return err
	}

	cs, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}
	ns, err := cs.CoreV1().Namespaces().Get("klstr", metav1.GetOptions{})
	log.Info(ns)
	log.Info(fmt.Sprintf("%+v", err))
	if err != nil {
		createdNs, err := cs.CoreV1().Namespaces().Create(&corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{Name: "klstr"},
		})
		if err != nil {
			return err
		}
		log.Infof("created namespace %v", createdNs)
	}
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("dbi-%s-%s", dbr.DBType, dbr.Name),
		},
		StringData: map[string]string{
			"dbtype":   "postgres",
			"host":     dbr.Host,
			"port":     strconv.Itoa(dbr.Port),
			"username": dbr.Username,
			"password": dbr.Password,
		},
	}
	createdSec, err := cs.CoreV1().Secrets("klstr").Create(secret)
	if err != nil {
		return err
	}
	log.Infof("Registered secret - %+v", createdSec)
	return nil
}
