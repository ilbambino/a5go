package core

import "math"

var Quaternions = func() []Quaternion {
	sqrt5 := math.Sqrt(5)
	invSqrt5 := math.Sqrt(0.2)

	sinAlpha := math.Sqrt((1 - invSqrt5) / 2)
	cosAlpha := math.Sqrt((1 + invSqrt5) / 2)

	a := 0.5
	b := math.Sqrt((2.5 - sqrt5) / 10)
	c := math.Sqrt((2.5 + sqrt5) / 10)
	d := math.Sqrt((1 + invSqrt5) / 8)
	e := math.Sqrt((1 - invSqrt5) / 8)
	f := math.Sqrt((3 - sqrt5) / 8)
	g := math.Sqrt((3 + sqrt5) / 8)

	faceCenters := [][2]float64{
		{0, 0},
		{sinAlpha, 0},
		{b, a},
		{-d, f},
		{-d, -f},
		{b, -a},
		{-cosAlpha, 0},
		{-e, -g},
		{c, -a},
		{c, a},
		{-e, g},
		{0, 0},
	}

	axes := make([][2]float64, len(faceCenters))
	for i, center := range faceCenters {
		axes[i] = [2]float64{-center[1], center[0]}
	}

	quaternions := make([]Quaternion, len(axes))
	for i, axis := range axes {
		switch i {
		case 0:
			quaternions[i] = Quaternion{0, 0, 0, 1}
		case 11:
			quaternions[i] = Quaternion{0, -1, 0, 0}
		default:
			if i < 6 {
				quaternions[i] = Quaternion{axis[0], axis[1], 0, cosAlpha}
			} else {
				quaternions[i] = Quaternion{axis[0], axis[1], 0, sinAlpha}
			}
		}
	}
	return quaternions
}()
