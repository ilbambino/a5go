package projections

import (
	"a5go/internal/core"
	"a5go/internal/geometry"
	"a5go/utils"
	"math"
)

type PolyhedralProjection struct{}

func (PolyhedralProjection) Forward(v core.Cartesian, sphericalTriangle core.SphericalTriangle, faceTriangle core.FaceTriangle) core.Face {
	a, b, c := sphericalTriangle[0], sphericalTriangle[1], sphericalTriangle[2]
	triangleShape := geometry.NewSphericalTriangleShape([]geometry.Cartesian{geometry.Cartesian(a), geometry.Cartesian(b), geometry.Cartesian(c)})
	z := normalize3(core.Cartesian{v[0] - a[0], v[1] - a[1], v[2] - a[2]})
	var p core.Cartesian
	utils.QuadrupleProduct(&p, a, z, b, c)
	p = normalize3(p)

	h := utils.VectorDifference(a, v) / utils.VectorDifference(a, p)
	areaABC := triangleShape.GetArea()
	scaledArea := h / areaABC
	shape1 := geometry.NewSphericalTriangleShape([]geometry.Cartesian{geometry.Cartesian(a), geometry.Cartesian(p), geometry.Cartesian(c)})
	shape2 := geometry.NewSphericalTriangleShape([]geometry.Cartesian{geometry.Cartesian(a), geometry.Cartesian(b), geometry.Cartesian(p)})
	bary := core.Barycentric{
		1 - h,
		scaledArea * shape1.GetArea(),
		scaledArea * shape2.GetArea(),
	}
	return core.BarycentricToFace(bary, faceTriangle)
}

func (p PolyhedralProjection) Inverse(facePoint core.Face, faceTriangle core.FaceTriangle, sphericalTriangle core.SphericalTriangle) core.Cartesian {
	a, b, c := sphericalTriangle[0], sphericalTriangle[1], sphericalTriangle[2]
	triangleShape := geometry.NewSphericalTriangleShape([]geometry.Cartesian{geometry.Cartesian(a), geometry.Cartesian(b), geometry.Cartesian(c)})
	bary := core.FaceToBarycentric(facePoint, faceTriangle)

	threshold := 1 - 1e-14
	if bary[0] > threshold {
		return a
	}
	if bary[1] > threshold {
		return b
	}
	if bary[2] > threshold {
		return c
	}

	c1 := cross3(b, c)
	areaABC := triangleShape.GetArea()
	h := 1 - bary[0]
	r := bary[2] / h
	alpha := r * areaABC
	s := math.Sin(alpha)
	halfC := math.Sin(alpha / 2)
	cc := 2 * halfC * halfC

	c01 := dot3(a, b)
	c12 := dot3(b, c)
	c20 := dot3(c, a)
	s12 := math.Sqrt(dot3(c1, c1))

	v := dot3(a, c1)
	f := s*v + cc*(c01*c12-c20)
	g := cc * s12 * (1 + c01)
	q := (2 / math.Acos(c12)) * math.Atan2(g, f)
	var pPoint core.Cartesian
	utils.Slerp(&pPoint, b, c, q)
	k := utils.VectorDifference(a, pPoint)
	t := p.safeAcos(h*k) / p.safeAcos(k)
	var out core.Cartesian
	utils.Slerp(&out, a, pPoint, t)
	return out
}

func (PolyhedralProjection) safeAcos(x float64) float64 {
	if x < 1e-3 {
		return 2*x + x*x*x/3
	}
	return math.Acos(1 - 2*x*x)
}

func normalize3(v core.Cartesian) core.Cartesian {
	length := math.Sqrt(v[0]*v[0] + v[1]*v[1] + v[2]*v[2])
	return core.Cartesian{v[0] / length, v[1] / length, v[2] / length}
}

func dot3(a, b core.Cartesian) float64 {
	return a[0]*b[0] + a[1]*b[1] + a[2]*b[2]
}

func cross3(a, b core.Cartesian) core.Cartesian {
	return core.Cartesian{
		a[1]*b[2] - a[2]*b[1],
		a[2]*b[0] - a[0]*b[2],
		a[0]*b[1] - a[1]*b[0],
	}
}
