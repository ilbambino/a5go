package core

import (
	"a5go/internal/testutil"
	"math"
	"testing"
)

func TestCoordinateTransforms(t *testing.T) {
	if DegToRad(180) != math.Pi || DegToRad(90) != math.Pi/2 || DegToRad(0) != 0 {
		t.Fatalf("DegToRad mismatch")
	}
	if RadToDeg(Radians(math.Pi)) != 180 || RadToDeg(Radians(math.Pi/2)) != 90 || RadToDeg(0) != 0 {
		t.Fatalf("RadToDeg mismatch")
	}

	testTriangle := FaceTriangle{{0, 0}, {1, 0}, {0, 1}}
	testPoints := []Face{{0.1, 0.1}, {0.25, 0.25}, {0.01, 0.8}}
	for _, point := range testPoints {
		bary := FaceToBarycentric(point, testTriangle)
		result := BarycentricToFace(bary, testTriangle)
		testutil.RequireCloseSlice(t, result[:], point[:], 1e-12)
		testutil.RequireClose(t, bary[0]+bary[1]+bary[2], 1, 1e-12)
	}

	vertices := []Face{testTriangle[0], testTriangle[1], testTriangle[2]}
	expectedBary := []Barycentric{{1, 0, 0}, {0, 1, 0}, {0, 0, 1}}
	for i, vertex := range vertices {
		bary := FaceToBarycentric(vertex, testTriangle)
		testutil.RequireCloseSlice(t, bary[:], expectedBary[i][:], 1e-12)
	}

	edgeMidpoints := []Face{{0.5, 0}, {0, 0.5}, {0.5, 0.5}}
	expectedMidpoints := []Barycentric{{0.5, 0.5, 0}, {0.5, 0, 0.5}, {0, 0.5, 0.5}}
	for i, point := range edgeMidpoints {
		bary := FaceToBarycentric(point, testTriangle)
		testutil.RequireCloseSlice(t, bary[:], expectedMidpoints[i][:], 1e-12)
		result := BarycentricToFace(bary, testTriangle)
		testutil.RequireCloseSlice(t, result[:], point[:], 1e-12)
	}

	northPole := ToCartesian(Spherical{0, 0})
	testutil.RequireCloseSlice(t, northPole[:], []float64{0, 0, 1}, 1e-12)
	equator0 := ToCartesian(Spherical{0, math.Pi / 2})
	testutil.RequireCloseSlice(t, equator0[:], []float64{1, 0, 0}, 1e-12)
	equator90 := ToCartesian(Spherical{math.Pi / 2, math.Pi / 2})
	testutil.RequireCloseSlice(t, equator90[:], []float64{0, 1, 0}, 1e-12)
	spherical := Spherical{math.Pi / 4, math.Pi / 6}
	roundTripSpherical := ToSpherical(ToCartesian(spherical))
	testutil.RequireCloseSlice(t, roundTripSpherical[:], spherical[:], 1e-12)

	greenwich := FromLonLat(LonLat{0, 0})
	testutil.RequireCloseSlice(t, greenwich[:], []float64{float64(DegToRad(93)), math.Pi / 2}, 1e-12)
	northPoleLonLat := FromLonLat(LonLat{0, 90})
	testutil.RequireCloseSlice(t, northPoleLonLat[:], []float64{float64(DegToRad(93)), 0}, 1e-12)
	southPoleLonLat := FromLonLat(LonLat{0, -90})
	testutil.RequireCloseSlice(t, southPoleLonLat[:], []float64{float64(DegToRad(93)), math.Pi}, 1e-12)

	lonlatPoints := []LonLat{{0, 0}, {90, 0}, {180, 0}, {0, 45}, {0, -45}, {-90, -45}, {180, 45}, {90, 45}, {0, 90}, {0, -90}, {123, 45}}
	for _, point := range lonlatPoints {
		result := ToLonLatFromSpherical(FromLonLat(point))
		testutil.RequireCloseSlice(t, result[:], point[:], 1e-9)
	}

	contour := Contour{{0, 0}, {10, 0}, {10, 10}, {0, 10}, {0, 0}}
	normalized := NormalizeLongitudes(contour)
	for i := range contour {
		testutil.RequireCloseSlice(t, normalized[i][:], contour[i][:], 1e-12)
	}

	wrapped := Contour{{-170, 0}, {-175, 0}, {-180, 0}, {175, 0}, {170, 0}}
	normalized = NormalizeLongitudes(wrapped)
	testutil.RequireClose(t, float64(normalized[3][0]), -185, 1e-12)
	testutil.RequireClose(t, float64(normalized[4][0]), -190, 1e-12)
}
