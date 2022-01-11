package twemproxyconfig

import (
	"encoding/json"
	"fmt"

	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

type TwemproxyServer struct {
	Address  string
	Priority int
	Name     string
}

func (tserver *TwemproxyServer) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s:%d %s\"", tserver.Address, tserver.Priority, tserver.Name)), nil
}

type TwemproxyConfigServerPool struct {
	Listen             string            `json:"listen"`
	Hash               string            `json:"hash"`
	HashTag            string            `json:"hash_tag"`
	Distribution       string            `json:"distribution"`
	Timeout            int               `json:"timeout"`
	Backlog            int               `json:"backlog"`
	PreConnect         bool              `json:"preconnect"`
	Redis              bool              `json:"redis"`
	AutoEjectHosts     bool              `json:"auto_eject_hosts"`
	ServerFailureLimit int               `json:"server_failure_limit"`
	Servers            []TwemproxyServer `json:"servers"`
}

// configMap returns a function that will return a ConfigMap
// resource when called. This ConfigMap holds the twemproxy config file.
func (gen *Generator) configMap(toYAML bool) func() *corev1.ConfigMap {

	return func() *corev1.ConfigMap {

		config := make(map[string]TwemproxyConfigServerPool, len(gen.Spec.ServerPools))
		for _, pool := range gen.Spec.ServerPools {
			config[pool.Name] = generateServerPool(pool, gen.monitoredShards)
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
}

func generateServerPool(pool saasv1alpha1.TwemproxyServerPool, monitoredShards map[string]TwemproxyServer) TwemproxyConfigServerPool {

	servers := make([]TwemproxyServer, 0, len(pool.Topology))
	for _, s := range pool.Topology {
		srv := monitoredShards[s.PhysicalShard]
		srv.Name = s.ShardName
		servers = append(servers, srv)
	}

	return TwemproxyConfigServerPool{
		// The following parameters cannot be changed
		Hash:           "fnv1a_64",
		HashTag:        "{}",
		Distribution:   "ketama",
		AutoEjectHosts: false,
		Redis:          true,
		// The following parameters could be safely modified or exposed in the CR
		Listen:     pool.BindAddress,
		Backlog:    pool.TCPBacklog,
		PreConnect: pool.PreConnect,
		Timeout:    pool.Timeout,
		// The list of servers is generated from the
		// list fo shards provided by the user in the Backend spec
		Servers: servers,
	}
}
