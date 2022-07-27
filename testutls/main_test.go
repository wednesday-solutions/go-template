package testutls

import (
	"os"
	"testing"
)

func TestIsInTests(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{
			name: "-test.paniconexit0 is in os.Args",
			want: true,
		},
		{
			name: "no -test.paniconexit0 in os.Args",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.want {
				os.Args = []string{"something"}
			}
			if got := IsInTests(); got != tt.want {
				t.Errorf("IsInTests() = %v, want %v", got, tt.want)
			}
		})
	}
}
