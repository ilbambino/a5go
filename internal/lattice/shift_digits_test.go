package lattice_test

import (
	"a5go/internal/lattice"
	"a5go/internal/testutil"
	"testing"
)

type shiftDigitsFixtures struct {
	ShiftDigits []struct {
		DigitsBefore []uint8 `json:"digitsBefore"`
		I            int     `json:"i"`
		Flips        [2]int8 `json:"flips"`
		InvertJ      bool    `json:"invertJ"`
		PatternName  string  `json:"patternName"`
		DigitsAfter  []uint8 `json:"digitsAfter"`
	} `json:"shiftDigits"`
}

func TestShiftDigitsFixtures(t *testing.T) {
	var fixtures shiftDigitsFixtures
	testutil.LoadJSON(t, "../../testdata/fixtures/lattice/shift-digits.json", &fixtures)

	patterns := map[string][]int{
		"PATTERN":                  lattice.Pattern,
		"PATTERN_FLIPPED":          lattice.PatternFlipped,
		"PATTERN_REVERSED":         lattice.PatternReversed,
		"PATTERN_FLIPPED_REVERSED": lattice.PatternFlippedReversed,
	}

	for _, f := range fixtures.ShiftDigits {
		digits := make([]lattice.Quaternary, len(f.DigitsBefore))
		for i, digit := range f.DigitsBefore {
			digits[i] = lattice.Quaternary(digit)
		}
		lattice.ShiftDigits(digits, f.I, [2]lattice.Flip{lattice.Flip(f.Flips[0]), lattice.Flip(f.Flips[1])}, f.InvertJ, patterns[f.PatternName])
		for i, got := range digits {
			if got != lattice.Quaternary(f.DigitsAfter[i]) {
				t.Fatalf("case %+v index %d: got %d want %d", f, i, got, f.DigitsAfter[i])
			}
		}
	}
}
