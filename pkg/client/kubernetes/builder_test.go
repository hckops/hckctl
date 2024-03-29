package kubernetes

import (
	"testing"

	"github.com/stretchr/testify/assert"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
        env:
        - name: TTYD_USERNAME
          value: username
        - name: TTYD_PASSWORD
          value: password
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
			Env: []KubeEnv{
				{Key: "TTYD_USERNAME", Value: "username"},
				{Key: "TTYD_PASSWORD", Value: "password"},
			},
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
	assert.YAMLEqf(t, expectedDeployment, ObjectToYaml(actualDeployment), "unexpected deployment")
	assert.YAMLEqf(t, expectedService, ObjectToYaml(actualService), "unexpected service")
}

func TestBuildJob(t *testing.T) {

	expectedJob := `
apiVersion: batch/v1
kind: Job
metadata:
  annotations:
    a.b.c: hello
    x.y.z: world
  creationTimestamp: null
  labels:
    app.kubernetes.io/instance: hckops-my-image
    app.kubernetes.io/managed-by: hckops
    app.kubernetes.io/name: my-task-name
    app.kubernetes.io/version: latest
    com.hckops.schema.kind: task-v1
  name: my-task-name
  namespace: my-namespace
spec:
  backoffLimit: 0
  template:
    metadata:
      creationTimestamp: null
    spec:
      containers:
      - args:
        - cmd1
        - cmd2
        image: hckops/my-image:latest
        imagePullPolicy: IfNotPresent
        name: hckops-my-image
        resources: {}
      restartPolicy: Never
status: {}
`

	namespace := "my-namespace"
	taskName := "my-task-name"
	annotations := map[string]string{
		"a.b.c": "hello",
		"x.y.z": "world",
	}
	extra := map[string]string{
		"com.hckops.schema.kind": "task-v1",
	}
	labels := BuildLabels(taskName, "hckops-my-image", "latest", extra)
	jobOpts := &JobOpts{
		Namespace:   namespace,
		Name:        taskName,
		Annotations: annotations,
		Labels:      labels,
		PodInfo: &PodInfo{
			Namespace:     namespace,
			PodName:       "INVALID_POD_NAME",
			ContainerName: "hckops/my-image",
			ImageName:     "hckops/my-image:latest",
			Arguments:     []string{"cmd1", "cmd2"},
			Env:           []KubeEnv{},
			Resource:      &KubeResource{},
		},
	}

	actualJob := BuildJob(jobOpts)
	// fix model
	actualJob.TypeMeta = metav1.TypeMeta{Kind: "Job", APIVersion: "batch/v1"}

	assert.YAMLEqf(t, expectedJob, ObjectToYaml(actualJob), "unexpected job")
}
