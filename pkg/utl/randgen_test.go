package utl

import (
	"testing"
)

func TestRandomSequence(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Success",
			args: args{
				n: 1,
			},
			want: "a", // first character of the alphabet
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Intn = func(n int) int {
				return 0
			}
			if got := RandomSequence(tt.args.n); got != tt.want {

				t.Errorf("RandomSequence() = %v, want %v", got, tt.want)
			}
		})
	}
}
