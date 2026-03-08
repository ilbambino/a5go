package core

import (
	"a5go/internal/testutil"
	"testing"
)

func TestPentagonConstants(t *testing.T) {
	if A != 72 || B != 127.94543761193603 || C != 108 || D != 82.29202980963508 || E != 149.7625318412527 {
		t.Fatalf("unexpected angle constants")
	}

	expectedPentagon := [][2]float64{
		{0, 0},
		{0.1993818474311588, 0.3754138223914238},
		{0.6180339887498949, 0.4490279765795854},
		{0.8174158361810537, 0.0736141541881617},
		{0.418652141318736, -0.07361415418816161},
	}
	gotPentagon := PentagonShapeDef.GetVertices()
	for i, vertex := range gotPentagon {
		testutil.RequireClose(t, vertex[0], expectedPentagon[i][0], 1e-15)
		testutil.RequireClose(t, vertex[1], expectedPentagon[i][1], 1e-15)
	}

	if u != (Face{0, 0}) {
		t.Fatalf("unexpected u: %v", u)
	}
	testutil.RequireClose(t, v[0], 0.6180339887498949, 1e-15)
	testutil.RequireClose(t, v[1], 0.4490279765795854, 1e-15)
	testutil.RequireClose(t, w[0], 0.6180339887498949, 1e-15)
	testutil.RequireClose(t, w[1], -0.4490279765795854, 1e-15)
	testutil.RequireClose(t, float64(V), 0.6283185307179586, 1e-15)

	expectedTriangle := [][2]float64{
		{0, 0},
		{0.6180339887498949, 0.4490279765795854},
		{0.6180339887498949, -0.4490279765795854},
	}
	gotTriangle := TriangleShapeDef.GetVertices()
	for i, vertex := range gotTriangle {
		testutil.RequireClose(t, vertex[0], expectedTriangle[i][0], 1e-15)
		testutil.RequireClose(t, vertex[1], expectedTriangle[i][1], 1e-15)
	}

	expectedBasis := []float64{
		0.6180339887498949,
		0.4490279765795854,
		0.6180339887498949,
		-0.4490279765795854,
	}
	expectedInverse := []float64{
		0.8090169943749475,
		0.8090169943749475,
		1.1135163644116068,
		-1.1135163644116068,
	}
	for i, value := range Basis {
		testutil.RequireClose(t, value, expectedBasis[i], 1e-15)
	}
	for i, value := range BasisInverse {
		testutil.RequireClose(t, value, expectedInverse[i], 1e-15)
	}

	product := Mat2Multiply(Basis, BasisInverse)
	testutil.RequireCloseSlice(t, product[:], []float64{1, 0, 0, 1}, 1e-10)
}
