package sentinel

import (
	basereconciler_types "github.com/3scale/saas-operator/pkg/basereconciler/types"
	"github.com/MakeNowJust/heredoc"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ConfigMap returns a basereconciler_types.GeneratorFunction function that will return a ConfigMap
// resource when called
func (gen *Generator) ConfigMap() basereconciler_types.GeneratorFunction {

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
					if [ ! -f $1 ]; then
						echo "dir /redis" >> $1
						echo "port 26379" >> $1
						echo "daemonize no" >> $1
						echo "logfile /dev/stdout" >> $1
						echo "sentinel deny-scripts-reconfig yes" >> $1
						echo "protected-mode no" >> $1
						echo "sentinel announce-ip ${POD_IP}" >> $1
						echo "sentinel announce-port 26379" >> $1
					else
						sed -i "s/^sentinel announce-ip.*/sentinel announce-ip ${POD_IP}/g" $1
					fi
				`),
			},
		}
	}
}