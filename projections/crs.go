package projections

import (
	"a5go/core"
	"math"
)

type CRS struct {
	vertices    []core.Cartesian
	invocations int
}

func NewCRS() *CRS {
	crs := &CRS{}
	crs.addFaceCenters()
	crs.addVertices()
	crs.addMidpoints()
	if len(crs.vertices) != 62 {
		panic("Failed to construct CRS: vertices length is not 62")
	}
	return crs
}

func (c *CRS) Vertices() []core.Cartesian {
	return c.vertices
}

func (c *CRS) GetVertex(point core.Cartesian) core.Cartesian {
	c.invocations++
	for _, vertex := range c.vertices {
		if distance3(point, vertex) < 1e-5 {
			return vertex
		}
	}
	panic("Failed to find vertex in CRS")
}

func (c *CRS) addFaceCenters() {
	for _, origin := range core.Origins {
		c.add(core.ToCartesian(origin.Axis))
	}
}

func (c *CRS) addVertices() {
	phiVertex := math.Atan(core.DistanceToVertex)
	for _, origin := range core.Origins {
		for i := 0; i < 5; i++ {
			thetaVertex := float64(2*i+1) * math.Pi / 5
			vertex := core.ToCartesian(core.Spherical{thetaVertex + float64(origin.Angle), phiVertex})
			c.add(core.TransformQuat(vertex, origin.Quat))
		}
	}
}

func (c *CRS) addMidpoints() {
	phiMidpoint := math.Atan(core.DistanceToEdge)
	for _, origin := range core.Origins {
		for i := 0; i < 5; i++ {
			thetaMidpoint := float64(2*i) * math.Pi / 5
			midpoint := core.ToCartesian(core.Spherical{thetaMidpoint + float64(origin.Angle), phiMidpoint})
			c.add(core.TransformQuat(midpoint, origin.Quat))
		}
	}
}

func (c *CRS) add(newVertex core.Cartesian) bool {
	normalized := normalize3(newVertex)
	for _, existing := range c.vertices {
		if distance3(normalized, existing) < 1e-5 {
			return false
		}
	}
	c.vertices = append(c.vertices, normalized)
	return true
}

func distance3(a, b core.Cartesian) float64 {
	dx := a[0] - b[0]
	dy := a[1] - b[1]
	dz := a[2] - b[2]
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}
