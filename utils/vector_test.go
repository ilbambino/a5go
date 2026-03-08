package utils_test

import (
	"a5go/core"
	"a5go/internal/testutil"
	"a5go/utils"
	"math"
	"testing"
)

func normalize(v core.Cartesian) core.Cartesian {
	length := math.Sqrt(v[0]*v[0] + v[1]*v[1] + v[2]*v[2])
	return core.Cartesian{v[0] / length, v[1] / length, v[2] / length}
}

func TestVectorDifference(t *testing.T) {
	if got := utils.VectorDifference(core.Cartesian{1, 0, 0}, core.Cartesian{1, 0, 0}); math.Abs(got) > 1e-6 {
		t.Fatalf("identical vectors: got %f", got)
	}

	got := utils.VectorDifference(core.Cartesian{1, 0, 0}, core.Cartesian{0, 1, 0})
	testutil.RequireClose(t, got, math.Sqrt(0.5), 1e-6)

	got = utils.VectorDifference(core.Cartesian{1, 0, 0}, normalize(core.Cartesian{0.999, 0.001, 0}))
	if !(got > 0 && got < 0.1) {
		t.Fatalf("small angle result out of range: %f", got)
	}
}

func TestQuadrupleProduct(t *testing.T) {
	a := core.Cartesian{1, 0, 0}
	b := core.Cartesian{0, 1, 0}
	c := core.Cartesian{0, 0, 1}
	d := normalize(core.Cartesian{1, 1, 1})

	var out core.Cartesian
	result := utils.QuadrupleProduct(&out, a, b, c, d)
	if result != &out {
		t.Fatalf("expected returned pointer to match out")
	}

	out = core.Cartesian{}
	result = utils.QuadrupleProduct(&out, a, b, c, a)
	if result != &out {
		t.Fatalf("expected returned pointer to match out")
	}
	if out[0] == 0 && out[1] == 0 && out[2] == 0 {
		t.Fatalf("expected non-zero quadruple product")
	}
}

func TestSlerp(t *testing.T) {
	a := core.Cartesian{1, 0, 0}
	b := core.Cartesian{0, 1, 0}

	var out core.Cartesian
	result := utils.Slerp(&out, a, b, 0.5)
	if result != &out {
		t.Fatalf("expected returned pointer to match out")
	}
	testutil.RequireCloseSlice(t, out[:], []float64{1 / math.Sqrt(2), 1 / math.Sqrt(2), 0}, 1e-6)

	utils.Slerp(&out, a, b, 0)
	testutil.RequireCloseSlice(t, out[:], []float64{1, 0, 0}, 1e-6)

	utils.Slerp(&out, a, b, 1)
	testutil.RequireCloseSlice(t, out[:], []float64{0, 1, 0}, 1e-6)

	utils.Slerp(&out, a, a, 0.5)
	testutil.RequireCloseSlice(t, out[:], []float64{1, 0, 0}, 1e-6)

	var out1, out2 core.Cartesian
	utils.Slerp(&out1, a, b, 0.25)
	utils.Slerp(&out2, a, b, 0.75)
	if out1[0] <= out1[1] {
		t.Fatalf("expected t=0.25 to be closer to A: %v", out1)
	}
	if out2[1] <= out2[0] {
		t.Fatalf("expected t=0.75 to be closer to B: %v", out2)
	}
}
