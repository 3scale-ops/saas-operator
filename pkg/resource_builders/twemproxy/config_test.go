package twemproxy

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestTwemproxyServer_MarshalJSON(t *testing.T) {
	type fields struct {
		Address  string
		Priority int
		Name     string
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			name: "marshal",
			fields: fields{
				Address:  "127.0.0.1:6379",
				Priority: 1,
				Name:     "server",
			},
			want:    []byte("\"127.0.0.1:6379:1 server\""),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := &Server{
				Address:  tt.fields.Address,
				Priority: tt.fields.Priority,
				Name:     tt.fields.Name,
			}
			got, err := srv.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("TwemproxyServer.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TwemproxyServer.MarshalJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTwemproxyServer_UnmarshalJSON(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *Server
		wantErr bool
	}{
		{
			name: "unmarshal",
			args: args{
				data: []byte("\"127.0.0.1:6379:2 server\""),
			},
			want: &Server{
				Address:  "127.0.0.1:6379",
				Priority: 2,
				Name:     "server",
			},
			wantErr: false,
		},
		{
			name: "error",
			args: args{
				data: []byte("\"127.0.0.1:6379:x server\""),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := &Server{}
			err := srv.UnmarshalJSON(tt.args.data)
			if err != nil {
				if tt.wantErr != true {
					t.Fatalf("TwemproxyServer.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			} else {
				if diff := cmp.Diff(srv, tt.want); len(diff) != 0 {
					t.Fatalf("TwemproxyServer.UnmarshalJSON() diff = %v", diff)
				}
			}

		})
	}
}
