package geometry

import "math"

type Cartesian [3]float64
type SphericalPolygon []Cartesian

type SphericalPolygonShape struct {
	vertices SphericalPolygon
	area     *float64
}

func NewSphericalPolygonShape(vertices SphericalPolygon) *SphericalPolygonShape {
	return &SphericalPolygonShape{vertices: vertices}
}

func (s *SphericalPolygonShape) GetBoundary(nSegments int, closedRing bool) SphericalPolygon {
	points := make(SphericalPolygon, 0, len(s.vertices)*nSegments+1)
	n := len(s.vertices)
	for seg := 0; seg < n*nSegments; seg++ {
		t := float64(seg) / float64(nSegments)
		points = append(points, s.Slerp(t))
	}
	if closedRing && len(points) > 0 {
		points = append(points, points[0])
	}
	return points
}

func (s *SphericalPolygonShape) Slerp(t float64) Cartesian {
	n := len(s.vertices)
	f := math.Mod(t, 1)
	i := int(math.Floor(math.Mod(t, float64(n))))
	j := (i + 1) % n
	return slerp3(s.vertices[i], s.vertices[j], f)
}

func (s *SphericalPolygonShape) getTransformedVertices(t int) (Cartesian, Cartesian, Cartesian) {
	n := len(s.vertices)
	i := t % n
	j := (i + 1) % n
	k := (i + n - 1) % n
	v := s.vertices[i]
	va := Cartesian{s.vertices[j][0] - v[0], s.vertices[j][1] - v[1], s.vertices[j][2] - v[2]}
	vb := Cartesian{s.vertices[k][0] - v[0], s.vertices[k][1] - v[1], s.vertices[k][2] - v[2]}
	return v, normalize(va), normalize(vb)
}

func (s *SphericalPolygonShape) ContainsPoint(point Cartesian) float64 {
	n := len(s.vertices)
	thetaDeltaMin := math.Inf(1)
	for i := 0; i < n; i++ {
		v, va, vb := s.getTransformedVertices(i)
		vp := normalize(Cartesian{point[0] - v[0], point[1] - v[1], point[2] - v[2]})
		crossAP := cross3(va, vp)
		crossPB := cross3(vp, vb)
		sinAP := dot3(v, crossAP)
		sinPB := dot3(v, crossPB)
		thetaDeltaMin = math.Min(thetaDeltaMin, math.Min(sinAP, sinPB))
	}
	return thetaDeltaMin
}

func (s *SphericalPolygonShape) getTriangleArea(v1, v2, v3 Cartesian) float64 {
	midA := normalize(lerp3(v2, v3, 0.5))
	midB := normalize(lerp3(v3, v1, 0.5))
	midC := normalize(lerp3(v1, v2, 0.5))
	stp := dot3(midA, cross3(midB, midC))
	clamped := math.Max(-1, math.Min(1, stp))
	if math.Abs(clamped) < 1e-8 {
		return 2 * clamped
	}
	return 2 * math.Asin(clamped)
}

func (s *SphericalPolygonShape) GetArea() float64 {
	if s.area != nil {
		return *s.area
	}
	area := s.computeArea()
	s.area = &area
	return area
}

func (s *SphericalPolygonShape) computeArea() float64 {
	if len(s.vertices) < 3 {
		return 0
	}
	if len(s.vertices) == 3 {
		return s.getTriangleArea(s.vertices[0], s.vertices[1], s.vertices[2])
	}
	center := Cartesian{}
	for _, vertex := range s.vertices {
		center[0] += vertex[0]
		center[1] += vertex[1]
		center[2] += vertex[2]
	}
	center = normalize(center)
	area := 0.0
	for i := 0; i < len(s.vertices); i++ {
		v1 := s.vertices[i]
		v2 := s.vertices[(i+1)%len(s.vertices)]
		triArea := s.getTriangleArea(center, v1, v2)
		if !math.IsNaN(triArea) {
			area += triArea
		}
	}
	return area
}

func normalize(v Cartesian) Cartesian {
	length := math.Sqrt(v[0]*v[0] + v[1]*v[1] + v[2]*v[2])
	if length == 0 {
		return v
	}
	return Cartesian{v[0] / length, v[1] / length, v[2] / length}
}

func dot3(a, b Cartesian) float64 {
	return a[0]*b[0] + a[1]*b[1] + a[2]*b[2]
}

func cross3(a, b Cartesian) Cartesian {
	return Cartesian{
		a[1]*b[2] - a[2]*b[1],
		a[2]*b[0] - a[0]*b[2],
		a[0]*b[1] - a[1]*b[0],
	}
}

func lerp3(a, b Cartesian, t float64) Cartesian {
	return Cartesian{
		a[0] + (b[0]-a[0])*t,
		a[1] + (b[1]-a[1])*t,
		a[2] + (b[2]-a[2])*t,
	}
}

func slerp3(a, b Cartesian, t float64) Cartesian {
	denominator := math.Sqrt(dot3(a, a)) * math.Sqrt(dot3(b, b))
	if denominator == 0 {
		return lerp3(a, b, t)
	}
	cosine := dot3(a, b) / denominator
	if cosine > 1 {
		cosine = 1
	}
	if cosine < -1 {
		cosine = -1
	}
	gamma := math.Acos(cosine)
	if gamma < 1e-12 {
		return lerp3(a, b, t)
	}
	weightA := math.Sin((1-t)*gamma) / math.Sin(gamma)
	weightB := math.Sin(t*gamma) / math.Sin(gamma)
	return Cartesian{
		a[0]*weightA + b[0]*weightB,
		a[1]*weightA + b[1]*weightB,
		a[2]*weightA + b[2]*weightB,
	}
}
