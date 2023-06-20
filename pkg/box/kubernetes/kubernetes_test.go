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
	resourceOptions := &kubernetes.KubeResource{
		Memory: "512Mi",
		Cpu:    "500m",
	}

	expectedDeployment :=
		`apiVersion: apps/v1
kind: Deployment
metadata:
 creationTimestamp: null
 labels:
   app.kubernetes.io/instance: hckops-my-image
   app.kubernetes.io/managed-by: hckops
   app.kubernetes.io/name: my-container-name
   app.kubernetes.io/version: latest
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
 strategy: {}
 template:
   metadata:
     creationTimestamp: null
     labels:
       app.kubernetes.io/instance: hckops-my-image
       app.kubernetes.io/managed-by: hckops
       app.kubernetes.io/name: my-container-name
       app.kubernetes.io/version: latest
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
 labels:
   app.kubernetes.io/instance: hckops-my-image
   app.kubernetes.io/managed-by: hckops
   app.kubernetes.io/name: my-container-name
   app.kubernetes.io/version: latest
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
 type: ClusterIP
status:
 loadBalancer: {}
`
	actualDeployment, actualService, err := buildSpec(containerName, namespace, template, resourceOptions)
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