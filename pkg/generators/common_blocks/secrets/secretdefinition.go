package secrets

import (
	"bytes"
	"fmt"
	"strings"

	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	secretsmanagerv1alpha1 "github.com/3scale/saas-operator/pkg/apis/secrets-manager/v1alpha1"
	"github.com/3scale/saas-operator/pkg/basereconciler"
	json "github.com/exponent-io/jsonpath"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// SecretConfiguration defines a configuration option that is stored as a
// key in a Secret resource
type SecretConfiguration struct {
	SecretName    string
	ConfigOptions map[string]string
}

// GenerateSecretDefinitionFn generates a SecretDefinition
func (sc *SecretConfiguration) GenerateSecretDefinitionFn(namespace string, labels map[string]string,
	basePath string, serializedConfig []byte) basereconciler.GeneratorFunction {

	return func() client.Object {
		return &secretsmanagerv1alpha1.SecretDefinition{
			TypeMeta: metav1.TypeMeta{
				Kind:       "SecretDefinition",
				APIVersion: secretsmanagerv1alpha1.GroupVersion.String(),
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      sc.SecretName,
				Namespace: namespace,
				Labels:    labels,
			},
			Spec: secretsmanagerv1alpha1.SecretDefinitionSpec{
				Name: sc.SecretName,
				Type: "opaque",
				KeysMap: func() map[string]secretsmanagerv1alpha1.DataSource {
					km, err := sc.keysMap(basePath, serializedConfig)
					if err != nil {
						// This is a code error, so panic is ok
						panic(err)
					}
					return km
				}(),
			},
		}
	}
}

func (sc *SecretConfiguration) keysMap(basePath string, serializedConfig []byte) (map[string]secretsmanagerv1alpha1.DataSource, error) {
	dsm := map[string]secretsmanagerv1alpha1.DataSource{}
	w := json.NewDecoder(bytes.NewReader([]byte(serializedConfig)))

	for configOption, path := range sc.ConfigOptions {
		vsr := &saasv1alpha1.VaultSecretReference{}
		relativePath := strings.TrimPrefix(path, basePath+"/")
		pathElems := func() []interface{} {
			pe := []interface{}{}
			for _, elem := range strings.Split(relativePath, "/") {
				pe = append(pe, elem)
			}
			pe = append(pe, "fromVault")
			return pe
		}()
		ok, err := w.SeekTo(pathElems...)
		if !ok || err != nil {
			return nil, fmt.Errorf("Couldn't find any VaultSecretReference under key '%s'", path)
		}
		w.Decode(vsr)
		dsm[configOption] = secretsmanagerv1alpha1.DataSource{
			Key:  vsr.Key,
			Path: vsr.Path,
		}
	}

	return dsm, nil
}

// SecretConfigurations is an slice of SecretConfiguration
type SecretConfigurations []SecretConfiguration

// LookupSecretName returns the name of the Secret resource
// where the configuration option is stored
func (sc *SecretConfigurations) LookupSecretName(config string) string {
	for _, secretConfig := range *sc {
		for name := range secretConfig.ConfigOptions {
			if name == config {
				return secretConfig.SecretName
			}
		}
	}
	panic("not found")
}

// LookupSecretConfiguration returns the SecretConfiguration resource
// that matches a given SecretName
func (sc *SecretConfigurations) LookupSecretConfiguration(secretName string) SecretConfiguration {
	for _, secretConfig := range *sc {
		if secretConfig.SecretName == secretName {
			return secretConfig
		}

	}
	panic("not found")
}
