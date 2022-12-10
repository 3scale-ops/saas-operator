package util

import (
	"reflect"
	"testing"
)

func TestUnique(t *testing.T) {
	type args struct {
		stringSlice []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "removes duplicates and sorts the list",
			args: args{
				stringSlice: []string{"z", "d", "r", "d", "z", "e"},
			},
			want: []string{"d", "e", "r", "z"},
		},
		{
			name: "just sorts the list",
			args: args{
				stringSlice: []string{"f", "c", "d", "x"},
			},
			want: []string{"c", "d", "f", "x"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Unique(tt.args.stringSlice); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Unique() = %v, want %v", got, tt.want)
			}
		})
	}
}
