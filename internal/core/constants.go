package core

import "math"

var (
	Phi              = (1 + math.Sqrt(5)) / 2
	TwoPi            = Radians(2 * math.Pi)
	TwoPiOver5       = Radians(2 * math.Pi / 5)
	PiOver5          = Radians(math.Pi / 5)
	PiOver10         = Radians(math.Pi / 10)
	DihedralAngle    = Radians(2 * math.Atan(Phi))
	InterhedralAngle = Radians(math.Pi - float64(DihedralAngle))
	FaceEdgeAngle    = Radians(-0.5*math.Pi + math.Acos(-1/math.Sqrt(3-Phi)))
	DistanceToEdge   = (math.Sqrt(5) - 1) / 2
	DistanceToVertex = 3 - math.Sqrt(5)
	RInscribed       = 1.0
	RMidEdge         = math.Sqrt(3 - Phi)
	RCircumscribed   = math.Sqrt(3) * RMidEdge / Phi
)

const (
	AuthalicRadiusEarth = 6371007.2
	// Match the TS/JS evaluated value exactly so downstream fixtures compare byte-for-byte.
	AuthalicAreaEarth = 510065624779439.1
)
