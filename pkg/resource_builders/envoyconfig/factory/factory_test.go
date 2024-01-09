package factory

import (
	"testing"

	"github.com/3scale-ops/marin3r/pkg/envoy"
	descriptor "github.com/3scale-ops/saas-operator/pkg/resource_builders/envoyconfig/descriptor"
	envoy_service_runtime_v3 "github.com/envoyproxy/go-control-plane/envoy/service/runtime/v3"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
)

type unregisteredType struct{ opts *opts }
type opts struct{}

func (x *unregisteredType) GetOptions() interface{}     { return x.opts }
func (x *unregisteredType) GetGeneratorVersion() string { return "" }
func (x *unregisteredType) GetName() string             { return "" }

type testDescriptor struct {
	name             string
	generatorVersion string
	opts             *testOptions
}

type testOptions struct {
	structpb *structpb.Struct
}

func (x *testDescriptor) GetOptions() interface{}     { return x.opts }
func (x *testDescriptor) GetGeneratorVersion() string { return x.generatorVersion }
func (x *testDescriptor) GetName() string             { return x.name }

func testTemplate(name string, opts interface{}) (envoy.Resource, error) {
	o := opts.(*testOptions)

	return &envoy_service_runtime_v3.Runtime{
		Name:  name,
		Layer: o.structpb,
	}, nil
}

var testFactory = EnvoyDynamicConfigFactory{
	"testOptions_v1": RegisterTemplate(testTemplate, &envoy_service_runtime_v3.Runtime{}),
}

func TestEnvoyDynamicConfigFactory_NewResource(t *testing.T) {
	type args struct {
		descriptor descriptor.EnvoyDynamicConfigDescriptor
	}
	tests := []struct {
		name    string
		factory EnvoyDynamicConfigFactory
		args    args
		want    envoy.Resource
		wantErr bool
	}{
		{
			name:    "Generates a runtime proto",
			factory: testFactory,
			args: args{
				descriptor: &testDescriptor{
					name:             "test",
					generatorVersion: "v1",
					opts: &testOptions{
						structpb: func() *structpb.Struct {
							l, _ := structpb.NewStruct(map[string]interface{}{
								"key": map[string]interface{}{
									"key": map[string]interface{}{},
								},
							})
							return l
						}(),
					},
				},
			},
			want: func() envoy.Resource {
				return &envoy_service_runtime_v3.Runtime{
					Name: "test",
					Layer: func() *structpb.Struct {
						l, _ := structpb.NewStruct(map[string]interface{}{
							"key": map[string]interface{}{
								"key": map[string]interface{}{},
							},
						})
						return l
					}(),
				}
			}(),
			wantErr: false,
		},
		{
			name:    "Unregistered class",
			factory: testFactory,
			args: args{
				descriptor: &unregisteredType{},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.factory.NewResource(tt.args.descriptor)
			if (err != nil) != tt.wantErr {
				t.Errorf("EnvoyDynamicConfigFactory.NewResource() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !proto.Equal(got, tt.want) {
				t.Errorf("EnvoyDynamicConfigFactory.NewResource() = %v, want %v", got, tt.want)
			}
		})
	}
}
