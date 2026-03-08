package lattice

var (
	kPos = KJ{1, 0}
	jPos = KJ{0, 1}
	kNeg = KJ{-1, 0}
	jNeg = KJ{0, -1}
	zero = KJ{0, 0}
)

func QuaternaryToKJ(n Quaternary, flips [2]Flip) KJ {
	p := zero
	q := zero

	switch {
	case flips[0] == NO && flips[1] == NO:
		p = kPos
		q = jPos
	case flips[0] == YES && flips[1] == NO:
		p = jNeg
		q = kNeg
	case flips[0] == NO && flips[1] == YES:
		p = jPos
		q = kPos
	case flips[0] == YES && flips[1] == YES:
		p = kNeg
		q = jNeg
	}

	switch n {
	case Quaternary0:
		return zero
	case Quaternary1:
		return p
	case Quaternary2:
		return KJ{q[0] + p[0], q[1] + p[1]}
	case Quaternary3:
		return KJ{q[0] + 2*p[0], q[1] + 2*p[1]}
	default:
		panic("invalid Quaternary value")
	}
}

func QuaternaryToFlips(n Quaternary) [2]Flip {
	return [...][2]Flip{
		{NO, NO},
		{NO, YES},
		{NO, NO},
		{YES, NO},
	}[n]
}

func IJToQuaternary(ij IJ, flips [2]Flip) Quaternary {
	i := ij[0]
	j := ij[1]
	digit := Quaternary0

	a := i + j
	if flips[0] == YES {
		a = -(i + j)
	}
	b := i
	if flips[1] == YES {
		b = -i
	}
	c := j
	if flips[0] == YES {
		c = -j
	}

	if flips[0]+flips[1] == 0 {
		if c < 1 {
			digit = Quaternary0
		} else if b > 1 {
			digit = Quaternary3
		} else if a > 1 {
			digit = Quaternary2
		} else {
			digit = Quaternary1
		}
	} else {
		if a < 1 {
			digit = Quaternary0
		} else if b > 1 {
			digit = Quaternary3
		} else if c > 1 {
			digit = Quaternary2
		} else {
			digit = Quaternary1
		}
	}

	return digit
}
