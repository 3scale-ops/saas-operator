package sentinel

import (
	"github.com/3scale/saas-operator/pkg/basereconciler"
	"github.com/MakeNowJust/heredoc"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ConfigMap returns a basereconciler.GeneratorFunction function that will return a ConfigMap
// resource when called
func (gen *Generator) ConfigMap() basereconciler.GeneratorFunction {

	return func() client.Object {

		return &corev1.ConfigMap{
			TypeMeta: metav1.TypeMeta{
				Kind:       "ConfigMap",
				APIVersion: corev1.SchemeGroupVersion.String(),
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      gen.GetComponent() + "-gen-config",
				Namespace: gen.GetNamespace(),
				Labels:    gen.GetLabels(),
			},
			Data: map[string]string{
				"generate-config.sh": heredoc.Doc(`
					echo "dir /redis" >> $1
					echo "port 26379" >> $1
					echo "daemonize no" >> $1
					echo "logfile /dev/stdout" >> $1
					echo "sentinel deny-scripts-reconfig yes" >> $1
					echo "protected-mode no" >> $1
					echo "sentinel announce-ip ${POD_IP}" >> $1
					echo "sentinel announce-port 26379" >> $1
				`),
			},
		}
	}
}
