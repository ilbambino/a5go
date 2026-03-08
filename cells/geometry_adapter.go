package cells

import "a5go/geometry"

type geometryAdapter struct {
	shape *geometry.PentagonShape
}

func wrapPentagon(shape *geometry.PentagonShape) *geometryAdapter {
	return &geometryAdapter{shape: shape}
}

func (g *geometryAdapter) GetCenter() [2]float64 {
	center := g.shape.GetCenter()
	return [2]float64{center[0], center[1]}
}

func (g *geometryAdapter) GetVertices() [][2]float64 {
	vertices := g.shape.GetVertices()
	out := make([][2]float64, len(vertices))
	for i, vertex := range vertices {
		out[i] = [2]float64{vertex[0], vertex[1]}
	}
	return out
}

func (g *geometryAdapter) SplitEdges(segments int) *geometryAdapter {
	return wrapPentagon(g.shape.SplitEdges(segments))
}

func (g *geometryAdapter) ContainsPoint(point [2]float64) float64 {
	return g.shape.ContainsPoint(geometry.Face(point))
}
