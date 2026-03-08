package core

import (
	"a5go/internal/geometry"
	"a5go/internal/lattice"
	"a5go/internal/testutil"
	"testing"
)

type tilingFixtures struct {
	GetPentagonVertices []struct {
		Input struct {
			Resolution int `json:"resolution"`
			Quintant   int `json:"quintant"`
			Anchor     struct {
				Q      uint8      `json:"q"`
				Offset [2]float64 `json:"offset"`
				Flips  [2]int8    `json:"flips"`
			} `json:"anchor"`
		} `json:"input"`
		Output struct {
			Vertices [][2]float64 `json:"vertices"`
			Area     float64      `json:"area"`
			Center   [2]float64   `json:"center"`
		} `json:"output"`
	} `json:"getPentagonVertices"`
	GetQuintantVertices []struct {
		Input struct {
			Quintant int `json:"quintant"`
		} `json:"input"`
		Output struct {
			Vertices [][2]float64 `json:"vertices"`
			Area     float64      `json:"area"`
			Center   [2]float64   `json:"center"`
		} `json:"output"`
	} `json:"getQuintantVertices"`
	GetFaceVertices struct {
		Vertices [][2]float64 `json:"vertices"`
		Area     float64      `json:"area"`
		Center   [2]float64   `json:"center"`
	} `json:"getFaceVertices"`
	GetQuintantPolar []struct {
		Input struct {
			Polar [2]float64 `json:"polar"`
		} `json:"input"`
		Output struct {
			Quintant int `json:"quintant"`
		} `json:"output"`
	} `json:"getQuintantPolar"`
}

func assertPolygon(t *testing.T, got [][2]float64, want [][2]float64, tolerance float64) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("vertex count mismatch: got %d want %d", len(got), len(want))
	}
	for i := range got {
		testutil.RequireClose(t, got[i][0], want[i][0], tolerance)
		testutil.RequireClose(t, got[i][1], want[i][1], tolerance)
	}
}

func toFloats(vertices geometry.Pentagon) [][2]float64 {
	result := make([][2]float64, len(vertices))
	for i, vertex := range vertices {
		result[i] = [2]float64{vertex[0], vertex[1]}
	}
	return result
}

func TestTilingFixtures(t *testing.T) {
	var fixtures tilingFixtures
	testutil.LoadJSON(t, "../../testdata/fixtures/tiling.json", &fixtures)

	for _, testCase := range fixtures.GetPentagonVertices {
		anchor := lattice.Anchor{
			Q:      lattice.Quaternary(testCase.Input.Anchor.Q),
			Offset: lattice.IJ(testCase.Input.Anchor.Offset),
			Flips:  [2]lattice.Flip{lattice.Flip(testCase.Input.Anchor.Flips[0]), lattice.Flip(testCase.Input.Anchor.Flips[1])},
		}
		pentagon := GetPentagonVertices(testCase.Input.Resolution, testCase.Input.Quintant, anchor)
		assertPolygon(t, toFloats(pentagon.GetVertices()), testCase.Output.Vertices, 1e-15)
		testutil.RequireClose(t, pentagon.GetArea(), testCase.Output.Area, 1e-15)
		center := pentagon.GetCenter()
		testutil.RequireClose(t, center[0], testCase.Output.Center[0], 1e-15)
		testutil.RequireClose(t, center[1], testCase.Output.Center[1], 1e-15)
	}

	for _, testCase := range fixtures.GetQuintantVertices {
		pentagon := GetQuintantVertices(testCase.Input.Quintant)
		assertPolygon(t, toFloats(pentagon.GetVertices()), testCase.Output.Vertices, 1e-15)
		testutil.RequireClose(t, pentagon.GetArea(), testCase.Output.Area, 1e-15)
		center := pentagon.GetCenter()
		testutil.RequireClose(t, center[0], testCase.Output.Center[0], 1e-15)
		testutil.RequireClose(t, center[1], testCase.Output.Center[1], 1e-15)
	}

	faceVertices := GetFaceVertices()
	assertPolygon(t, toFloats(faceVertices.GetVertices()), fixtures.GetFaceVertices.Vertices, 1e-15)
	testutil.RequireClose(t, faceVertices.GetArea(), fixtures.GetFaceVertices.Area, 1e-15)
	center := faceVertices.GetCenter()
	testutil.RequireClose(t, center[0], fixtures.GetFaceVertices.Center[0], 1e-15)
	testutil.RequireClose(t, center[1], fixtures.GetFaceVertices.Center[1], 1e-15)

	for _, testCase := range fixtures.GetQuintantPolar {
		got := GetQuintantPolar(Polar(testCase.Input.Polar))
		if got != testCase.Output.Quintant {
			t.Fatalf("GetQuintantPolar(%v) = %d want %d", testCase.Input.Polar, got, testCase.Output.Quintant)
		}
	}
}
