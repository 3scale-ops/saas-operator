package pod

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/3scale-ops/basereconciler/mutators"
	"github.com/3scale-ops/basereconciler/resource"
	"github.com/3scale-ops/basereconciler/util"
	saasv1alpha1 "github.com/3scale-ops/saas-operator/api/v1alpha1"
	"github.com/3scale-ops/saas-operator/pkg/resource_builders/externalsecret"
	operatorutil "github.com/3scale-ops/saas-operator/pkg/util"
	externalsecretsv1beta1 "github.com/external-secrets/external-secrets/apis/externalsecrets/v1beta1"
	clone "github.com/huandu/go-clone/generic"
	"github.com/samber/lo"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

type Option struct {
	value       *string
	valueFrom   *corev1.EnvVarSource
	envVariable string
	secretName  string
	seedKey     string
	vaultKey    string
	vaultPath   string
	isSet       bool
	isEmpty     bool
}

func (o *Option) IntoEnvvar(e string) *Option          { o.envVariable = e; return o }
func (o *Option) AsSecretRef(s fmt.Stringer) *Option   { o.secretName = s.String(); return o }
func (o *Option) WithSeedKey(key fmt.Stringer) *Option { o.seedKey = key.String(); return o }
func (o *Option) EmptyIf(empty bool) *Option {
	if empty {
		o.isEmpty = true
	}
	return o
}

// Unpack retrieves the value specified from the API and adds a matching option to the
// list of options. It handles both values and pointers seamlessly.
// Considers a nil value as an unset option.
// It always unpacks into an string representation of the value so it can be stored as
// an environment variable.
// A parameter indicating the format (as in a call to fmt.Sprintf()) can be optionally passed.
func (opt *Option) Unpack(o any, params ...string) *Option {
	if len(params) > 1 {
		panic(fmt.Errorf("too many params in call to Unpack"))
	}

	if opt.isEmpty {
		opt.isSet = true
		opt.value = util.Pointer("")
		return opt
	}

	var val any

	if reflect.ValueOf(o).Kind() == reflect.Ptr {
		if lo.IsNil(o) {
			// underlying value is nil so option is unset
			return &Option{isSet: false}
		} else {
			val = reflect.ValueOf(o).Elem().Interface()
		}
	} else {
		val = o
	}

	switch v := val.(type) {

	case saasv1alpha1.SecretReference:
		if opt.envVariable == "" {
			panic("AddEnvvar must be invoked to add a new option")
		}
		opt.isSet = true

		// is a secret with override
		if v.Override != nil {
			opt.value = v.Override

			// is a secret with value from vault
		} else if v.FromVault != nil {
			if opt.secretName == "" {
				panic("AsSecretRef must be invoked when using 'SecretReference.FromVault'")
			}
			opt.vaultKey = v.FromVault.Key
			opt.vaultPath = v.FromVault.Path
			opt.valueFrom = &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					Key: opt.envVariable,
					LocalObjectReference: corev1.LocalObjectReference{
						Name: opt.secretName,
					},
				}}

			// is a secret retrieved ffom the default seed Secret
		} else if v.FromSeed != nil {
			if opt.seedKey == "" {
				panic("WithSeedKey must be invoked when using 'SecretReference.FromSeed'")
			}
			opt.valueFrom = &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: saasv1alpha1.DefaultSeedSecret,
					},
					Key: opt.seedKey,
				}}
		}

	default:
		opt.isSet = true
		opt.value = unpackValue(v, params...)
	}

	return opt
}

func unpackValue(o any, params ...string) *string {
	var format string
	if len(params) > 0 {
		format = params[0]
	} else {
		format = "%v"
	}
	return util.Pointer(fmt.Sprintf(format, o))
}

type Options []*Option

func NewOptions() *Options { return &Options{} }

// DeepCopy traveses the struct and returns
// a deep copy of it
func (options *Options) DeepCopy() *Options {
	return clone.Clone(options)
}

// FilterSecretOptions returns a list of options that will generate a Secret resource
func (options *Options) FilterSecretOptions() Options {
	return lo.Filter(*options, func(item *Option, index int) bool {
		return item.valueFrom != nil && item.valueFrom.SecretKeyRef != nil
	})
}

// FilterSecretOptions returns a list of options that will generate a Secret resource
// with a Vault secret store as its source (via an ExternalSecret)
func (options *Options) FilterFromVaultOptions() Options {
	return lo.Filter(*options, func(item *Option, index int) bool {
		return item.vaultKey != "" && item.vaultPath != ""
	})
}

func (options *Options) ListSecretResourceNames() []string {
	list := lo.Reduce(options.FilterSecretOptions(), func(agg []string, item *Option, _ int) []string {
		return append(agg, item.valueFrom.SecretKeyRef.Name)
	}, []string{})

	return lo.Uniq(list)
}

func (options *Options) GenerateRolloutTriggers(additionalSecrets ...string) []resource.TemplateMutationFunction {
	secrets := options.ListSecretResourceNames()
	triggers := make([]resource.TemplateMutationFunction, 0, len(secrets))
	for _, secret := range append(secrets, additionalSecrets...) {
		triggers = append(
			triggers,
			mutators.RolloutTrigger{Name: secret, SecretName: util.Pointer(secret)}.Add(),
		)
	}
	return triggers
}

func (options *Options) AddEnvvar(e string) *Option {
	opt := &Option{envVariable: e}
	*options = append(*options, opt)
	return opt
}

// WithExtraEnv returns a copy of the Options list with the added extra envvars
func (options *Options) WithExtraEnv(extra []corev1.EnvVar) *Options {

	out := options.DeepCopy()
	for _, envvar := range extra {
		o, exists := lo.Find(*out, func(o *Option) bool {
			return o.envVariable == envvar.Name
		})

		if exists {
			o.value = util.Pointer(envvar.Value)
			o.valueFrom = envvar.ValueFrom
			o.isSet = true
			o.secretName = ""
		} else {
			var v *string
			if envvar.Value != "" {
				v = util.Pointer(envvar.Value)
			}
			*out = append(*out, &Option{
				value:       v,
				valueFrom:   envvar.ValueFrom,
				envVariable: envvar.Name,
				isSet:       true,
			})
		}
	}
	return out
}

// BuildEnvironment generates a list of corev1.Envvar that matches the
// list of options
func (opts *Options) BuildEnvironment() []corev1.EnvVar {

	env := []corev1.EnvVar{}
	for _, opt := range *opts {

		if !opt.isSet {
			continue
		}

		// Direct value (if both value and valueFrom are set, value takes precedence and
		// valueFrom will be ignored)
		if opt.value != nil {
			env = append(env, corev1.EnvVar{
				Name:  opt.envVariable,
				Value: *opt.value,
			})
			continue

		}

		// ValueFrom
		if opt.valueFrom != nil {
			env = append(env, corev1.EnvVar{
				Name:      opt.envVariable,
				ValueFrom: opt.valueFrom,
			})
			continue
		}
	}

	return env
}

// GenerateExternalSecrets generates the external secret templates required to match the list of options
func (opts *Options) GenerateExternalSecrets(namespace string, labels map[string]string, secretStoreName, secretStoreKind string, refreshInterval metav1.Duration) []resource.TemplateInterface {
	list := []resource.TemplateInterface{}

	for _, group := range lo.PartitionBy(opts.FilterFromVaultOptions(), func(item *Option) string { return item.secretName }) {
		data := []externalsecretsv1beta1.ExternalSecretData{}
		name := group[0].secretName
		for _, opt := range group {
			data = append(data, externalsecretsv1beta1.ExternalSecretData{
				SecretKey: opt.envVariable,
				RemoteRef: externalsecretsv1beta1.ExternalSecretDataRemoteRef{
					Key:                strings.TrimPrefix(opt.vaultPath, "secret/data/"),
					Property:           opt.vaultKey,
					ConversionStrategy: "Default",
					DecodingStrategy:   "None",
				},
			})
		}
		list = append(list, resource.NewTemplateFromObjectFunction(
			func() *externalsecretsv1beta1.ExternalSecret {
				return externalsecret.New(types.NamespacedName{Name: name, Namespace: namespace}, labels, secretStoreName, secretStoreKind, refreshInterval, data)
			}))
	}
	return list
}

func Union(lists ...[]*Option) *Options {
	all := operatorutil.ConcatSlices(lists...)
	all = lo.UniqBy(all, func(item *Option) string {
		return item.envVariable
	})
	return util.Pointer[Options](all)
}
