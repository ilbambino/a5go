package core

import "a5go/internal/lattice"

type Quaternion [4]float64

type Origin struct {
	ID            int
	Axis          Spherical
	Quat          Quaternion
	InverseQuat   Quaternion
	Angle         Radians
	Orientation   []lattice.Orientation
	FirstQuintant int
}

type A5Cell struct {
	Origin     *Origin
	Segment    int
	S          uint64
	Resolution int
}
