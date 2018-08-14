package util

import (
	"testing"

	prometheusopv1 "github.com/coreos/prometheus-operator/pkg/client/monitoring/v1"
	corev1 "k8s.io/api/core/v1"
)

func TestDecodeWithNoArg(t *testing.T) {
	serviceYaml := `
apiVersion: v1
kind: Service
metadata:
  name: grafana
spec:
  selector:
    app: grafana
  ports:
  - name: grafana
    port: 3000
    targetPort: 3000
`
	sd := NewSchemaDecoder([]byte(serviceYaml))
	obj, err := sd.Decode()
	if err != nil {
		t.Error("error decoding ", err)
	}
	_, ok := obj.(*corev1.Service)
	if !ok {
		t.Error("object is of wrong type")
	}
}

func TestDecodeWithArg(t *testing.T) {
	serviceYaml := `
apiVersion: v1
kind: Service
metadata:
  name: grafana
spec:
  selector:
    app: grafana
  ports:
  - name: grafana
    port: 3000
    targetPort: 3000
`
	sd := NewSchemaDecoder([]byte(serviceYaml))
	obj := &prometheusopv1.Prometheus{}
	_, err := sd.Decode(obj)
	if err != nil {
		t.Error("error decoding ", err)
	}
}
