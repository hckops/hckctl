package client

import (
	"bytes"
	"log"
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/printers"
	"k8s.io/client-go/kubernetes/scheme"

	"github.com/hckops/hckctl/internal/config"
	"github.com/hckops/hckctl/pkg/schema"
	"github.com/stretchr/testify/assert"
)

// https://github.com/kubernetes/client-go/issues/193
// https://medium.com/@harshjniitr/reading-and-writing-k8s-resource-as-yaml-in-golang-81dc8c7ea800

func TestBuildSpec(t *testing.T) {
	namespaceName := "my-namespace"
	containerName := "my-container-name"
	template := &schema.BoxV1{
		Kind: "box/v1",
		Name: "mybox",
		Tags: []string{"my-test"},
		Image: struct {
			Repository string
			Version    string
		}{
			Repository: "hckops/box-mybox",
		},
		Network: struct{ Ports []string }{Ports: []string{
			"aaa:123",
			"bbb:456:789",
		}},
	}
	kubeConfig := &config.KubeConfig{
		Namespace:  "labs",
		ConfigPath: "~/.kube/config",
		Resources: config.KubeResources{
			Memory: "512Mi",
			Cpu:    "500m",
		},
	}

	expectedDeployment :=
		`apiVersion: apps/v1
kind: Deployment
metadata:
  creationTimestamp: null
  labels:
    app.kubernetes.io/instance: hckops-box-mybox
    app.kubernetes.io/managed-by: hckops
    app.kubernetes.io/name: my-container-name
    app.kubernetes.io/version: latest
  name: my-container-name
  namespace: my-namespace
spec:
  replicas: 1
  selector:
    matchLabels:
      app.kubernetes.io/instance: hckops-box-mybox
      app.kubernetes.io/managed-by: hckops
      app.kubernetes.io/name: my-container-name
      app.kubernetes.io/version: latest
  strategy: {}
  template:
    metadata:
      creationTimestamp: null
      labels:
        app.kubernetes.io/instance: hckops-box-mybox
        app.kubernetes.io/managed-by: hckops
        app.kubernetes.io/name: my-container-name
        app.kubernetes.io/version: latest
      name: my-container-name
      namespace: my-namespace
    spec:
      containers:
      - image: hckops/box-mybox:latest
        imagePullPolicy: IfNotPresent
        name: hckops-box-mybox
        ports:
        - containerPort: 123
          name: aaa-svc
          protocol: TCP
        - containerPort: 789
          name: bbb-svc
          protocol: TCP
        resources:
          limits:
            memory: 512Mi
          requests:
            cpu: 500m
            memory: 512Mi
        stdin: true
        tty: true
status: {}
`
	actualDeployment, _ := BuildSpec(namespaceName, containerName, template, kubeConfig)
	actualDeployment.TypeMeta = metav1.TypeMeta{
		Kind: "Deployment", APIVersion: "apps/v1",
	}
	assert.YAMLEqf(t, expectedDeployment, deploymentToYaml(actualDeployment), "deployments are not equal")
}

func yamlToDeployment(data string) *appsv1.Deployment {
	decoder := scheme.Codecs.UniversalDeserializer().Decode
	object, _, err := decoder([]byte(data), nil, nil)
	if err != nil {
		log.Fatalf("yamlToDeployment: %#v\n", err)
	}
	return object.(*appsv1.Deployment)
}

func deploymentToYaml(deployment *appsv1.Deployment) string {
	buffer := new(bytes.Buffer)
	printer := printers.YAMLPrinter{}
	if err := printer.PrintObj(deployment, buffer); err != nil {
		log.Fatalf("deploymentToYaml: %#v\n", err)
	}
	return buffer.String()
}
