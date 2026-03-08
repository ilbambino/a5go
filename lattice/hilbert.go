package lattice

import "math"

var flipShift = IJ{-1, 1}

func SToAnchor(s uint64, resolution int, orientation Orientation, doShiftDigits ...bool) Anchor {
	shiftDigitsEnabled := true
	if len(doShiftDigits) > 0 {
		shiftDigitsEnabled = doShiftDigits[0]
	}

	input := s
	reverse := orientation == OrientationVU || orientation == OrientationWU || orientation == OrientationVW
	invertJ := orientation == OrientationWV || orientation == OrientationVW
	flipIJ := orientation == OrientationWU || orientation == OrientationUW
	if reverse {
		input = (uint64(1) << uint(2*resolution)) - input - 1
	}
	anchor := sToAnchorInternal(input, resolution, invertJ, flipIJ, shiftDigitsEnabled)

	if flipIJ {
		i, j := anchor.Offset[0], anchor.Offset[1]
		flipX, flipY := anchor.Flips[0], anchor.Flips[1]
		anchor.Offset = IJ{j, i}
		if flipX == YES {
			anchor.Offset = addIJ(anchor.Offset, flipShift)
		}
		if flipY == YES {
			anchor.Offset = subIJ(anchor.Offset, flipShift)
		}
	}
	if invertJ {
		i, j := anchor.Offset[0], anchor.Offset[1]
		flips := anchor.Flips
		newJ := float64(uint64(1)<<uint(resolution)) - (i + j)
		flips[0] = -flips[0]
		anchor.Offset[1] = newJ
		anchor.Flips = flips
	}
	return anchor
}

func sToAnchorInternal(s uint64, resolution int, invertJ, flipIJ, doShiftDigits bool) Anchor {
	offset := KJ{0, 0}
	flips := [2]Flip{NO, NO}
	input := s

	digits := make([]Quaternary, 0, resolution)
	for input > 0 || len(digits) < resolution {
		digits = append(digits, Quaternary(input%4))
		input >>= 2
	}

	pattern := Pattern
	if flipIJ {
		pattern = PatternFlipped
	}

	for i := len(digits) - 1; i >= 0; i-- {
		if doShiftDigits {
			ShiftDigits(digits, i, flips, invertJ, pattern)
		}
		qFlips := QuaternaryToFlips(digits[i])
		flips = [2]Flip{flips[0] * qFlips[0], flips[1] * qFlips[1]}
	}

	flips = [2]Flip{NO, NO}
	for i := len(digits) - 1; i >= 0; i-- {
		offset = KJ{offset[0] * 2, offset[1] * 2}
		childOffset := QuaternaryToKJ(digits[i], flips)
		offset = KJ{offset[0] + childOffset[0], offset[1] + childOffset[1]}
		qFlips := QuaternaryToFlips(digits[i])
		flips = [2]Flip{flips[0] * qFlips[0], flips[1] * qFlips[1]}
	}

	q := Quaternary0
	if len(digits) > 0 {
		q = digits[0]
	}
	return Anchor{Q: q, Offset: KJToIJ(offset), Flips: flips}
}

func IJToS(input IJ, resolution int, orientation Orientation, doShiftDigits ...bool) uint64 {
	shiftDigitsEnabled := true
	if len(doShiftDigits) > 0 {
		shiftDigitsEnabled = doShiftDigits[0]
	}

	reverse := orientation == OrientationVU || orientation == OrientationWU || orientation == OrientationVW
	invertJ := orientation == OrientationWV || orientation == OrientationVW
	flipIJ := orientation == OrientationWU || orientation == OrientationUW

	ij := input
	if flipIJ {
		ij[0], ij[1] = input[1], input[0]
	}
	if invertJ {
		i, j := ij[0], ij[1]
		ij[1] = float64(uint64(1)<<uint(resolution)) - (i + j)
	}

	s := ijToSInternal(ij, invertJ, flipIJ, resolution, shiftDigitsEnabled)
	if reverse {
		s = (uint64(1) << uint(2*resolution)) - s - 1
	}
	return s
}

func ijToSInternal(input IJ, invertJ, flipIJ bool, resolution int, doShiftDigits bool) uint64 {
	numDigits := resolution
	digits := make([]Quaternary, numDigits)
	flips := [2]Flip{NO, NO}
	pivot := IJ{0, 0}

	for i := numDigits - 1; i >= 0; i-- {
		relativeOffset := IJ{input[0] - pivot[0], input[1] - pivot[1]}
		scale := float64(uint64(1) << uint(i))
		scaledOffset := IJ{relativeOffset[0] / scale, relativeOffset[1] / scale}
		digit := IJToQuaternary(scaledOffset, flips)
		digits[i] = digit

		childOffset := KJToIJ(QuaternaryToKJ(digit, flips))
		upscaledChildOffset := IJ{childOffset[0] * scale, childOffset[1] * scale}
		pivot = addIJ(pivot, upscaledChildOffset)
		qFlips := QuaternaryToFlips(digit)
		flips = [2]Flip{flips[0] * qFlips[0], flips[1] * qFlips[1]}
	}

	pattern := PatternReversed
	if flipIJ {
		pattern = PatternFlippedReversed
	}

	for i := 0; i < len(digits); i++ {
		qFlips := QuaternaryToFlips(digits[i])
		flips = [2]Flip{flips[0] * qFlips[0], flips[1] * qFlips[1]}
		if doShiftDigits {
			ShiftDigits(digits, i, flips, invertJ, pattern)
		}
	}

	var output uint64
	for i := numDigits - 1; i >= 0; i-- {
		scale := uint(2 * i)
		output += uint64(digits[i]) << scale
	}
	return output
}

func IJToFlips(input IJ, resolution int) [2]Flip {
	numDigits := resolution
	flips := [2]Flip{NO, NO}
	pivot := IJ{0, 0}

	for i := numDigits - 1; i >= 0; i-- {
		relativeOffset := IJ{input[0] - pivot[0], input[1] - pivot[1]}
		scale := float64(uint64(1) << uint(i))
		scaledOffset := IJ{relativeOffset[0] / scale, relativeOffset[1] / scale}
		digit := IJToQuaternary(scaledOffset, flips)
		childOffset := KJToIJ(QuaternaryToKJ(digit, flips))
		upscaledChildOffset := IJ{childOffset[0] * scale, childOffset[1] * scale}
		pivot = addIJ(pivot, upscaledChildOffset)
		qFlips := QuaternaryToFlips(digit)
		flips = [2]Flip{flips[0] * qFlips[0], flips[1] * qFlips[1]}
	}
	return flips
}

var probeOffsets = [][2]float64{
	{0.1 * math.Cos(45*math.Pi/180), 0.1 * math.Sin(45*math.Pi/180)},
	{0.1 * math.Cos(113*math.Pi/180), 0.1 * math.Sin(113*math.Pi/180)},
	{0.1 * math.Cos(293*math.Pi/180), 0.1 * math.Sin(293*math.Pi/180)},
	{0.1 * math.Cos(225*math.Pi/180), 0.1 * math.Sin(225*math.Pi/180)},
}

func AnchorToS(anchor Anchor, resolution int, orientation Orientation) uint64 {
	i, j := anchor.Offset[0], anchor.Offset[1]
	index := int((1 - anchor.Flips[0]) + (1-anchor.Flips[1])/2)
	probeOffset := probeOffsets[index]
	return IJToS(IJ{i + probeOffset[0], j + probeOffset[1]}, resolution, orientation)
}

func addIJ(a, b IJ) IJ {
	return IJ{a[0] + b[0], a[1] + b[1]}
}

func subIJ(a, b IJ) IJ {
	return IJ{a[0] - b[0], a[1] - b[1]}
}
