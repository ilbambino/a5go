package internal

import (
	"a5go/internal/core"
	"a5go/internal/projections"
)

var (
	helperDodecahedron = projections.NewDodecahedronProjection()
)

func LonLatToFace(lonLat core.LonLat, resolution int, centroid ...core.LonLat) core.Face {
	_ = resolution
	spherical := core.FromLonLat(lonLat)
	var origin *core.Origin
	if len(centroid) > 0 {
		origin = core.FindNearestOrigin(core.FromLonLat(centroid[0]))
	} else {
		origin = core.FindNearestOrigin(spherical)
	}

	rotations := []int{0, 0, 9, 6, 7, 6, 5, 4, 7, 7, 9, 0}
	helperRotation := core.Mat2FromRotation(float64(rotations[origin.ID])*float64(core.PiOver5) + float64(origin.Angle))

	dodecPoint := helperDodecahedron.Forward(spherical, origin.ID)
	dodecPoint = core.Mat2Transform(helperRotation, dodecPoint)

	shift := core.Face{0, 0}
	path := []int{0, 0, 1, 2, 5, 2, 3, 2, 9, 2, 3, 4}
	for i := 0; i <= origin.ID; i++ {
		offset := core.Face{
			1.232 * cos(float64(path[i])*float64(core.PiOver5)),
			1.232 * sin(float64(path[i])*float64(core.PiOver5)),
		}
		shift[0] += offset[0]
		shift[1] += offset[1]
	}

	dodecPoint[0] += shift[0]
	dodecPoint[1] += shift[1]
	return dodecPoint
}
