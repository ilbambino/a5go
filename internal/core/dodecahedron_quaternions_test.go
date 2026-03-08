package core

import (
	"a5go/internal/testutil"
	"math"
	"testing"
)

type quaternionFixture struct {
	Metadata struct {
		TotalQuaternions int `json:"totalQuaternions"`
	} `json:"metadata"`
	Quaternions []struct {
		Magnitude float64 `json:"magnitude"`
	} `json:"quaternions"`
	ValidationTests struct {
		AllNormalized     bool `json:"allNormalized"`
		NorthPoleIdentity bool `json:"northPoleIdentity"`
		SouthPoleCorrect  bool `json:"southPoleCorrect"`
		AllFinite         bool `json:"allFinite"`
	} `json:"validationTests"`
	Constants struct {
		INVSQRT5              float64 `json:"INV_SQRT5"`
		ExpectedPentagonAngle float64 `json:"expectedPentagonAngle"`
	} `json:"constants"`
}

func distance(a, b Cartesian) float64 {
	dx := a[0] - b[0]
	dy := a[1] - b[1]
	dz := a[2] - b[2]
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}

func dot(a, b Cartesian) float64 {
	return a[0]*b[0] + a[1]*b[1] + a[2]*b[2]
}

func TestDodecahedronQuaternions(t *testing.T) {
	var fixture quaternionFixture
	testutil.LoadJSON(t, "../../testdata/fixtures/dodecahedron-quaternions.json", &fixture)

	if len(Quaternions) != fixture.Metadata.TotalQuaternions {
		t.Fatalf("quaternion count mismatch")
	}
	for i, q := range Quaternions {
		testutil.RequireClose(t, QuaternionLength(q), fixture.Quaternions[i].Magnitude, 1e-10)
	}
	if Quaternions[0] != (Quaternion{0, 0, 0, 1}) || Quaternions[11] != (Quaternion{0, -1, 0, 0}) {
		t.Fatalf("pole quaternions mismatch")
	}

	cosAlpha := math.Sqrt((1 + math.Sqrt(0.2)) / 2)
	for i := 1; i <= 5; i++ {
		testutil.RequireClose(t, Quaternions[i][2], 0, 1e-15)
		testutil.RequireClose(t, Quaternions[i][3], cosAlpha, 1e-10)
	}
	sinAlpha := math.Sqrt((1 - math.Sqrt(0.2)) / 2)
	for i := 6; i <= 10; i++ {
		testutil.RequireClose(t, Quaternions[i][2], 0, 1e-15)
		testutil.RequireClose(t, Quaternions[i][3], sinAlpha, 1e-10)
	}

	northPole := Cartesian{0, 0, 1}
	faceCenters := make([]Cartesian, len(Quaternions))
	for i, q := range Quaternions {
		rotated := TransformQuat(northPole, q)
		faceCenters[i] = rotated
		testutil.RequireClose(t, math.Sqrt(dot(rotated, rotated)), 1, 1e-10)
		if i != 0 && distance(rotated, northPole) <= 0.1 {
			t.Fatalf("rotation %d too close to north pole", i)
		}
	}

	for i := 0; i < len(faceCenters); i++ {
		for j := i + 1; j < len(faceCenters); j++ {
			if distance(faceCenters[i], faceCenters[j]) <= 0.1 {
				t.Fatalf("face centers %d and %d overlap", i, j)
			}
		}
	}

	testVector := Cartesian{1, 0, 0}
	for _, q := range Quaternions {
		backRotated := TransformQuat(TransformQuat(testVector, q), QuaternionConjugate(q))
		testutil.RequireClose(t, distance(testVector, backRotated), 0, 1e-10)
	}

	v1 := Cartesian{1, 0, 0}
	v2 := Cartesian{0, 1, 0}
	for _, q := range Quaternions {
		rotated1 := TransformQuat(v1, q)
		rotated2 := TransformQuat(v2, q)
		testutil.RequireClose(t, dot(rotated1, rotated2), 0, 1e-10)
	}

	testVectors := []Cartesian{{1, 0, 0}, {0, 1, 0}, {0, 0, 1}, {1, 1, 1}}
	for _, q := range Quaternions {
		for _, v := range testVectors {
			testutil.RequireClose(t, math.Sqrt(dot(TransformQuat(v, q), TransformQuat(v, q))), math.Sqrt(dot(v, v)), 1e-10)
		}
	}

	zValues := make([]float64, len(faceCenters))
	for i, center := range faceCenters {
		zValues[i] = center[2]
	}
	sortFloatDesc(zValues)
	testutil.RequireClose(t, zValues[0], 1, 1e-10)
	testutil.RequireClose(t, zValues[11], -1, 1e-10)
	invSqrt5 := math.Sqrt(0.2)
	for _, z := range zValues[1:6] {
		testutil.RequireClose(t, z, invSqrt5, 1e-5)
	}
	for _, z := range zValues[6:11] {
		testutil.RequireClose(t, z, -invSqrt5, 1e-5)
	}

	firstRing := faceCenters[1:6]
	for i := 0; i < 5; i++ {
		next := (i + 1) % 5
		angle1 := math.Atan2(firstRing[i][1], firstRing[i][0])
		angle2 := math.Atan2(firstRing[next][1], firstRing[next][0])
		angleDiff := angle2 - angle1
		if angleDiff < 0 {
			angleDiff += 2 * math.Pi
		}
		if angleDiff > math.Pi {
			angleDiff = 2*math.Pi - angleDiff
		}
		testutil.RequireClose(t, angleDiff, 2*math.Pi/5, 1e-1)
	}
}

func sortFloatDesc(values []float64) {
	for i := 0; i < len(values); i++ {
		for j := i + 1; j < len(values); j++ {
			if values[j] > values[i] {
				values[i], values[j] = values[j], values[i]
			}
		}
	}
}
