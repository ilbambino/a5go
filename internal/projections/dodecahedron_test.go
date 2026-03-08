package projections_test

import (
	"a5go/internal/core"
	"a5go/internal/projections"
	"a5go/internal/testutil"
	"testing"
)

type dodecahedronFixture struct {
	Static struct {
		OriginID int `json:"ORIGIN_ID"`
	} `json:"static"`
	Forward []struct {
		Input    [2]float64 `json:"input"`
		Expected [2]float64 `json:"expected"`
	} `json:"forward"`
	Inverse []struct {
		Input    [2]float64 `json:"input"`
		Expected [2]float64 `json:"expected"`
	} `json:"inverse"`
}

func TestDodecahedronProjection(t *testing.T) {
	var fixture dodecahedronFixture
	testutil.LoadJSON(t, "../../testdata/fixtures/projections/dodecahedron.json", &fixture)

	dodecahedron := projections.NewDodecahedronProjection()
	originID := fixture.Static.OriginID

	for _, testCase := range fixture.Forward {
		result := dodecahedron.Forward(core.Spherical(testCase.Input), originID)
		testutil.RequireCloseSlice(t, result[:], testCase.Expected[:], 1e-9)
		roundTrip := dodecahedron.Inverse(result, originID)
		testutil.RequireCloseSlice(t, roundTrip[:], testCase.Input[:], 1e-9)
	}

	for _, testCase := range fixture.Inverse {
		result := dodecahedron.Inverse(core.Face(testCase.Input), originID)
		testutil.RequireCloseSlice(t, result[:], testCase.Expected[:], 1e-9)
		roundTrip := dodecahedron.Forward(result, originID)
		testutil.RequireCloseSlice(t, roundTrip[:], testCase.Input[:], 1e-9)
	}
}
