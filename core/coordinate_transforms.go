package core

import (
	"a5go/projections"
	"math"
)

type Contour []LonLat

var authalic = projections.AuthalicProjection{}

func DegToRad(deg Degrees) Radians {
	return Radians(float64(deg) * (math.Pi / 180))
}

func RadToDeg(rad Radians) Degrees {
	return Degrees(float64(rad) * (180 / math.Pi))
}

func ToPolar(xy Face) Polar {
	return Polar{math.Hypot(xy[0], xy[1]), math.Atan2(xy[1], xy[0])}
}

func ToFace(polar Polar) Face {
	return Face{
		polar[0] * math.Cos(polar[1]),
		polar[0] * math.Sin(polar[1]),
	}
}

func FaceToIJ(face Face) IJ {
	return IJ(Mat2Transform(BasisInverse, face))
}

func IJToFace(ij IJ) Face {
	return Mat2Transform(Basis, Face(ij))
}

func FaceToBarycentric(p Face, triangle FaceTriangle) Barycentric {
	p1, p2, p3 := triangle[0], triangle[1], triangle[2]
	d31 := [2]float64{p1[0] - p3[0], p1[1] - p3[1]}
	d23 := [2]float64{p3[0] - p2[0], p3[1] - p2[1]}
	d3p := [2]float64{p[0] - p3[0], p[1] - p3[1]}

	det := d23[0]*d31[1] - d23[1]*d31[0]
	b0 := (d23[0]*d3p[1] - d23[1]*d3p[0]) / det
	b1 := (d31[0]*d3p[1] - d31[1]*d3p[0]) / det
	b2 := 1 - (b0 + b1)
	return Barycentric{b0, b1, b2}
}

func BarycentricToFace(b Barycentric, triangle FaceTriangle) Face {
	p1, p2, p3 := triangle[0], triangle[1], triangle[2]
	return Face{
		b[0]*p1[0] + b[1]*p2[0] + b[2]*p3[0],
		b[0]*p1[1] + b[1]*p2[1] + b[2]*p3[1],
	}
}

func ToSpherical(xyz Cartesian) Spherical {
	theta := math.Atan2(xyz[1], xyz[0])
	r := math.Sqrt(xyz[0]*xyz[0] + xyz[1]*xyz[1] + xyz[2]*xyz[2])
	phi := math.Acos(xyz[2] / r)
	return Spherical{theta, phi}
}

func ToCartesian(spherical Spherical) Cartesian {
	theta, phi := spherical[0], spherical[1]
	sinPhi := math.Sin(phi)
	return Cartesian{
		sinPhi * math.Cos(theta),
		sinPhi * math.Sin(theta),
		math.Cos(phi),
	}
}

const longitudeOffset = Degrees(93)

func FromLonLat(lonLat LonLat) Spherical {
	longitude, latitude := Degrees(lonLat[0]), Degrees(lonLat[1])
	theta := DegToRad(longitude + longitudeOffset)
	geodeticLat := DegToRad(latitude)
	authalicLat := authalic.Forward(float64(geodeticLat))
	phi := Radians(math.Pi/2 - authalicLat)
	return Spherical{float64(theta), float64(phi)}
}

func ToLonLatFromSpherical(spherical Spherical) LonLat {
	theta, phi := Radians(spherical[0]), Radians(spherical[1])
	longitude := RadToDeg(theta) - longitudeOffset
	authalicLat := Radians(math.Pi/2 - float64(phi))
	geodeticLat := authalic.Inverse(float64(authalicLat))
	latitude := RadToDeg(Radians(geodeticLat))
	return LonLat{float64(longitude), float64(latitude)}
}

func NormalizeLongitudes(contour Contour) Contour {
	points := make([]Cartesian, len(contour))
	center := Cartesian{}
	for i, lonLat := range contour {
		points[i] = ToCartesian(FromLonLat(lonLat))
		center[0] += points[i][0]
		center[1] += points[i][1]
		center[2] += points[i][2]
	}

	length := math.Sqrt(center[0]*center[0] + center[1]*center[1] + center[2]*center[2])
	center = Cartesian{center[0] / length, center[1] / length, center[2] / length}

	centerLonLat := ToLonLatFromSpherical(ToSpherical(center))
	centerLon := Degrees(centerLonLat[0])
	centerLat := Degrees(centerLonLat[1])
	if centerLat > 89.99 || centerLat < -89.99 {
		centerLon = Degrees(contour[0][0])
	}
	centerLon = Degrees(math.Mod(math.Mod(float64(centerLon)+180, 360)+360, 360) - 180)

	normalized := make(Contour, len(contour))
	for i, point := range contour {
		longitude, latitude := Degrees(point[0]), Degrees(point[1])
		for longitude-centerLon > 180 {
			longitude -= 360
		}
		for longitude-centerLon < -180 {
			longitude += 360
		}
		normalized[i] = LonLat{float64(longitude), float64(latitude)}
	}
	return normalized
}
