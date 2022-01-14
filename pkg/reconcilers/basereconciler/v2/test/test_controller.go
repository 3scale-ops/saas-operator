/*
Copyright 2021.

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

package test

import (
	"context"
	"time"

	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	externalsecretsv1alpha1 "github.com/3scale/saas-operator/pkg/apis/externalsecrets/v1alpha1"
	secretsmanagerv1alpha1 "github.com/3scale/saas-operator/pkg/apis/secrets-manager/v1alpha1"
	"github.com/3scale/saas-operator/pkg/generators/common_blocks/marin3r"
	basereconciler "github.com/3scale/saas-operator/pkg/reconcilers/basereconciler/v2"
	"github.com/3scale/saas-operator/pkg/reconcilers/basereconciler/v2/resources"
	"github.com/3scale/saas-operator/pkg/reconcilers/basereconciler/v2/test/api/v1alpha1"
	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/pointer"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// Reconciler reconciles a Test object
// +kubebuilder:object:generate=false
type Reconciler struct {
	basereconciler.Reconciler
	Log logr.Logger
}

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("name", req.Name, "namespace", req.Namespace)

	instance := &v1alpha1.Test{}
	key := types.NamespacedName{Name: req.Name, Namespace: req.Namespace}
	result, err := r.GetInstance(ctx, key, instance, "finalizer.example.com", []func(){}, log)
	if result != nil || err != nil {
		return *result, err
	}

	err = r.ReconcileOwnedResources(ctx, instance, r.GetScheme(), []basereconciler.Resource{
		resources.DeploymentTemplate{
			Template: deployment(req.Namespace, instance.Spec.Marin3r),
			RolloutTriggers: []resources.RolloutTrigger{{
				Name:       "secret",
				SecretName: pointer.String("secret"),
			}},
			EnforceReplicas: true,
			IsEnabled:       true,
		},
		resources.SecretDefinitionTemplate{
			Template:  secretDefinition(req.Namespace),
			IsEnabled: true,
		},
		resources.ExternalSecretTemplate{
			Template:  externalSecret(req.Namespace),
			IsEnabled: true,
		},
		resources.ServiceTemplate{
			Template:  service(req.Namespace, instance.Spec.ServiceAnnotations),
			IsEnabled: true,
		},
	})

	if err != nil {
		log.Error(err, "unable to reconcile owned resources")
		return r.ManageError(ctx, instance, err)
	}

	return r.ManageSuccess(ctx, instance)
}

// SetupWithManager sets up the controller with the Manager.
func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.Test{}).
		Watches(&source.Channel{Source: r.GetStatusChangeChannel()}, &handler.EnqueueRequestForObject{}).
		Watches(&source.Kind{Type: &corev1.Secret{TypeMeta: metav1.TypeMeta{Kind: "Secret"}}},
			r.SecretEventHandler(&v1alpha1.TestList{}, r.Log)).
		Complete(r)
}

func deployment(namespace string, marin3rSpec *saasv1alpha1.Marin3rSidecarSpec) func() *appsv1.Deployment {
	return func() *appsv1.Deployment {
		dep := &appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "deployment",
				Namespace: namespace,
			},
			Spec: appsv1.DeploymentSpec{
				Replicas: pointer.Int32Ptr(1),
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{"selector": "deployment"},
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{"selector": "deployment"},
					},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:      "container",
								Image:     "example.com:latest",
								Resources: corev1.ResourceRequirements{},
							},
						},
					},
				},
			},
		}

		if marin3rSpec != nil && !marin3rSpec.IsDeactivated() {
			dep = marin3r.EnableSidecar(*dep, *marin3rSpec)
		}

		return dep
	}
}

func service(namespace string, annotations map[string]string) func() *corev1.Service {
	return func() *corev1.Service {
		return &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:        "service",
				Namespace:   namespace,
				Annotations: annotations,
			},
			Spec: corev1.ServiceSpec{
				Type:                  corev1.ServiceTypeLoadBalancer,
				ExternalTrafficPolicy: corev1.ServiceExternalTrafficPolicyTypeCluster,
				SessionAffinity:       corev1.ServiceAffinityNone,
				Ports: []corev1.ServicePort{{
					Name: "port", Port: 80, TargetPort: intstr.FromInt(80), Protocol: corev1.ProtocolTCP}},
				Selector: map[string]string{"selector": "deployment"},
			},
		}
	}
}

func secretDefinition(namespace string) func() *secretsmanagerv1alpha1.SecretDefinition {

	return func() *secretsmanagerv1alpha1.SecretDefinition {
		return &secretsmanagerv1alpha1.SecretDefinition{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "secret",
				Namespace: namespace,
			},
			Spec: secretsmanagerv1alpha1.SecretDefinitionSpec{
				Name: "secret",
				Type: "opaque",
				KeysMap: map[string]secretsmanagerv1alpha1.DataSource{
					"KEY": {Key: "vault-key", Path: "vault-path"},
				},
			},
		}
	}
}

func externalSecret(namespace string) func() *externalsecretsv1alpha1.ExternalSecret {

	return func() *externalsecretsv1alpha1.ExternalSecret {
		return &externalsecretsv1alpha1.ExternalSecret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "secret",
				Namespace: namespace,
			},
			Spec: externalsecretsv1alpha1.ExternalSecretSpec{
				SecretStoreRef:  externalsecretsv1alpha1.SecretStoreRef{Name: "vault-mgmt", Kind: "ClusterSecretStore"},
				Target:          externalsecretsv1alpha1.ExternalSecretTarget{Name: "secret"},
				RefreshInterval: &metav1.Duration{Duration: 60 * time.Second},
				Data: []externalsecretsv1alpha1.ExternalSecretData{
					{
						SecretKey: "KEY",
						RemoteRef: externalsecretsv1alpha1.ExternalSecretDataRemoteRef{
							Key:      "vault-path",
							Property: "vault-key",
						},
					},
				},
			},
		}
	}
}
