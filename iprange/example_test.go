package iprange_test

import (
	"testing"

	"gitlab.com/slon/shad-go/iprange"
)

func FuzzIPRangeParsing(f *testing.F) {
	inputs := []string{
		"10.0.0.1",
		"10.0.0.5-10",
		"192.168.1.*",
		"192.168.10.0/24",
		"172.16.0.0/12",
		"255.255.255.255",
		"0.0.0.0",
		"127.0.0.1",
		"10.0.0.1-10.0.0.10",
		"10.0.0.1/30",
		"192.168.0.1-192.168.0.50",
		"192.168.100.0/25",
	}
	for _, input := range inputs {
		f.Add(input)
	}

	f.Fuzz(func(t *testing.T, data string) {
		defer getError(t, data)
		if _, err := iprange.Parse(data); err != nil {
			t.Logf("Failed to parse %q: %v", data, err)
		}
	})
}

func getError(t *testing.T, data string) {
	if recoverInfo := recover(); recoverInfo != nil {
		t.Fatalf("Unexpected panic when parsing %q: %v", data, recoverInfo)
	}
}
