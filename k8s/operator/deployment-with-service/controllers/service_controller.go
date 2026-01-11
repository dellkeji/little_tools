package controllers

import (
	"context"
	"fmt"
	"reflect"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	appsv1alpha1 "service-operator/api/v1"
)

// ServiceReconciler reconciles a Service object
type ServiceReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=apps.example.com,resources=services,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps.example.com,resources=services/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=apps.example.com,resources=services/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=networking.k8s.io,resources=ingresses,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *ServiceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Fetch the Service instance
	service := &appsv1alpha1.Service{}
	err := r.Get(ctx, req.NamespacedName, service)
	if err != nil {
		if errors.IsNotFound(err) {
			logger.Info("Service resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		logger.Error(err, "Failed to get Service")
		return ctrl.Result{}, err
	}

	// Set default values
	r.setDefaults(service)

	// Reconcile ConfigMap
	if err := r.reconcileConfigMap(ctx, service); err != nil {
		logger.Error(err, "Failed to reconcile ConfigMap")
		return ctrl.Result{}, err
	}

	// Reconcile Deployment
	if err := r.reconcileDeployment(ctx, service); err != nil {
		logger.Error(err, "Failed to reconcile Deployment")
		return ctrl.Result{}, err
	}

	// Reconcile Service
	if err := r.reconcileService(ctx, service); err != nil {
		logger.Error(err, "Failed to reconcile Service")
		return ctrl.Result{}, err
	}

	// Reconcile Ingress
	if err := r.reconcileIngress(ctx, service); err != nil {
		logger.Error(err, "Failed to reconcile Ingress")
		return ctrl.Result{}, err
	}

	// Update status
	if err := r.updateStatus(ctx, service); err != nil {
		logger.Error(err, "Failed to update Service status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *ServiceReconciler) setDefaults(service *appsv1alpha1.Service) {
	if service.Spec.Replicas == nil {
		replicas := int32(1)
		service.Spec.Replicas = &replicas
	}
	if service.Spec.Port == 0 {
		service.Spec.Port = 8080
	}
	if service.Spec.ServiceType == "" {
		service.Spec.ServiceType = "ClusterIP"
	}
	if service.Spec.Ingress != nil && service.Spec.Ingress.Path == "" {
		service.Spec.Ingress.Path = "/"
	}
}

func (r *ServiceReconciler) reconcileConfigMap(ctx context.Context, service *appsv1alpha1.Service) error {
	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      service.Name + "-config",
			Namespace: service.Namespace,
		},
	}

	if len(service.Spec.ConfigData) == 0 {
		// Delete ConfigMap if no config data
		err := r.Delete(ctx, configMap)
		if err != nil && !errors.IsNotFound(err) {
			return err
		}
		return nil
	}

	_, err := controllerutil.CreateOrUpdate(ctx, r.Client, configMap, func() error {
		configMap.Data = service.Spec.ConfigData
		return controllerutil.SetControllerReference(service, configMap, r.Scheme)
	})

	return err
}

func (r *ServiceReconciler) reconcileDeployment(ctx context.Context, service *appsv1alpha1.Service) error {
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      service.Name,
			Namespace: service.Namespace,
		},
	}

	_, err := controllerutil.CreateOrUpdate(ctx, r.Client, deployment, func() error {
		// Set labels
		labels := map[string]string{
			"app":     service.Name,
			"version": "v1",
		}

		deployment.Spec = appsv1.DeploymentSpec{
			Replicas: service.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  service.Name,
							Image: service.Spec.Image,
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: service.Spec.Port,
									Protocol:      corev1.ProtocolTCP,
								},
							},
							Env: r.buildEnvVars(service),
						},
					},
				},
			},
		}

		// Set resource requirements if specified
		if service.Spec.Resources != nil {
			container := &deployment.Spec.Template.Spec.Containers[0]
			container.Resources = corev1.ResourceRequirements{}

			if service.Spec.Resources.Limits != nil {
				container.Resources.Limits = corev1.ResourceList{}
				for k, v := range service.Spec.Resources.Limits {
					container.Resources.Limits[corev1.ResourceName(k)] = resource.MustParse(v)
				}
			}

			if service.Spec.Resources.Requests != nil {
				container.Resources.Requests = corev1.ResourceList{}
				for k, v := range service.Spec.Resources.Requests {
					container.Resources.Requests[corev1.ResourceName(k)] = resource.MustParse(v)
				}
			}
		}

		// Add ConfigMap volume if config data exists
		if len(service.Spec.ConfigData) > 0 {
			deployment.Spec.Template.Spec.Volumes = []corev1.Volume{
				{
					Name: "config",
					VolumeSource: corev1.VolumeSource{
						ConfigMap: &corev1.ConfigMapVolumeSource{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: service.Name + "-config",
							},
						},
					},
				},
			}

			deployment.Spec.Template.Spec.Containers[0].VolumeMounts = []corev1.VolumeMount{
				{
					Name:      "config",
					MountPath: "/etc/config",
					ReadOnly:  true,
				},
			}
		}

		return controllerutil.SetControllerReference(service, deployment, r.Scheme)
	})

	return err
}

func (r *ServiceReconciler) buildEnvVars(service *appsv1alpha1.Service) []corev1.EnvVar {
	var envVars []corev1.EnvVar

	// Add custom environment variables
	for _, env := range service.Spec.Env {
		envVars = append(envVars, corev1.EnvVar{
			Name:  env.Name,
			Value: env.Value,
		})
	}

	// Add default environment variables
	envVars = append(envVars, corev1.EnvVar{
		Name:  "SERVICE_NAME",
		Value: service.Name,
	})

	envVars = append(envVars, corev1.EnvVar{
		Name:  "SERVICE_NAMESPACE",
		Value: service.Namespace,
	})

	return envVars
}

func (r *ServiceReconciler) reconcileService(ctx context.Context, service *appsv1alpha1.Service) error {
	k8sService := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      service.Name,
			Namespace: service.Namespace,
		},
	}

	_, err := controllerutil.CreateOrUpdate(ctx, r.Client, k8sService, func() error {
		labels := map[string]string{
			"app": service.Name,
		}

		k8sService.Spec = corev1.ServiceSpec{
			Selector: labels,
			Type:     corev1.ServiceType(service.Spec.ServiceType),
			Ports: []corev1.ServicePort{
				{
					Port:       service.Spec.Port,
					TargetPort: intstr.FromInt(int(service.Spec.Port)),
					Protocol:   corev1.ProtocolTCP,
				},
			},
		}

		return controllerutil.SetControllerReference(service, k8sService, r.Scheme)
	})

	return err
}

func (r *ServiceReconciler) reconcileIngress(ctx context.Context, service *appsv1alpha1.Service) error {
	ingress := &networkingv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      service.Name,
			Namespace: service.Namespace,
		},
	}

	if service.Spec.Ingress == nil || !service.Spec.Ingress.Enabled {
		// Delete ingress if not enabled
		err := r.Delete(ctx, ingress)
		if err != nil && !errors.IsNotFound(err) {
			return err
		}
		return nil
	}

	_, err := controllerutil.CreateOrUpdate(ctx, r.Client, ingress, func() error {
		// Set annotations
		if ingress.Annotations == nil {
			ingress.Annotations = make(map[string]string)
		}
		for k, v := range service.Spec.Ingress.Annotations {
			ingress.Annotations[k] = v
		}

		pathType := networkingv1.PathTypePrefix
		ingress.Spec = networkingv1.IngressSpec{
			Rules: []networkingv1.IngressRule{
				{
					Host: service.Spec.Ingress.Host,
					IngressRuleValue: networkingv1.IngressRuleValue{
						HTTP: &networkingv1.HTTPIngressRuleValue{
							Paths: []networkingv1.HTTPIngressPath{
								{
									Path:     service.Spec.Ingress.Path,
									PathType: &pathType,
									Backend: networkingv1.IngressBackend{
										Service: &networkingv1.IngressServiceBackend{
											Name: service.Name,
											Port: networkingv1.ServiceBackendPort{
												Number: service.Spec.Port,
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

		// Add TLS if enabled
		if service.Spec.Ingress.TLS != nil && service.Spec.Ingress.TLS.Enabled {
			ingress.Spec.TLS = []networkingv1.IngressTLS{
				{
					Hosts:      []string{service.Spec.Ingress.Host},
					SecretName: service.Spec.Ingress.TLS.SecretName,
				},
			}
		}

		return controllerutil.SetControllerReference(service, ingress, r.Scheme)
	})

	return err
}

func (r *ServiceReconciler) updateStatus(ctx context.Context, service *appsv1alpha1.Service) error {
	// Get deployment status
	deployment := &appsv1.Deployment{}
	err := r.Get(ctx, types.NamespacedName{
		Name:      service.Name,
		Namespace: service.Namespace,
	}, deployment)
	if err != nil {
		return err
	}

	// Update service status
	service.Status.ReadyReplicas = deployment.Status.ReadyReplicas

	if deployment.Status.ReadyReplicas == *service.Spec.Replicas {
		service.Status.Phase = "Ready"
	} else {
		service.Status.Phase = "Pending"
	}

	// Set URL if ingress is enabled
	if service.Spec.Ingress != nil && service.Spec.Ingress.Enabled && service.Spec.Ingress.Host != "" {
		protocol := "http"
		if service.Spec.Ingress.TLS != nil && service.Spec.Ingress.TLS.Enabled {
			protocol = "https"
		}
		service.Status.URL = fmt.Sprintf("%s://%s%s", protocol, service.Spec.Ingress.Host, service.Spec.Ingress.Path)
	}

	// Update conditions
	condition := metav1.Condition{
		Type:               "Ready",
		Status:             metav1.ConditionFalse,
		LastTransitionTime: metav1.Now(),
		Reason:             "DeploymentNotReady",
		Message:            fmt.Sprintf("Deployment has %d/%d ready replicas", deployment.Status.ReadyReplicas, *service.Spec.Replicas),
	}

	if deployment.Status.ReadyReplicas == *service.Spec.Replicas {
		condition.Status = metav1.ConditionTrue
		condition.Reason = "DeploymentReady"
		condition.Message = "All replicas are ready"
	}

	// Update or add condition
	updated := false
	for i, existingCondition := range service.Status.Conditions {
		if existingCondition.Type == condition.Type {
			if !reflect.DeepEqual(existingCondition.Status, condition.Status) ||
				!reflect.DeepEqual(existingCondition.Reason, condition.Reason) ||
				!reflect.DeepEqual(existingCondition.Message, condition.Message) {
				service.Status.Conditions[i] = condition
			}
			updated = true
			break
		}
	}
	if !updated {
		service.Status.Conditions = append(service.Status.Conditions, condition)
	}

	return r.Status().Update(ctx, service)
}

// SetupWithManager sets up the controller with the Manager.
func (r *ServiceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1alpha1.Service{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Owns(&corev1.ConfigMap{}).
		Owns(&networkingv1.Ingress{}).
		Complete(r)
}
