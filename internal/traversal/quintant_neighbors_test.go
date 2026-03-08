package traversal

import (
	"a5go/internal/lattice"
	"a5go/internal/testutil"
	"testing"
)

type quintantNeighborFixture struct {
	Input struct {
		S           uint64 `json:"s"`
		Resolution  int    `json:"resolution"`
		Orientation string `json:"orientation"`
	} `json:"input"`
	Output struct {
		Neighbors []uint64 `json:"neighbors"`
	} `json:"output"`
}

func TestQuintantNeighbors(t *testing.T) {
	var fixtures []quintantNeighborFixture
	testutil.LoadJSON(t, "../../testdata/fixtures/traversal/quintant-neighbors.json", &fixtures)

	for _, f := range fixtures {
		result := GetCellNeighbors(f.Input.S, f.Input.Resolution, lattice.Orientation(f.Input.Orientation), false)
		if len(result) != len(f.Output.Neighbors) {
			t.Fatalf("neighbor count mismatch")
		}
		for i := range result {
			if result[i] != f.Output.Neighbors[i] {
				t.Fatalf("neighbor mismatch at %d", i)
			}
		}
	}
}
