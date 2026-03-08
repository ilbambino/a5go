package projections_test

import (
	"a5go/internal/projections"
	"a5go/internal/testutil"
	"math"
	"testing"
)

func TestCRS(t *testing.T) {
	var expectedVertices [][3]float64
	testutil.LoadJSON(t, "../../testdata/fixtures/crs-vertices.json", &expectedVertices)

	crs := projections.NewCRS()
	vertices := crs.Vertices()
	if len(vertices) != 62 || len(vertices) != len(expectedVertices) {
		t.Fatalf("crs vertex count mismatch")
	}
	for i, vertex := range vertices {
		testutil.RequireCloseSlice(t, vertex[:], expectedVertices[i][:], 1e-12)
		length := math.Sqrt(vertex[0]*vertex[0] + vertex[1]*vertex[1] + vertex[2]*vertex[2])
		testutil.RequireClose(t, length, 1, 1e-15)
	}

	defer func() {
		if recover() == nil {
			t.Fatalf("expected GetVertex to panic for non-vertex")
		}
	}()
	_ = crs.GetVertex([3]float64{1, 0, 0})
}
