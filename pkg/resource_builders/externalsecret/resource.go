package externalsecret

import (
	externalsecretsv1beta1 "github.com/external-secrets/external-secrets/apis/externalsecrets/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

func New(key types.NamespacedName, labels map[string]string,
	secretStoreName, secretStoreKind string, refreshInterval metav1.Duration,
	data []externalsecretsv1beta1.ExternalSecretData) *externalsecretsv1beta1.ExternalSecret {

	return &externalsecretsv1beta1.ExternalSecret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      key.Name,
			Namespace: key.Namespace,
			Labels:    labels,
		},
		Spec: externalsecretsv1beta1.ExternalSecretSpec{
			SecretStoreRef: externalsecretsv1beta1.SecretStoreRef{
				Name: secretStoreName,
				Kind: secretStoreKind,
			},
			Target: externalsecretsv1beta1.ExternalSecretTarget{
				Name:           key.Name,
				CreationPolicy: "Owner",
				DeletionPolicy: "Retain",
			},
			RefreshInterval: &refreshInterval,
			Data:            data,
		},
	}
}
