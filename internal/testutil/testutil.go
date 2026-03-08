package testutil

import (
	"encoding/json"
	"math"
	"os"
	"testing"
)

func LoadJSON(t *testing.T, path string, target any) {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	if err := json.Unmarshal(data, target); err != nil {
		t.Fatalf("unmarshal %s: %v", path, err)
	}
}

func RequireClose(t *testing.T, got, want, tolerance float64) {
	t.Helper()
	if math.Abs(got-want) > tolerance {
		t.Fatalf("got %.16f want %.16f tolerance %.16f", got, want, tolerance)
	}
}

func RequireCloseSlice(t *testing.T, got, want []float64, tolerance float64) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("length mismatch: got %d want %d", len(got), len(want))
	}
	for i := range got {
		if math.Abs(got[i]-want[i]) > tolerance {
			t.Fatalf("index %d: got %.16f want %.16f tolerance %.16f", i, got[i], want[i], tolerance)
		}
	}
}
