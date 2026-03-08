package traversal

import (
	"a5go/internal/core"
	"a5go/internal/testutil"
	"testing"
)

type globalNeighborFixture struct {
	Input struct {
		CellID string `json:"cellId"`
	} `json:"input"`
	Output struct {
		Neighbors     []string `json:"neighbors"`
		EdgeNeighbors []string `json:"edgeNeighbors"`
	} `json:"output"`
}

func hexStrings(values []uint64) []string {
	result := make([]string, len(values))
	for i, value := range values {
		result[i] = core.U64ToHex(value)
	}
	return result
}

func parseHexLocal(hex string) uint64 {
	value, err := core.HexToU64(hex)
	if err != nil {
		panic(err)
	}
	return value
}

func TestGlobalNeighbors(t *testing.T) {
	var fixtures []globalNeighborFixture
	testutil.LoadJSON(t, "../../testdata/fixtures/traversal/global-neighbors.json", &fixtures)

	for _, f := range fixtures {
		cellID := parseHexLocal(f.Input.CellID)
		result := hexStrings(GetGlobalCellNeighbors(cellID, false))
		if len(result) != len(f.Output.Neighbors) {
			t.Fatalf("global neighbor count mismatch")
		}
		for i := range result {
			if result[i] != f.Output.Neighbors[i] {
				t.Fatalf("global neighbor mismatch at %d", i)
			}
		}

		edgeResult := hexStrings(GetGlobalCellNeighbors(cellID, true))
		if len(edgeResult) != len(f.Output.EdgeNeighbors) {
			t.Fatalf("edge neighbor count mismatch")
		}
		for i := range edgeResult {
			if edgeResult[i] != f.Output.EdgeNeighbors[i] {
				t.Fatalf("edge neighbor mismatch at %d", i)
			}
		}
		if len(edgeResult) != 5 {
			t.Fatalf("expected 5 edge neighbors")
		}
	}
}
