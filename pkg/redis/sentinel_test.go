package redis

import (
	"reflect"
	"testing"

	redistypes "github.com/3scale/saas-operator/pkg/redis/types"
)

func Test_sliceCmdToStruct(t *testing.T) {
	type args struct {
		in  []interface{}
		out interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "holy crap",
			args: args{
				in: []interface{}{
					"name",
					"10.244.0.7:6379",
					"ip",
					"10.244.0.7",
					"port",
					"6379",
					"runid",
					"4009bfff0cb1115e8b5d06e00730808a051fa2c5",
					"flags",
					"slave",
					"link-pending-commands",
					"0",
					"link-refcount",
					"1",
					"last-ping-sent",
					"0",
					"last-ok-ping-reply",
					"68",
					"last-ping-reply",
					"68",
					"down-after-milliseconds",
					"5000",
					"info-refresh",
					"2035",
					"role-reported",
					"slave",
					"role-reported-time",
					"91031035",
					"master-link-down-time",
					"0",
					"master-link-status",
					"ok",
					"master-host",
					"10.244.0.8",
					"master-port",
					"6379",
					"slave-priority",
					"100",
					"slave-repl-offset",
					"12707378",
				},
				out: &redistypes.SentinelSlaveCmdResult{},
			},
			want: &redistypes.SentinelSlaveCmdResult{
				Name:                  "10.244.0.7:6379",
				IP:                    "10.244.0.7",
				Port:                  6379,
				RunID:                 "4009bfff0cb1115e8b5d06e00730808a051fa2c5",
				Flags:                 "slave",
				LinkPendingCommands:   0,
				LinkRefcount:          1,
				LastPingSet:           0,
				LastOkPingReply:       68,
				LastPingReply:         68,
				DownAfterMilliseconds: 5000,
				InfoRefresh:           2035,
				RoleReported:          "slave",
				RoleReportedTime:      91031035,
				MasterLinkDownTime:    0,
				MasterLinkStatus:      "ok",
				MasterHost:            "10.244.0.8",
				MasterPort:            "6379",
				SlavePriority:         "100",
				SlaveReplOffset:       12707378,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := sliceCmdToStruct(tt.args.in, tt.args.out)
			if (err != nil) != tt.wantErr {
				t.Errorf("sliceCmdToStruct() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("sliceCmdToStruct() = %v, want %v", got, tt.want)
			}
		})
	}
}
