package lattice_test

import (
	"a5go/internal/testutil"
	"a5go/lattice"
	"testing"
)

type hilbertFixture struct {
	SToAnchor []struct {
		S           uint64     `json:"s"`
		Resolution  int        `json:"resolution"`
		Orientation string     `json:"orientation"`
		Q           uint8      `json:"q"`
		Offset      [2]float64 `json:"offset"`
		Flips       [2]int8    `json:"flips"`
	} `json:"sToAnchor"`
}

func TestHilbertFixtures(t *testing.T) {
	var fixture hilbertFixture
	testutil.LoadJSON(t, "../testdata/fixtures/lattice/hilbert.json", &fixture)

	for _, f := range fixture.SToAnchor {
		anchor := lattice.SToAnchor(f.S, f.Resolution, lattice.Orientation(f.Orientation))
		if anchor.Q != lattice.Quaternary(f.Q) {
			t.Fatalf("q mismatch for s=%d res=%d ori=%s", f.S, f.Resolution, f.Orientation)
		}
		if anchor.Offset[0] != f.Offset[0] || anchor.Offset[1] != f.Offset[1] {
			t.Fatalf("offset mismatch for s=%d res=%d ori=%s: got %v want %v", f.S, f.Resolution, f.Orientation, anchor.Offset, f.Offset)
		}
		if anchor.Flips[0] != lattice.Flip(f.Flips[0]) || anchor.Flips[1] != lattice.Flip(f.Flips[1]) {
			t.Fatalf("flips mismatch for s=%d res=%d ori=%s", f.S, f.Resolution, f.Orientation)
		}

		s := lattice.AnchorToS(lattice.Anchor{
			Q:      lattice.Quaternary(f.Q),
			Offset: lattice.IJ(f.Offset),
			Flips:  [2]lattice.Flip{lattice.Flip(f.Flips[0]), lattice.Flip(f.Flips[1])},
		}, f.Resolution, lattice.Orientation(f.Orientation))
		if s != f.S {
			t.Fatalf("AnchorToS mismatch for res=%d ori=%s: got %d want %d", f.Resolution, f.Orientation, s, f.S)
		}
	}
}
