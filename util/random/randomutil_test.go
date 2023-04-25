package random

import "testing"

func TestRandomInt(t *testing.T) {
	type args struct {
		limit int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
		{"intTest", args{100}, 100},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RandomInt(tt.args.limit); got >= tt.want {
				t.Errorf("RandomInt() = %v, want %v", got, tt.want)
			}
		})
	}
}
