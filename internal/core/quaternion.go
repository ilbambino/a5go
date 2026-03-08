package core

import "math"

func QuaternionConjugate(q Quaternion) Quaternion {
	return Quaternion{-q[0], -q[1], -q[2], q[3]}
}

func QuaternionLength(q Quaternion) float64 {
	return math.Sqrt(q[0]*q[0] + q[1]*q[1] + q[2]*q[2] + q[3]*q[3])
}

func QuaternionMultiply(a, b Quaternion) Quaternion {
	return Quaternion{
		a[3]*b[0] + a[0]*b[3] + a[1]*b[2] - a[2]*b[1],
		a[3]*b[1] - a[0]*b[2] + a[1]*b[3] + a[2]*b[0],
		a[3]*b[2] + a[0]*b[1] - a[1]*b[0] + a[2]*b[3],
		a[3]*b[3] - a[0]*b[0] - a[1]*b[1] - a[2]*b[2],
	}
}

func TransformQuat(v Cartesian, q Quaternion) Cartesian {
	p := Quaternion{v[0], v[1], v[2], 0}
	rotated := QuaternionMultiply(QuaternionMultiply(q, p), QuaternionConjugate(q))
	return Cartesian{rotated[0], rotated[1], rotated[2]}
}
