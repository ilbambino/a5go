package projections_test

import (
	"a5go/core"
	"a5go/internal/testutil"
	"a5go/projections"
	"math"
	"testing"
)

type gnomonicFixture struct {
	Forward []struct {
		Input    [2]float64 `json:"input"`
		Expected [2]float64 `json:"expected"`
	} `json:"forward"`
	Inverse []struct {
		Input    [2]float64 `json:"input"`
		Expected [2]float64 `json:"expected"`
	} `json:"inverse"`
}

func TestGnomonicProjection(t *testing.T) {
	var fixture gnomonicFixture
	testutil.LoadJSON(t, "../testdata/fixtures/projections/gnomonic.json", &fixture)
	gnomonic := projections.GnomonicProjection{}
	for _, testCase := range fixture.Forward {
		result := gnomonic.Forward(core.Spherical(testCase.Input))
		requireClosePair(t, result[:], testCase.Expected[:])
		roundTrip := gnomonic.Inverse(result)
		requireClosePair(t, roundTrip[:], testCase.Input[:])
	}
	for _, testCase := range fixture.Inverse {
		result := gnomonic.Inverse(core.Polar(testCase.Input))
		requireClosePair(t, result[:], testCase.Expected[:])
		roundTrip := gnomonic.Forward(result)
		requireClosePair(t, roundTrip[:], testCase.Input[:])
	}
}

func requireClosePair(t *testing.T, got, want []float64) {
	t.Helper()
	for i := range got {
		diff := math.Abs(got[i] - want[i])
		limit := 1e-12
		if math.Abs(want[i]) > 1 {
			limit = math.Abs(want[i]) * 1e-12
		}
		if diff > limit {
			t.Fatalf("index %d: got %.16f want %.16f tolerance %.16f", i, got[i], want[i], limit)
		}
	}
}
