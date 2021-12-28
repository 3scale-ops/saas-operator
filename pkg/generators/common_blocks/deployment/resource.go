package deployment

import (
	"fmt"

	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	basereconciler_types "github.com/3scale/saas-operator/pkg/basereconciler/types"
	"github.com/3scale/saas-operator/pkg/util"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	TrafficSwitch = fmt.Sprintf("%s/traffic", saasv1alpha1.GroupVersion.Group)
)

// New returns a basereconciler_types.GeneratorFunction function that will return a Service
// resource when called
func New(key types.NamespacedName, labels map[string]string, selector map[string]string, trafficSelector map[string]string,
	fn basereconciler_types.GeneratorFunction) basereconciler_types.GeneratorFunction {

	return func() client.Object {

		dep := fn().(*appsv1.Deployment)

		// Set the Pod selector
		dep.Spec.Selector = &metav1.LabelSelector{MatchLabels: selector}

		// Set the Pod labels including the traffic selector
		dep.Spec.Template.ObjectMeta.Labels = util.MergeMaps(map[string]string{}, labels, selector, trafficSelector)

		return &appsv1.Deployment{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Deployment",
				APIVersion: appsv1.SchemeGroupVersion.String(),
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:        key.Name,
				Namespace:   key.Namespace,
				Labels:      labels,
				Annotations: dep.GetAnnotations(),
			},
			Spec: dep.Spec,
		}
	}
}
