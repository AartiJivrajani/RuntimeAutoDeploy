package generateK8S

/**
Pods can be created given some environment variables/config read from a file.
This way is simpler than replacing the variables in a templatized yaml since this
leaves some room for extending the project and allows other developer flexibility
in terms of version changes in k8s APIs.
*/

import (
	"RuntimeAutoDeploy/common"
	"RuntimeAutoDeploy/config"
	"context"
	"flag"
	"fmt"
	"path/filepath"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"k8s.io/client-go/util/homedir"

	apiv1 "k8s.io/api/core/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	ClientSet   *kubernetes.Clientset
	ClientReady = false
	ClientLock  sync.Mutex
)

func GetK8sClient(ctx context.Context) error {
	ClientLock.Lock()
	defer ClientLock.Unlock()

	var kubeconfig *string
	common.AddToStatusList(ctx.Value(common.TRACE_ID).(string),
		fmt.Sprintf(common.STAGE_FORMAT,
			common.GetTimestampFormat(fmt.Sprintf("%d", time.Now().Unix()), "", ""),
			common.STAGE_STATUS_WIP,
			common.STAGE_K8S_BOOTSTRAP), true)

	if !ClientReady {
		if home := homedir.HomeDir(); home != "" {
			kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute "+
				"path to the kubeconfig file")
		} else {
			kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
		}
		flag.Parse()

		k8sConfig, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			common.AddToStatusList(ctx.Value(common.TRACE_ID).(string),
				fmt.Sprintf(common.STAGE_ERROR_FORMAT,
					common.GetTimestampFormat(fmt.Sprintf("%d", time.Now().Unix()), "", ""),
					common.STAGE_STATUS_ERROR,
					common.STAGE_K8S_BOOTSTRAP,
					"error starting k8s client", err.Error()), false)
			return err
		}
		ClientSet, err = kubernetes.NewForConfig(k8sConfig)
		if err != nil {
			common.AddToStatusList(ctx.Value(common.TRACE_ID).(string),
				fmt.Sprintf(common.STAGE_ERROR_FORMAT,
					common.GetTimestampFormat(fmt.Sprintf("%d", time.Now().Unix()), "", ""),
					common.STAGE_STATUS_ERROR,
					common.STAGE_K8S_BOOTSTRAP,
					"error starting k8s client", err.Error()), false)
			return err
		}
		ClientReady = true
	}
	common.AddToStatusList(ctx.Value(common.TRACE_ID).(string),
		fmt.Sprintf(common.STAGE_FORMAT,
			common.GetTimestampFormat(fmt.Sprintf("%d", time.Now().Unix()), "", ""),
			common.STAGE_STATUS_DONE,
			common.STAGE_K8S_BOOTSTRAP), false)

	return nil
}

func CreateService(ctx context.Context, conf *config.Application) error {
	serviceName := fmt.Sprintf("%s-%s", conf.AppName, "svc")

	common.AddToStatusList(ctx.Value(common.TRACE_ID).(string),
		fmt.Sprintf(common.STAGE_FORMAT,
			common.GetTimestampFormat(fmt.Sprintf("%d", time.Now().Unix()), "", ""),
			common.STAGE_STATUS_WIP,
			fmt.Sprintf(common.STAGE_CREATING_SERVICE, conf.AppName)), false)

	serviceClient := ClientSet.CoreV1().Services(corev1.NamespaceDefault)

	service := &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   serviceName,
			Labels: map[string]string{"app": conf.AppName},
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeNodePort,
			Ports: []corev1.ServicePort{
				{
					Name: "http",
					Port: int32(conf.Port),
					//NodePort: int32(conf.Port),
				},
			},
		},
	}
	log.WithFields(log.Fields{
		"name": serviceName,
	}).Info("Creating service")

	result, err := serviceClient.Create(context.TODO(), service, metav1.CreateOptions{})
	if err != nil {
		common.AddToStatusList(ctx.Value(common.TRACE_ID).(string),
			fmt.Sprintf(common.STAGE_ERROR_FORMAT,
				common.GetTimestampFormat(fmt.Sprintf("%d", time.Now().Unix()), "", ""),
				common.STAGE_STATUS_ERROR,
				fmt.Sprintf(common.STAGE_CREATING_SERVICE, serviceName),
				fmt.Sprintf("%s-%s", "error creating k8s service", err.Error())), false)
		return err
	}
	log.WithFields(log.Fields{
		"name": result.GetObjectMeta().GetName(),
	}).Info("Created service")

	common.AddToStatusList(ctx.Value(common.TRACE_ID).(string),
		fmt.Sprintf(common.STAGE_FORMAT,
			common.GetTimestampFormat(fmt.Sprintf("%d", time.Now().Unix()), "", ""),
			common.STAGE_STATUS_DONE,
			fmt.Sprintf(common.STAGE_CREATING_SERVICE, serviceName)), false)
	return nil
}

func int32Ptr(i int32) *int32 { return &i }

func CreateDeployment(ctx context.Context, conf *config.Application) error {

	common.AddToStatusList(ctx.Value(common.TRACE_ID).(string),
		fmt.Sprintf(common.STAGE_FORMAT,
			common.GetTimestampFormat(fmt.Sprintf("%d", time.Now().Unix()), "", ""),
			common.STAGE_STATUS_WIP,
			fmt.Sprintf(common.STAGE_CREATING_DEPLOYMENT, conf.AppName)), false)

	deploymentClient := ClientSet.AppsV1().Deployments(apiv1.NamespaceDefault)

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("%s-%s", conf.AppName, "deployment"),
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(int32(conf.ReplicaCount)),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": conf.AppName,
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": conf.AppName,
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  conf.AppName,
							Image: fmt.Sprintf("%s:%s", config.UserConfig.Reg.Address, conf.AppName),
							Ports: []apiv1.ContainerPort{
								{
									Name:          "http",
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: int32(conf.Port), //TODO: Make this configurable
								},
							},
							//Resources: apiv1.ResourceRequirements{
							//Limits: apiv1.ResourceList{
							//	"cpu":    resource.MustParse("1"),
							//	"memory": resource.MustParse("100Mi"),
							//},
							//Requests: apiv1.ResourceList{
							//	"cpu":    resource.MustParse("0.5"),
							//	"memory": resource.MustParse("100Mi"),
							//},
							//},
						},
					},
				},
			},
		},
	}
	log.WithFields(log.Fields{
		"name": fmt.Sprintf("%s-%s", conf.AppName, "deployment"),
	}).Info("Creating deployment")

	result, err := deploymentClient.Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		common.AddToStatusList(ctx.Value(common.TRACE_ID).(string),
			fmt.Sprintf(common.STAGE_ERROR_FORMAT,
				common.GetTimestampFormat(fmt.Sprintf("%d", time.Now().Unix()), "", ""),
				common.STAGE_STATUS_ERROR,
				fmt.Sprintf(common.STAGE_CREATING_DEPLOYMENT, conf.AppName),
				fmt.Sprintf("%s-%s", "error creating k8s deployment", err.Error())), false)
		return err
	}
	log.WithFields(log.Fields{
		"name": result.GetObjectMeta().GetName(),
	}).Info("Created deployment")

	common.AddToStatusList(ctx.Value(common.TRACE_ID).(string),
		fmt.Sprintf(common.STAGE_FORMAT,
			common.GetTimestampFormat(fmt.Sprintf("%d", time.Now().Unix()), "", ""),
			common.STAGE_STATUS_DONE,
			fmt.Sprintf(common.STAGE_CREATING_DEPLOYMENT, conf.AppName)), false)
	return nil
}
