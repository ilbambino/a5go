package core

import (
	"a5go/geometry"
	"math"
)

const (
	A = Degrees(72)
	B = Degrees(127.94543761193603)
	C = Degrees(108)
	D = Degrees(82.29202980963508)
	E = Degrees(149.7625318412527)
)

var (
	a = Face{0, 0}
	b = Face{0, 1}
	c = Face{0.7885966681787006, 1.6149108024237764}
	d = Face{1.6171013659387945, 1.054928690397459}
	e = Face{math.Cos(float64(PiOver10)), math.Sin(float64(PiOver10))}

	edgeMidpointD = 2 * faceLength(c) * math.Cos(float64(PiOver5))
	basisRotation = float64(PiOver5) - math.Atan2(c[1], c[0])
	scaleFactor   = 2 * DistanceToEdge / edgeMidpointD

	PentagonShapeDef *geometry.PentagonShape
	u                = Face{0, 0}
	v                Face
	w                Face
	V                Radians
	TriangleShapeDef *geometry.PentagonShape
	Basis            Mat2
	BasisInverse     Mat2
)

func faceLength(v Face) float64 {
	return math.Sqrt(v[0]*v[0] + v[1]*v[1])
}

func rotateFace(v Face, angle float64) Face {
	cosine := math.Cos(angle)
	sine := math.Sin(angle)
	return Face{
		v[0]*cosine - v[1]*sine,
		v[0]*sine + v[1]*cosine,
	}
}

func init() {
	a = rotateFace(Face{a[0] * scaleFactor, a[1] * scaleFactor}, basisRotation)
	b = rotateFace(Face{b[0] * scaleFactor, b[1] * scaleFactor}, basisRotation)
	c = rotateFace(Face{c[0] * scaleFactor, c[1] * scaleFactor}, basisRotation)
	d = rotateFace(Face{d[0] * scaleFactor, d[1] * scaleFactor}, basisRotation)
	e = rotateFace(Face{e[0] * scaleFactor, e[1] * scaleFactor}, basisRotation)
	PentagonShapeDef = geometry.NewPentagonShape(geometry.Pentagon{
		geometry.Face(a),
		geometry.Face(b),
		geometry.Face(c),
		geometry.Face(d),
		geometry.Face(e),
	})

	bisectorAngle := math.Atan2(c[1], c[0]) - float64(PiOver5)
	l := DistanceToEdge / math.Cos(float64(PiOver5))

	V = Radians(bisectorAngle + float64(PiOver5))
	v = Face{l * math.Cos(float64(V)), l * math.Sin(float64(V))}

	wAngle := bisectorAngle - float64(PiOver5)
	w = Face{l * math.Cos(wAngle), l * math.Sin(wAngle)}
	TriangleShapeDef = geometry.NewPentagonShape(geometry.Pentagon{
		geometry.Face(u),
		geometry.Face(v),
		geometry.Face(w),
	})

	Basis = Mat2FromValues(v[0], v[1], w[0], w[1])
	BasisInverse = Mat2Invert(Basis)
}
