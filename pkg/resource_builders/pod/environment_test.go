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

type TSecret string

func (s TSecret) String() string { return string(s) }

type TSeedKey string

func (s TSeedKey) String() string { return string(s) }

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
				}}).IntoEnvvar("envvar").AsSecretRef(TSecret("secret"))
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
				}}).IntoEnvvar("envvar").AsSecretRef(TSecret("secret"))
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
				o.Unpack(v).IntoEnvvar("envvar").AsSecretRef(TSecret("secret"))
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
				}}).IntoEnvvar("envvar").AsSecretRef(TSecret("secret")).EmptyIf(true)
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
				}}).IntoEnvvar("envvar1").AsSecretRef(TSecret("secret")).EmptyIf(true)
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
				{
					Name: "envvar4", ValueFrom: &corev1.EnvVarSource{FieldRef: &corev1.ObjectFieldSelector{
						APIVersion: "v1",
						FieldPath:  "metadata.name",
					}},
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
				{
					Name: "envvar4", ValueFrom: &corev1.EnvVarSource{FieldRef: &corev1.ObjectFieldSelector{
						APIVersion: "v1",
						FieldPath:  "metadata.name",
					}},
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
		{
			name: "SecretReference from seed",
			opts: func() *Options {
				o := &Options{}
				o.Unpack(saasv1alpha1.SecretReference{Override: util.Pointer("value1")}).IntoEnvvar("envvar1")
				o.Unpack(saasv1alpha1.SecretReference{FromSeed: &saasv1alpha1.SeedSecretReference{}}).IntoEnvvar("envvar2").
					AsSecretRef(TSecret("some-secret")).
					WithSeedKey(TSeedKey("seed-key"))
				return o
			}(),
			args: args{extra: []corev1.EnvVar{}},
			want: []corev1.EnvVar{
				{
					Name:  "envvar1",
					Value: "value1",
				},
				{
					Name: "envvar2",
					ValueFrom: &corev1.EnvVarSource{
						SecretKeyRef: &corev1.SecretKeySelector{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: saasv1alpha1.DefaultSeedSecret,
							},
							Key: "seed-key",
						},
					},
				},
			},
		},
		{
			name: "SecretReference from vault, but with seed configured",
			opts: func() *Options {
				o := &Options{}
				o.Unpack(saasv1alpha1.SecretReference{Override: util.Pointer("value1")}).IntoEnvvar("envvar1")
				o.Unpack(saasv1alpha1.SecretReference{FromVault: &saasv1alpha1.VaultSecretReference{
					Path: "path",
					Key:  "key",
				}}).IntoEnvvar("envvar2").
					AsSecretRef(TSecret("some-secret")).
					WithSeedKey(TSeedKey("seed-key"))
				return o
			}(),
			args: args{extra: []corev1.EnvVar{}},
			want: []corev1.EnvVar{
				{
					Name:  "envvar1",
					Value: "value1",
				},
				{
					Name: "envvar2",
					ValueFrom: &corev1.EnvVarSource{
						SecretKeyRef: &corev1.SecretKeySelector{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: "some-secret",
							},
							Key: "envvar2",
						},
					},
				},
			},
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
				}}).IntoEnvvar("envvar1").AsSecretRef(TSecret("secret1"))
				o.Unpack(&saasv1alpha1.SecretReference{FromVault: &saasv1alpha1.VaultSecretReference{
					Path: "path2",
					Key:  "key2",
				}}).IntoEnvvar("envvar2").AsSecretRef(TSecret("secret1"))
				o.Unpack(&saasv1alpha1.SecretReference{FromVault: &saasv1alpha1.VaultSecretReference{
					Path: "path3",
					Key:  "key3",
				}}).IntoEnvvar("envvar3").AsSecretRef(TSecret("secret2"))
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
				o.Unpack(&saasv1alpha1.SecretReference{Override: util.Pointer("override")}).IntoEnvvar("envvar1").AsSecretRef(TSecret("secret"))
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
				o.Unpack(v).IntoEnvvar("envvar1").AsSecretRef(TSecret("secret"))
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
			options: &Options{
				{
					value:       util.Pointer("value1"),
					envVariable: "envvar1",
					isSet:       true,
				},
				{
					value:       util.Pointer("value2"),
					envVariable: "envvar2",
					isSet:       true,
				},
			},
			args: args{
				extra: []corev1.EnvVar{
					{Name: "envvar1", Value: "aaaa"},
					{Name: "envvar3", Value: "bbbb"},
					{Name: "envvar4", ValueFrom: &corev1.EnvVarSource{FieldRef: &corev1.ObjectFieldSelector{
						APIVersion: "v1",
						FieldPath:  "metadata.name",
					}}},
				},
			},
			want: &Options{
				{
					value:       util.Pointer("aaaa"),
					envVariable: "envvar1",
					isSet:       true,
				},
				{
					value:       util.Pointer("value2"),
					envVariable: "envvar2",
					isSet:       true,
				},
				{
					value:       util.Pointer("bbbb"),
					envVariable: "envvar3",
					isSet:       true,
				},
				{
					valueFrom:   &corev1.EnvVarSource{FieldRef: &corev1.ObjectFieldSelector{APIVersion: "v1", FieldPath: "metadata.name"}},
					envVariable: "envvar4",
					isSet:       true,
				},
			},
			wantOld: &Options{
				{
					value:       util.Pointer("value1"),
					envVariable: "envvar1",
					isSet:       true,
				},
				{
					value:       util.Pointer("value2"),
					envVariable: "envvar2",
					isSet:       true,
				},
			},
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
				o.Unpack(&saasv1alpha1.SecretReference{FromVault: &saasv1alpha1.VaultSecretReference{}}).IntoEnvvar("envvar1").AsSecretRef(TSecret("secret1"))
				// not ok: not a secret value
				o.Unpack("value").IntoEnvvar("envvar2")
				// not ok: secret value with override
				o.Unpack(&saasv1alpha1.SecretReference{Override: util.Pointer("value")}).IntoEnvvar("envvar3").AsSecretRef(TSecret("secret2"))
				var v *saasv1alpha1.SecretReference
				// not ok: secret value is nil
				o.Unpack(v).IntoEnvvar("envvar1").AsSecretRef(TSecret("secret3"))
				// ok
				o.Unpack(&saasv1alpha1.SecretReference{FromVault: &saasv1alpha1.VaultSecretReference{}}).IntoEnvvar("envvar2").AsSecretRef(TSecret("secret1"))
				// ok
				o.Unpack(&saasv1alpha1.SecretReference{FromVault: &saasv1alpha1.VaultSecretReference{}}).IntoEnvvar("envvar3").AsSecretRef(TSecret("secret2"))
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

func TestUnion(t *testing.T) {
	type args struct {
		lists [][]*Option
	}
	tests := []struct {
		name string
		args args
		want Options
	}{
		{
			name: "",
			args: args{
				lists: [][]*Option{
					{
						{
							value:       util.Pointer("value1"),
							envVariable: "ENVVAR1",
							isSet:       false,
						},
						{
							value:       util.Pointer("value2"),
							envVariable: "ENVVAR2",
							isSet:       false,
						},
					},
					{
						{
							value:       util.Pointer("value1"),
							envVariable: "ENVVAR1",
							isSet:       false,
						},
						{
							value:       util.Pointer("value3"),
							envVariable: "ENVVAR3",
							isSet:       false,
						},
						{
							value:       util.Pointer("value4"),
							envVariable: "ENVVAR4",
							isSet:       false,
						},
					},
				},
			},
			want: []*Option{
				{
					value:       util.Pointer("value1"),
					envVariable: "ENVVAR1",
					isSet:       false,
				},
				{
					value:       util.Pointer("value2"),
					envVariable: "ENVVAR2",
					isSet:       false,
				},
				{
					value:       util.Pointer("value3"),
					envVariable: "ENVVAR3",
					isSet:       false,
				},
				{
					value:       util.Pointer("value4"),
					envVariable: "ENVVAR4",
					isSet:       false,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Union(tt.args.lists...)
			if diff := cmp.Diff(*got, tt.want, cmp.AllowUnexported(Option{})); len(diff) > 0 {
				t.Errorf("Union() got diff %v", diff)
			}
		})
	}
}
