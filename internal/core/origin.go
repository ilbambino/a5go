package core

import "a5go/internal/lattice"

var (
	clockwiseFan = []lattice.Orientation{
		lattice.OrientationVU, lattice.OrientationUW, lattice.OrientationVW, lattice.OrientationVW, lattice.OrientationVW,
	}
	clockwiseStep = []lattice.Orientation{
		lattice.OrientationWU, lattice.OrientationUW, lattice.OrientationVW, lattice.OrientationVU, lattice.OrientationUW,
	}
	counterStep = []lattice.Orientation{
		lattice.OrientationWU, lattice.OrientationUV, lattice.OrientationWV, lattice.OrientationWU, lattice.OrientationUW,
	}
	counterJump = []lattice.Orientation{
		lattice.OrientationVU, lattice.OrientationUV, lattice.OrientationWV, lattice.OrientationWU, lattice.OrientationUW,
	}
	quintantOrientations = [][]lattice.Orientation{
		clockwiseFan, counterJump, counterStep,
		clockwiseStep, counterStep, counterJump,
		counterStep, clockwiseStep, clockwiseStep,
		clockwiseStep, counterJump, counterJump,
	}
	quintantFirst = []int{4, 2, 3, 2, 0, 4, 3, 2, 2, 0, 3, 0}
	originOrder   = []int{0, 1, 2, 4, 3, 5, 7, 8, 6, 11, 10, 9}
)

var Origins = func() []*Origin {
	origins := make([]*Origin, 0, 12)
	nextID := 0
	addOrigin := func(axis Spherical, angle Radians, quaternion Quaternion) {
		origin := &Origin{
			ID:            nextID,
			Axis:          axis,
			Quat:          quaternion,
			InverseQuat:   QuaternionConjugate(quaternion),
			Angle:         angle,
			Orientation:   quintantOrientations[nextID],
			FirstQuintant: quintantFirst[nextID],
		}
		origins = append(origins, origin)
		nextID++
	}

	addOrigin(Spherical{0, 0}, 0, Quaternions[0])
	for i := 0; i < 5; i++ {
		alpha := float64(i) * float64(TwoPiOver5)
		alpha2 := alpha + float64(PiOver5)
		addOrigin(Spherical{alpha, float64(InterhedralAngle)}, PiOver5, Quaternions[i+1])
		addOrigin(Spherical{alpha2, mathPi - float64(InterhedralAngle)}, PiOver5, Quaternions[(i+3)%5+6])
	}
	addOrigin(Spherical{0, mathPi}, 0, Quaternions[11])

	for i := 0; i < len(origins); i++ {
		for j := i + 1; j < len(origins); j++ {
			if originOrderIndex(origins[j].ID) < originOrderIndex(origins[i].ID) {
				origins[i], origins[j] = origins[j], origins[i]
			}
		}
	}
	for i, origin := range origins {
		origin.ID = i
	}
	return origins
}()

const mathPi = 3.14159265358979323846264338327950288419716939937510

func QuintantToSegment(quintant int, origin *Origin) (int, lattice.Orientation) {
	layout := origin.Orientation
	step := 1
	if sameOrientationLayout(layout, clockwiseFan) || sameOrientationLayout(layout, clockwiseStep) {
		step = -1
	}
	delta := (quintant - origin.FirstQuintant + 5) % 5
	faceRelativeQuintant := (step*delta + 5) % 5
	orientation := layout[faceRelativeQuintant]
	segment := (origin.FirstQuintant + faceRelativeQuintant) % 5
	return segment, orientation
}

func SegmentToQuintant(segment int, origin *Origin) (int, lattice.Orientation) {
	layout := origin.Orientation
	step := 1
	if sameOrientationLayout(layout, clockwiseFan) || sameOrientationLayout(layout, clockwiseStep) {
		step = -1
	}
	faceRelativeQuintant := (segment - origin.FirstQuintant + 5) % 5
	orientation := layout[faceRelativeQuintant]
	quintant := (origin.FirstQuintant + step*faceRelativeQuintant + 5) % 5
	return quintant, orientation
}

func sameOrientationLayout(a, b []lattice.Orientation) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func originOrderIndex(id int) int {
	for i, candidate := range originOrder {
		if candidate == id {
			return i
		}
	}
	return len(originOrder)
}

func FindNearestOrigin(point Spherical) *Origin {
	minDistance := 1e9
	nearest := Origins[0]
	for _, origin := range Origins {
		distance := Haversine(point, origin.Axis)
		if distance < minDistance {
			minDistance = distance
			nearest = origin
		}
	}
	return nearest
}

func IsNearestOrigin(point Spherical, origin *Origin) bool {
	return Haversine(point, origin.Axis) > 0.49999999
}

func Haversine(point, axis Spherical) float64 {
	theta, phi := point[0], point[1]
	theta2, phi2 := axis[0], axis[1]
	dtheta := theta2 - theta
	dphi := phi2 - phi
	a1 := sinHalf(dphi)
	a2 := sinHalf(dtheta)
	return a1*a1 + a2*a2*mathSin(phi)*mathSin(phi2)
}

func sinHalf(v float64) float64 { return mathSin(v / 2) }
