package core

type Degrees float64
type Radians float64

// 2D coordinate systems.
type Face [2]float64
type Polar [2]float64
type IJ [2]float64
type KJ [2]float64

// 3D coordinate systems.
type Cartesian [3]float64
type Spherical [2]float64
type LonLat [2]float64

// Barycentric coordinates.
type Barycentric [3]float64
type FaceTriangle [3]Face
type SphericalTriangle [3]Cartesian
