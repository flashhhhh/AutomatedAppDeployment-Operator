/*
Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	automationv1 "local.io/automated-app-deployment/api/v1"
)

// AutomatedAppDeploymentReconciler reconciles a AutomatedAppDeployment object
type AutomatedAppDeploymentReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=automation.local.io,resources=automatedappdeployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=automation.local.io,resources=automatedappdeployments/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=automation.local.io,resources=automatedappdeployments/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the AutomatedAppDeployment object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.21.0/pkg/reconcile
func (r *AutomatedAppDeploymentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := logf.FromContext(ctx)

	var automatedAppDeployment automationv1.AutomatedAppDeployment
	if err := r.Get(ctx, req.NamespacedName, &automatedAppDeployment); err != nil {
		logger.Info("AutomatedAppDeployment not found", "name", req.NamespacedName)
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	logger.Info("Reconciling AutomatedAppDeployment", "name", automatedAppDeployment.Name)

	labels := automatedAppDeployment.GetLabels()
	if labels == nil {
		labels = map[string]string{
			"app": automatedAppDeployment.Name,
		}
	}

	deployment := &appsv1.Deployment{}
	err := r.Get(ctx, client.ObjectKey{
		Namespace: automatedAppDeployment.Namespace,
		Name:      automatedAppDeployment.Name,
	}, deployment)

	if k8serrors.IsNotFound(err) {
		deployment = &appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      automatedAppDeployment.Name,
				Namespace: automatedAppDeployment.Namespace,
				Labels:    labels,
			},
			Spec: appsv1.DeploymentSpec{
				Replicas: &automatedAppDeployment.Spec.Replicas,
				Selector: &metav1.LabelSelector{
					MatchLabels: labels,
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: labels,
					},
					Spec: corev1.PodSpec{
						Containers: func() []corev1.Container {
							var containers []corev1.Container

							for _, deploymentSpec := range automatedAppDeployment.Spec.Deployments {
								containers = append(containers, corev1.Container{
									Name:  automatedAppDeployment.Name,
									Image: deploymentSpec.Image,
									Ports: func() []corev1.ContainerPort {
										var containerPorts []corev1.ContainerPort
										for _, port := range deploymentSpec.Ports {
											containerPorts = append(containerPorts, corev1.ContainerPort{
												ContainerPort: port,
											})
										}
										return containerPorts
									}(),
									Env: func() []corev1.EnvVar {
										var envVars []corev1.EnvVar
										for key, value := range deploymentSpec.EnvVars {
											envVars = append(envVars, corev1.EnvVar{
												Name:  key,
												Value: value,
											})
										}
										return envVars
									}(),
								})
							}

							return containers
						}(),
					},
				},
			},
		}

		if err := ctrl.SetControllerReference(&automatedAppDeployment, deployment, r.Scheme); err != nil {
			logger.Error(err, "unable to set controller reference for Deployment")
			return ctrl.Result{}, err
		}
		if err := r.Create(ctx, deployment); err != nil {
			logger.Error(err, "unable to create Deployment", "deployment", deployment.Name)
			return ctrl.Result{}, err
		}
		logger.Info("Created Deployment", "deployment", deployment.Name)
	}

	service := &corev1.Service{}
	err = r.Get(ctx, client.ObjectKey{
		Namespace: automatedAppDeployment.Namespace,
		Name:      automatedAppDeployment.Name,
	}, service)

	if k8serrors.IsNotFound(err) {
		service = &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      automatedAppDeployment.Name,
				Namespace: automatedAppDeployment.Namespace,
				Labels:    labels,
			},
			Spec: corev1.ServiceSpec{
				Selector: labels,
				Ports: func() []corev1.ServicePort {
					var servicePorts []corev1.ServicePort
					for _, deploymentSpec := range automatedAppDeployment.Spec.Deployments {
						for _, port := range deploymentSpec.Ports {
							servicePorts = append(servicePorts, corev1.ServicePort{
								Port: port,
								TargetPort: intstr.IntOrString{
									IntVal: port,
								},
							})
						}
					}
					return servicePorts
				}(),
				Type: corev1.ServiceTypeClusterIP,
			},
		}

		if err := ctrl.SetControllerReference(&automatedAppDeployment, service, r.Scheme); err != nil {
			logger.Error(err, "unable to set controller reference for Service")
			return ctrl.Result{}, err
		}
		if err := r.Create(ctx, service); err != nil {
			logger.Error(err, "unable to create Service", "service", service.Name)
			return ctrl.Result{}, err
		}
		logger.Info("Created Service", "service", service.Name)
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *AutomatedAppDeploymentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&automationv1.AutomatedAppDeployment{}).
		Owns(&automationv1.AutomatedAppDeployment{}).
		Complete(r)
}
