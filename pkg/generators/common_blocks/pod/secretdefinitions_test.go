package pod

import (
	"reflect"
	"testing"

	saasv1alpha1 "github.com/3scale/saas-operator/api/v1alpha1"
	secretsmanagerv1alpha1 "github.com/3scale/saas-operator/pkg/apis/secrets-manager/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestGenerateSecretDefinitionFn(t *testing.T) {
	type args struct {
		name      string
		namespace string
		labels    map[string]string
		opts      interface{}
	}
	tests := []struct {
		name string
		args args
		want *secretsmanagerv1alpha1.SecretDefinition
	}{
		{
			name: "Generates a new SecretDefinition from an Options struct",
			args: args{
				name:      "my-secret",
				namespace: "test",
				labels:    map[string]string{},
				opts: struct {
					Option1 EnvVarValue `env:"OPTION1"`
					Option2 EnvVarValue `env:"OPTION2" secret:"my-secret"`
					Option3 EnvVarValue `env:"OPTION3" secret:"my-secret"`
					Option4 EnvVarValue `env:"OPTION4" secret:"other-secret"`
				}{
					Option1: &ClearTextValue{Value: "value1"},
					Option2: &SecretValue{Value: saasv1alpha1.SecretReference{
						FromVault: saasv1alpha1.VaultSecretReference{Key: "key2", Path: "path2"}}},
					Option3: &SecretValue{Value: saasv1alpha1.SecretReference{
						FromVault: saasv1alpha1.VaultSecretReference{Key: "key3", Path: "path3"}}},
					Option4: &SecretValue{Value: saasv1alpha1.SecretReference{
						FromVault: saasv1alpha1.VaultSecretReference{Key: "key3", Path: "path3"}}},
				},
			},
			want: &secretsmanagerv1alpha1.SecretDefinition{
				TypeMeta: metav1.TypeMeta{
					Kind:       "SecretDefinition",
					APIVersion: secretsmanagerv1alpha1.GroupVersion.String(),
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-secret",
					Namespace: "test",
					Labels:    map[string]string{},
				},
				Spec: secretsmanagerv1alpha1.SecretDefinitionSpec{
					Name: "my-secret",
					Type: "opaque",
					KeysMap: map[string]secretsmanagerv1alpha1.DataSource{
						"OPTION2": {Key: "key2", Path: "path2"},
						"OPTION3": {Key: "key3", Path: "path3"}},
				},
			},
		},
		{
			name: "Generates other SecretDefinition from the same Options struct (see previous test)",
			args: args{
				name:      "other-secret",
				namespace: "test",
				labels:    map[string]string{},
				opts: struct {
					Option1 EnvVarValue `env:"OPTION1"`
					Option2 EnvVarValue `env:"OPTION2" secret:"my-secret"`
					Option3 EnvVarValue `env:"OPTION3" secret:"my-secret"`
					Option4 EnvVarValue `env:"OPTION4" secret:"other-secret"`
				}{
					Option1: &ClearTextValue{Value: "value1"},
					Option2: &SecretValue{Value: saasv1alpha1.SecretReference{
						FromVault: saasv1alpha1.VaultSecretReference{Key: "key2", Path: "path2"}}},
					Option3: &SecretValue{Value: saasv1alpha1.SecretReference{
						FromVault: saasv1alpha1.VaultSecretReference{Key: "key3", Path: "path3"}}},
					Option4: &SecretValue{Value: saasv1alpha1.SecretReference{
						FromVault: saasv1alpha1.VaultSecretReference{Key: "key4", Path: "path4"}}},
				},
			},
			want: &secretsmanagerv1alpha1.SecretDefinition{
				TypeMeta: metav1.TypeMeta{
					Kind:       "SecretDefinition",
					APIVersion: secretsmanagerv1alpha1.GroupVersion.String(),
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "other-secret",
					Namespace: "test",
					Labels:    map[string]string{},
				},
				Spec: secretsmanagerv1alpha1.SecretDefinitionSpec{
					Name: "other-secret",
					Type: "opaque",
					KeysMap: map[string]secretsmanagerv1alpha1.DataSource{
						"OPTION4": {Key: "key4", Path: "path4"},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fn := GenerateSecretDefinitionFn(tt.args.name, tt.args.namespace, tt.args.labels, tt.args.opts)
			if got := fn(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GenerateSecretDefinitionFn() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_keysMap(t *testing.T) {
	type args struct {
		name string
		opts interface{}
	}
	tests := []struct {
		name      string
		args      args
		want      map[string]secretsmanagerv1alpha1.DataSource
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
						FromVault: saasv1alpha1.VaultSecretReference{Key: "key2", Path: "path2"}}},
					Option3: &SecretValue{Value: saasv1alpha1.SecretReference{
						FromVault: saasv1alpha1.VaultSecretReference{Key: "key3", Path: "path3"}}},
					Option4: &SecretValue{Value: saasv1alpha1.SecretReference{
						FromVault: saasv1alpha1.VaultSecretReference{Key: "key4", Path: "path4"}}},
				},
			},
			want: map[string]secretsmanagerv1alpha1.DataSource{
				"OPTION2": {Key: "key2", Path: "path2"},
				"OPTION3": {Key: "key3", Path: "path3"},
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
			want:      map[string]secretsmanagerv1alpha1.DataSource{},
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
			want:      map[string]secretsmanagerv1alpha1.DataSource{},
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

			if got := keysMap(tt.args.name, tt.args.opts); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("keysMap() = %v, want %v", got, tt.want)
			}
		})
	}
}
