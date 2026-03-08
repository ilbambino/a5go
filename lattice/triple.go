package lattice

import "math"

type Triple struct {
	X int `json:"x"`
	Y int `json:"y"`
	Z int `json:"z"`
}

func TripleParity(t Triple) int {
	return t.X + t.Y + t.Z
}

func TripleInBounds(t Triple, maxRow int) bool {
	sum := t.X + t.Y + t.Z
	if sum != 0 && sum != 1 {
		return false
	}
	limit := t.Y - sum
	return t.X <= 0 && t.Z <= 0 && t.Y >= 0 && t.Y <= maxRow && t.X >= -limit && t.Z >= -limit
}

func TripleToS(t Triple, resolution int, orientation Orientation) *uint64 {
	anchor := TripleToAnchor(t, resolution, orientation)
	if anchor == nil {
		return nil
	}
	s := AnchorToS(*anchor, resolution, orientation)
	return &s
}

func AnchorToTriple(anchor Anchor) Triple {
	shiftI := 0.25
	shiftJ := 0.25
	flip0, flip1 := anchor.Flips[0], anchor.Flips[1]

	if flip0 == NO && flip1 == YES {
		shiftI = -shiftI
		shiftJ = -shiftJ
	}

	if flip0 == YES && flip1 == YES {
		shiftI = -shiftI
		shiftJ = -shiftJ
	} else if flip0 == YES {
		shiftJ -= 1
	} else if flip1 == YES {
		shiftJ += 1
	}

	i := anchor.Offset[0] + shiftI
	j := anchor.Offset[1] + shiftJ

	r := (i + j) - 0.5
	c := (i - j) + r

	x := int(math.Floor((c+1)/2 - r))
	y := int(r)
	z := int(math.Floor((1 - c) / 2))
	return Triple{X: x, Y: y, Z: z}
}

func TripleToAnchor(t Triple, resolution int, orientation Orientation) *Anchor {
	x, y, z := t.X, t.Y, t.Z
	sum := x + y + z
	if sum != 0 && sum != 1 {
		return nil
	}

	r := float64(y)
	cMin := math.Max(float64(2*x+2*y-1), float64(-2*z-1)+0.0001)
	cMax := math.Min(float64(2*x+2*y+1)-0.0001, float64(1-2*z))
	c := math.Round((cMin + cMax) / 2)

	centerI := (c + 0.5) / 2
	centerJ := r - c/2 + 0.25

	if orientation == OrientationUV || orientation == OrientationVU {
		flips := IJToFlips(IJ{centerI, centerJ}, resolution)
		shiftI := 0.25
		shiftJ := 0.25
		if flips[0] == NO && flips[1] == YES {
			shiftI = -shiftI
			shiftJ = -shiftJ
		}
		if flips[0] == YES && flips[1] == YES {
			shiftI = -shiftI
			shiftJ = -shiftJ
		} else if flips[0] == YES {
			shiftJ -= 1
		} else if flips[1] == YES {
			shiftJ += 1
		}
		offset := IJ{math.Round(centerI - shiftI), math.Round(centerJ - shiftJ)}
		anchor := OffsetFlipsToAnchor(offset, flips, orientation)
		return &anchor
	}

	s := IJToS(IJ{centerI, centerJ}, resolution, orientation)
	anchor := SToAnchor(s, resolution, orientation)
	return &anchor
}
