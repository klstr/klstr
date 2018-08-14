package util

import (
	"testing"

	prometheusopv1 "github.com/coreos/prometheus-operator/pkg/client/monitoring/v1"
	appsv1 "k8s.io/api/apps/v1"
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

func testMultiDecode(t *testing.T) {
	multiYaml := `
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: grafana
  name: grafana
spec:
  replicas: 1
  selector:
    matchLabels:
      app: grafana
  revisionHistoryLimit: 2
  template:
    metadata:
      labels:
        app: grafana
    spec:
      containers:
      - image: grafana/grafana:5.2.2
        name: grafana
        imagePullPolicy: Always
        ports:
        - containerPort: 3000
        env:
          - name: GF_AUTH_BASIC_ENABLED
            value: "false"
          - name: GF_AUTH_ANONYMOUS_ENABLED
            value: "true"
          - name: GF_AUTH_ANONYMOUS_ORG_ROLE
            value: Admin
---
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
	sd := NewSchemaDecoder([]byte(multiYaml))
	objs, err := sd.MultiDecode()
	if err != nil {
		t.Error("error decoding multi yaml ", err)
	}
	_, ok := objs[0].(*appsv1.Deployment)
	if !ok {
		t.Error("first object is not of type deployment ", err)
	}
	_, ok = objs[1].(*corev1.Service)
	if !ok {
		t.Error("second object is not of type service ", err)
	}
}
