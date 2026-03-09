package cells_test

import (
	"a5go/internal/cells"
	"a5go/internal/core"
	"a5go/internal/testutil"
	"testing"
)

type populatedPlacesFixture struct {
	Features []struct {
		Properties struct {
			Name string `json:"name"`
		} `json:"properties"`
		Geometry struct {
			Coordinates [2]float64 `json:"coordinates"`
		} `json:"geometry"`
	} `json:"features"`
}

func TestCellValidation(t *testing.T) {
	got, err := cells.LonLatToCell(core.LonLat{0, 0}, -1)
	if err != nil {
		t.Fatalf("world cell lookup error: %v", err)
	}
	if got != 0 {
		t.Fatalf("world cell mismatch: got %d", got)
	}
	if got := cells.CellToLonLat(0); got != (core.LonLat{0, 0}) {
		t.Fatalf("world cell center mismatch: got %v", got)
	}
	if got := cells.CellToBoundary(0); len(got) != 0 {
		t.Fatalf("world cell boundary mismatch: got %v", got)
	}
}

func TestAntimeridianCellBoundaries(t *testing.T) {
	cellIDs := []string{"eb60000000000000", "2e00000000000000"}
	for _, cellIDHex := range cellIDs {
		cellID, err := core.HexToU64(cellIDHex)
		if err != nil {
			t.Fatalf("parse hex %s: %v", cellIDHex, err)
		}
		boundaries := [][]core.LonLat{
			cells.CellToBoundary(cellID, cells.CellBoundaryOptions{ClosedRing: true, Segments: 1}),
			cells.CellToBoundary(cellID, cells.CellBoundaryOptions{ClosedRing: true, Segments: 10}),
			cells.CellToBoundary(cellID, cells.CellBoundaryOptions{ClosedRing: true, AutoSegments: true}),
		}
		for _, boundary := range boundaries {
			minLon, maxLon := boundary[0][0], boundary[0][0]
			for _, point := range boundary[1:] {
				if point[0] < minLon {
					minLon = point[0]
				}
				if point[0] > maxLon {
					maxLon = point[0]
				}
			}
			if maxLon-minLon >= 180 {
				t.Fatalf("antimeridian span too wide for %s: %.6f", cellIDHex, maxLon-minLon)
			}
		}
	}
}

func TestCellContainsOriginalPointForAllResolutions(t *testing.T) {
	var populatedPlaces populatedPlacesFixture
	testutil.LoadJSON(t, "../../testdata/data/ne_50m_populated_places_nameonly.json", &populatedPlaces)

	for _, feature := range populatedPlaces.Features {
		testLonLat := core.LonLat{feature.Geometry.Coordinates[0], feature.Geometry.Coordinates[1]}
		for resolution := 1; resolution <= core.MaxResolution; resolution++ {
			if resolution == core.MaxResolution || abs(testLonLat[1]) > 80 {
				continue
			}
			cellID, err := cells.LonLatToCell(testLonLat, resolution)
			if err != nil {
				t.Fatalf("LonLatToCell error at resolution %d: %v", resolution, err)
			}
			_ = cells.CellToBoundary(cellID)
			cell := core.Deserialize(cellID)
			if cells.A5CellContainsPoint(cell, testLonLat) < 0 {
				t.Fatalf("cell %s at resolution %d does not contain point %v (%s)", core.U64ToHex(cellID), resolution, testLonLat, feature.Properties.Name)
			}
		}
	}
}

func abs(v float64) float64 {
	if v < 0 {
		return -v
	}
	return v
}
