package time

import (
	"reflect"
	"testing"
	"time"
)

func TestNow(t *testing.T) {
	tests := []struct {
		name string
		want time.Time
	}{
		// TODO: Add test cases.
		{"0", time.Now()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Now(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Now() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTimeToMs(t *testing.T) {
	type args struct {
		t0 time.Time
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		// TODO: Add test cases.
		{"time2", args{Now()}, time.Now().UTC().UnixMilli()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TimeToMs(tt.args.t0); got != tt.want {
				t.Errorf("TimeToMs() = %v, want %v", got, tt.want)
			} else {
				t.Log(got, tt.args.t0)
			}
		})
	}
}

func TestMsToTime(t *testing.T) {
	type args struct {
		epochMilli int64
	}
	w := func() time.Time {
		t0, _ := ParseDateTime("2023-05-04 02:32:33")
		return t0
	}()
	tests := []struct {
		name string
		args args
		want time.Time
	}{
		// TODO: Add test cases.
		{"time", args{1683167553720}, w},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MsToTime(tt.args.epochMilli); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MsToTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormatNormal(t *testing.T) {
	type args struct {
		t0 time.Time
	}
	t02 := func() time.Time {
		t01, _ := ParseDateTime("2023-05-04 10:43:03")
		return t01
	}()
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"t", args{t02}, "2023-05-04 10:43:03",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatNormal(tt.args.t0); got != tt.want {
				t.Errorf("FormatNormal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPlusSeconds(t *testing.T) {
	type args struct {
		t0      time.Time
		seconds int
	}
	tests := []struct {
		name string
		args args
		want time.Time
	}{
		// TODO: Add test cases.
		{"t0", args{t0: func() time.Time {
			t01, _ := ParseDateTime("2023-05-04 10:43:03")
			return t01
		}(), seconds: 60}, func() time.Time {
			t01, _ := ParseDateTime("2023-05-04 10:44:03")
			return t01
		}()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PlusSeconds(tt.args.t0, tt.args.seconds); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PlusSeconds() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPlusMinutes(t *testing.T) {
	type args struct {
		t0      time.Time
		minutes int
	}
	tests := []struct {
		name string
		args args
		want time.Time
	}{
		// TODO: Add test cases.
		{"t0", args{t0: func() time.Time {
			t01, _ := ParseDateTime("2023-05-04 10:43:03")
			return t01
		}(), minutes: 1}, func() time.Time {
			t01, _ := ParseDateTime("2023-05-04 10:44:03")
			return t01
		}()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PlusMinutes(tt.args.t0, tt.args.minutes); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PlusMinutes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPlusMinutes1(t *testing.T) {
	type args struct {
		t0      time.Time
		minutes int
	}
	tests := []struct {
		name string
		args args
		want time.Time
	}{
		// TODO: Add test cases.
		{"t0", args{t0: func() time.Time {
			t01, _ := ParseDateTime("2023-05-04 10:44:03")
			return t01
		}(), minutes: -1}, func() time.Time {
			t01, _ := ParseDateTime("2023-05-04 10:43:03")
			return t01
		}()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PlusMinutes(tt.args.t0, tt.args.minutes); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PlusMinutes() = %v, want %v", got, tt.want)
			}
		})
	}
}
