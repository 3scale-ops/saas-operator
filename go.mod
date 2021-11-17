module github.com/3scale/saas-operator

go 1.16

require (
	github.com/blang/semver v3.5.1+incompatible
	github.com/davecgh/go-spew v1.1.1
	github.com/go-logr/logr v0.4.0
	github.com/goombaio/namegenerator v0.0.0-20181006234301-989e774b106e
	github.com/onsi/ginkgo v1.16.4
	github.com/onsi/gomega v1.13.0
	github.com/openshift/api v3.9.0+incompatible
	github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring v0.42.1
	github.com/redhat-cop/operator-utils v1.2.2
	k8s.io/api v0.21.2
	k8s.io/apimachinery v0.21.2
	k8s.io/client-go v0.21.2
	k8s.io/utils v0.0.0-20210527160623-6fdb442a123b
	sigs.k8s.io/controller-runtime v0.9.2
)
