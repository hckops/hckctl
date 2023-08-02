package kubernetes

import (
	"bytes"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/cli-runtime/pkg/printers"
	"k8s.io/client-go/kubernetes/scheme"
)

func TestBuildResources(t *testing.T) {
	expectedDeployment := `
apiVersion: apps/v1
kind: Deployment
metadata:
  creationTimestamp: null
  annotations:
    a.b.c: hello
    x.y.z: world
  labels:
    app.kubernetes.io/instance: hckops-my-image
    app.kubernetes.io/managed-by: hckops
    app.kubernetes.io/name: my-box-name
    app.kubernetes.io/version: latest
    com.hckops.schema.kind: box-v1
  name: my-box-name
  namespace: my-namespace
spec:
  replicas: 1
  selector:
    matchLabels:
      app.kubernetes.io/instance: hckops-my-image
      app.kubernetes.io/managed-by: hckops
      app.kubernetes.io/name: my-box-name
      app.kubernetes.io/version: latest
      com.hckops.schema.kind: box-v1
  strategy: {}
  template:
    metadata:
      creationTimestamp: null
      annotations:
        a.b.c: hello
        x.y.z: world
      labels:
        app.kubernetes.io/instance: hckops-my-image
        app.kubernetes.io/managed-by: hckops
        app.kubernetes.io/name: my-box-name
        app.kubernetes.io/version: latest
        com.hckops.schema.kind: box-v1
      name: my-box-name
      namespace: my-namespace
    spec:
      containers:
      - image: hckops/my-image:latest
        imagePullPolicy: IfNotPresent
        name: hckops-my-image
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

	expectedService := `
apiVersion: v1
kind: Service
metadata:
  creationTimestamp: null
  annotations:
    a.b.c: hello
    x.y.z: world
  labels:
    app.kubernetes.io/instance: hckops-my-image
    app.kubernetes.io/managed-by: hckops
    app.kubernetes.io/name: my-box-name
    app.kubernetes.io/version: latest
    com.hckops.schema.kind: box-v1
  name: my-box-name
  namespace: my-namespace
spec:
  ports:
  - name: aaa
    port: 123
    protocol: TCP
    targetPort: aaa-svc
  - name: bbb
    port: 789
    protocol: TCP
    targetPort: bbb-svc
  selector:
    app.kubernetes.io/instance: hckops-my-image
    app.kubernetes.io/managed-by: hckops
    app.kubernetes.io/name: my-box-name
    app.kubernetes.io/version: latest
    com.hckops.schema.kind: box-v1
  type: ClusterIP
status:
  loadBalancer: {}
`

	namespace := "my-namespace"
	boxName := "my-box-name"
	annotations := map[string]string{
		"a.b.c": "hello",
		"x.y.z": "world",
	}
	extra := map[string]string{
		"com.hckops.schema.kind": "box-v1",
	}
	labels := BuildLabels(boxName, "hckops-my-image", "latest", extra)
	specOpts := &ResourcesOpts{
		Namespace:   namespace,
		Name:        boxName,
		Annotations: annotations,
		Labels:      labels,
		Ports: []KubePort{
			{Name: "aaa", Port: "123"},
			{Name: "bbb", Port: "789"},
		},
		PodInfo: &PodInfo{
			Namespace:     namespace,
			PodName:       "INVALID_POD_NAME",
			ContainerName: "hckops/my-image",
			ImageName:     "hckops/my-image:latest",
			Env:           nil,
			Resource: &KubeResource{
				Memory: "512Mi",
				Cpu:    "500m",
			},
		},
	}
	actualDeployment, actualService, err := BuildResources(specOpts)
	// fix models
	actualDeployment.TypeMeta = metav1.TypeMeta{Kind: "Deployment", APIVersion: "apps/v1"}
	actualService.TypeMeta = metav1.TypeMeta{Kind: "Service", APIVersion: "v1"}

	assert.NoError(t, err)
	assert.YAMLEqf(t, expectedDeployment, objectToYaml(actualDeployment), "unexpected deployment")
	assert.YAMLEqf(t, expectedService, objectToYaml(actualService), "unexpected service")
}

func objectToYaml(object runtime.Object) string {
	buffer := new(bytes.Buffer)
	printer := printers.YAMLPrinter{}
	if err := printer.PrintObj(object, buffer); err != nil {
		log.Fatalf("objectToYaml: %#v\n", err)
	}
	return buffer.String()
}

// https://github.com/kubernetes/client-go/issues/193
// https://medium.com/@harshjniitr/reading-and-writing-k8s-resource-as-yaml-in-golang-81dc8c7ea800
func yamlToDeployment(data string) *appsv1.Deployment {
	decoder := scheme.Codecs.UniversalDeserializer().Decode
	object, _, err := decoder([]byte(data), nil, nil)
	if err != nil {
		log.Fatalf("yamlToDeployment: %#v\n", err)
	}
	return object.(*appsv1.Deployment)
}
