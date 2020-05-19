package main

/**
This is a sample of how the pods can be created given some environment variables
This way is simpler than replacing the variables in a templatized yaml since this
leaves some room for extending the project and allows other developer flexibility
in terms of version changes in k8s APIs.
*/

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/json"
)

func main() {
	pod := createPod("dev")
	podBytes, err := json.Marshal(pod)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(podBytes))
	service := createService()
	serviceBytes, err := json.Marshal(service)
	fmt.Println(string(serviceBytes))

}

var (
	ServiceName     = "sample_service"
	LabelName       = "sample_label"
	ServicePort     = int32(3000)
	ServiceNodePort = int32(3080)
)

func createService() *corev1.Service {
	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   ServiceName,
			Labels: map[string]string{"app": LabelName},
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeNodePort,
			Ports: []corev1.ServicePort{
				{
					Name:     "http",
					Port:     ServicePort,
					NodePort: ServiceNodePort,
				},
			},
		},
	}
}

func createPod(environment string) *corev1.Pod {
	return &corev1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-pod",
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "test-container",
					Image: "nginx",
					Env: []corev1.EnvVar{
						{
							Name:  "ENV",
							Value: environment,
						},
					},
				},
			},
		},
	}
}
