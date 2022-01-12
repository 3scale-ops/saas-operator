package resources

import (
	"context"
	"fmt"

	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	"github.com/3scale/saas-operator/pkg/util"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// RolloutTrigger defines a configuration source that should trigger a
// rollout whenever the data within that configuration source changes
type RolloutTrigger struct {
	Name          string
	ConfigMapName *string
	SecretName    *string
}

// GetHash returns the hash of the data container in the RolloutTrigger
// config source
func (rt *RolloutTrigger) GetHash(ctx context.Context, cl client.Client, namespace string) (string, error) {

	if rt.SecretName != nil {
		secret := &corev1.Secret{}
		key := types.NamespacedName{Name: *rt.SecretName, Namespace: namespace}
		if err := cl.Get(ctx, key, secret); err != nil {
			if errors.IsNotFound(err) {
				return "", nil
			}
			return "", err
		}
		return util.Hash(secret.Data), nil

	} else if rt.ConfigMapName != nil {
		cm := &corev1.ConfigMap{}
		key := types.NamespacedName{Name: *rt.ConfigMapName, Namespace: namespace}
		if err := cl.Get(ctx, key, cm); err != nil {
			if errors.IsNotFound(err) {
				return "", nil
			}
			return "", err
		}
		return util.Hash(cm.Data), nil

	} else {
		return "", fmt.Errorf("empty rollout trigger")
	}
}

// GetAnnotationKey returns the annotation key to be used in the Pods that read
// from the config source defined in the RolloutTrigger
func (rt *RolloutTrigger) GetAnnotationKey() string {
	if rt.SecretName != nil {
		return fmt.Sprintf("%s/%s.%s", saasv1alpha1.AnnotationsDomain, rt.Name, "secret-hash")
	}
	return fmt.Sprintf("%s/%s.%s", saasv1alpha1.AnnotationsDomain, rt.Name, "configmap-hash")
}
