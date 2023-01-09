package auto

import (
	marin3rv1alpha1 "github.com/3scale-ops/marin3r/apis/marin3r/v1alpha1"
	"github.com/3scale-ops/marin3r/pkg/envoy"
	"github.com/3scale/saas-operator/pkg/util"
	envoy_config_listener_v3 "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	envoy_extensions_transport_sockets_tls_v3 "github.com/envoyproxy/go-control-plane/envoy/extensions/transport_sockets/tls/v3"
)

func GenerateSecrets(resources []envoy.Resource) ([]marin3rv1alpha1.EnvoySecretResource, error) {

	refs := []string{}

	for _, res := range resources {

		switch o := res.(type) {

		case *envoy_config_listener_v3.Listener:
			secrets, err := secretRefsFromListener(o)
			if err != nil {
				return nil, err
			}
			refs = append(refs, secrets...)

		}
	}

	secrets := []marin3rv1alpha1.EnvoySecretResource{}
	for _, ref := range util.Unique(refs) {
		secrets = append(secrets, marin3rv1alpha1.EnvoySecretResource{Name: ref})
	}

	return secrets, nil
}

func secretRefsFromListener(listener *envoy_config_listener_v3.Listener) ([]string, error) {

	if listener.FilterChains[0].TransportSocket == nil {
		return nil, nil
	}

	secrets := []string{}
	proto, err := listener.FilterChains[0].TransportSocket.GetTypedConfig().UnmarshalNew()
	if err != nil {
		return nil, err
	}
	tlsContext := proto.(*envoy_extensions_transport_sockets_tls_v3.DownstreamTlsContext)
	for _, sdsConfig := range tlsContext.CommonTlsContext.TlsCertificateSdsSecretConfigs {
		secrets = append(secrets, sdsConfig.Name)
	}

	return util.Unique(secrets), nil
}
