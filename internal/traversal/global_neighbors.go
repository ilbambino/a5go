package traversal

import (
	"a5go/internal/core"
	"a5go/internal/lattice"
	"sort"
)

type neighborDelta [4]int

var leftEdgeDeltas = [][]neighborDelta{
	{{0, 0, 0, 1}, {0, 0, 1, 0}},
	{{0, 0, 0, 1}, {0, 1, 0, 1}, {0, -1, 1, 0}, {0, 1, -1, 0}},
	{},
	{{0, -1, 0, 1}, {0, 0, -1, 0}},
}

var rightEdgeDeltas = [][]neighborDelta{
	{{0, 0, 0, 1}, {0, 1, 0, 1}, {-1, 1, 0, 0}, {1, -1, 0, 0}},
	{{0, 0, 0, 1}, {1, 0, 0, 0}},
	{{0, -1, 0, 1}, {-1, 0, 0, 0}},
	{},
}

var crossFaceDeltas = [][]neighborDelta{
	{{0, 0, 0, 1}, {1, 0, 0, 1}, {1, 0, -1, 0}},
	{{0, 0, -1, 1}, {0, 0, 0, 0}},
}

type neighborContext struct {
	hilbertRes  int
	resolution  int
	maxS        uint64
	maxRow      int
	edgeOnly    bool
	neighborSet map[uint64]struct{}
}

func addNeighbor(ctx *neighborContext, neighborTriple lattice.Triple, orientation lattice.Orientation, neighborOrigin *core.Origin, neighborSegment int) {
	s := lattice.TripleToS(neighborTriple, ctx.hilbertRes, orientation)
	if s == nil || *s >= ctx.maxS {
		return
	}
	ctx.neighborSet[core.Serialize(core.A5Cell{Origin: neighborOrigin, Segment: neighborSegment, S: *s, Resolution: ctx.resolution})] = struct{}{}
}

func addDeltaNeighbors(ctx *neighborContext, base lattice.Triple, deltas []neighborDelta, orientation lattice.Orientation, neighborOrigin *core.Origin, neighborSegment int) {
	for _, delta := range deltas {
		if ctx.edgeOnly && delta[3] == 0 {
			continue
		}
		neighborTriple := lattice.Triple{X: base.X + delta[0], Y: base.Y + delta[1], Z: base.Z + delta[2]}
		if !lattice.TripleInBounds(neighborTriple, ctx.maxRow) {
			continue
		}
		addNeighbor(ctx, neighborTriple, orientation, neighborOrigin, neighborSegment)
	}
}

func GetGlobalCellNeighbors(cellID uint64, edgeOnly bool) []uint64 {
	cell := core.Deserialize(cellID)
	if cell.Resolution < core.FirstHilbertResolution {
		return []uint64{}
	}

	hilbertRes := cell.Resolution - core.FirstHilbertResolution + 1
	sourceQuintant, sourceOrientation := core.SegmentToQuintant(cell.Segment, cell.Origin)
	anchor := lattice.SToAnchor(cell.S, hilbertRes, sourceOrientation)
	triple := lattice.AnchorToTriple(anchor)
	uvSourceAnchor := lattice.TripleToAnchor(triple, hilbertRes, lattice.OrientationUV)

	ctx := &neighborContext{
		hilbertRes:  hilbertRes,
		resolution:  cell.Resolution,
		maxS:        uint64(1) << uint(2*hilbertRes),
		maxRow:      (1 << hilbertRes) - 1,
		edgeOnly:    edgeOnly,
		neighborSet: make(map[uint64]struct{}),
	}

	for _, neighborS := range FindQuintantNeighborS(triple, uvSourceAnchor, cell.S, hilbertRes, sourceOrientation, ctx.edgeOnly) {
		ctx.neighborSet[core.Serialize(core.A5Cell{Origin: cell.Origin, Segment: cell.Segment, S: neighborS, Resolution: cell.Resolution})] = struct{}{}
	}

	parity := lattice.TripleParity(triple)
	yOdd := triple.Y%2 != 0
	deltaIndex := parity * 2
	if yOdd {
		deltaIndex++
	}

	if triple.Z == 0 {
		targetQuintant := (sourceQuintant - 1 + 5) % 5
		targetSegment, targetOrientation := core.QuintantToSegment(targetQuintant, cell.Origin)
		swappedBase := lattice.Triple{X: 0, Y: triple.Y, Z: triple.X}
		addDeltaNeighbors(ctx, swappedBase, leftEdgeDeltas[deltaIndex], targetOrientation, cell.Origin, targetSegment)
	}
	if triple.X == 0 {
		targetQuintant := (sourceQuintant + 1) % 5
		targetSegment, targetOrientation := core.QuintantToSegment(targetQuintant, cell.Origin)
		swappedBase := lattice.Triple{X: triple.Z, Y: triple.Y, Z: 0}
		addDeltaNeighbors(ctx, swappedBase, rightEdgeDeltas[deltaIndex], targetOrientation, cell.Origin, targetSegment)
	}

	if triple.Y == ctx.maxRow {
		adj := core.FaceAdjacency[cell.Origin.ID][sourceQuintant]
		adjOrigin := core.Origins[adj[0]]
		adjSegment, adjOrientation := core.QuintantToSegment(adj[1], adjOrigin)
		mirroredBase := lattice.Triple{X: triple.Z, Y: ctx.maxRow, Z: triple.X}
		addDeltaNeighbors(ctx, mirroredBase, crossFaceDeltas[parity], adjOrientation, adjOrigin, adjSegment)
	}

	if triple.X == 0 && triple.Y == 0 && triple.Z == 0 {
		for q := 0; q < 5; q++ {
			if q == sourceQuintant {
				continue
			}
			distance := min((q-sourceQuintant+5)%5, (sourceQuintant-q+5)%5)
			if ctx.edgeOnly && distance != 1 {
				continue
			}
			targetSegment, targetOrientation := core.QuintantToSegment(q, cell.Origin)
			addNeighbor(ctx, triple, targetOrientation, cell.Origin, targetSegment)
		}
	}

	if triple.X == -ctx.maxRow && triple.Y == ctx.maxRow && triple.Z == 0 {
		prevQuintant := (sourceQuintant - 1 + 5) % 5
		prevAdj := core.FaceAdjacency[cell.Origin.ID][prevQuintant]
		prevAdjOrigin := core.Origins[prevAdj[0]]
		prevAdjSegment, prevAdjOrientation := core.QuintantToSegment(prevAdj[1], prevAdjOrigin)
		addNeighbor(ctx, triple, prevAdjOrientation, prevAdjOrigin, prevAdjSegment)

		crossFace := core.FaceAdjacency[cell.Origin.ID][sourceQuintant]
		crossOrigin := core.Origins[crossFace[0]]
		nextCrossQuintant := (crossFace[1] + 1) % 5
		crossSegment, crossOrientation := core.QuintantToSegment(nextCrossQuintant, crossOrigin)
		addNeighbor(ctx, triple, crossOrientation, crossOrigin, crossSegment)
	}

	result := make([]uint64, 0, len(ctx.neighborSet))
	for neighbor := range ctx.neighborSet {
		result = append(result, neighbor)
	}
	sort.Slice(result, func(i, j int) bool { return result[i] < result[j] })
	return result
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
