package utils

import (
	"a5go/core"
	"math"
)

func dot(a, b core.Cartesian) float64 {
	return a[0]*b[0] + a[1]*b[1] + a[2]*b[2]
}

func cross(a, b core.Cartesian) core.Cartesian {
	return core.Cartesian{
		a[1]*b[2] - a[2]*b[1],
		a[2]*b[0] - a[0]*b[2],
		a[0]*b[1] - a[1]*b[0],
	}
}

func add(a, b core.Cartesian) core.Cartesian {
	return core.Cartesian{a[0] + b[0], a[1] + b[1], a[2] + b[2]}
}

func sub(a, b core.Cartesian) core.Cartesian {
	return core.Cartesian{a[0] - b[0], a[1] - b[1], a[2] - b[2]}
}

func scale(v core.Cartesian, s float64) core.Cartesian {
	return core.Cartesian{v[0] * s, v[1] * s, v[2] * s}
}

func length(v core.Cartesian) float64 {
	return math.Sqrt(dot(v, v))
}

func normalize(v core.Cartesian) core.Cartesian {
	l := length(v)
	if l == 0 {
		return v
	}
	return scale(v, 1/l)
}

func lerp(a, b core.Cartesian, t float64) core.Cartesian {
	return core.Cartesian{
		a[0] + (b[0]-a[0])*t,
		a[1] + (b[1]-a[1])*t,
		a[2] + (b[2]-a[2])*t,
	}
}

func angle(a, b core.Cartesian) float64 {
	denominator := length(a) * length(b)
	if denominator == 0 {
		return 0
	}
	cosine := dot(a, b) / denominator
	if cosine > 1 {
		cosine = 1
	}
	if cosine < -1 {
		cosine = -1
	}
	return math.Acos(cosine)
}

func VectorDifference(a, b core.Cartesian) float64 {
	midpointAB := normalize(lerp(a, b, 0.5))
	d := length(cross(midpointAB, a))
	if d < 1e-8 {
		ab := sub(a, b)
		return 0.5 * length(ab)
	}
	return d
}

func TripleProduct(a, b, c core.Cartesian) float64 {
	return dot(a, cross(b, c))
}

func QuadrupleProduct(out *core.Cartesian, a, b, c, d core.Cartesian) *core.Cartesian {
	crossCD := cross(c, d)
	tripleProductACD := dot(a, crossCD)
	tripleProductBCD := dot(b, crossCD)
	scaledA := scale(a, tripleProductBCD)
	scaledB := scale(b, tripleProductACD)
	*out = sub(scaledB, scaledA)
	return out
}

func Slerp(out *core.Cartesian, a, b core.Cartesian, t float64) *core.Cartesian {
	gamma := angle(a, b)
	if gamma < 1e-12 {
		*out = lerp(a, b, t)
		return out
	}
	weightA := math.Sin((1-t)*gamma) / math.Sin(gamma)
	weightB := math.Sin(t*gamma) / math.Sin(gamma)
	*out = add(scale(a, weightA), scale(b, weightB))
	return out
}
