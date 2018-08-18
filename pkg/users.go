package klstr

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	certsv1beta1 "k8s.io/api/certificates/v1beta1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	typedcertsv1beta1 "k8s.io/client-go/kubernetes/typed/certificates/v1beta1"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

func NewUser(username, kubeConfig string) error {

	cs, err := getKubeClientSet(kubeConfig)
	if err != nil {
		return err
	}
	csr, err := newCSR(username)
	if err != nil {
		return err
	}

	kubecsr := &certsv1beta1.CertificateSigningRequest{
		ObjectMeta: metav1.ObjectMeta{
			Name:   username,
			Labels: map[string]string{"name": "username"},
		},
		Spec: certsv1beta1.CertificateSigningRequestSpec{
			Request: csr.CSR,
			Groups:  []string{"system:authenticated"},
			Usages:  []certsv1beta1.KeyUsage{certsv1beta1.UsageClientAuth},
		},
	}
	createdCsr, err := cs.CertificatesV1beta1().CertificateSigningRequests().Create(kubecsr)
	if err != nil {
		return err
	}

	log.Infof("Created CSR : %+v", createdCsr)

	createdCsr.Status.Conditions = append(createdCsr.Status.Conditions, certsv1beta1.CertificateSigningRequestCondition{
		Type:           certsv1beta1.CertificateApproved,
		Reason:         "automatically approved by klstr",
		Message:        "This CSR was generated and automatically approved by klstr",
		LastUpdateTime: metav1.Now(),
	})

	log.Info(createdCsr)
	approvedCsr, err := cs.CertificatesV1beta1().CertificateSigningRequests().UpdateApproval(createdCsr)
	if err != nil {
		return err
	}

	waitForIssue(cs.CertificatesV1beta1().CertificateSigningRequests(), username)
	approvedCsr, err = cs.CertificatesV1beta1().CertificateSigningRequests().Get(username, metav1.GetOptions{})
	if err != nil {
		return err
	}
	log.Infof("Approved CSR : %+v", approvedCsr)

	err = createPrivateNS(cs, username)
	if err != nil {
		return err
	}

	config, err := clientcmd.LoadFromFile(kubeConfig)
	if err != nil {
		return err
	}

	clusterName := config.Contexts[config.CurrentContext].Cluster
	cfg := clientcmdapi.NewConfig()
	currentCluster := config.Clusters[clusterName].DeepCopy()
	cfg.Clusters[config.Contexts[config.CurrentContext].Cluster] = currentCluster

	ai := clientcmdapi.NewAuthInfo()
	ai.ClientCertificateData = approvedCsr.Status.Certificate
	ai.ClientKeyData = csr.PrivateKey
	cfg.AuthInfos[username] = ai

	ctxName := fmt.Sprintf("%s@%s", username, clusterName)
	ctx := clientcmdapi.NewContext()
	ctx.Namespace = username
	ctx.Cluster = config.Contexts[config.CurrentContext].Cluster
	ctx.AuthInfo = username
	cfg.Contexts[ctxName] = ctx
	cfg.CurrentContext = ctxName

	return clientcmd.WriteToFile(*cfg, fmt.Sprintf("%s-config.yaml", username))

}

func waitForIssue(ci typedcertsv1beta1.CertificateSigningRequestInterface, certificateName string) {
	ch := make(chan bool)
	go func() {
		for {
			csr, _ := ci.Get(certificateName, metav1.GetOptions{}) // TODO: handle this cleanly. may be setup watcher
			if len(csr.Status.Certificate) == 0 {
				time.Sleep(5 * time.Second)
				continue
			}
			break
		}
		ch <- true
	}()
	<-ch
}

type CSR struct {
	PrivateKey []byte
	CSR        []byte
}

func newCSR(username string) (*CSR, error) {
	var name pkix.Name
	name.CommonName = username

	pkey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}

	tpl := &x509.CertificateRequest{
		Subject:            name,
		SignatureAlgorithm: x509.ECDSAWithSHA256,
	}

	csr, err := x509.CreateCertificateRequest(rand.Reader, tpl, pkey)
	if err != nil {
		return nil, err
	}

	pkeyM, err := x509.MarshalECPrivateKey(pkey)
	if err != nil {
		return nil, err
	}
	pkeyBlock := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: pkeyM,
	}
	csrBlock := &pem.Block{
		Type:  "CERTIFICATE REQUEST",
		Bytes: csr,
	}

	csrBytes := pem.EncodeToMemory(csrBlock)
	pKeyBytes := pem.EncodeToMemory(pkeyBlock)
	return &CSR{CSR: csrBytes, PrivateKey: pKeyBytes}, nil
}

func getKubeClientSet(kubeConfig string) (*kubernetes.Clientset, error) {

	config, err := clientcmd.BuildConfigFromFlags("", kubeConfig)
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(config)

}

func createPrivateNS(cs *kubernetes.Clientset, username string) error {
	ns := corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: username,
		},
	}
	createdNs, err := cs.Core().Namespaces().Create(&ns)
	if err != nil {
		return err
	}
	log.Infof("Created Namespace - %v", createdNs.Name)

	role := rbacv1.Role{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-manager", username),
			Namespace: username,
		},
		Rules: []rbacv1.PolicyRule{
			rbacv1.PolicyRule{
				APIGroups: []string{"*"},
				Verbs:     []string{"*"},
				Resources: []string{"*"},
			},
		},
	}

	newRole, err := cs.RbacV1().Roles(username).Create(&role)
	if err != nil {
		return err
	}
	log.Infof("Created Role %v", newRole)

	rb := rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("%s-manager-binding", username),
		},
		Subjects: []rbacv1.Subject{
			rbacv1.Subject{Kind: "User", Name: username},
		},
		RoleRef: rbacv1.RoleRef{
			Kind: "Role",
			Name: fmt.Sprintf("%s-manager", username),
		},
	}
	newRb, err := cs.RbacV1().RoleBindings(username).Create(&rb)
	if err != nil {
		return err
	}
	log.Infof("Created Role binding - %v", newRb)
	return nil
}
