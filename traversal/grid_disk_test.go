package traversal

import (
	"a5go/core"
	"a5go/internal/testutil"
	"sort"
	"testing"
)

type gridDiskFixture struct {
	CellID           string   `json:"cellId"`
	K                int      `json:"k"`
	Cells            []string `json:"cells"`
	ExtraVertexCells []string `json:"extraVertexCells"`
}

func TestGridDisk(t *testing.T) {
	var fixtures []gridDiskFixture
	testutil.LoadJSON(t, "../testdata/fixtures/traversal/grid-disk.json", &fixtures)

	for _, f := range fixtures {
		cellID := parseHexLocal(f.CellID)
		targetRes := core.GetResolution(cellID)
		result := hexStrings(core.Uncompact(GridDisk(cellID, f.K), targetRes))
		if len(result) != len(f.Cells) {
			t.Fatalf("gridDisk length mismatch")
		}
		for i := range result {
			if result[i] != f.Cells[i] {
				t.Fatalf("gridDisk mismatch at %d", i)
			}
		}
	}

	cellID := parseHexLocal(fixtures[0].CellID)
	result := hexStrings(GridDisk(cellID, 0))
	if len(result) != 1 || result[0] != fixtures[0].CellID {
		t.Fatalf("gridDisk k=0 mismatch")
	}
}

func TestGridDiskVertex(t *testing.T) {
	var fixtures []gridDiskFixture
	testutil.LoadJSON(t, "../testdata/fixtures/traversal/grid-disk.json", &fixtures)

	for _, f := range fixtures {
		cellID := parseHexLocal(f.CellID)
		targetRes := core.GetResolution(cellID)
		expected := append(append([]string{}, f.Cells...), f.ExtraVertexCells...)
		sort.Slice(expected, func(i, j int) bool {
			return pad(expected[i], 20) < pad(expected[j], 20)
		})
		result := hexStrings(core.Uncompact(GridDiskVertex(cellID, f.K), targetRes))
		if len(result) != len(expected) {
			t.Fatalf("gridDiskVertex length mismatch")
		}
		for i := range result {
			if result[i] != expected[i] {
				t.Fatalf("gridDiskVertex mismatch at %d", i)
			}
		}
	}

	cellID := parseHexLocal(fixtures[0].CellID)
	result := hexStrings(GridDiskVertex(cellID, 0))
	if len(result) != 1 || result[0] != fixtures[0].CellID {
		t.Fatalf("gridDiskVertex k=0 mismatch")
	}
}

func pad(value string, width int) string {
	for len(value) < width {
		value = "0" + value
	}
	return value
}
