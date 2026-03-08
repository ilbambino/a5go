package lattice

func reversePattern(pattern []int) []int {
	reversed := make([]int, len(pattern))
	for i := range pattern {
		reversed[i] = indexOf(pattern, i)
	}
	return reversed
}

func indexOf(values []int, needle int) int {
	for i, value := range values {
		if value == needle {
			return i
		}
	}
	return -1
}

var (
	Pattern                = []int{0, 1, 3, 4, 5, 6, 7, 2}
	PatternFlipped         = []int{0, 1, 2, 7, 3, 4, 5, 6}
	PatternReversed        = reversePattern(Pattern)
	PatternFlippedReversed = reversePattern(PatternFlipped)
)

func ShiftDigits(digits []Quaternary, i int, flips [2]Flip, invertJ bool, pattern []int) {
	if i <= 0 {
		return
	}

	parentK := digits[i]
	childK := digits[i-1]
	f := flips[0] + flips[1]

	needsShift := true
	first := true

	if invertJ != (f == 0) {
		needsShift = parentK == Quaternary1 || parentK == Quaternary2
		first = parentK == Quaternary1
	} else {
		needsShift = parentK < Quaternary2
		first = parentK == Quaternary0
	}
	if !needsShift {
		return
	}

	src := int(childK)
	if !first {
		src += 4
	}
	dst := pattern[src]
	digits[i-1] = Quaternary(dst % 4)
	digits[i] = Quaternary((int(parentK) + 4 + dst/4 - src/4) % 4)
}
