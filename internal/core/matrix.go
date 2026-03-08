package core

import "math"

type Mat2 [4]float64
type Mat2D [6]float64

func Mat2FromValues(m00, m01, m10, m11 float64) Mat2 {
	return Mat2{m00, m01, m10, m11}
}

func Mat2FromRotation(angle float64) Mat2 {
	c := math.Cos(angle)
	s := math.Sin(angle)
	return Mat2{c, s, -s, c}
}

func Mat2Transform(m Mat2, v Face) Face {
	return Face{
		m[0]*v[0] + m[2]*v[1],
		m[1]*v[0] + m[3]*v[1],
	}
}

func Mat2DTransform(m Mat2D, v Face) Face {
	return Face{
		m[0]*v[0] + m[2]*v[1] + m[4],
		m[1]*v[0] + m[3]*v[1] + m[5],
	}
}

func Mat2Multiply(a, b Mat2) Mat2 {
	return Mat2{
		a[0]*b[0] + a[2]*b[1],
		a[1]*b[0] + a[3]*b[1],
		a[0]*b[2] + a[2]*b[3],
		a[1]*b[2] + a[3]*b[3],
	}
}

func Mat2Invert(m Mat2) Mat2 {
	det := m[0]*m[3] - m[1]*m[2]
	return Mat2{
		m[3] / det,
		-m[1] / det,
		-m[2] / det,
		m[0] / det,
	}
}
