package traversal

import "a5go/internal/core"

func gridDiskBFS(cellID uint64, k int, edgeOnly bool) []uint64 {
	if k == 0 {
		return []uint64{cellID}
	}

	interior := make([]uint64, 0)
	prevFrontier := map[uint64]struct{}{}
	frontier := map[uint64]struct{}{cellID: {}}

	for ring := 1; ring <= k; ring++ {
		nextFrontier := make(map[uint64]struct{})
		for id := range frontier {
			for _, neighbor := range GetGlobalCellNeighbors(id, edgeOnly) {
				if _, ok := prevFrontier[neighbor]; ok {
					continue
				}
				if _, ok := frontier[neighbor]; ok {
					continue
				}
				if _, ok := nextFrontier[neighbor]; ok {
					continue
				}
				nextFrontier[neighbor] = struct{}{}
			}
		}

		for id := range prevFrontier {
			interior = append(interior, id)
		}
		if len(interior) > 100 {
			interior = core.Compact(interior)
		}

		prevFrontier = frontier
		frontier = nextFrontier
	}

	for id := range prevFrontier {
		interior = append(interior, id)
	}
	for id := range frontier {
		interior = append(interior, id)
	}

	return core.Compact(interior)
}

func GridDisk(cellID uint64, k int) []uint64 {
	return gridDiskBFS(cellID, k, true)
}

func GridDiskVertex(cellID uint64, k int) []uint64 {
	return gridDiskBFS(cellID, k, false)
}
