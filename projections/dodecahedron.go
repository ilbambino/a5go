package projections

import (
	"a5go/core"
	"math"
)

type faceTriangleIndex int

type DodecahedronProjection struct {
	faceTriangles      []core.FaceTriangle
	sphericalTriangles []core.SphericalTriangle
	polyhedral         PolyhedralProjection
	gnomonic           GnomonicProjection
	crs                *CRS
}

func NewDodecahedronProjection() *DodecahedronProjection {
	return &DodecahedronProjection{
		faceTriangles:      make([]core.FaceTriangle, 30),
		sphericalTriangles: make([]core.SphericalTriangle, 240),
		polyhedral:         PolyhedralProjection{},
		gnomonic:           GnomonicProjection{},
		crs:                NewCRS(),
	}
}

func (d *DodecahedronProjection) Forward(spherical core.Spherical, originID int) core.Face {
	origin := core.Origins[originID]
	unprojected := core.ToCartesian(spherical)
	out := core.TransformQuat(unprojected, origin.InverseQuat)
	projectedSpherical := core.ToSpherical(out)
	polar := d.gnomonic.Forward(projectedSpherical)
	polar[1] -= float64(origin.Angle)
	faceTriangleIndex := d.getFaceTriangleIndex(polar)
	reflect := d.shouldReflect(polar)
	faceTriangle := d.getFaceTriangle(faceTriangleIndex, reflect, false)
	sphericalTriangle := d.getSphericalTriangle(faceTriangleIndex, originID, reflect)
	return d.polyhedral.Forward(unprojected, sphericalTriangle, faceTriangle)
}

func (d *DodecahedronProjection) Inverse(face core.Face, originID int) core.Spherical {
	polar := core.ToPolar(face)
	faceTriangleIndex := d.getFaceTriangleIndex(polar)
	reflect := d.shouldReflect(polar)
	faceTriangle := d.getFaceTriangle(faceTriangleIndex, reflect, false)
	sphericalTriangle := d.getSphericalTriangle(faceTriangleIndex, originID, reflect)
	unprojected := d.polyhedral.Inverse(face, faceTriangle, sphericalTriangle)
	return core.ToSpherical(unprojected)
}

func (d *DodecahedronProjection) shouldReflect(polar core.Polar) bool {
	rho, gamma := polar[0], polar[1]
	x := core.ToFace(core.Polar{rho, d.NormalizeGamma(gamma)})[0]
	return x > core.DistanceToEdge
}

func (d *DodecahedronProjection) getFaceTriangleIndex(polar core.Polar) faceTriangleIndex {
	return faceTriangleIndex((int(math.Floor(polar[1]/float64(core.PiOver5))) + 10) % 10)
}

func (d *DodecahedronProjection) getFaceTriangle(index faceTriangleIndex, reflected, squashed bool) core.FaceTriangle {
	cacheIndex := int(index)
	if reflected {
		if squashed {
			cacheIndex += 20
		} else {
			cacheIndex += 10
		}
	}
	if triangle := d.faceTriangles[cacheIndex]; triangle != (core.FaceTriangle{}) {
		return triangle
	}
	if reflected {
		d.faceTriangles[cacheIndex] = d.getReflectedFaceTriangle(index, squashed)
	} else {
		d.faceTriangles[cacheIndex] = d.buildFaceTriangle(index)
	}
	return d.faceTriangles[cacheIndex]
}

func (d *DodecahedronProjection) buildFaceTriangle(index faceTriangleIndex) core.FaceTriangle {
	quintant := (int(index) + 1) / 2 % 5
	vertices := core.GetQuintantVertices(quintant).GetVertices()
	vCenter := core.Face(vertices[0])
	vCorner1 := core.Face(vertices[1])
	vCorner2 := core.Face(vertices[2])
	vEdgeMidpoint := core.Face{
		(vCorner1[0] + vCorner2[0]) * 0.5,
		(vCorner1[1] + vCorner2[1]) * 0.5,
	}
	even := int(index)%2 == 0
	if even {
		return core.FaceTriangle{vCenter, vEdgeMidpoint, vCorner1}
	}
	return core.FaceTriangle{vCenter, vCorner2, vEdgeMidpoint}
}

func (d *DodecahedronProjection) getReflectedFaceTriangle(index faceTriangleIndex, squashed bool) core.FaceTriangle {
	triangle := d.buildFaceTriangle(index)
	a := core.Face{triangle[0][0], triangle[0][1]}
	b := core.Face{triangle[1][0], triangle[1][1]}
	c := core.Face{triangle[2][0], triangle[2][1]}
	even := int(index)%2 == 0
	a[0] = -a[0]
	a[1] = -a[1]
	midpoint := b
	if !even {
		midpoint = c
	}
	scale := 2.0
	if squashed {
		scale = 1 + 1/math.Cos(float64(core.InterhedralAngle))
	}
	a[0] += midpoint[0] * scale
	a[1] += midpoint[1] * scale
	return core.FaceTriangle{a, c, b}
}

func (d *DodecahedronProjection) getSphericalTriangle(index faceTriangleIndex, originID int, reflected bool) core.SphericalTriangle {
	cacheIndex := 10*originID + int(index)
	if reflected {
		cacheIndex += 120
	}
	if triangle := d.sphericalTriangles[cacheIndex]; triangle != (core.SphericalTriangle{}) {
		return triangle
	}
	d.sphericalTriangles[cacheIndex] = d.buildSphericalTriangle(index, originID, reflected)
	return d.sphericalTriangles[cacheIndex]
}

func (d *DodecahedronProjection) buildSphericalTriangle(index faceTriangleIndex, originID int, reflected bool) core.SphericalTriangle {
	origin := core.Origins[originID]
	faceTriangle := d.getFaceTriangle(index, reflected, true)
	var sphericalTriangle core.SphericalTriangle
	for i, face := range faceTriangle {
		rhoGamma := core.ToPolar(face)
		rotatedPolar := core.Polar{rhoGamma[0], rhoGamma[1] + float64(origin.Angle)}
		rotated := core.ToCartesian(d.gnomonic.Inverse(rotatedPolar))
		rotated = core.TransformQuat(rotated, origin.Quat)
		sphericalTriangle[i] = d.crs.GetVertex(rotated)
	}
	return sphericalTriangle
}

func (d *DodecahedronProjection) NormalizeGamma(gamma float64) float64 {
	segment := gamma / float64(core.TwoPiOver5)
	sCenter := math.Round(segment)
	sOffset := segment - sCenter
	return sOffset * float64(core.TwoPiOver5)
}
