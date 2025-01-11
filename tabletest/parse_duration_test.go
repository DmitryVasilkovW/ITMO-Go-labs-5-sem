package tabletest

import (
	"math/rand"
	"testing"
	"time"
)

var parseDurationTests = []struct {
	input    string
	duration time.Duration
	ok       bool
}{
	{"0", 0, true},
	{"2s", 2 * time.Second, true},
	{"10m", 10 * time.Minute, true},
	{"3h", 3 * time.Hour, true},
	{"1234ms", 1234 * time.Millisecond, true},
	{"555ns", 555 * time.Nanosecond, true},
	{"789us", 789 * time.Microsecond, true},
	{"456µs", 456 * time.Microsecond, true},
	{"239μs", 239 * time.Microsecond, true},

	{"-2m", -2 * time.Minute, true},
	{"+3s", 3 * time.Second, true},
	{"-0", 0, true},
	{"+0", 0, true},

	{"3.0m", 3 * time.Minute, true},
	{"2.8s", 2*time.Second + 800*time.Millisecond, true},
	{"4.s", 4 * time.Second, true},
	{".6s", 600 * time.Millisecond, true},
	{"1.50s", 1*time.Second + 500*time.Millisecond, true},
	{"7.000s", 7 * time.Second, true},
	{"9.009s", 9*time.Second + 9*time.Millisecond, true},
	{"2.1230s", 2*time.Second + 123*time.Millisecond, true},
	{"2h45m", 2*time.Hour + 45*time.Minute, true},
	{"15m1h", 15*time.Minute + 1*time.Hour, true},
	{"7m12.5s", 7*time.Minute + 12*time.Second + 500*time.Millisecond, true},
	{"12.5s7m", 7*time.Minute + 12*time.Second + 500*time.Millisecond, true},
	{"-2m12.5s", -(2*time.Minute + 12*time.Second + 500*time.Millisecond), true},
	{"-12.5s2m", -(2*time.Minute + 12*time.Second + 500*time.Millisecond), true},
	{"2h30m10s20ms30ns40µs", 2*time.Hour + 30*time.Minute + 10*time.Second + 20*time.Millisecond + 40*time.Microsecond + 30*time.Nanosecond, true},
	{"1h10m5.25s", 1*time.Hour + 10*time.Minute + 5*time.Second + 250*time.Millisecond, true},
	{"123456789012ns", 123456789012 * time.Nanosecond, true},
	{"9223372036854775807ns", (1<<63 - 1) * time.Nanosecond, true},
	{"0.1111111111111111111h", 6*time.Minute + 40*time.Second, true},
	{"0.05000000000000000000h", 3 * time.Minute, true},

	{"9223372036854775808ns", 0, false},
	{"10000000000000000000ns", 0, false},
	{"9223372036854775808ms", 0, false},
	{"10000000000000000000ms", 0, false},
	{"9223372036854775808s", 0, false},
	{"10000000000000000000s", 0, false},
	{"9223372036854775808m", 0, false},
	{"10000000000000000000m", 0, false},
	{"9223372036854775808h", 0, false},
	{"10000000000000000000h", 0, false},
	{"9223372036854775807.0000000001ms", 0, false},
	{"9223372036854775.808000000001us", 0, false},

	{"", 0, false},
	{"12", 0, false},
	{"-12", 0, false},
	{"z", 0, false},
	{".", 0, false},
	{"-.", 0, false},
	{"+.", 0, false},
	{"x.s", 0, false},
	{"+.t", 0, false},

	{"92233720368ns", 92233720368 * time.Nanosecond, true},
	{"9223372036854775807ns", (1<<63 - 1) * time.Nanosecond, true},
	{"0.3333333333333333333h", 20 * time.Minute, true},
	{"0.100000000000000000000h", 6 * time.Minute, true},
	{"0.9223372036854775807h", 55*time.Minute + 20*time.Second + 413933267*time.Nanosecond, true},

	{"5000000h", 0, false},
	{"-9223372036854775808ns", 0, false},
	{"9223372036854775808ns", 0, false},
	{"3000000h", 0, false},
	{"9223372036854775808ns", 0, false},
	{"-9223372036854775808ns", 0, false},
	{"9223372036854775.808us", 0, false},
	{"9223372036854ms775μs808ns", 0, false},
}

func TestParseDuration(t *testing.T) {
	for _, tc := range parseDurationTests {
		parseDuration, err := ParseDuration(tc.input)
		if tc.ok && (err != nil || parseDuration != tc.duration) {
			t.Errorf("ParseDuration(%q) = (%v, %v), want (%v, nil)", tc.input, parseDuration, err, tc.duration)
		} else if !tc.ok && err == nil {
			t.Errorf("ParseDuration(%q) = (_, nil), want (_, non-nil)", tc.input)
		}
	}
}

func TestParseDurationRoundTrip(t *testing.T) {
	for i := 0; i < 100; i++ {
		first := time.Duration(rand.Int31()) * time.Millisecond
		s := first.String()
		second, err := ParseDuration(s)
		if err != nil || first != second {
			t.Errorf("round-trip failed: %d => %q => %d, %v", first, s, second, err)
		}
	}
}
