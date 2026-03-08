package core

import "math"

func GetNumCells(resolution int) float64 {
	if resolution < 0 {
		return 0
	}
	if resolution == 0 {
		return 12
	}
	return 60 * math.Pow(4, float64(resolution-1))
}

func GetNumChildren(parentResolution, childResolution int) float64 {
	if childResolution < parentResolution {
		return 0
	}
	if childResolution == parentResolution {
		return 1
	}
	if parentResolution >= FirstHilbertResolution {
		return math.Pow(4, float64(childResolution-parentResolution))
	}
	parentCount := GetNumCells(parentResolution)
	if parentCount == 0 {
		parentCount = 1
	}
	return GetNumCells(childResolution) / parentCount
}

func CellArea(resolution int) float64 {
	if resolution < 0 {
		return AuthalicAreaEarth
	}
	return AuthalicAreaEarth / GetNumCells(resolution)
}
