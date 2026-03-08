package core_test

import (
	"a5go/core"
	"testing"
)

func TestHexToU64(t *testing.T) {
	cases := map[string]uint64{
		"1a2b3c":   1715004,
		"0":        0,
		"ff":       255,
		"ffffffff": 4294967295,
	}

	for input, want := range cases {
		got, err := core.HexToU64(input)
		if err != nil {
			t.Fatalf("HexToU64(%q) error: %v", input, err)
		}
		if got != want {
			t.Fatalf("HexToU64(%q) = %d want %d", input, got, want)
		}
	}
}

func TestU64ToHex(t *testing.T) {
	cases := map[uint64]string{
		1715004:    "1a2b3c",
		0:          "0",
		255:        "ff",
		4294967295: "ffffffff",
	}

	for input, want := range cases {
		if got := core.U64ToHex(input); got != want {
			t.Fatalf("U64ToHex(%d) = %q want %q", input, got, want)
		}
	}
}

func TestHexRoundTrip(t *testing.T) {
	values := []string{"1a2b3c", "0", "ff", "ffffffff"}
	for _, value := range values {
		index, err := core.HexToU64(value)
		if err != nil {
			t.Fatalf("HexToU64(%q) error: %v", value, err)
		}
		if got := core.U64ToHex(index); got != value {
			t.Fatalf("round trip %q => %q", value, got)
		}
	}
}
