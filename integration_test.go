package a5go_test

import (
	"a5go"
	"a5go/internal/testutil"
	"testing"
)

type integrationGeoJSON struct {
	Features []struct {
		Properties struct {
			CellIDHex string `json:"cellIdHex"`
		} `json:"properties"`
		Geometry struct {
			Coordinates [][][]float64 `json:"coordinates"`
		} `json:"geometry"`
	} `json:"features"`
}

func TestWireframeIntegration(t *testing.T) {
	paths := []struct {
		path     string
		segments int
		auto     bool
	}{
		{"testdata/integration/wireframe-0.json", 1, false},
		{"testdata/integration/wireframe-1.json", 1, false},
		{"testdata/integration/wireframe-2.json", 1, false},
		{"testdata/integration/wireframe-3.json", 1, false},
		{"testdata/integration/wireframe-auto-edges-0.json", 0, true},
		{"testdata/integration/wireframe-auto-edges-1.json", 0, true},
		{"testdata/integration/wireframe-auto-edges-2.json", 0, true},
		{"testdata/integration/wireframe-auto-edges-3.json", 0, true},
	}

	for _, testCase := range paths {
		var fixture integrationGeoJSON
		testutil.LoadJSON(t, testCase.path, &fixture)
		for _, feature := range fixture.Features {
			cellID, err := a5go.HexToU64(feature.Properties.CellIDHex)
			if err != nil {
				t.Fatalf("parse hex %s: %v", feature.Properties.CellIDHex, err)
			}
			opts := a5go.CellBoundaryOptions{ClosedRing: true, Segments: testCase.segments, AutoSegments: testCase.auto}
			actualBoundary := a5go.CellToBoundary(cellID, opts)
			expectedBoundary := feature.Geometry.Coordinates[0]
			if len(actualBoundary) != len(expectedBoundary) {
				t.Fatalf("boundary length mismatch for %s: got %d want %d", feature.Properties.CellIDHex, len(actualBoundary), len(expectedBoundary))
			}
			for i := range actualBoundary {
				testutil.RequireCloseSlice(t, actualBoundary[i][:], expectedBoundary[i], 1e-6)
			}
		}
	}
}
