package traversal

import (
	"a5go/internal/cells"
	"a5go/internal/core"
	"math"
	"sort"
)

const (
	cellRadiusSafetyFactor = 2.0
	minCellsForSubdivision = 20
)

var cellRadius = func() []float64 {
	radii := make([]float64, 31)
	radii[0] = cellRadiusSafetyFactor * core.AuthalicRadiusEarth / math.Sqrt(3)
	baseCellRadius := cellRadiusSafetyFactor * core.AuthalicRadiusEarth / math.Sqrt(15)
	for r := 1; r <= 30; r++ {
		radii[r] = baseCellRadius / float64(uint64(1)<<uint(r-1))
	}
	return radii
}()

func MetersToH(meters float64) float64 {
	s := math.Sin(meters / (2 * core.AuthalicRadiusEarth))
	return s * s
}

func EstimateCellRadius(resolution int) float64 {
	return cellRadius[resolution]
}

func PickCoarseResolution(radius float64, targetRes int) int {
	capAreaM2 := 2 * math.Pi * core.AuthalicRadiusEarth * core.AuthalicRadiusEarth * (1 - math.Cos(radius/core.AuthalicRadiusEarth))
	for res := core.FirstHilbertResolution; res <= targetRes; res++ {
		cArea := core.CellArea(res)
		if capAreaM2/cArea >= minCellsForSubdivision {
			return res
		}
	}
	return targetRes
}

func SphericalCap(cellID uint64, radius float64) []uint64 {
	targetRes := core.GetResolution(cellID)
	coarseRes := PickCoarseResolution(radius, targetRes)
	center := cells.CellToSpherical(cellID)
	hRadius := MetersToH(radius)

	startCell := cellID
	if coarseRes < targetRes {
		startCell = core.CellToParent(cellID, coarseRes)
	}
	coarseCellRadius := EstimateCellRadius(coarseRes)
	hExpanded := MetersToH(radius + coarseCellRadius)
	coarseVisited := map[uint64]struct{}{startCell: {}}
	coarseFrontier := map[uint64]struct{}{startCell: {}}

	for len(coarseFrontier) > 0 {
		nextFrontier := map[uint64]struct{}{}
		for id := range coarseFrontier {
			for _, neighbor := range GetGlobalCellNeighbors(id, false) {
				if _, seen := coarseVisited[neighbor]; seen {
					continue
				}
				coarseVisited[neighbor] = struct{}{}
				if core.Haversine(center, cells.CellToSpherical(neighbor)) <= hExpanded {
					nextFrontier[neighbor] = struct{}{}
				}
			}
		}
		coarseFrontier = nextFrontier
	}

	result := make([]uint64, 0)
	boundary := make([]uint64, 0, len(coarseVisited))
	for cell := range coarseVisited {
		boundary = append(boundary, cell)
	}

	for res := coarseRes; res < targetRes; res++ {
		cellRadius := EstimateCellRadius(res)
		hInner := -1.0
		if radius > cellRadius {
			hInner = MetersToH(radius - cellRadius)
		}
		hOuter := MetersToH(radius + cellRadius)
		nextBoundary := make([]uint64, 0)

		for _, cell := range boundary {
			h := core.Haversine(center, cells.CellToSpherical(cell))
			if h <= hInner {
				result = append(result, cell)
			} else if h > hOuter {
				continue
			} else {
				nextBoundary = append(nextBoundary, core.CellToChildren(cell, res+1)...)
			}
		}
		boundary = nextBoundary
	}

	for _, cell := range boundary {
		if core.Haversine(center, cells.CellToSpherical(cell)) <= hRadius {
			result = append(result, cell)
		}
	}

	sort.Slice(result, func(i, j int) bool { return result[i] < result[j] })
	return result
}
