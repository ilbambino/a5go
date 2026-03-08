package geometry_test

import (
	"a5go/internal/geometry"
	"a5go/internal/testutil"
	"math"
	"testing"
)

type sphericalTriangleFixture = sphericalPolygonFixture

func mustPanicTriangle(t *testing.T, fn func()) {
	t.Helper()
	defer func() {
		if recover() == nil {
			t.Fatalf("expected panic")
		}
	}()
	fn()
}

func TestSphericalTriangleFixtures(t *testing.T) {
	mustPanicTriangle(t, func() { geometry.NewSphericalTriangleShape(nil) })
	mustPanicTriangle(t, func() { geometry.NewSphericalTriangleShape([]geometry.Cartesian{{1, 0, 0}, {0, 1, 0}}) })
	mustPanicTriangle(t, func() {
		geometry.NewSphericalTriangleShape([]geometry.Cartesian{{1, 0, 0}, {0, 1, 0}, {0, 0, 1}, {1, 1, 1}})
	})
	_ = geometry.NewSphericalTriangleShape([]geometry.Cartesian{{1, 0, 0}, {0, 1, 0}, {0, 0, 1}})

	var fixtures []sphericalTriangleFixture
	testutil.LoadJSON(t, "../../testdata/fixtures/geometry/spherical-triangle.json", &fixtures)

	for _, fixture := range fixtures {
		triangle := geometry.NewSphericalTriangleShape(toCartesianSlice(fixture.Vertices))
		assertCartesianList(t, triangle.GetBoundary(1, true), fixture.Boundary1)
		assertCartesianList(t, triangle.GetBoundary(2, true), fixture.Boundary2)
		assertCartesianList(t, triangle.GetBoundary(3, true), fixture.Boundary3)

		for _, testCase := range fixture.SlerpTests {
			actual := triangle.Slerp(testCase.T)
			testutil.RequireCloseSlice(t, actual[:], testCase.Result[:], 1e-6)
			length := math.Sqrt(actual[0]*actual[0] + actual[1]*actual[1] + actual[2]*actual[2])
			if math.Abs(length-1) >= 1e-10 {
				t.Fatalf("slerp result not normalized")
			}
		}
		for _, testCase := range fixture.ContainsPointTests {
			actual := triangle.ContainsPoint(geometry.Cartesian(testCase.Point))
			testutil.RequireClose(t, actual, testCase.Result, 1e-6)
		}
		area := triangle.GetArea()
		testutil.RequireClose(t, area, fixture.Area, 1e-6)
		if math.Abs(area) <= 0 || math.Abs(area) > 2*math.Pi {
			t.Fatalf("triangle area out of bounds")
		}
	}
}
