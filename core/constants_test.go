package core_test

import (
	"a5go/core"
	"a5go/internal/testutil"
	"math"
	"testing"
)

type constantsFixture struct {
	Phi struct {
		Value         float64 `json:"value"`
		ExpectedValue float64 `json:"expectedValue"`
		Properties    struct {
			GoldenRatioSquared float64 `json:"goldenRatioSquared"`
			GoldenRatioPlusOne float64 `json:"goldenRatioPlusOne"`
			Reciprocal         float64 `json:"reciprocal"`
			ReciprocalMinusOne float64 `json:"reciprocalMinusOne"`
		} `json:"properties"`
	} `json:"φ"`
	Angles struct {
		TwoPi struct {
			Value         float64 `json:"value"`
			ExpectedValue float64 `json:"expectedValue"`
		} `json:"TWO_PI"`
		TwoPiOver5 struct {
			Value         float64 `json:"value"`
			ExpectedValue float64 `json:"expectedValue"`
		} `json:"TWO_PI_OVER_5"`
		PiOver5 struct {
			Value         float64 `json:"value"`
			ExpectedValue float64 `json:"expectedValue"`
		} `json:"PI_OVER_5"`
		PiOver10 struct {
			Value         float64 `json:"value"`
			ExpectedValue float64 `json:"expectedValue"`
		} `json:"PI_OVER_10"`
	} `json:"angles"`
	DodecahedronAngles struct {
		DihedralAngle struct {
			Value         float64 `json:"value"`
			ExpectedValue float64 `json:"expectedValue"`
		} `json:"dihedralAngle"`
		InterhedralAngle struct {
			Value         float64 `json:"value"`
			ExpectedValue float64 `json:"expectedValue"`
		} `json:"interhedralAngle"`
		FaceEdgeAngle struct {
			Value         float64 `json:"value"`
			ExpectedValue float64 `json:"expectedValue"`
		} `json:"faceEdgeAngle"`
		AngleSum float64 `json:"angleSum"`
	} `json:"dodecahedronAngles"`
	Distances struct {
		DistanceToEdge struct {
			Value              float64 `json:"value"`
			ExpectedValue      float64 `json:"expectedValue"`
			AlternativeFormula float64 `json:"alternativeFormula"`
		} `json:"distanceToEdge"`
		DistanceToVertex struct {
			Value              float64 `json:"value"`
			ExpectedValue      float64 `json:"expectedValue"`
			AlternativeFormula float64 `json:"alternativeFormula"`
		} `json:"distanceToVertex"`
	} `json:"distances"`
	SphereRadii struct {
		RInscribed struct {
			Value         float64 `json:"value"`
			ExpectedValue float64 `json:"expectedValue"`
		} `json:"Rinscribed"`
		RMidedge struct {
			Value         float64 `json:"value"`
			ExpectedValue float64 `json:"expectedValue"`
		} `json:"Rmidedge"`
		RCircumscribed struct {
			Value         float64 `json:"value"`
			ExpectedValue float64 `json:"expectedValue"`
		} `json:"Rcircumscribed"`
		Relationships struct {
			InscribedLessThanMidedge     bool `json:"inscribedLessThanMidedge"`
			MidedgeLessThanCircumscribed bool `json:"midedgeLessThanCircumscribed"`
		} `json:"relationships"`
	} `json:"sphereRadii"`
	ValidationTests struct {
		FiniteNumbers []struct {
			IsFinite bool `json:"isFinite"`
			IsNaN    bool `json:"isNaN"`
		} `json:"finiteNumbers"`
		PositiveConstants []struct {
			IsPositive bool `json:"isPositive"`
		} `json:"positiveConstants"`
	} `json:"validationTests"`
}

func TestConstants(t *testing.T) {
	var fixture constantsFixture
	testutil.LoadJSON(t, "../testdata/fixtures/constants.json", &fixture)

	testutil.RequireClose(t, core.Phi, fixture.Phi.ExpectedValue, 1e-15)
	if core.Phi != fixture.Phi.Value {
		t.Fatalf("phi mismatch: got %.16f want %.16f", core.Phi, fixture.Phi.Value)
	}
	testutil.RequireClose(t, fixture.Phi.Properties.GoldenRatioSquared, fixture.Phi.Properties.GoldenRatioPlusOne, 1e-15)
	testutil.RequireClose(t, fixture.Phi.Properties.Reciprocal, fixture.Phi.Properties.ReciprocalMinusOne, 1e-15)

	if float64(core.TwoPi) != fixture.Angles.TwoPi.Value {
		t.Fatalf("TwoPi mismatch")
	}
	testutil.RequireClose(t, float64(core.TwoPi), fixture.Angles.TwoPi.ExpectedValue, 1e-15)
	if float64(core.TwoPiOver5) != fixture.Angles.TwoPiOver5.Value {
		t.Fatalf("TwoPiOver5 mismatch")
	}
	testutil.RequireClose(t, float64(core.TwoPiOver5), fixture.Angles.TwoPiOver5.ExpectedValue, 1e-15)
	if float64(core.PiOver5) != fixture.Angles.PiOver5.Value {
		t.Fatalf("PiOver5 mismatch")
	}
	testutil.RequireClose(t, float64(core.PiOver5), fixture.Angles.PiOver5.ExpectedValue, 1e-15)
	if float64(core.PiOver10) != fixture.Angles.PiOver10.Value {
		t.Fatalf("PiOver10 mismatch")
	}
	testutil.RequireClose(t, float64(core.PiOver10), fixture.Angles.PiOver10.ExpectedValue, 1e-15)
	testutil.RequireClose(t, float64(core.TwoPiOver5), 2*float64(core.PiOver5), 1e-15)
	testutil.RequireClose(t, float64(core.PiOver5), 2*float64(core.PiOver10), 1e-15)

	if float64(core.DihedralAngle) != fixture.DodecahedronAngles.DihedralAngle.Value {
		t.Fatalf("DihedralAngle mismatch")
	}
	testutil.RequireClose(t, float64(core.DihedralAngle), fixture.DodecahedronAngles.DihedralAngle.ExpectedValue, 1e-15)
	if float64(core.InterhedralAngle) != fixture.DodecahedronAngles.InterhedralAngle.Value {
		t.Fatalf("InterhedralAngle mismatch")
	}
	testutil.RequireClose(t, float64(core.InterhedralAngle), fixture.DodecahedronAngles.InterhedralAngle.ExpectedValue, 1e-15)
	testutil.RequireClose(t, fixture.DodecahedronAngles.AngleSum, math.Pi, 1e-15)
	if float64(core.FaceEdgeAngle) != fixture.DodecahedronAngles.FaceEdgeAngle.Value {
		t.Fatalf("FaceEdgeAngle mismatch")
	}
	testutil.RequireClose(t, float64(core.FaceEdgeAngle), fixture.DodecahedronAngles.FaceEdgeAngle.ExpectedValue, 1e-15)

	if core.DistanceToEdge != fixture.Distances.DistanceToEdge.Value {
		t.Fatalf("DistanceToEdge mismatch")
	}
	testutil.RequireClose(t, core.DistanceToEdge, fixture.Distances.DistanceToEdge.ExpectedValue, 1e-15)
	testutil.RequireClose(t, core.DistanceToEdge, fixture.Distances.DistanceToEdge.AlternativeFormula, 1e-15)
	if core.DistanceToVertex != fixture.Distances.DistanceToVertex.Value {
		t.Fatalf("DistanceToVertex mismatch")
	}
	testutil.RequireClose(t, core.DistanceToVertex, fixture.Distances.DistanceToVertex.ExpectedValue, 1e-15)
	testutil.RequireClose(t, core.DistanceToVertex, fixture.Distances.DistanceToVertex.AlternativeFormula, 1e-15)

	if core.RInscribed != fixture.SphereRadii.RInscribed.Value || core.RInscribed != fixture.SphereRadii.RInscribed.ExpectedValue {
		t.Fatalf("RInscribed mismatch")
	}
	if core.RMidEdge != fixture.SphereRadii.RMidedge.Value {
		t.Fatalf("RMidEdge mismatch")
	}
	testutil.RequireClose(t, core.RMidEdge, fixture.SphereRadii.RMidedge.ExpectedValue, 1e-15)
	if core.RCircumscribed != fixture.SphereRadii.RCircumscribed.Value {
		t.Fatalf("RCircumscribed mismatch")
	}
	testutil.RequireClose(t, core.RCircumscribed, fixture.SphereRadii.RCircumscribed.ExpectedValue, 1e-15)
	if !fixture.SphereRadii.Relationships.InscribedLessThanMidedge || !fixture.SphereRadii.Relationships.MidedgeLessThanCircumscribed {
		t.Fatalf("fixture relationship sanity check failed")
	}

	for _, test := range fixture.ValidationTests.FiniteNumbers {
		if !test.IsFinite || test.IsNaN {
			t.Fatalf("expected finite non-NaN constant")
		}
	}
	for _, test := range fixture.ValidationTests.PositiveConstants {
		if !test.IsPositive {
			t.Fatalf("expected positive constant")
		}
	}
}
