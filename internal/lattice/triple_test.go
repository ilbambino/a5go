package lattice_test

import (
	"a5go/internal/lattice"
	"a5go/internal/testutil"
	"testing"
)

type tripleFixtures struct {
	AnchorToTriple []struct {
		S           uint64 `json:"s"`
		Resolution  int    `json:"resolution"`
		Orientation string `json:"orientation"`
		X           int    `json:"x"`
		Y           int    `json:"y"`
		Z           int    `json:"z"`
		Parity      int    `json:"parity"`
	} `json:"anchorToTriple"`
	TripleInBounds []struct {
		X        int  `json:"x"`
		Y        int  `json:"y"`
		Z        int  `json:"z"`
		MaxRow   int  `json:"maxRow"`
		Expected bool `json:"expected"`
	} `json:"tripleInBounds"`
}

func TestTripleFixtures(t *testing.T) {
	var fixture tripleFixtures
	testutil.LoadJSON(t, "../../testdata/fixtures/lattice/triple.json", &fixture)

	for _, f := range fixture.AnchorToTriple {
		anchor := lattice.SToAnchor(f.S, f.Resolution, lattice.Orientation(f.Orientation))
		triple := lattice.AnchorToTriple(anchor)
		if triple.X != f.X || triple.Y != f.Y || triple.Z != f.Z {
			t.Fatalf("AnchorToTriple mismatch for s=%d res=%d ori=%s: got %+v", f.S, f.Resolution, f.Orientation, triple)
		}

		if lattice.TripleParity(lattice.Triple{X: f.X, Y: f.Y, Z: f.Z}) != f.Parity {
			t.Fatalf("TripleParity mismatch for (%d,%d,%d)", f.X, f.Y, f.Z)
		}

		s := lattice.TripleToS(lattice.Triple{X: f.X, Y: f.Y, Z: f.Z}, f.Resolution, lattice.Orientation(f.Orientation))
		if s == nil || *s != f.S {
			t.Fatalf("TripleToS mismatch for (%d,%d,%d) res=%d ori=%s", f.X, f.Y, f.Z, f.Resolution, f.Orientation)
		}

		expected := lattice.SToAnchor(f.S, f.Resolution, lattice.Orientation(f.Orientation))
		actual := lattice.TripleToAnchor(lattice.Triple{X: f.X, Y: f.Y, Z: f.Z}, f.Resolution, lattice.Orientation(f.Orientation))
		if actual == nil {
			t.Fatalf("TripleToAnchor returned nil for (%d,%d,%d)", f.X, f.Y, f.Z)
		}
		if actual.Offset[0] != expected.Offset[0] || actual.Offset[1] != expected.Offset[1] || actual.Flips[0] != expected.Flips[0] || actual.Flips[1] != expected.Flips[1] {
			t.Fatalf("TripleToAnchor mismatch for (%d,%d,%d)", f.X, f.Y, f.Z)
		}
	}

	for _, f := range fixture.TripleInBounds {
		got := lattice.TripleInBounds(lattice.Triple{X: f.X, Y: f.Y, Z: f.Z}, f.MaxRow)
		if got != f.Expected {
			t.Fatalf("TripleInBounds mismatch for (%d,%d,%d) maxRow=%d", f.X, f.Y, f.Z, f.MaxRow)
		}
	}
}
