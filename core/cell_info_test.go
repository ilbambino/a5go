package core

import (
	"a5go/internal/testutil"
	"testing"
)

type cellInfoFixture struct {
	NumCells []struct {
		Resolution  int     `json:"resolution"`
		Count       float64 `json:"count"`
		CountBigInt string  `json:"countBigInt"`
	} `json:"numCells"`
	NumChildren []struct {
		ParentResolution int     `json:"parentResolution"`
		ChildResolution  int     `json:"childResolution"`
		NumChildren      float64 `json:"numChildren"`
	} `json:"numChildren"`
	CellArea []struct {
		Resolution int     `json:"resolution"`
		AreaM2     float64 `json:"areaM2"`
	} `json:"cellArea"`
}

func TestCellInfo(t *testing.T) {
	var fixture cellInfoFixture
	testutil.LoadJSON(t, "../testdata/fixtures/cell-info.json", &fixture)

	for _, testCase := range fixture.NumCells {
		if GetNumCells(testCase.Resolution) != testCase.Count {
			t.Fatalf("GetNumCells(%d) mismatch", testCase.Resolution)
		}
	}
	for _, testCase := range fixture.NumChildren {
		if GetNumChildren(testCase.ParentResolution, testCase.ChildResolution) != testCase.NumChildren {
			t.Fatalf("GetNumChildren(%d, %d) mismatch", testCase.ParentResolution, testCase.ChildResolution)
		}
	}
	for _, testCase := range fixture.CellArea {
		if CellArea(testCase.Resolution) != testCase.AreaM2 {
			t.Fatalf("CellArea(%d) mismatch", testCase.Resolution)
		}
	}
}
