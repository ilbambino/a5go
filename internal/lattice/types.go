package lattice

type Orientation string

const (
	OrientationUV Orientation = "uv"
	OrientationVU Orientation = "vu"
	OrientationUW Orientation = "uw"
	OrientationWU Orientation = "wu"
	OrientationVW Orientation = "vw"
	OrientationWV Orientation = "wv"
)

type Quaternary uint8

const (
	Quaternary0 Quaternary = 0
	Quaternary1 Quaternary = 1
	Quaternary2 Quaternary = 2
	Quaternary3 Quaternary = 3
)

type Flip int8

const (
	YES Flip = -1
	NO  Flip = 1
)

type Anchor struct {
	Q      Quaternary
	Offset IJ
	Flips  [2]Flip
}

type IJ [2]float64
type KJ [2]float64
