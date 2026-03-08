package projections_test

import (
	"a5go/internal/projections"
	"a5go/internal/testutil"
	"testing"
)

type authalicFixture struct {
	Forward []struct {
		Input    float64 `json:"input"`
		Expected float64 `json:"expected"`
	} `json:"forward"`
	Inverse []struct {
		Input    float64 `json:"input"`
		Expected float64 `json:"expected"`
	} `json:"inverse"`
}

func TestAuthalicProjection(t *testing.T) {
	var fixture authalicFixture
	testutil.LoadJSON(t, "../../testdata/fixtures/projections/authalic.json", &fixture)
	authalic := projections.AuthalicProjection{}

	for _, testCase := range fixture.Forward {
		result := authalic.Forward(testCase.Input)
		testutil.RequireClose(t, result, testCase.Expected, 1e-10)
		roundTrip := authalic.Inverse(result)
		testutil.RequireClose(t, roundTrip, testCase.Input, 1e-15)
	}
	for _, testCase := range fixture.Inverse {
		result := authalic.Inverse(testCase.Input)
		testutil.RequireClose(t, result, testCase.Expected, 1e-10)
		roundTrip := authalic.Forward(result)
		testutil.RequireClose(t, roundTrip, testCase.Input, 1e-15)
	}
}
