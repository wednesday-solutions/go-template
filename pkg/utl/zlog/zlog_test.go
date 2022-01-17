package zlog

import (
	"os"
	"reflect"
	"testing"

	"github.com/rs/zerolog"
)

func TestNew(t *testing.T) {
	z := zerolog.New(os.Stdout)
	tests := []struct {
		name string
		want *Log
	}{
		{
			name: "Success",
			want: &Log{
				logger: &z,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
