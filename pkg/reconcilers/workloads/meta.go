package workloads

import (
	"github.com/3scale/saas-operator/pkg/util"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func applyKey(o client.Object, w WithKey) {
	o.SetName(w.GetKey().Name)
	o.SetNamespace(w.GetKey().Namespace)
}

func applyLabels(o client.Object, w WithLabels) {
	o.SetLabels(util.MergeMaps(map[string]string{}, o.GetLabels(), w.GetLabels()))
}
