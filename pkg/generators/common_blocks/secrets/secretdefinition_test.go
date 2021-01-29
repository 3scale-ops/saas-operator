package secrets

import (
	"reflect"
	"testing"

	secretsmanagerv1alpha1 "github.com/3scale/saas-operator/pkg/apis/secrets-manager/v1alpha1"
	"github.com/MakeNowJust/heredoc"
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
