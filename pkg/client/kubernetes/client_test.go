package kubernetes

import (
	"k8s.io/apimachinery/pkg/api/resource"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestNewDeploymentInfo(t *testing.T) {
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "myDeploymentNamespace",
			Name:      "myDeploymentName",
		},
		Status: appsv1.DeploymentStatus{
			Conditions: []appsv1.DeploymentCondition{
				{Status: corev1.ConditionTrue},
			},
		},
	}
	podInfo := &PodInfo{
		Namespace:     "myPodNamespace",
		PodName:       "myPodName",
		ContainerName: "myContainerName",
		Env: map[string]string{
			"MY_KEY_1": "MY_VALUE_1",
			"MY_KEY_2": "MY_VALUE_2",
		},
	}
	expected := DeploymentInfo{
		Namespace: "myDeploymentNamespace",
		Name:      "myDeploymentName",
		Healthy:   true,
		PodInfo:   podInfo,
	}

	assert.Equal(t, expected, newDeploymentInfo(deployment, podInfo))
}

func TestIsDeploymentHealthy(t *testing.T) {
	statusHealthy := appsv1.DeploymentStatus{
		Conditions: []appsv1.DeploymentCondition{
			{Status: corev1.ConditionTrue},
		},
	}
	assert.True(t, isDeploymentHealthy(statusHealthy))

	statusNotHealthy := appsv1.DeploymentStatus{
		Conditions: []appsv1.DeploymentCondition{
			{Status: corev1.ConditionTrue},
			{Status: corev1.ConditionFalse},
			{Status: corev1.ConditionTrue},
		},
	}
	assert.False(t, isDeploymentHealthy(statusNotHealthy))
}

func TestNewDeploymentDetails(t *testing.T) {
	createdTime, _ := time.Parse(time.RFC3339, "2042-12-08T10:30:05.265113665Z")

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Namespace:         "myDeploymentNamespace",
			Name:              "myDeploymentName",
			CreationTimestamp: metav1.Time{Time: createdTime},
			Annotations: map[string]string{
				"com.hckops.schema.kind": "box/v1",
			},
		},
		Status: appsv1.DeploymentStatus{
			Conditions: []appsv1.DeploymentCondition{
				{Status: corev1.ConditionTrue},
			},
		},
	}
	podInfo := &PodInfo{
		Namespace:     "myPodNamespace",
		PodName:       "myPodName",
		ContainerName: "myContainerName",
		Env: map[string]string{
			"MY_KEY_1": "MY_VALUE_1",
			"MY_KEY_2": "MY_VALUE_2",
		},
	}
	expected := &DeploymentDetails{
		Info: &DeploymentInfo{
			Namespace: "myDeploymentNamespace",
			Name:      "myDeploymentName",
			Healthy:   true,
			PodInfo:   podInfo,
		},
		Created: createdTime,
		Annotations: map[string]string{
			"com.hckops.schema.kind": "box/v1",
		},
	}

	assert.Equal(t, expected, newDeploymentDetails(deployment, podInfo))
}

func TestNewPodInfo(t *testing.T) {
	pods := &corev1.PodList{
		Items: []corev1.Pod{
			{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "myPodNamespace",
					Name:      "myPodName",
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "myContainerName",
							Image: "myImageName",
							Env: []corev1.EnvVar{
								{Name: "MY_KEY_1", Value: "MY_VALUE_1"},
								{Name: "MY_KEY_2", Value: "MY_VALUE_2"},
							},
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									corev1.ResourceMemory: resource.MustParse("512Mi"),
									corev1.ResourceCPU:    resource.MustParse("500m"),
								},
							},
						},
					},
				},
			},
		},
	}
	result, err := newPodInfo("myNamespace", pods)
	expected := &PodInfo{
		Namespace:     "myPodNamespace",
		PodName:       "myPodName",
		ContainerName: "myContainerName",
		ImageName:     "myImageName",
		Env: map[string]string{
			"MY_KEY_1": "MY_VALUE_1",
			"MY_KEY_2": "MY_VALUE_2",
		},
		Resource: &KubeResource{
			Memory: "512Mi",
			Cpu:    "500m",
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestNewPodInfoErrorReplica(t *testing.T) {
	pods := &corev1.PodList{
		Items: []corev1.Pod{
			{ObjectMeta: metav1.ObjectMeta{Name: "pod-1"}},
			{ObjectMeta: metav1.ObjectMeta{Name: "pod-2"}},
		},
	}
	result, err := newPodInfo("myPodNamespace", pods)

	assert.EqualError(t, err, "found 2 pods, expected only 1 pod for deployment: namespace=myPodNamespace")
	assert.Nil(t, result)
}

func TestNewPodInfoErrorContainer(t *testing.T) {
	pods := &corev1.PodList{
		Items: []corev1.Pod{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: "pod-1",
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{Name: "myContainer-1"},
						{Name: "myContainer-2"},
					},
				},
			},
		},
	}
	result, err := newPodInfo("myPodNamespace", pods)

	assert.EqualError(t, err, "found 2 containers, expected only 1 container for pod: namespace=myPodNamespace")
	assert.Nil(t, result)
}

func TestNewServiceInfo(t *testing.T) {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "myServiceNamespace",
			Name:      "myServiceName",
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{Name: "alias-1", Port: 123},
				{Name: "alias-2", Port: 456},
			},
		},
	}
	serviceInfo := &ServiceInfo{
		Namespace: "myServiceNamespace",
		Name:      "myServiceName",
		Ports: []KubePort{
			{Name: "alias-1", Port: "123"},
			{Name: "alias-2", Port: "456"},
		},
	}

	assert.Equal(t, serviceInfo, newServiceInfo(service))
}
