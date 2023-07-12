package client

import (
	"testing"
	"time"
)

func TestSentinelInfoCache_GetValue(t *testing.T) {
	type args struct {
		shard       string
		runID       string
		key         string
		maxCacheAge time.Duration
	}
	tests := []struct {
		name    string
		sic     SentinelInfoCache
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Returns the value",
			sic: map[string]map[string]RedisServerInfoCache{
				"shard01": {
					"xxxxxx": RedisServerInfoCache{
						CacheAge: 500 * time.Millisecond,
						Info: map[string]string{
							"key": "value",
						},
					},
				},
			},
			args: args{
				shard:       "shard01",
				runID:       "xxxxxx",
				key:         "key",
				maxCacheAge: 1 * time.Second,
			},
			want:    "value",
			wantErr: false,
		},
		{
			name: "Error, shard not found",
			sic:  map[string]map[string]RedisServerInfoCache{"zzzz": {}},
			args: args{
				shard:       "shard01",
				runID:       "xxxxxx",
				key:         "key",
				maxCacheAge: 1 * time.Second,
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "Error, run_id not found",
			sic: map[string]map[string]RedisServerInfoCache{
				"shard01": {
					"zzzzz": RedisServerInfoCache{},
				},
			},
			args: args{
				shard:       "shard01",
				runID:       "xxxxxx",
				key:         "key",
				maxCacheAge: 1 * time.Second,
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "Error, cache too old",
			sic: map[string]map[string]RedisServerInfoCache{
				"shard01": {
					"xxxxxx": RedisServerInfoCache{
						CacheAge: 2 * time.Second,
						Info: map[string]string{
							"key": "value",
						},
					},
				},
			},
			args: args{
				shard:       "shard01",
				runID:       "xxxxxx",
				key:         "key",
				maxCacheAge: 1 * time.Second,
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "Error, key not found",
			sic: map[string]map[string]RedisServerInfoCache{
				"shard01": {
					"xxxxxx": RedisServerInfoCache{
						CacheAge: 1 * time.Millisecond,
						Info:     map[string]string{},
					},
				},
			},
			args: args{
				shard:       "shard01",
				runID:       "xxxxxx",
				key:         "key",
				maxCacheAge: 1 * time.Second,
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.sic.GetValue(tt.args.shard, tt.args.runID, tt.args.key, tt.args.maxCacheAge)
			if (err != nil) != tt.wantErr {
				t.Errorf("SentinelInfoCache.GetValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SentinelInfoCache.GetValue() = %v, want %v", got, tt.want)
			}
		})
	}
}
