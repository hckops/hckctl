package kubernetes

import (
	"bytes"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/cli-runtime/pkg/printers"
	"k8s.io/client-go/kubernetes/scheme"

	"github.com/hckops/hckctl/pkg/box/model"
	"github.com/hckops/hckctl/pkg/client/kubernetes"
)

func TestBuildSpec(t *testing.T) {
	template := &model.BoxV1{
		Kind: "box/v1",
		Name: "my-name",
		Tags: []string{"my-tag"},
		Image: struct {
			Repository string
			Version    string
		}{
			Repository: "hckops/my-image",
		},
		Shell: "/bin/bash",
		Network: struct{ Ports []string }{Ports: []string{
			"aaa:123",
			"bbb:456:789",
		}},
	}
	containerName := "my-container-name"
	namespace := "my-namespace"

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
    app.kubernetes.io/name: my-container-name
    app.kubernetes.io/version: latest
    com.hckops.schema.kind: box-v1
  name: my-container-name
  namespace: my-namespace
spec:
  replicas: 1
  selector:
    matchLabels:
      app.kubernetes.io/instance: hckops-my-image
      app.kubernetes.io/managed-by: hckops
      app.kubernetes.io/name: my-container-name
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
        app.kubernetes.io/name: my-container-name
        app.kubernetes.io/version: latest
        com.hckops.schema.kind: box-v1
      name: my-container-name
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
    app.kubernetes.io/name: my-container-name
    app.kubernetes.io/version: latest
    com.hckops.schema.kind: box-v1
  name: my-container-name
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
    app.kubernetes.io/name: my-container-name
    app.kubernetes.io/version: latest
    com.hckops.schema.kind: box-v1
  type: ClusterIP
status:
  loadBalancer: {}
`
	templateOpts := &model.TemplateOptions{
		Template: template,
		Size:     model.ExtraSmall,
		Labels: map[string]string{
			"a.b.c": "hello",
			"x.y.z": "world",
		},
	}
	actualDeployment, actualService, err := buildSpec(containerName, namespace, templateOpts)
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

func TestBoxLabel(t *testing.T) {
	expected := "com.hckops.schema.kind=box-v1"
	assert.Equal(t, expected, boxLabel())
}

// TODO flaky test due to order of map/string
func TestToBoxDetails(t *testing.T) {
	createdTime, _ := time.Parse(time.RFC3339, "2042-12-08T10:30:05.265113665Z")

	deployment := &kubernetes.DeploymentDetails{
		Info: &kubernetes.DeploymentInfo{
			Namespace: "myDeploymentNamespace",
			Name:      "myDeploymentName",
			Healthy:   false,
			PodInfo: &kubernetes.PodInfo{
				Namespace:     "myPodNamespace",
				PodName:       "myPodName",
				ContainerName: "myContainerName",
				Env: map[string]string{
					"MY_KEY_1": "MY_VALUE_1",
					"MY_KEY_2": "MY_VALUE_2",
				},
			},
		},
		Created: createdTime,
		Annotations: map[string]string{
			"com.hckops.template.git":          "true",
			"com.hckops.template.git.url":      "myUrl",
			"com.hckops.template.git.revision": "myRevision",
			"com.hckops.template.git.commit":   "myCommit",
			"com.hckops.template.git.name":     "box/base/arch",
			"com.hckops.template.cache.path":   "/tmp/cache/myUuid",
			"com.hckops.box.size":              "m",
		},
	}
	serviceInfo := &kubernetes.ServiceInfo{
		Namespace: "myServiceNamespace",
		Name:      "myServiceName",
		Ports: []kubernetes.ServicePort{
			{Name: "portName", Port: "remotePort"},
		},
	}
	expected := &model.BoxDetails{
		Info: model.BoxInfo{
			Id:      "myPodName",
			Name:    "myDeploymentName",
			Healthy: false,
		},
		TemplateInfo: &model.BoxTemplateInfo{
			GitTemplate: &model.GitTemplateInfo{
				Url:      "myUrl",
				Revision: "myRevision",
				Commit:   "myCommit",
				Name:     "box/base/arch",
			},
		},
		ProviderInfo: &model.BoxProviderInfo{
			Provider: model.Kubernetes,
			KubeProvider: &model.KubeProviderInfo{
				Namespace: "myDeploymentNamespace",
			},
		},
		Size: model.Medium,
		Env: []model.BoxEnv{
			{Key: "MY_KEY_1", Value: "MY_VALUE_1"},
			{Key: "MY_KEY_2", Value: "MY_VALUE_2"},
		},
		Ports: []model.BoxPort{
			{Alias: "portName", Local: "TODO", Remote: "remotePort", Public: false},
		},
		Created: createdTime,
	}
	result, err := toBoxDetails(deployment, serviceInfo)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}
