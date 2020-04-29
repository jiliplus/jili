package bar

import (
	"testing"
)

func Test_maxFloat64(t *testing.T) {
	type args struct {
		a float64
		b float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{

		{
			"a 大",
			args{
				a: 1,
				b: 0,
			},
			1,
		},

		{
			"b 大",
			args{
				a: 1,
				b: 2,
			},
			2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := maxFloat64(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("maxFloat64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_minFloat64(t *testing.T) {
	type args struct {
		a float64
		b float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			"b 大",
			args{
				a: 1,
				b: 2,
			},
			1,
		},

		{
			"b 小",
			args{
				a: 1,
				b: 0,
			},
			0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := minFloat64(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("minFloat64() = %v, want %v", got, tt.want)
			}
		})
	}
}
