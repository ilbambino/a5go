package projections_test

import (
	"a5go/internal/core"
	"a5go/internal/projections"
	"a5go/internal/testutil"
	"math"
	"testing"
)

type polyhedralFixture struct {
	Static struct {
		TestSphericalTriangle [][3]float64 `json:"TEST_SPHERICAL_TRIANGLE"`
		TestFaceTriangle      [][2]float64 `json:"TEST_FACE_TRIANGLE"`
	} `json:"static"`
	Forward []struct {
		Input    [3]float64 `json:"input"`
		Expected [2]float64 `json:"expected"`
	} `json:"forward"`
	Inverse []struct {
		Input    [2]float64 `json:"input"`
		Expected [3]float64 `json:"expected"`
	} `json:"inverse"`
}

func maxAngle(triangle []core.Cartesian) float64 {
	angles := []float64{
		math.Acos(dot3(triangle[0], triangle[1])),
		math.Acos(dot3(triangle[1], triangle[2])),
		math.Acos(dot3(triangle[2], triangle[0])),
	}
	max := angles[0]
	for _, angle := range angles[1:] {
		if angle > max {
			max = angle
		}
	}
	return max
}

func dot3(a, b core.Cartesian) float64 {
	return a[0]*b[0] + a[1]*b[1] + a[2]*b[2]
}

func distance3(a, b core.Cartesian) float64 {
	dx := a[0] - b[0]
	dy := a[1] - b[1]
	dz := a[2] - b[2]
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}

func TestPolyhedralProjection(t *testing.T) {
	var fixture polyhedralFixture
	testutil.LoadJSON(t, "../../testdata/fixtures/projections/polyhedral.json", &fixture)
	polyhedral := projections.PolyhedralProjection{}

	sphericalTriangle := make([]core.Cartesian, len(fixture.Static.TestSphericalTriangle))
	for i, point := range fixture.Static.TestSphericalTriangle {
		sphericalTriangle[i] = core.Cartesian(point)
	}
	faceTriangle := core.FaceTriangle{core.Face(fixture.Static.TestFaceTriangle[0]), core.Face(fixture.Static.TestFaceTriangle[1]), core.Face(fixture.Static.TestFaceTriangle[2])}

	authalicRadius := 6371.0072
	maxArcLengthMM := authalicRadius * maxAngle(sphericalTriangle) * 1e9
	largestError := 0.0

	for _, testCase := range fixture.Forward {
		result := polyhedral.Forward(core.Cartesian(testCase.Input), [3]core.Cartesian{sphericalTriangle[0], sphericalTriangle[1], sphericalTriangle[2]}, faceTriangle)
		testutil.RequireCloseSlice(t, result[:], testCase.Expected[:], 1e-9)
		roundTrip := polyhedral.Inverse(result, faceTriangle, [3]core.Cartesian{sphericalTriangle[0], sphericalTriangle[1], sphericalTriangle[2]})
		d := distance3(roundTrip, core.Cartesian(testCase.Input))
		if d > largestError {
			largestError = d
		}
		testutil.RequireCloseSlice(t, roundTrip[:], testCase.Input[:], 1e-9)
	}
	if largestError*maxArcLengthMM >= 0.01 {
		t.Fatalf("polyhedral error too large: %f mm", largestError*maxArcLengthMM)
	}

	for _, testCase := range fixture.Inverse {
		result := polyhedral.Inverse(core.Face(testCase.Input), faceTriangle, [3]core.Cartesian{sphericalTriangle[0], sphericalTriangle[1], sphericalTriangle[2]})
		testutil.RequireCloseSlice(t, result[:], testCase.Expected[:], 1e-9)
		roundTrip := polyhedral.Forward(result, [3]core.Cartesian{sphericalTriangle[0], sphericalTriangle[1], sphericalTriangle[2]}, faceTriangle)
		testutil.RequireCloseSlice(t, roundTrip[:], testCase.Input[:], 1e-9)
	}
}
