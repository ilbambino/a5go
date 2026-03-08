package lattice_test

import (
	"a5go/internal/lattice"
	"a5go/internal/testutil"
	"testing"
)

type quaternaryFixtures struct {
	IJToQuaternary []struct {
		IJ    [2]float64 `json:"ij"`
		Flips [2]int8    `json:"flips"`
		Digit uint8      `json:"digit"`
	} `json:"IJToQuaternary"`
	QuaternaryToKJ []struct {
		Q     uint8      `json:"q"`
		Flips [2]int8    `json:"flips"`
		KJ    [2]float64 `json:"kj"`
	} `json:"quaternaryToKJ"`
	QuaternaryToFlips []struct {
		Q     uint8   `json:"q"`
		Flips [2]int8 `json:"flips"`
	} `json:"quaternaryToFlips"`
}

func TestQuaternaryFixtures(t *testing.T) {
	var fixtures quaternaryFixtures
	testutil.LoadJSON(t, "../../testdata/fixtures/lattice/quaternary.json", &fixtures)

	for _, f := range fixtures.IJToQuaternary {
		got := lattice.IJToQuaternary(lattice.IJ(f.IJ), [2]lattice.Flip{lattice.Flip(f.Flips[0]), lattice.Flip(f.Flips[1])})
		if got != lattice.Quaternary(f.Digit) {
			t.Fatalf("IJToQuaternary(%v, %v) = %d want %d", f.IJ, f.Flips, got, f.Digit)
		}
	}

	for _, f := range fixtures.QuaternaryToKJ {
		got := lattice.QuaternaryToKJ(lattice.Quaternary(f.Q), [2]lattice.Flip{lattice.Flip(f.Flips[0]), lattice.Flip(f.Flips[1])})
		if got[0] != f.KJ[0] || got[1] != f.KJ[1] {
			t.Fatalf("QuaternaryToKJ(%d, %v) = %v want %v", f.Q, f.Flips, got, f.KJ)
		}
	}

	for _, f := range fixtures.QuaternaryToFlips {
		got := lattice.QuaternaryToFlips(lattice.Quaternary(f.Q))
		if got[0] != lattice.Flip(f.Flips[0]) || got[1] != lattice.Flip(f.Flips[1]) {
			t.Fatalf("QuaternaryToFlips(%d) = %v want %v", f.Q, got, f.Flips)
		}
	}
}
