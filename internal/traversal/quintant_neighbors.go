package traversal

import "a5go/internal/lattice"

func FindQuintantNeighborS(sourceTriple lattice.Triple, uvSourceAnchor *lattice.Anchor, sourceS uint64, resolution int, orientation lattice.Orientation, edgeOnly bool) []uint64 {
	maxS := uint64(1) << uint(2*resolution)
	maxRow := (1 << resolution) - 1
	neighbors := make([]uint64, 0)

	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			for dz := -1; dz <= 1; dz++ {
				if dx == 0 && dy == 0 && dz == 0 {
					continue
				}
				if abs(dx)+abs(dy)+abs(dz) > 3 {
					continue
				}
				if edgeOnly && abs(dx)+abs(dy)+abs(dz) > 2 {
					continue
				}
				neighborTriple := lattice.Triple{X: sourceTriple.X + dx, Y: sourceTriple.Y + dy, Z: sourceTriple.Z + dz}
				if !lattice.TripleInBounds(neighborTriple, maxRow) {
					continue
				}
				uvNeighborAnchor := lattice.TripleToAnchor(neighborTriple, resolution, lattice.OrientationUV)
				if uvNeighborAnchor == nil || uvSourceAnchor == nil {
					continue
				}
				if !IsNeighbor(*uvSourceAnchor, *uvNeighborAnchor) {
					continue
				}
				neighborS := lattice.TripleToS(neighborTriple, resolution, orientation)
				if neighborS != nil && *neighborS < maxS && *neighborS != sourceS {
					neighbors = append(neighbors, *neighborS)
				}
			}
		}
	}

	return neighbors
}

func GetCellNeighbors(s uint64, resolution int, orientation lattice.Orientation, edgeOnly bool) []uint64 {
	anchor := lattice.SToAnchor(s, resolution, orientation)
	triple := lattice.AnchorToTriple(anchor)
	uvSourceAnchor := lattice.TripleToAnchor(triple, resolution, lattice.OrientationUV)
	return sortedUint64(FindQuintantNeighborS(triple, uvSourceAnchor, s, resolution, orientation, edgeOnly))
}

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}
