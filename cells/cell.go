package cells

import (
	"a5go/core"
	"a5go/geometry"
	"a5go/lattice"
	"a5go/projections"
	"math"
)

type CellBoundaryOptions struct {
	ClosedRing   bool
	Segments     int
	AutoSegments bool
}

var (
	dodecahedron = projections.NewDodecahedronProjection()
)

func LonLatToCell(lonLat core.LonLat, resolution int) uint64 {
	if resolution == -1 {
		return core.WorldCell
	}
	if resolution < core.FirstHilbertResolution {
		return core.Serialize(lonLatToEstimate(lonLat, resolution))
	}

	hilbertResolution := 1 + resolution - core.FirstHilbertResolution
	samples := []core.LonLat{lonLat}
	const n = 25
	scale := 50 / math.Pow(2, float64(hilbertResolution))
	for i := 0; i < n; i++ {
		r := (float64(i) / n) * scale
		sample := core.LonLat{
			lonLat[0] + math.Cos(float64(i))*r,
			lonLat[1] + math.Sin(float64(i))*r,
		}
		samples = append(samples, sample)
	}

	estimateSet := map[uint64]struct{}{}
	cells := make([]struct {
		cell     core.A5Cell
		distance float64
	}, 0, len(samples))
	for _, sample := range samples {
		estimate := lonLatToEstimate(sample, resolution)
		estimateKey := core.Serialize(estimate)
		if _, exists := estimateSet[estimateKey]; exists {
			continue
		}
		estimateSet[estimateKey] = struct{}{}
		distance := A5CellContainsPoint(estimate, lonLat)
		if distance > 0 {
			return estimateKey
		}
		cells = append(cells, struct {
			cell     core.A5Cell
			distance float64
		}{cell: estimate, distance: distance})
	}

	best := cells[0]
	for _, candidate := range cells[1:] {
		if candidate.distance > best.distance {
			best = candidate
		}
	}
	return core.Serialize(best.cell)
}

func lonLatToEstimate(lonLat core.LonLat, resolution int) core.A5Cell {
	spherical := core.FromLonLat(lonLat)
	origin := *core.FindNearestOrigin(spherical)

	dodecPoint := dodecahedron.Forward(spherical, origin.ID)
	polar := core.ToPolar(dodecPoint)
	quintant := core.GetQuintantPolar(polar)
	segment, orientation := core.QuintantToSegment(quintant, &origin)
	if resolution < core.FirstHilbertResolution {
		return core.A5Cell{S: 0, Segment: segment, Origin: &origin, Resolution: resolution}
	}

	if quintant != 0 {
		extraAngle := 2 * float64(core.PiOver5) * float64(quintant)
		rotation := core.Mat2FromRotation(-extraAngle)
		dodecPoint = core.Mat2Transform(rotation, dodecPoint)
	}

	hilbertResolution := 1 + resolution - core.FirstHilbertResolution
	scale := math.Pow(2, float64(hilbertResolution))
	dodecPoint[0] *= scale
	dodecPoint[1] *= scale

	ij := core.FaceToIJ(dodecPoint)
	s := lattice.IJToS(lattice.IJ(ij), hilbertResolution, orientation)
	return core.A5Cell{S: s, Segment: segment, Origin: &origin, Resolution: resolution}
}

func getPentagon(cell core.A5Cell) *geometry.PentagonShape {
	quintant, orientation := core.SegmentToQuintant(cell.Segment, cell.Origin)
	if cell.Resolution == core.FirstHilbertResolution-1 {
		return core.GetQuintantVertices(quintant)
	}
	if cell.Resolution == core.FirstHilbertResolution-2 {
		return core.GetFaceVertices()
	}

	hilbertResolution := cell.Resolution - core.FirstHilbertResolution + 1
	anchor := lattice.SToAnchor(cell.S, hilbertResolution, orientation)
	return core.GetPentagonVertices(hilbertResolution, quintant, anchor)
}

func CellToSpherical(cellID uint64) core.Spherical {
	cell := core.Deserialize(cellID)
	pentagon := getPentagon(cell)
	center := pentagon.Center()
	return dodecahedron.Inverse(core.Face(center), cell.Origin.ID)
}

func CellToLonLat(cellID uint64) core.LonLat {
	if cellID == core.WorldCell {
		return core.LonLat{0, 0}
	}
	return core.ToLonLatFromSpherical(CellToSpherical(cellID))
}

func CellToBoundary(cellID uint64, options ...CellBoundaryOptions) []core.LonLat {
	if cellID == core.WorldCell {
		return []core.LonLat{}
	}

	opts := CellBoundaryOptions{ClosedRing: true, AutoSegments: true}
	if len(options) > 0 {
		opts = options[0]
		if !opts.AutoSegments && opts.Segments == 0 {
			opts.Segments = 1
		}
	}

	cell := core.Deserialize(cellID)
	segments := resolveBoundarySegments(cell.Resolution, opts)

	pentagon := getPentagon(cell).SplitEdges(segments)
	vertices := pentagon.Vertices()
	boundary := make([]core.LonLat, len(vertices))
	for i, vertex := range vertices {
		unprojected := dodecahedron.Inverse(core.Face(vertex), cell.Origin.ID)
		boundary[i] = core.ToLonLatFromSpherical(unprojected)
	}
	normalizedBoundary := core.NormalizeLongitudes(boundary)
	if opts.ClosedRing && len(normalizedBoundary) > 0 {
		normalizedBoundary = append(normalizedBoundary, normalizedBoundary[0])
	}
	for i, j := 0, len(normalizedBoundary)-1; i < j; i, j = i+1, j-1 {
		normalizedBoundary[i], normalizedBoundary[j] = normalizedBoundary[j], normalizedBoundary[i]
	}
	return normalizedBoundary
}

func A5CellContainsPoint(cell core.A5Cell, point core.LonLat) float64 {
	pentagon := getPentagon(cell)
	spherical := core.FromLonLat(point)
	projectedPoint := dodecahedron.Forward(spherical, cell.Origin.ID)
	return pentagon.ContainsPoint([2]float64(projectedPoint))
}

func resolveBoundarySegments(resolution int, opts CellBoundaryOptions) int {
	if opts.AutoSegments {
		return int(math.Max(1, math.Pow(2, float64(6-resolution))))
	}
	if opts.Segments <= 0 {
		return 1
	}
	return opts.Segments
}
