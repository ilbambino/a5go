package traversal_test

import (
	"a5go/internal/core"
	"a5go/internal/testutil"
	"a5go/internal/traversal"
	"testing"
)

type capFixture struct {
	SphericalCap []struct {
		CellID string   `json:"cellId"`
		Radius float64  `json:"radius"`
		Cells  []string `json:"cells"`
	} `json:"sphericalCap"`
	SphericalCapCompact []struct {
		CellID         string   `json:"cellId"`
		Radius         float64  `json:"radius"`
		CompactedCells []string `json:"compactedCells"`
	} `json:"sphericalCapCompact"`
	Helpers struct {
		MetersToH []struct {
			Meters    float64 `json:"meters"`
			ExpectedH float64 `json:"expectedH"`
		} `json:"metersToH"`
		EstimateCellRadius []struct {
			Resolution     int     `json:"resolution"`
			ExpectedMeters float64 `json:"expectedMeters"`
		} `json:"estimateCellRadius"`
		PickCoarseResolution []struct {
			Radius            float64 `json:"radius"`
			TargetRes         int     `json:"targetRes"`
			ExpectedCoarseRes int     `json:"expectedCoarseRes"`
		} `json:"pickCoarseResolution"`
	} `json:"helpers"`
}

func TestSphericalCap(t *testing.T) {
	var fixtures capFixture
	testutil.LoadJSON(t, "../../testdata/fixtures/traversal/cap.json", &fixtures)

	for _, f := range fixtures.SphericalCap {
		cellID := parseHexCap(t, f.CellID)
		targetRes := core.GetResolution(cellID)
		result := core.Uncompact(traversal.SphericalCap(cellID, f.Radius), targetRes)
		got := make([]string, len(result))
		for i, cell := range result {
			got[i] = core.U64ToHex(cell)
		}
		assertStringSlicesEqual(t, got, f.Cells)
	}

	for _, f := range fixtures.SphericalCapCompact {
		cellID := parseHexCap(t, f.CellID)
		result := traversal.SphericalCap(cellID, f.Radius)
		got := make([]string, len(result))
		for i, cell := range result {
			got[i] = core.U64ToHex(cell)
		}
		assertStringSlicesEqual(t, got, f.CompactedCells)
	}
}

func TestCapHelpers(t *testing.T) {
	var fixtures capFixture
	testutil.LoadJSON(t, "../../testdata/fixtures/traversal/cap.json", &fixtures)

	for _, f := range fixtures.Helpers.MetersToH {
		testutil.RequireClose(t, traversal.MetersToH(f.Meters), f.ExpectedH, 1e-15)
	}

	for i, f := range fixtures.Helpers.EstimateCellRadius {
		testutil.RequireClose(t, traversal.EstimateCellRadius(f.Resolution), f.ExpectedMeters, 1e-9)
		if i > 0 && f.ExpectedMeters >= fixtures.Helpers.EstimateCellRadius[i-1].ExpectedMeters {
			t.Fatalf("cell radius did not decrease at index %d", i)
		}
	}

	for _, f := range fixtures.Helpers.PickCoarseResolution {
		got := traversal.PickCoarseResolution(f.Radius, f.TargetRes)
		if got != f.ExpectedCoarseRes {
			t.Fatalf("coarse resolution mismatch: got %d want %d", got, f.ExpectedCoarseRes)
		}
		if got > f.TargetRes {
			t.Fatalf("coarse resolution exceeded target: got %d target %d", got, f.TargetRes)
		}
	}
}

func parseHexCap(t *testing.T, hex string) uint64 {
	t.Helper()
	value, err := core.HexToU64(hex)
	if err != nil {
		t.Fatalf("parse hex %s: %v", hex, err)
	}
	return value
}

func assertStringSlicesEqual(t *testing.T, got, want []string) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("length mismatch: got %d want %d", len(got), len(want))
	}
	for i := range got {
		if got[i] != want[i] {
			t.Fatalf("index %d: got %s want %s", i, got[i], want[i])
		}
	}
}
