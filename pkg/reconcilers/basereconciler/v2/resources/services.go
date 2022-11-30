package resources

import (
	"context"
	"fmt"

	basereconciler "github.com/3scale/saas-operator/pkg/reconcilers/basereconciler/v2"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

var _ basereconciler.Resource = ServiceTemplate{}

// ServiceTemplate has methods to generate and reconcile a Service
type ServiceTemplate struct {
	Template  func() *corev1.Service
	IsEnabled bool
}

// Build returns a Service resource
func (st ServiceTemplate) Build(ctx context.Context, cl client.Client) (client.Object, error) {

	svc := st.Template()

	if err := populateServiceSpecRuntimeValues(ctx, cl, svc); err != nil {
		return nil, err
	}

	return svc.DeepCopy(), nil
}

// Enabled indicates if the resource should be present or not
func (dt ServiceTemplate) Enabled() bool {
	return dt.IsEnabled
}

// ResourceReconciler implements a generic reconciler for Service resources
func (st ServiceTemplate) ResourceReconciler(ctx context.Context, cl client.Client, obj client.Object) error {
	logger := log.FromContext(ctx, "ResourceReconciler", "Service")

	needsUpdate := false
	desired := obj.(*corev1.Service)

	instance := &corev1.Service{}
	err := cl.Get(ctx, types.NamespacedName{Name: desired.GetName(), Namespace: desired.GetNamespace()}, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			err = cl.Create(ctx, desired)
			if err != nil {
				return fmt.Errorf("unable to create object: " + err.Error())
			}
			logger.Info("Resource created")
			return nil
		}
		return err
	}

	/* Reconcile metadata */
	if !equality.Semantic.DeepEqual(instance.GetAnnotations(), desired.GetAnnotations()) {
		instance.ObjectMeta.Annotations = desired.GetAnnotations()
		needsUpdate = true
	}
	if !equality.Semantic.DeepEqual(instance.GetLabels(), desired.GetLabels()) {
		instance.ObjectMeta.Labels = desired.GetLabels()
		needsUpdate = true
	}

	/* Reconcile the ports */
	if !equality.Semantic.DeepEqual(instance.Spec.Ports, desired.Spec.Ports) {
		instance.Spec.Ports = desired.Spec.Ports
		needsUpdate = true
	}

	/* Reconcile label selector */
	if !equality.Semantic.DeepEqual(instance.Spec.Selector, desired.Spec.Selector) {
		instance.Spec.Selector = desired.Spec.Selector
		needsUpdate = true
	}

	if needsUpdate {
		err := cl.Update(ctx, instance)
		if err != nil {
			return err
		}
		logger.Info("Resource updated")
	}

	return nil
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
