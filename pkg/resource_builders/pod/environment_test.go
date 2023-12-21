package pod

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/3scale-ops/basereconciler/util"
	saasv1alpha1 "github.com/3scale-ops/saas-operator/api/v1alpha1"
	externalsecretsv1beta1 "github.com/external-secrets/external-secrets/apis/externalsecrets/v1beta1"
	"github.com/google/go-cmp/cmp"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestOptions_BuildEnvironment(t *testing.T) {
	type args struct {
		extra []corev1.EnvVar
	}
	tests := []struct {
		name string
		opts *Options
		args args
		want []corev1.EnvVar
	}{
		{
			name: "Text value",
			opts: func() *Options {
				o := NewOptions()
				o.Unpack("value").IntoEnvvar("envvar")
				return o
			}(),
			args: args{extra: []corev1.EnvVar{}},
			want: []corev1.EnvVar{{
				Name:  "envvar",
				Value: "value",
			}},
		},
		{
			name: "Text value with custom format",
			opts: func() *Options {
				o := NewOptions()
				o.Unpack(8080, ":%d").IntoEnvvar("envvar")
				return o
			}(),
			args: args{extra: []corev1.EnvVar{}},
			want: []corev1.EnvVar{{
				Name:  "envvar",
				Value: ":8080",
			}},
		},
		{
			name: "Pointer to text value",
			opts: func() *Options {
				o := NewOptions()
				o.Unpack(util.Pointer("value")).IntoEnvvar("envvar")
				return o
			}(),
			args: args{extra: []corev1.EnvVar{}},
			want: []corev1.EnvVar{{
				Name:  "envvar",
				Value: "value",
			}},
		},
		{
			name: "Don't panic on nil pointer to text value",
			opts: func() *Options {
				o := NewOptions()
				var v *string
				o.Unpack(v).IntoEnvvar("envvar")
				return o
			}(),
			args: args{extra: []corev1.EnvVar{}},
			want: []corev1.EnvVar{},
		},
		{
			name: "SecretReference",
			opts: func() *Options {
				o := &Options{}
				o.Unpack(saasv1alpha1.SecretReference{FromVault: &saasv1alpha1.VaultSecretReference{
					Path: "path",
					Key:  "key",
				}}).IntoEnvvar("envvar").AsSecretRef("secret")
				return o
			}(),
			args: args{extra: []corev1.EnvVar{}},
			want: []corev1.EnvVar{{
				Name: "envvar",
				ValueFrom: &corev1.EnvVarSource{
					SecretKeyRef: &corev1.SecretKeySelector{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: "secret",
						},
						Key: "envvar",
					},
				},
			}},
		},
		{
			name: "Pointer to SecretReference",
			opts: func() *Options {
				o := &Options{}
				o.Unpack(&saasv1alpha1.SecretReference{FromVault: &saasv1alpha1.VaultSecretReference{
					Path: "path",
					Key:  "key",
				}}).IntoEnvvar("envvar").AsSecretRef("secret")
				return o
			}(),
			args: args{extra: []corev1.EnvVar{}},
			want: []corev1.EnvVar{{
				Name: "envvar",
				ValueFrom: &corev1.EnvVarSource{
					SecretKeyRef: &corev1.SecretKeySelector{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: "secret",
						},
						Key: "envvar",
					},
				},
			}},
		},
		{
			name: "Don't panic on nil pointer to SecretReference",
			opts: func() *Options {
				o := &Options{}
				var v *saasv1alpha1.SecretReference
				o.Unpack(v).IntoEnvvar("envvar").AsSecretRef("secret")
				return o
			}(),
			args: args{extra: []corev1.EnvVar{}},
			want: []corev1.EnvVar{},
		},
		{
			name: "SecretReference with override",
			opts: func() *Options {
				o := &Options{}
				o.Unpack(saasv1alpha1.SecretReference{Override: util.Pointer("value")}).IntoEnvvar("envvar")
				return o
			}(),
			args: args{extra: []corev1.EnvVar{}},
			want: []corev1.EnvVar{{
				Name:  "envvar",
				Value: "value",
			}},
		},
		{
			name: "EmptyIf",
			opts: func() *Options {
				o := &Options{}
				o.Unpack(&saasv1alpha1.SecretReference{FromVault: &saasv1alpha1.VaultSecretReference{
					Path: "path",
					Key:  "key",
				}}).IntoEnvvar("envvar").AsSecretRef("secret").EmptyIf(true)
				return o
			}(),
			args: args{extra: []corev1.EnvVar{}},
			want: []corev1.EnvVar{{
				Name:  "envvar",
				Value: "",
			}},
		},
		{
			name: "Adds/overwrites extra envvars",
			opts: func() *Options {
				o := &Options{}
				o.Unpack(&saasv1alpha1.SecretReference{FromVault: &saasv1alpha1.VaultSecretReference{
					Path: "path",
					Key:  "key",
				}}).IntoEnvvar("envvar1").AsSecretRef("secret").EmptyIf(true)
				o.Unpack("value2").IntoEnvvar("envvar2")
				return o
			}(),
			args: args{extra: []corev1.EnvVar{
				{
					Name:  "envvar1",
					Value: "value1",
				},
				{
					Name:  "envvar3",
					Value: "value3",
				},
			}},
			want: []corev1.EnvVar{
				{
					Name:  "envvar1",
					Value: "value1",
				},
				{
					Name:  "envvar2",
					Value: "value2",
				},
				{
					Name:  "envvar3",
					Value: "value3",
				},
			},
		},
		{
			name: "bool value",
			opts: func() *Options {
				o := NewOptions()
				o.Unpack(true).IntoEnvvar("envvar")
				return o
			}(),
			args: args{extra: []corev1.EnvVar{}},
			want: []corev1.EnvVar{{
				Name:  "envvar",
				Value: "true",
			}},
		},
		{
			name: "Pointer to int value",
			opts: func() *Options {
				o := NewOptions()
				o.Unpack(util.Pointer(100)).IntoEnvvar("envvar")
				return o
			}(),
			args: args{extra: []corev1.EnvVar{}},
			want: []corev1.EnvVar{{
				Name:  "envvar",
				Value: "100",
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.opts.WithExtraEnv(tt.args.extra).BuildEnvironment()
			if diff := cmp.Diff(got, tt.want); len(diff) > 0 {
				t.Errorf("Options.BuildEnvironment() got diff %v", diff)
			}
		})
	}
}

func TestOptions_GenerateExternalSecrets(t *testing.T) {
	type args struct {
		namespace       string
		labels          map[string]string
		secretStoreName string
		secretStoreKind string
		refreshInterval metav1.Duration
	}
	tests := []struct {
		name string
		opts *Options
		args args
		want []client.Object
	}{
		{
			name: "Does not generate any external secret",
			opts: func() *Options {
				o := NewOptions()
				o.Unpack("value1").IntoEnvvar("envvar1")
				o.Unpack("value2").IntoEnvvar("envvar2")
				return o
			}(),
			args: args{},
			want: []client.Object{},
		},
		{
			name: "Generates external secrets for the secret options",
			opts: func() *Options {
				o := NewOptions()
				o.Unpack(&saasv1alpha1.SecretReference{FromVault: &saasv1alpha1.VaultSecretReference{
					Path: "path1",
					Key:  "key1",
				}}).IntoEnvvar("envvar1").AsSecretRef("secret1")
				o.Unpack(&saasv1alpha1.SecretReference{FromVault: &saasv1alpha1.VaultSecretReference{
					Path: "path2",
					Key:  "key2",
				}}).IntoEnvvar("envvar2").AsSecretRef("secret1")
				o.Unpack(&saasv1alpha1.SecretReference{FromVault: &saasv1alpha1.VaultSecretReference{
					Path: "path3",
					Key:  "key3",
				}}).IntoEnvvar("envvar3").AsSecretRef("secret2")
				return o
			}(),
			args: args{
				namespace:       "ns",
				labels:          map[string]string{"label-key": "label-value"},
				secretStoreName: "vault",
				secretStoreKind: "SecretStore",
				refreshInterval: metav1.Duration{Duration: 60 * time.Second},
			},
			want: []client.Object{
				&externalsecretsv1beta1.ExternalSecret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "secret1",
						Namespace: "ns",
						Labels:    map[string]string{"label-key": "label-value"},
					},
					Spec: externalsecretsv1beta1.ExternalSecretSpec{
						SecretStoreRef: externalsecretsv1beta1.SecretStoreRef{
							Name: "vault",
							Kind: "SecretStore",
						},
						Target: externalsecretsv1beta1.ExternalSecretTarget{
							Name:           "secret1",
							CreationPolicy: "Owner",
							DeletionPolicy: "Retain",
						},
						RefreshInterval: util.Pointer(metav1.Duration{Duration: 60 * time.Second}),
						Data: []externalsecretsv1beta1.ExternalSecretData{
							{
								SecretKey: "envvar1",
								RemoteRef: externalsecretsv1beta1.ExternalSecretDataRemoteRef{
									Key:                "path1",
									Property:           "key1",
									ConversionStrategy: "Default",
									DecodingStrategy:   "None",
								},
							},
							{
								SecretKey: "envvar2",
								RemoteRef: externalsecretsv1beta1.ExternalSecretDataRemoteRef{
									Key:                "path2",
									Property:           "key2",
									ConversionStrategy: "Default",
									DecodingStrategy:   "None",
								},
							},
						},
					},
				},
				&externalsecretsv1beta1.ExternalSecret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "secret2",
						Namespace: "ns",
						Labels:    map[string]string{"label-key": "label-value"},
					},
					Spec: externalsecretsv1beta1.ExternalSecretSpec{
						SecretStoreRef: externalsecretsv1beta1.SecretStoreRef{
							Name: "vault",
							Kind: "SecretStore",
						},
						Target: externalsecretsv1beta1.ExternalSecretTarget{
							Name:           "secret2",
							CreationPolicy: "Owner",
							DeletionPolicy: "Retain",
						},
						RefreshInterval: util.Pointer(metav1.Duration{Duration: 60 * time.Second}),
						Data: []externalsecretsv1beta1.ExternalSecretData{
							{
								SecretKey: "envvar3",
								RemoteRef: externalsecretsv1beta1.ExternalSecretDataRemoteRef{
									Key:                "path3",
									Property:           "key3",
									ConversionStrategy: "Default",
									DecodingStrategy:   "None",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Skips secret options with override",
			opts: func() *Options {
				o := NewOptions()
				o.Unpack(&saasv1alpha1.SecretReference{Override: util.Pointer("override")}).IntoEnvvar("envvar1").AsSecretRef("secret")
				return o
			}(),
			args: args{},
			want: []client.Object{},
		},
		{
			name: "Skips secret options with nil value",
			opts: func() *Options {
				o := NewOptions()
				var v *saasv1alpha1.SecretReference
				o.Unpack(v).IntoEnvvar("envvar1").AsSecretRef("secret")
				return o
			}(),
			args: args{},
			want: []client.Object{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			templates := tt.opts.GenerateExternalSecrets(tt.args.namespace, tt.args.labels, tt.args.secretStoreName, tt.args.secretStoreKind, tt.args.refreshInterval)
			got := []client.Object{}
			for _, tplt := range templates {
				es, _ := tplt.Build(context.TODO(), fake.NewClientBuilder().Build(), nil)
				got = append(got, es)
			}
			if diff := cmp.Diff(got, tt.want); len(diff) > 0 {
				t.Errorf("Options.GenerateExternalSecrets() got diff %v", diff)
			}
		})
	}
}

func TestOptions_WithExtraEnv(t *testing.T) {
	type args struct {
		extra []corev1.EnvVar
	}
	tests := []struct {
		name    string
		options *Options
		args    args
		want    *Options
		wantOld *Options
	}{
		{
			name: "",
			options: &Options{{
				value:       util.Pointer("value1"),
				envVariable: "envvar1",
				set:         true,
			}},
			args: args{
				extra: []corev1.EnvVar{
					{Name: "envvar1", Value: "aaaa"},
					{Name: "envvar2", Value: "bbbb"},
				},
			},
			want: &Options{
				{
					value:       util.Pointer("aaaa"),
					envVariable: "envvar1",
					set:         true,
				},
				{
					value:       util.Pointer("bbbb"),
					envVariable: "envvar2",
					set:         true,
				},
			},
			wantOld: &Options{{
				value:       util.Pointer("value1"),
				envVariable: "envvar1",
				set:         true,
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.options.WithExtraEnv(tt.args.extra)
			if diff := cmp.Diff(got, tt.want, cmp.AllowUnexported(Option{})); len(diff) > 0 {
				t.Errorf("Options.WithExtraEnv() got diff %v", diff)
			}
			if diff := cmp.Diff(tt.options, tt.wantOld, cmp.AllowUnexported(Option{})); len(diff) > 0 {
				t.Errorf("Options.WithExtraEnv() gotOld diff %v", diff)
			}
		})
	}
}

func TestOptions_ListSecretResourceNames(t *testing.T) {
	tests := []struct {
		name    string
		options *Options
		want    []string
	}{
		{
			name: "",
			options: func() *Options {
				o := &Options{}
				// ok
				o.Unpack(&saasv1alpha1.SecretReference{FromVault: &saasv1alpha1.VaultSecretReference{}}).IntoEnvvar("envvar1").AsSecretRef("secret1")
				// not ok: not a secret value
				o.Unpack("value").IntoEnvvar("envvar2")
				// not ok: secret value with override
				o.Unpack(&saasv1alpha1.SecretReference{Override: util.Pointer("value")}).IntoEnvvar("envvar3").AsSecretRef("secret2")
				var v *saasv1alpha1.SecretReference
				// not ok: secret value is nil
				o.Unpack(v).IntoEnvvar("envvar1").AsSecretRef("secret3")
				// ok
				o.Unpack(&saasv1alpha1.SecretReference{FromVault: &saasv1alpha1.VaultSecretReference{}}).IntoEnvvar("envvar2").AsSecretRef("secret1")
				// ok
				o.Unpack(&saasv1alpha1.SecretReference{FromVault: &saasv1alpha1.VaultSecretReference{}}).IntoEnvvar("envvar3").AsSecretRef("secret2")
				return o
			}(),
			want: []string{"secret1", "secret2"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.options.ListSecretResourceNames(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Options.ListSecretResourceNames() = %v, want %v", got, tt.want)
			}
		})
	}
}
