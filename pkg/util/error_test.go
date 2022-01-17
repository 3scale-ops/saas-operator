package util

import (
	"errors"
	"testing"
)

func TestMultiError_Error(t *testing.T) {
	tests := []struct {
		name string
		me   MultiError
		want string
	}{
		{
			name: "Returns several errors",
			me: []error{
				errors.New("error1"),
				errors.New("error2"),
				errors.New("error3"),
			},
			want: `["error1","error2","error3"]`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.me.Error(); got != tt.want {
				t.Errorf("MultiError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}
