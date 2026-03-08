package geometry

import (
	"math"
)

type Face [2]float64
type Mat2 [4]float64
type Mat2D [6]float64
type Pentagon []Face

type PentagonShape struct {
	vertices Pentagon
}

func NewPentagonShape(vertices Pentagon) *PentagonShape {
	p := &PentagonShape{vertices: vertices}
	if !p.isWindingCorrect() {
		reverseFaces(p.vertices)
	}
	return p
}

func reverseFaces(values []Face) {
	for i, j := 0, len(values)-1; i < j; i, j = i+1, j-1 {
		values[i], values[j] = values[j], values[i]
	}
}

func (p *PentagonShape) GetArea() float64 {
	signedArea := 0.0
	n := len(p.vertices)
	for i := 0; i < n; i++ {
		j := (i + 1) % n
		signedArea += (p.vertices[j][0] - p.vertices[i][0]) * (p.vertices[j][1] + p.vertices[i][1])
	}
	return signedArea
}

func (p *PentagonShape) isWindingCorrect() bool {
	return p.GetArea() >= 0
}

func (p *PentagonShape) GetVertices() Pentagon {
	return p.vertices
}

func (p *PentagonShape) Scale(scale float64) *PentagonShape {
	for i := range p.vertices {
		p.vertices[i][0] *= scale
		p.vertices[i][1] *= scale
	}
	return p
}

func (p *PentagonShape) Rotate180() *PentagonShape {
	for i := range p.vertices {
		p.vertices[i][0] = -p.vertices[i][0]
		p.vertices[i][1] = -p.vertices[i][1]
	}
	return p
}

func (p *PentagonShape) ReflectY() *PentagonShape {
	for i := range p.vertices {
		p.vertices[i][1] = -p.vertices[i][1]
	}
	reverseFaces(p.vertices)
	return p
}

func (p *PentagonShape) Translate(translation Face) *PentagonShape {
	for i := range p.vertices {
		p.vertices[i][0] += translation[0]
		p.vertices[i][1] += translation[1]
	}
	return p
}

func (p *PentagonShape) Transform(transform Mat2) *PentagonShape {
	for i := range p.vertices {
		p.vertices[i] = Face{
			transform[0]*p.vertices[i][0] + transform[2]*p.vertices[i][1],
			transform[1]*p.vertices[i][0] + transform[3]*p.vertices[i][1],
		}
	}
	return p
}

func (p *PentagonShape) Transform2D(transform Mat2D) *PentagonShape {
	for i := range p.vertices {
		p.vertices[i] = Face{
			transform[0]*p.vertices[i][0] + transform[2]*p.vertices[i][1] + transform[4],
			transform[1]*p.vertices[i][0] + transform[3]*p.vertices[i][1] + transform[5],
		}
	}
	return p
}

func (p *PentagonShape) Clone() *PentagonShape {
	vertices := make(Pentagon, len(p.vertices))
	for i, vertex := range p.vertices {
		vertices[i] = Face{vertex[0], vertex[1]}
	}
	return NewPentagonShape(vertices)
}

func (p *PentagonShape) GetCenter() Face {
	n := float64(len(p.vertices))
	center := Face{0, 0}
	for _, vertex := range p.vertices {
		center[0] += vertex[0] / n
		center[1] += vertex[1] / n
	}
	return center
}

func (p *PentagonShape) ContainsPoint(point Face) float64 {
	if !p.isWindingCorrect() {
		panic("Pentagon is not counter-clockwise")
	}

	n := len(p.vertices)
	dMax := 1.0
	for i := 0; i < n; i++ {
		v1 := p.vertices[i]
		v2 := p.vertices[(i+1)%n]

		dx := v1[0] - v2[0]
		dy := v1[1] - v2[1]
		px := point[0] - v1[0]
		py := point[1] - v1[1]

		crossProduct := dx*py - dy*px
		if crossProduct < 0 {
			pLength := math.Sqrt(px*px + py*py)
			dMax = min(dMax, crossProduct/pLength)
		}
	}

	return dMax
}

func (p *PentagonShape) SplitEdges(segments int) *PentagonShape {
	if segments <= 1 {
		return p
	}

	newVertices := make(Pentagon, 0, len(p.vertices)*segments)
	n := len(p.vertices)
	for i := 0; i < n; i++ {
		v1 := p.vertices[i]
		v2 := p.vertices[(i+1)%n]
		newVertices = append(newVertices, Face{v1[0], v1[1]})
		for j := 1; j < segments; j++ {
			t := float64(j) / float64(segments)
			newVertices = append(newVertices, Face{
				v1[0] + (v2[0]-v1[0])*t,
				v1[1] + (v2[1]-v1[1])*t,
			})
		}
	}

	return NewPentagonShape(newVertices)
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
