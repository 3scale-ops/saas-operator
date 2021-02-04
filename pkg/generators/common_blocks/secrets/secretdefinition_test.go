package secrets

import (
	"reflect"
	"testing"

	secretsmanagerv1alpha1 "github.com/3scale/saas-operator/pkg/apis/secrets-manager/v1alpha1"
	"github.com/MakeNowJust/heredoc"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestSecretConfiguration_keysMap(t *testing.T) {
	type fields struct {
		SecretName    string
		ConfigOptions map[string]string
	}
	type args struct {
		basePath         string
		serializedConfig []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[string]secretsmanagerv1alpha1.DataSource
		wantErr bool
	}{
		{
			name: "Builds SecretDefinition DataSources from serialized config",
			fields: fields{
				SecretName: "secret-name",
				ConfigOptions: map[string]string{
					"SOME_CONFIG":  "/spec/config/someConfig",
					"OTHER_CONFIG": "/spec/config/otherConfig",
				},
			},
			args: args{
				basePath: "/spec",
				serializedConfig: []byte(heredoc.Doc(`
					{
						"config": {
							"someConfig": {
								"fromVault": {
									"path": "vault-path1",
									"key": "vault-key1"
								}
							},
							"otherConfig": {
								"fromVault": {
									"path": "vault-path2",
									"key": "vault-key2"
								}
							}
						}
					}
				`)),
			},
			want: map[string]secretsmanagerv1alpha1.DataSource{
				"SOME_CONFIG":  {Key: "vault-key1", Path: "vault-path1"},
				"OTHER_CONFIG": {Key: "vault-key2", Path: "vault-path2"},
			},
			wantErr: false,
		},
		{
			name: "Builds SecretDefinition DataSources from serialized config, reversed order",
			fields: fields{
				SecretName: "secret-name",
				ConfigOptions: map[string]string{
					"SOME_CONFIG":  "/spec/config/someConfig",
					"OTHER_CONFIG": "/spec/config/otherConfig",
				},
			},
			args: args{
				basePath: "/spec",
				serializedConfig: []byte(heredoc.Doc(`
					{
						"config": {
							"otherConfig": {
								"fromVault": {
									"path": "vault-path2",
									"key": "vault-key2"
								}
							},
							"someConfig": {
								"fromVault": {
									"path": "vault-path1",
									"key": "vault-key1"
								}
							}
						}
					}
				`)),
			},
			want: map[string]secretsmanagerv1alpha1.DataSource{
				"SOME_CONFIG":  {Key: "vault-key1", Path: "vault-path1"},
				"OTHER_CONFIG": {Key: "vault-key2", Path: "vault-path2"},
			},
			wantErr: false,
		},
		{
			name: "Fails if path is not found",
			fields: fields{
				SecretName: "secret-name",
				ConfigOptions: map[string]string{
					"SOME_CONFIG": "/spec/config/someConfig",
				},
			},
			args: args{
				basePath:         "/spec",
				serializedConfig: []byte(`{}`),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sc := &SecretConfiguration{
				SecretName:    tt.fields.SecretName,
				ConfigOptions: tt.fields.ConfigOptions,
			}
			got, err := sc.keysMap(tt.args.basePath, tt.args.serializedConfig)
			if (err != nil) != tt.wantErr {
				t.Errorf("SecretConfiguration.keysMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SecretConfiguration.keysMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSecretConfiguration_GenerateSecretDefinitionFn(t *testing.T) {
	type fields struct {
		SecretName    string
		ConfigOptions map[string]string
	}
	type args struct {
		namespace        string
		labels           map[string]string
		basePath         string
		serializedConfig []byte
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *secretsmanagerv1alpha1.SecretDefinition
	}{
		{
			name: "Returns a basereconciler.GeneratorFunction that returns a SecretDefinition when called",
			fields: fields{
				SecretName: "secret-name",
				ConfigOptions: map[string]string{
					"SOME_CONFIG": "/spec/config/someConfig",
				},
			},
			args: args{
				namespace: "namespace",
				labels:    map[string]string{},
				basePath:  "/spec",
				serializedConfig: []byte(heredoc.Doc(`
				{
					"config": {
						"someConfig": {
							"fromVault": {
								"path": "vault-path",
								"key": "vault-key"
							}
						}
					}
				}
			`)),
			},
			want: &secretsmanagerv1alpha1.SecretDefinition{
				TypeMeta: metav1.TypeMeta{
					Kind:       "SecretDefinition",
					APIVersion: secretsmanagerv1alpha1.GroupVersion.String(),
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "secret-name",
					Namespace: "namespace",
					Labels:    map[string]string{},
				},
				Spec: secretsmanagerv1alpha1.SecretDefinitionSpec{
					Name: "secret-name",
					Type: "opaque",
					KeysMap: map[string]secretsmanagerv1alpha1.DataSource{
						"SOME_CONFIG": {Path: "vault-path", Key: "vault-key"},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sc := &SecretConfiguration{
				SecretName:    tt.fields.SecretName,
				ConfigOptions: tt.fields.ConfigOptions,
			}
			fn := sc.GenerateSecretDefinitionFn(tt.args.namespace, tt.args.labels, tt.args.basePath, tt.args.serializedConfig)
			if got := fn(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SecretConfiguration.GenerateSecretDefinitionFn() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSecretConfigurations_LookupSecretName(t *testing.T) {
	type args struct {
		config string
	}
	tests := []struct {
		name string
		sc   *SecretConfigurations
		args args
		want string
	}{
		{
			name: "Returns the secret name given a config option name",
			sc: &SecretConfigurations{
				{
					SecretName: "secret1",
					ConfigOptions: map[string]string{
						"config1": "/spec/config/config1",
					},
				},
				{
					SecretName: "secret2",
					ConfigOptions: map[string]string{
						"config2": "/spec/config/config2",
					},
				},
			},
			args: args{config: "config2"},
			want: "secret2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.sc.LookupSecretName(tt.args.config); got != tt.want {
				t.Errorf("SecretConfigurations.LookupSecretName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSecretConfigurations_LookupSecretConfiguration(t *testing.T) {
	type args struct {
		secretName string
	}
	tests := []struct {
		name string
		sc   *SecretConfigurations
		args args
		want SecretConfiguration
	}{
		{
			name: "Returns a SecretConfiguration for the given Secret/SecretDefinition name",
			sc: &SecretConfigurations{
				{
					SecretName: "secret1",
					ConfigOptions: map[string]string{
						"config1": "/spec/config/config1",
					},
				},
				{
					SecretName: "secret2",
					ConfigOptions: map[string]string{
						"config2": "/spec/config/config2",
					},
				},
			},
			args: args{secretName: "secret2"},
			want: SecretConfiguration{
				SecretName: "secret2",
				ConfigOptions: map[string]string{
					"config2": "/spec/config/config2",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.sc.LookupSecretConfiguration(tt.args.secretName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SecretConfigurations.LookupSecretConfiguration() = %v, want %v", got, tt.want)
			}
		})
	}
}
