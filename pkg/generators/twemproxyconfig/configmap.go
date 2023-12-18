package twemproxyconfig

import (
	"encoding/json"

	saasv1alpha1 "github.com/3scale-ops/saas-operator/api/v1alpha1"
	"github.com/3scale-ops/saas-operator/pkg/resource_builders/twemproxy"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

const (
	HealthPoolName    string = "health"
	HealthBindAddress string = "127.0.0.1:22333"
)

// configMap returns a function that will return a ConfigMap
// resource when called. This ConfigMap holds the twemproxy config file.
func (gen *Generator) configMap(toYAML bool) *corev1.ConfigMap {
	config := make(map[string]twemproxy.ServerPoolConfig, len(gen.Spec.ServerPools)+1)
	for _, pool := range gen.Spec.ServerPools {
		if *pool.Target == saasv1alpha1.Masters {
			config[pool.Name] = twemproxy.GenerateServerPool(pool, gen.masterTargets)
		} else {
			config[pool.Name] = twemproxy.GenerateServerPool(pool, gen.slaverwTargets)
		}
	}

	config[HealthPoolName] = twemproxy.ServerPoolConfig{
		Listen: HealthBindAddress,
		Redis:  true,
		Servers: []twemproxy.Server{{
			Address:  "127.0.0.1:6379",
			Priority: 1,
			Name:     "dummy",
		}},
	}

	var b []byte
	var err error

	if toYAML {
		b, err = yaml.Marshal(config)

	} else {
		b, err = json.Marshal(config)
	}
	if err != nil {
		panic(err)
	}

	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      gen.GetInstanceName(),
			Namespace: gen.GetNamespace(),
			Labels:    gen.GetLabels(),
		},
		Data: map[string]string{
			"nutcracker.yml": string(b),
		},
	}
}
