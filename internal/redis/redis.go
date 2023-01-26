package redis

import (
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// GenerateRedisSvc used to setup the service resource
func GenerateRedisSvc(name, ns, role string) *v1.Service {
	return &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      generateName(name, role),
			Namespace: ns,
			Labels:    getLabels(name, role),
		},
		Spec: v1.ServiceSpec{
			Selector: getLabels(name, role),
			Ports: []v1.ServicePort{
				{
					Name:       "redis",
					Protocol:   v1.ProtocolTCP,
					TargetPort: intstr.FromString("redis"),
					Port:       6392,
				},
			},
		},
	}
}

// GenerateRedisDeploy used to setup the deployment resource
func GenerateRedisDeploy(name, ns, role string, replicas int, args []string) *appsv1.Deployment {
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      generateName(name, role),
			Namespace: ns,
			Labels:    getLabels(name, role),
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: IntPtr(int32(replicas)),
			Selector: &metav1.LabelSelector{
				MatchLabels: getLabels(name, role),
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: getLabels(name, role),
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name: "redis",
							//TODO Make this configurable
							Image: "redis:6.2.3-alpine",
							Args:  args,
							Ports: []v1.ContainerPort{
								{
									Name:          "redis",
									ContainerPort: 6392,
									Protocol:      v1.ProtocolTCP,
								},
							},
							LivenessProbe: &v1.Probe{
								InitialDelaySeconds: 30,
								TimeoutSeconds:      5,
								ProbeHandler: v1.ProbeHandler{
									Exec: &v1.ExecAction{
										Command: []string{
											"redis-cli",
											"ping",
										},
									},
								},
							},
							ReadinessProbe: &v1.Probe{
								InitialDelaySeconds: 30,
								TimeoutSeconds:      5,
								ProbeHandler: v1.ProbeHandler{
									Exec: &v1.ExecAction{
										Command: []string{
											"redis-cli",
											"ping",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func generateName(name, role string) string {
	return fmt.Sprintf("%v-%v", name, role)
}

func getLabels(name, role string) map[string]string {
	return map[string]string{
		"simple.simple.redis/name": name,
		"simple.simple.redis/role": role,
	}
}

func IntPtr(i int32) *int32 {
	return &i
}
