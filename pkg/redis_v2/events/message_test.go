package events

import (
	"testing"

	goredis "github.com/go-redis/redis/v8"
	"github.com/go-test/deep"
)

func init() {
	deep.CompareUnexportedFields = true
}

func TestNewRedisEventMessage(t *testing.T) {
	tests := []struct {
		name    string
		msg     *goredis.Message
		want    RedisEventMessage
		wantErr bool
	}{
		// empty event
		{
			name: "Returns an error for empty event messages",
			msg: &goredis.Message{
				Channel: "",
				Payload: "",
			},
			want:    RedisEventMessage{},
			wantErr: true,
		},
		// tilt event
		{
			name: "Returns a tilt RedisEventMessage object",
			msg: &goredis.Message{
				Channel: "+tilt",
				Payload: "",
			},
			want: RedisEventMessage{
				event:  "+tilt",
				target: RedisInstanceDetails{},
				master: RedisInstanceDetails{},
			},
			wantErr: false,
		},
		{
			name: "Returns an error for invalid tilt messages",
			msg: &goredis.Message{
				Channel: "+tilt",
				Payload: "text",
			},
			want:    RedisEventMessage{},
			wantErr: true,
		},
		// switch-master event
		{
			name: "Returns a switch-master RedisEventMessage object",
			msg: &goredis.Message{
				Channel: "+switch-master",
				Payload: "shard01 10.244.0.36 6379 10.244.0.38 6379",
			},
			want: RedisEventMessage{
				event: "+switch-master",
				target: RedisInstanceDetails{
					name: "shard01",
					ip:   "10.244.0.36",
					port: "6379",
					role: "master",
				},
				master: RedisInstanceDetails{
					name: "shard01",
					ip:   "10.244.0.38",
					port: "6379",
					role: "master",
				},
			},
			wantErr: false,
		},
		// set event
		{
			name: "Returns a valid set RedisEventMessage object",
			msg: &goredis.Message{
				Channel: "+set",
				Payload: "master shard02 10.244.0.39 6379 down-after-milliseconds 5000",
			},
			want: RedisEventMessage{
				event: "+set",
				config: RedisConfig{
					name:  "down-after-milliseconds",
					value: "5000",
				},
				target: RedisInstanceDetails{
					name: "shard02",
					ip:   "10.244.0.39",
					port: "6379",
					role: "master",
				},
				master: RedisInstanceDetails{
					name: "shard02",
					ip:   "10.244.0.39",
					port: "6379",
					role: "master",
				},
			},
			wantErr: false,
		},
		// monitor event
		{
			name: "Returns a monitor RedisEventMessage object",
			msg: &goredis.Message{
				Channel: "+monitor",
				Payload: "master shard01 10.244.0.24 6379 quorum 2",
			},
			want: RedisEventMessage{
				event: "+monitor",
				config: RedisConfig{
					name:  "quorum",
					value: "2",
				},
				target: RedisInstanceDetails{
					name: "shard01",
					ip:   "10.244.0.24",
					port: "6379",
					role: "master",
				},
				master: RedisInstanceDetails{
					name: "shard01",
					ip:   "10.244.0.24",
					port: "6379",
					role: "master",
				},
			},
			wantErr: false,
		},
		{
			name: "Returns an error for invalid  monitor event",
			msg: &goredis.Message{
				Channel: "+monitor",
				Payload: "master shard01 10.244.0.24 6379",
			},
			want:    RedisEventMessage{},
			wantErr: true,
		},
		// new-epoch event
		{
			name: "Returns a new-epoch RedisEventMessage object",
			msg: &goredis.Message{
				Channel: "+new-epoch",
				Payload: "1",
			},
			want: RedisEventMessage{
				event: "+new-epoch",
				config: RedisConfig{
					value: "1",
				},
			},
			wantErr: false,
		},
		// vote-for-leader event
		{
			name: "Returns a vote-for-leader RedisEventMessage object",
			msg: &goredis.Message{
				Channel: "+vote-for-leader",
				Payload: "43561253f764f487f959025d119e9d63354dc399 1",
			},
			want: RedisEventMessage{
				event: "+vote-for-leader",
				config: RedisConfig{
					name:  "43561253f764f487f959025d119e9d63354dc399",
					value: "1",
				},
			},
			wantErr: false,
		},
		// master instance with instance details event
		{
			name: "Returns a valid master <instance details> type RedisEventMessage",
			msg: &goredis.Message{
				Channel: "+sdown",
				Payload: "master shard01 10.244.0.24 6379",
			},
			want: RedisEventMessage{
				event: "+sdown",
				target: RedisInstanceDetails{
					name: "shard01",
					ip:   "10.244.0.24",
					port: "6379",
					role: "master",
				},
				master: RedisInstanceDetails{
					name: "shard01",
					ip:   "10.244.0.24",
					port: "6379",
					role: "master",
				},
			},
		},
		{
			name: "Returns an error for invalid master <instance details> type messages",
			msg: &goredis.Message{
				Channel: "+reset-master",
				Payload: "master 10.244.0.20:6379",
			},
			want:    RedisEventMessage{},
			wantErr: true,
		},
		// non-master instance with instance details event
		{
			name: "Returns a valid non-master <instance details> type RedisEventMessage",
			msg: &goredis.Message{
				Channel: "+sentinel",
				Payload: "sentinel 4b9b24bb72e032080232362b907678a5b0e1ec0b 10.244.0.46 26379 @ shard02 10.244.0.39 6379",
			},
			want: RedisEventMessage{
				event: "+sentinel",
				target: RedisInstanceDetails{
					name: "4b9b24bb72e032080232362b907678a5b0e1ec0b",
					ip:   "10.244.0.46",
					port: "26379",
					role: "sentinel",
				},
				master: RedisInstanceDetails{
					name: "shard02",
					ip:   "10.244.0.39",
					port: "6379",
					role: "master",
				},
			},
			wantErr: false,
		},
		{
			name: "Returns an error for non-master <instance details> type message without master info",
			msg: &goredis.Message{
				Channel: "+sentinel",
				Payload: "slave 10.244.0.20:6379 10.244.0.20 6379 @ missing-info",
			},
			want:    RedisEventMessage{},
			wantErr: true,
		},
		{
			name: "Returns an error for invalid <instance details> type messages",
			msg: &goredis.Message{
				Channel: "+slave-reconf-done",
				Payload: "slave 10.244.0.20:6379 missing-info",
			},
			want:    RedisEventMessage{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewRedisEventMessage(tt.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewRedisEventMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := deep.Equal(got, tt.want); len(diff) > 0 {
				t.Errorf("NewRedisEventMessage() got diff: %v", diff)
			}
		})
	}
}
