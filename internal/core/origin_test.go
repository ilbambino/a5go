package core

import (
	"a5go/internal/testutil"
	"math"
	"testing"
)

type originFixture struct {
	ID            int        `json:"id"`
	Axis          [2]float64 `json:"axis"`
	Quat          [4]float64 `json:"quat"`
	Angle         float64    `json:"angle"`
	Orientation   []string   `json:"orientation"`
	FirstQuintant int        `json:"firstQuintant"`
}

func TestOrigins(t *testing.T) {
	var expected []originFixture
	testutil.LoadJSON(t, "../../testdata/fixtures/origins.json", &expected)

	if len(Origins) != 12 || len(Origins) != len(expected) {
		t.Fatalf("origin count mismatch")
	}

	for i, origin := range Origins {
		exp := expected[i]
		if origin.ID != exp.ID {
			t.Fatalf("origin id mismatch")
		}
		testutil.RequireCloseSlice(t, origin.Axis[:], exp.Axis[:], 1e-12)
		testutil.RequireCloseSlice(t, origin.Quat[:], exp.Quat[:], 1e-12)
		testutil.RequireClose(t, float64(origin.Angle), exp.Angle, 1e-12)
		if origin.FirstQuintant != exp.FirstQuintant {
			t.Fatalf("first quintant mismatch")
		}
		cartesian := ToCartesian(origin.Axis)
		testutil.RequireClose(t, math.Sqrt(cartesian[0]*cartesian[0]+cartesian[1]*cartesian[1]+cartesian[2]*cartesian[2]), 1, 1e-12)
		testutil.RequireClose(t, QuaternionLength(origin.Quat), 1, 1e-12)
	}
}

func TestFindNearestOrigin(t *testing.T) {
	for _, origin := range Origins {
		if nearest := FindNearestOrigin(origin.Axis); nearest != origin {
			t.Fatalf("nearest origin mismatch for %d", origin.ID)
		}
	}

	boundaryPoints := []struct {
		Point           Spherical
		ExpectedOrigins []int
	}{
		{Spherical{0, float64(PiOver5) / 2}, []int{0, 1}},
		{Spherical{2 * float64(PiOver5), float64(PiOver5)}, []int{3, 4}},
		{Spherical{0, math.Pi - float64(PiOver5)/2}, []int{9, 10}},
	}
	for _, testCase := range boundaryPoints {
		nearest := FindNearestOrigin(testCase.Point)
		found := false
		for _, expected := range testCase.ExpectedOrigins {
			if nearest.ID == expected {
				found = true
			}
		}
		if !found {
			t.Fatalf("unexpected boundary nearest origin: %d", nearest.ID)
		}
	}
}

func TestHaversineAndConversions(t *testing.T) {
	point := Spherical{0, 0}
	if Haversine(point, point) != 0 {
		t.Fatalf("expected zero haversine")
	}
	point2 := Spherical{math.Pi / 4, math.Pi / 3}
	if Haversine(point2, point2) != 0 {
		t.Fatalf("expected zero haversine")
	}

	p1 := Spherical{0, math.Pi / 4}
	p2 := Spherical{math.Pi / 2, math.Pi / 3}
	testutil.RequireClose(t, Haversine(p1, p2), Haversine(p2, p1), 1e-12)

	distances := []Spherical{{0, math.Pi / 6}, {0, math.Pi / 4}, {0, math.Pi / 3}, {0, math.Pi / 2}}
	lastDistance := 0.0
	for _, p := range distances {
		distance := Haversine(point, p)
		if distance <= lastDistance {
			t.Fatalf("expected increasing haversine")
		}
		lastDistance = distance
	}

	lat := math.Pi / 4
	d1 := Haversine(Spherical{0, lat}, Spherical{math.Pi, lat})
	d2 := Haversine(Spherical{0, lat}, Spherical{math.Pi / 2, lat})
	if d1 <= d2 {
		t.Fatalf("expected larger longitudinal separation to increase distance")
	}

	testutil.RequireClose(t, Haversine(Spherical{0, 0}, Spherical{0, math.Pi / 2}), 0.5, 1e-4)
	testutil.RequireClose(t, Haversine(Spherical{0, math.Pi / 4}, Spherical{math.Pi / 2, math.Pi / 4}), 0.25, 1e-4)

	origin := Origins[0]
	for quintant := 0; quintant < 5; quintant++ {
		segment, _ := QuintantToSegment(quintant, origin)
		roundTripQuintant, _ := SegmentToQuintant(segment, origin)
		if roundTripQuintant != quintant {
			t.Fatalf("round trip quintant mismatch")
		}
	}

	for _, origin := range Origins {
		nearest := FindNearestOrigin(origin.Axis)
		if nearest != origin {
			t.Fatalf("expected origin to be nearest to itself")
		}
	}

	boundaryChecks := []struct {
		Point  Spherical
		Origin *Origin
	}{
		{Spherical{0, float64(PiOver5) / 2}, Origins[0]},
		{Spherical{2 * float64(PiOver5), float64(PiOver5)}, Origins[3]},
		{Spherical{0, math.Pi - float64(PiOver5)/2}, Origins[9]},
	}
	for _, testCase := range boundaryChecks {
		if IsNearestOrigin(testCase.Point, testCase.Origin) {
			t.Fatalf("expected boundary point not to be nearest origin")
		}
	}
}
