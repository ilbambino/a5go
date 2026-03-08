package geometry_test

import (
	"a5go/internal/geometry"
	"a5go/internal/testutil"
	"testing"
)

type pentagonFixture struct {
	Vertices           [][2]float64 `json:"vertices"`
	ContainsPointTests []struct {
		Point  [2]float64 `json:"point"`
		Result float64    `json:"result"`
	} `json:"containsPointTests"`
	Area           float64    `json:"area"`
	Center         [2]float64 `json:"center"`
	TransformTests struct {
		Scale     [][2]float64 `json:"scale"`
		Rotate180 [][2]float64 `json:"rotate180"`
		ReflectY  [][2]float64 `json:"reflectY"`
		Translate [][2]float64 `json:"translate"`
	} `json:"transformTests"`
	SplitEdgesTests struct {
		Segments2 [][2]float64 `json:"segments2"`
		Segments3 [][2]float64 `json:"segments3"`
	} `json:"splitEdgesTests"`
}

func toPentagon(vertices [][2]float64) geometry.Pentagon {
	result := make(geometry.Pentagon, len(vertices))
	for i, vertex := range vertices {
		result[i] = geometry.Face(vertex)
	}
	return result
}

func assertVertices(t *testing.T, got geometry.Pentagon, want [][2]float64) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("vertex count mismatch: got %d want %d", len(got), len(want))
	}
	for i, vertex := range got {
		testutil.RequireClose(t, vertex[0], want[i][0], 1e-6)
		testutil.RequireClose(t, vertex[1], want[i][1], 1e-6)
	}
}

func TestPentagonShapeFixtures(t *testing.T) {
	var fixtures []pentagonFixture
	testutil.LoadJSON(t, "../../testdata/fixtures/geometry/pentagon.json", &fixtures)

	for _, fixture := range fixtures {
		pentagon := geometry.NewPentagonShape(toPentagon(fixture.Vertices))

		for _, contains := range fixture.ContainsPointTests {
			got := pentagon.ContainsPoint(geometry.Face(contains.Point))
			testutil.RequireClose(t, got, contains.Result, 1e-6)
		}

		testutil.RequireClose(t, pentagon.GetArea(), fixture.Area, 1e-6)

		center := pentagon.GetCenter()
		testutil.RequireClose(t, center[0], fixture.Center[0], 1e-6)
		testutil.RequireClose(t, center[1], fixture.Center[1], 1e-6)

		assertVertices(t, pentagon.Clone().Scale(2).GetVertices(), fixture.TransformTests.Scale)
		assertVertices(t, pentagon.Clone().Rotate180().GetVertices(), fixture.TransformTests.Rotate180)
		assertVertices(t, pentagon.Clone().ReflectY().GetVertices(), fixture.TransformTests.ReflectY)
		assertVertices(t, pentagon.Clone().Translate(geometry.Face{1, 1}).GetVertices(), fixture.TransformTests.Translate)
		assertVertices(t, pentagon.Clone().SplitEdges(2).GetVertices(), fixture.SplitEdgesTests.Segments2)
		assertVertices(t, pentagon.Clone().SplitEdges(3).GetVertices(), fixture.SplitEdgesTests.Segments3)
	}
}
