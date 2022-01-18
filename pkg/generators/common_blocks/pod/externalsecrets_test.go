package pod

import (
	"reflect"
	"testing"
	"time"

	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	externalsecretsv1alpha1 "github.com/3scale/saas-operator/pkg/apis/externalsecrets/v1alpha1"
	"github.com/go-test/deep"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
)

func TestGenerateExternalSecretFn(t *testing.T) {
	type args struct {
		name            string
		namespace       string
		secretStoreName string
		secretStoreKind string
		refreshInterval metav1.Duration
		labels          map[string]string
		opts            interface{}
	}
	tests := []struct {
		name string
		args args
		want *externalsecretsv1alpha1.ExternalSecret
	}{
		{
			name: "Generates a new ExternalSecret from an Options struct",
			args: args{
				name:            "my-secret",
				namespace:       "test",
				secretStoreName: "vault-mgmt",
				secretStoreKind: "ClusterSecretStore",
				refreshInterval: metav1.Duration{Duration: 120 * time.Second},
				labels:          map[string]string{},
				opts: struct {
					Option1 EnvVarValue `env:"OPTION1"`
					Option2 EnvVarValue `env:"OPTION2" secret:"my-secret"`
					Option3 EnvVarValue `env:"OPTION3" secret:"my-secret"`
					Option4 EnvVarValue `env:"OPTION4" secret:"other-secret"`
					Option5 EnvVarValue `env:"OPTION5" secret:"not-set"`
				}{
					Option1: &ClearTextValue{Value: "value1"},
					Option2: &SecretValue{Value: saasv1alpha1.SecretReference{
						FromVault: &saasv1alpha1.VaultSecretReference{Key: "key2", Path: "path2"}}},
					Option3: &SecretValue{Value: saasv1alpha1.SecretReference{
						FromVault: &saasv1alpha1.VaultSecretReference{Key: "key3", Path: "path3"}}},
					Option4: &SecretValue{Value: saasv1alpha1.SecretReference{
						FromVault: &saasv1alpha1.VaultSecretReference{Key: "key3", Path: "path3"}}},
				},
			},
			want: &externalsecretsv1alpha1.ExternalSecret{
				TypeMeta: metav1.TypeMeta{
					Kind:       externalsecretsv1alpha1.ExtSecretKind,
					APIVersion: externalsecretsv1alpha1.ExtSecretGroupVersionKind.GroupVersion().String(),
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-secret",
					Namespace: "test",
					Labels:    map[string]string{},
				},
				Spec: externalsecretsv1alpha1.ExternalSecretSpec{
					SecretStoreRef: externalsecretsv1alpha1.SecretStoreRef{
						Name: "vault-mgmt",
						Kind: "ClusterSecretStore",
					},
					Target:          externalsecretsv1alpha1.ExternalSecretTarget{Name: "my-secret"},
					RefreshInterval: &metav1.Duration{Duration: 120 * time.Second},
					Data: []externalsecretsv1alpha1.ExternalSecretData{
						{
							SecretKey: "OPTION2",
							RemoteRef: externalsecretsv1alpha1.ExternalSecretDataRemoteRef{
								Key:      "path2",
								Property: "key2",
							},
						},
						{
							SecretKey: "OPTION3",
							RemoteRef: externalsecretsv1alpha1.ExternalSecretDataRemoteRef{
								Key:      "path3",
								Property: "key3",
							},
						},
					},
				},
			},
		},
		{
			name: "Generates other ExternalSecret from the same Options struct (see previous test)",
			args: args{
				name:            "other-secret",
				namespace:       "test",
				secretStoreName: "vault-mgmt",
				secretStoreKind: "ClusterSecretStore",
				refreshInterval: metav1.Duration{Duration: 2 * time.Minute},
				labels:          map[string]string{},
				opts: struct {
					Option1 EnvVarValue `env:"OPTION1"`
					Option2 EnvVarValue `env:"OPTION2" secret:"my-secret"`
					Option3 EnvVarValue `env:"OPTION3" secret:"my-secret"`
					Option4 EnvVarValue `env:"OPTION4" secret:"other-secret"`
				}{
					Option1: &ClearTextValue{Value: "value1"},
					Option2: &SecretValue{Value: saasv1alpha1.SecretReference{
						FromVault: &saasv1alpha1.VaultSecretReference{Key: "key2", Path: "path2"}}},
					Option3: &SecretValue{Value: saasv1alpha1.SecretReference{
						FromVault: &saasv1alpha1.VaultSecretReference{Key: "key3", Path: "path3"}}},
					Option4: &SecretValue{Value: saasv1alpha1.SecretReference{
						FromVault: &saasv1alpha1.VaultSecretReference{Key: "key4", Path: "path4"}}},
				},
			},
			want: &externalsecretsv1alpha1.ExternalSecret{
				TypeMeta: metav1.TypeMeta{
					Kind:       externalsecretsv1alpha1.ExtSecretKind,
					APIVersion: externalsecretsv1alpha1.ExtSecretGroupVersionKind.GroupVersion().String(),
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "other-secret",
					Namespace: "test",
					Labels:    map[string]string{},
				},
				Spec: externalsecretsv1alpha1.ExternalSecretSpec{
					SecretStoreRef: externalsecretsv1alpha1.SecretStoreRef{
						Name: "vault-mgmt",
						Kind: "ClusterSecretStore",
					},
					Target:          externalsecretsv1alpha1.ExternalSecretTarget{Name: "other-secret"},
					RefreshInterval: &metav1.Duration{Duration: 120 * time.Second},
					Data: []externalsecretsv1alpha1.ExternalSecretData{
						{
							SecretKey: "OPTION4",
							RemoteRef: externalsecretsv1alpha1.ExternalSecretDataRemoteRef{
								Key:      "path4",
								Property: "key4",
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateExternalSecretFn(tt.args.name, tt.args.namespace, tt.args.secretStoreName, tt.args.secretStoreKind, tt.args.refreshInterval, tt.args.labels, tt.args.opts)()
			if diff := deep.Equal(got, tt.want); len(diff) > 0 {
				t.Errorf("GenerateExternalSecretFn() = diff %v", diff)
			}
		})
	}
}

func Test_keysSlice(t *testing.T) {
	type args struct {
		name string
		opts interface{}
	}
	tests := []struct {
		name      string
		args      args
		want      []externalsecretsv1alpha1.ExternalSecretData
		wantPanic bool
	}{
		{
			name: "Generates a DataSources map",
			args: args{
				name: "my-secret",
				opts: struct {
					Option1 EnvVarValue `env:"OPTION1"`
					Option2 EnvVarValue `env:"OPTION2" secret:"my-secret"`
					Option3 EnvVarValue `env:"OPTION3" secret:"my-secret"`
					Option4 EnvVarValue `env:"OPTION4" secret:"other-secret"`
				}{
					Option1: &ClearTextValue{Value: "value1"},
					Option2: &SecretValue{Value: saasv1alpha1.SecretReference{
						FromVault: &saasv1alpha1.VaultSecretReference{Key: "key2", Path: "path2"}}},
					Option3: &SecretValue{Value: saasv1alpha1.SecretReference{
						FromVault: &saasv1alpha1.VaultSecretReference{Key: "key3", Path: "path3"}}},
					Option4: &SecretValue{Value: saasv1alpha1.SecretReference{
						FromVault: &saasv1alpha1.VaultSecretReference{Key: "key4", Path: "path4"}}},
				},
			},
			want: []externalsecretsv1alpha1.ExternalSecretData{
				{
					SecretKey: "OPTION2",
					RemoteRef: externalsecretsv1alpha1.ExternalSecretDataRemoteRef{
						Key:      "path2",
						Property: "key2",
					},
				},
				{
					SecretKey: "OPTION3",
					RemoteRef: externalsecretsv1alpha1.ExternalSecretDataRemoteRef{
						Key:      "path3",
						Property: "key3",
					},
				},
			},
			wantPanic: false,
		},
		{
			name: "Generates a DataSources map, with secret overrides",
			args: args{
				name: "my-secret",
				opts: struct {
					Option1 EnvVarValue `env:"OPTION1"`
					Option2 EnvVarValue `env:"OPTION2" secret:"my-secret"`
					Option3 EnvVarValue `env:"OPTION3" secret:"my-secret"`
					Option4 EnvVarValue `env:"OPTION4" secret:"other-secret"`
				}{
					Option1: &ClearTextValue{Value: "value1"},
					Option2: &SecretValue{Value: saasv1alpha1.SecretReference{
						FromVault: &saasv1alpha1.VaultSecretReference{Key: "key2", Path: "path2"}}},
					Option3: &SecretValue{Value: saasv1alpha1.SecretReference{
						Override: pointer.StringPtr("override")}},
					Option4: &SecretValue{Value: saasv1alpha1.SecretReference{
						FromVault: &saasv1alpha1.VaultSecretReference{Key: "key4", Path: "path4"}}},
				},
			},
			want: []externalsecretsv1alpha1.ExternalSecretData{
				{
					SecretKey: "OPTION2",
					RemoteRef: externalsecretsv1alpha1.ExternalSecretDataRemoteRef{
						Key:      "path2",
						Property: "key2",
					},
				},
			},
			wantPanic: false,
		},
		{
			name: "Panics if value is not a SecretValue",
			args: args{
				name: "my-secret",
				opts: struct {
					Option1 EnvVarValue `env:"OPTION1" secret:"my-secret"`
				}{
					Option1: &ClearTextValue{Value: "xxxx"},
				},
			},
			want:      []externalsecretsv1alpha1.ExternalSecretData{},
			wantPanic: true,
		},
		{
			name: "Panics if 'env' tag is missing",
			args: args{
				name: "my-secret",
				opts: struct {
					Option1 EnvVarValue `secret:"my-secret"`
				}{
					Option1: &SecretValue{Value: saasv1alpha1.SecretReference{}},
				},
			},
			want:      []externalsecretsv1alpha1.ExternalSecretData{},
			wantPanic: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if r == nil && tt.wantPanic {
					t.Errorf("code did not panic")
				}
				if r != nil && !tt.wantPanic {
					t.Errorf("code caused a panic")
				}
			}()

			if got := keysSlice(tt.args.name, tt.args.opts); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("keysSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}
