package internal

import (
	"a5go/internal/cells"
	"a5go/internal/core"
)

func GenerateWireframe(resolution int, segments ...int) [][]core.LonLat {
	segmentCount := 1
	if len(segments) > 0 {
		segmentCount = segments[0]
	}

	wireframe := make([][]core.LonLat, 0)
	baseCells := 1
	var stamp uint64
	if resolution == 0 {
		baseCells = 12
		stamp = 0b10 << 56
	} else {
		baseCells = 60
		stamp = 0b01 << 56
	}

	for i := 0; i < baseCells; i++ {
		segment := uint64(i) << 58
		index := segment | stamp
		if resolution < core.FirstHilbertResolution {
			wireframe = append(wireframe, cells.CellToBoundary(index, cells.CellBoundaryOptions{ClosedRing: false, Segments: segmentCount}))
		} else {
			children, err := core.CellToChildren(index, resolution)
			if err != nil {
				panic(err)
			}
			for _, child := range children {
				wireframe = append(wireframe, cells.CellToBoundary(child, cells.CellBoundaryOptions{ClosedRing: false, Segments: segmentCount}))
			}
		}
	}
	return wireframe
}
