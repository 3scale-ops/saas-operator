package resources

import (
	"context"
	"fmt"

	basereconciler "github.com/3scale/saas-operator/pkg/reconcilers/basereconciler/v2"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ basereconciler.Resource = ServiceTemplate{}

type ServiceTemplate struct {
	Template  func() *corev1.Service
	IsEnabled bool
}

func (st ServiceTemplate) Build(ctx context.Context, cl client.Client) (client.Object, []string, error) {

	svc := st.Template()
	svc.GetObjectKind().SetGroupVersionKind(corev1.SchemeGroupVersion.WithKind("Service"))

	if err := populateServiceSpecRuntimeValues(ctx, cl, svc); err != nil {
		return nil, nil, err
	}

	return svc.DeepCopy(), serviceExcludes(svc), nil
}

func (dt ServiceTemplate) Enabled() bool {
	return dt.IsEnabled
}

func populateServiceSpecRuntimeValues(ctx context.Context, cl client.Client, svc *corev1.Service) error {

	instance := &corev1.Service{}
	err := cl.Get(ctx, types.NamespacedName{Name: svc.GetName(), Namespace: svc.GetNamespace()}, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Resource not found, return the template as is
			// because there are not runtime values yet
			return nil
		}
		return err
	}

	// Set runtime values in the resource:
	// "/spec/clusterIP", "/spec/clusterIPs", "/spec/ipFamilies", "/spec/ipFamilyPolicy", "/spec/ports/*/nodePort"
	svc.Spec.ClusterIP = instance.Spec.ClusterIP
	svc.Spec.ClusterIPs = instance.Spec.ClusterIPs
	svc.Spec.IPFamilies = instance.Spec.IPFamilies
	svc.Spec.IPFamilyPolicy = instance.Spec.IPFamilyPolicy

	// For services that are not ClusterIP we need to populate the runtime values
	// of NodePort for each port
	if svc.Spec.Type != corev1.ServiceTypeClusterIP {
		for idx, port := range svc.Spec.Ports {
			runtimePort := findPort(port.Port, port.Protocol, instance.Spec.Ports)
			if runtimePort != nil {
				svc.Spec.Ports[idx].NodePort = runtimePort.NodePort
			}
		}
	}

	return nil
}

func findPort(pNumber int32, pProtocol corev1.Protocol, ports []corev1.ServicePort) *corev1.ServicePort {
	// Ports within a svc are uniquely identified by
	// the "port" and "protocol" fields. This is documented in
	// k8s API reference
	for _, port := range ports {
		if pNumber == port.Port && pProtocol == port.Protocol {
			return &port
		}
	}
	// not found
	return nil
}

// serviceExcludes generates the list of excluded paths for a Service resource
func serviceExcludes(svc *corev1.Service) []string {
	paths := append(DefaultExcludedPaths, "/spec/clusterIP", "/spec/clusterIPs", "/spec/ipFamilies", "/spec/ipFamilyPolicy")
	for idx := range svc.Spec.Ports {
		paths = append(paths, fmt.Sprintf("/spec/ports/%d/nodePort", idx))
	}
	return paths
}
