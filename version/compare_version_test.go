package main

import "testing"

func Test_compareVersion(t *testing.T) {
	type args struct {
		versionA string
		versionB string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"", args{"1.0.1", "1.0.0"}, VersionBig},
		{"", args{"1.0.1", "1.1.0"}, VersionSmall},
		{"", args{"1.2.0", "1.2.0"}, VersionEqual},
		{"", args{"1.2.0.beta", "1.2.0.alpha"}, VersionBig},
		{"", args{"3.17.0.rc-1", "3.17.0.beta-1"}, VersionBig},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Log("compare version:", tt.args.versionA, tt.args.versionB)
			if got := compareVersion(tt.args.versionA, tt.args.versionB); got != tt.want {
				t.Errorf("compareVersion() = %v, want %v", got, tt.want)
			} else {
				t.Logf("compareVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}
