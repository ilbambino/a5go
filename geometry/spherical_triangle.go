package geometry

type SphericalTriangleShape struct {
	*SphericalPolygonShape
}

func NewSphericalTriangleShape(vertices []Cartesian) *SphericalTriangleShape {
	if len(vertices) != 3 {
		panic("SphericalTriangleShape requires exactly 3 vertices")
	}
	return &SphericalTriangleShape{SphericalPolygonShape: NewSphericalPolygonShape(vertices)}
}
