package convert

import (
	"reflect"
	"testing"

	"github.com/volatiletech/null/v8"
)

const SuccessCase = "Success"

func TestStringToPointerString(t *testing.T) {
	type args struct {
		v string
	}
	tests := []struct {
		name string
		args args
		want *string
	}{
		{
			name: SuccessCase,
			args: args{
				v: "test",
			},
			want: null.StringFrom("test").Ptr(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StringToPointerString(tt.args.v); *got != *tt.want {
				t.Errorf("StringToPointerString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringToInt(t *testing.T) {
	type args struct {
		v string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: SuccessCase,
			args: args{
				v: "1",
			},
			want: 1,
		},
		{
			name: "Failure",
			args: args{
				v: "a",
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StringToInt(tt.args.v); got != tt.want {
				t.Errorf("StringToInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringToBool(t *testing.T) {
	type args struct {
		v string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: SuccessCase,
			args: args{
				v: "true",
			},
			want: true,
		},
		{
			name: "Failure",
			args: args{
				v: "falasdse",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StringToBool(tt.args.v); got != tt.want {
				t.Errorf("StringToBool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNullDotStringToPointerString(t *testing.T) {
	type args struct {
		v null.String
	}
	tests := []struct {
		name string
		args args
		want *string
	}{
		{
			name: SuccessCase,
			args: args{
				v: null.StringFrom("true"),
			},
			want: null.StringFrom("true").Ptr(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NullDotStringToPointerString(tt.args.v); *got != *tt.want {
				t.Errorf("NullDotStringToPointerString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNullDotStringToString(t *testing.T) {
	type args struct {
		v null.String
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: SuccessCase,
			args: args{
				v: null.StringFrom("true"),
			},
			want: "true",
		},
		{
			name: "Success_Nil",
			args: args{
				v: null.String{},
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NullDotStringToString(tt.args.v); got != tt.want {
				t.Errorf("NullDotStringToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNullDotIntToInt(t *testing.T) {
	type args struct {
		v null.Int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: SuccessCase,
			args: args{
				v: null.IntFrom(1),
			},
			want: 1,
		},
		{
			name: "Success_Nil",
			args: args{
				v: null.Int{},
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NullDotIntToInt(tt.args.v); got != tt.want {
				t.Errorf("NullDotIntToInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNullDotBoolToPointerBool(t *testing.T) {
	type args struct {
		v null.Bool
	}
	boolean := null.BoolFrom(true)
	tests := []struct {
		name string
		args args
		want *bool
	}{
		{
			name: SuccessCase,
			args: args{
				v: boolean,
			},
			want: boolean.Ptr(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NullDotBoolToPointerBool(tt.args.v); *got != *tt.want {
				t.Errorf("NullDotBoolToPointerBool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPointerStringToNullDotInt(t *testing.T) {
	type args struct {
		s *string
	}
	tests := []struct {
		name string
		args args
		want null.Int
	}{
		{
			name: SuccessCase,
			args: args{
				s: null.StringFrom("1").Ptr(),
			},
			want: null.IntFrom(1),
		},
		{
			name: "Success_InvalidValue",
			args: args{
				s: null.StringFrom("asd1").Ptr(),
			},
			want: null.IntFrom(0),
		},
		{
			name: "Success_Nil",
			args: args{
				s: null.String{}.Ptr(),
			},
			want: null.IntFrom(0),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PointerStringToNullDotInt(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PointerStringToNullDotInt() = %v, want %v", got, tt.want)
			}
		})
	}
}
