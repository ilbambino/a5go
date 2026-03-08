package core

import (
	"a5go/internal/geometry"
	"a5go/internal/lattice"
	"math"
)

var quintantRotations = func() []Mat2 {
	rotations := make([]Mat2, 5)
	for quintant := 0; quintant < 5; quintant++ {
		rotations[quintant] = Mat2FromRotation(float64(TwoPiOver5) * float64(quintant))
	}
	return rotations
}()

func GetPentagonVertices(resolution int, quintant int, anchor lattice.Anchor) *geometry.PentagonShape {
	pentagon := PentagonShapeDef.Clone()

	translation := Mat2Transform(Basis, Face{anchor.Offset[0], anchor.Offset[1]})

	if anchor.Flips[0] == lattice.NO && anchor.Flips[1] == lattice.YES {
		pentagon.Rotate180()
	}

	q := anchor.Q
	f := anchor.Flips[0] + anchor.Flips[1]
	if ((f == -2 || f == 2) && q > lattice.Quaternary1) || (f == 0 && (q == lattice.Quaternary0 || q == lattice.Quaternary3)) {
		pentagon.ReflectY()
	}
	shiftLeft := geometry.Face{-w[0], -w[1]}
	shiftRight := geometry.Face{w[0], w[1]}
	if anchor.Flips[0] == lattice.YES && anchor.Flips[1] == lattice.YES {
		pentagon.Rotate180()
	} else if anchor.Flips[0] == lattice.YES {
		pentagon.Translate(shiftLeft)
	} else if anchor.Flips[1] == lattice.YES {
		pentagon.Translate(shiftRight)
	}

	pentagon.Translate(geometry.Face(translation))
	pentagon.Scale(1 / math.Pow(2, float64(resolution)))
	pentagon.Transform(geometry.Mat2(quintantRotations[quintant]))

	return pentagon
}

type PentagonFlavor uint8

func GetPentagonFlavor(anchor lattice.Anchor) PentagonFlavor {
	f := 0
	if anchor.Flips[1] == lattice.YES {
		f += 2
	}

	q := anchor.Q
	flipSum := anchor.Flips[0] + anchor.Flips[1]
	if ((flipSum == -2 || flipSum == 2) && q > lattice.Quaternary1) || (flipSum == 0 && (q == lattice.Quaternary0 || q == lattice.Quaternary3)) {
		f++
	}
	if flipSum == -2 || flipSum == 2 {
		f += 4
	}

	return PentagonFlavor(f)
}

func GetQuintantVertices(quintant int) *geometry.PentagonShape {
	triangle := TriangleShapeDef.Clone()
	triangle.Transform(geometry.Mat2(quintantRotations[quintant]))
	return triangle
}

func GetFaceVertices() *geometry.PentagonShape {
	vertices := make([]Face, 0, len(quintantRotations))
	for _, rotation := range quintantRotations {
		vertices = append(vertices, Mat2Transform(rotation, v))
	}
	for i, j := 0, len(vertices)-1; i < j; i, j = i+1, j-1 {
		vertices[i], vertices[j] = vertices[j], vertices[i]
	}
	polygon := make(geometry.Pentagon, len(vertices))
	for i, vertex := range vertices {
		polygon[i] = geometry.Face(vertex)
	}
	return geometry.NewPentagonShape(polygon)
}

func GetQuintantPolar(polar Polar) int {
	return (jsRound(polar[1]/float64(TwoPiOver5)) + 5) % 5
}

func jsRound(v float64) int {
	return int(math.Floor(v + 0.5))
}
