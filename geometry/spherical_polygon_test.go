package geometry_test

import (
	"a5go/geometry"
	"a5go/internal/testutil"
	"math"
	"testing"
)

type sphericalPolygonFixture struct {
	Vertices   [][3]float64 `json:"vertices"`
	Boundary1  [][3]float64 `json:"boundary1"`
	Boundary2  [][3]float64 `json:"boundary2"`
	Boundary3  [][3]float64 `json:"boundary3"`
	SlerpTests []struct {
		T      float64    `json:"t"`
		Result [3]float64 `json:"result"`
	} `json:"slerpTests"`
	ContainsPointTests []struct {
		Point  [3]float64 `json:"point"`
		Result float64    `json:"result"`
	} `json:"containsPointTests"`
	Area float64 `json:"area"`
}

func toCartesianSlice(values [][3]float64) []geometry.Cartesian {
	result := make([]geometry.Cartesian, len(values))
	for i, value := range values {
		result[i] = geometry.Cartesian(value)
	}
	return result
}

func assertCartesianList(t *testing.T, got []geometry.Cartesian, want [][3]float64) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("length mismatch")
	}
	for i := range got {
		testutil.RequireCloseSlice(t, got[i][:], want[i][:], 1e-6)
	}
}

func TestSphericalPolygonFixtures(t *testing.T) {
	var fixtures []sphericalPolygonFixture
	testutil.LoadJSON(t, "../testdata/fixtures/geometry/spherical-polygon.json", &fixtures)

	for _, fixture := range fixtures {
		polygon := geometry.NewSphericalPolygonShape(toCartesianSlice(fixture.Vertices))
		assertCartesianList(t, polygon.GetBoundary(1, true), fixture.Boundary1)
		assertCartesianList(t, polygon.GetBoundary(2, true), fixture.Boundary2)
		assertCartesianList(t, polygon.GetBoundary(3, true), fixture.Boundary3)

		for _, testCase := range fixture.SlerpTests {
			actual := polygon.Slerp(testCase.T)
			testutil.RequireCloseSlice(t, actual[:], testCase.Result[:], 1e-6)
			length := math.Sqrt(actual[0]*actual[0] + actual[1]*actual[1] + actual[2]*actual[2])
			if math.Abs(length-1) >= 1e-10 {
				t.Fatalf("slerp result not normalized")
			}
		}

		for _, testCase := range fixture.ContainsPointTests {
			actual := polygon.ContainsPoint(geometry.Cartesian(testCase.Point))
			testutil.RequireClose(t, actual, testCase.Result, 1e-6)
		}

		area := polygon.GetArea()
		testutil.RequireClose(t, area, fixture.Area, 1e-6)
		if math.Abs(area) <= 0 || math.Abs(area) > 2*math.Pi {
			t.Fatalf("polygon area out of bounds")
		}
	}

	if geometry.NewSphericalPolygonShape(nil).GetArea() != 0 {
		t.Fatalf("expected empty polygon area 0")
	}
	if geometry.NewSphericalPolygonShape([]geometry.Cartesian{{1, 0, 0}}).GetArea() != 0 {
		t.Fatalf("expected single-point polygon area 0")
	}
	if geometry.NewSphericalPolygonShape([]geometry.Cartesian{{1, 0, 0}, {0, 1, 0}}).GetArea() != 0 {
		t.Fatalf("expected line polygon area 0")
	}
}
