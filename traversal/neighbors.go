package traversal

import (
	"a5go/lattice"
	"sort"
)

type neighborPattern [4]int

var neighbors = map[uint8][]neighborPattern{
	0: {{0, -2, -1, 1}, {0, -2, -1, -1}, {0, -1, 1, -1}, {0, -1, -1, -1}, {0, -1, 1, 1}, {1, -2, -1, -1}, {1, -1, -1, 1}, {1, -1, 1, -1}, {1, 0, 1, -1}, {2, -1, 1, -1}, {2, -2, -1, -1}},
	1: {{-1, -1, -1, 1}, {0, -2, -1, -1}, {0, -1, -1, -1}, {0, -1, 1, -1}, {0, 0, -1, 1}, {0, 0, -1, -1}, {0, 1, 1, -1}, {0, 1, 1, 1}, {1, -2, -1, -1}, {1, -1, 1, -1}, {1, -1, -1, -1}, {1, 0, 1, -1}},
	2: {{-2, 2, -1, -1}, {-2, 1, 1, -1}, {-1, 0, 1, -1}, {-1, 1, 1, -1}, {-1, 1, -1, 1}, {-1, 2, -1, -1}, {0, 1, -1, -1}, {0, 1, 1, -1}, {0, 1, 1, 1}, {0, 2, -1, -1}, {0, 2, -1, 1}},
	3: {{-1, 0, 1, -1}, {-1, 1, 1, -1}, {-1, 1, -1, -1}, {-1, 2, -1, -1}, {0, -1, 1, -1}, {0, -1, 1, 1}, {0, 0, -1, -1}, {0, 0, -1, 1}, {0, 1, -1, -1}, {0, 1, 1, -1}, {0, 2, -1, -1}, {1, 1, -1, 1}},
	4: {{0, -1, 1, -1}, {0, -1, 1, 1}, {0, 0, -1, -1}, {0, 0, -1, 1}, {0, 1, -1, -1}, {1, 0, -1, -1}, {1, 0, 1, -1}, {1, -1, 1, -1}, {1, 1, -1, 1}, {2, -1, 1, -1}, {2, 0, -1, -1}},
	5: {{-1, 1, -1, 1}, {0, -1, 1, -1}, {0, 0, -1, -1}, {0, 1, -1, -1}, {0, 1, 1, -1}, {0, 1, 1, 1}, {0, 2, -1, -1}, {0, 2, -1, 1}, {1, -1, 1, -1}, {1, 0, -1, -1}, {1, 0, 1, -1}, {1, 1, -1, -1}},
	6: {{-2, 0, -1, -1}, {-2, 1, 1, -1}, {-1, -1, -1, 1}, {-1, 0, -1, -1}, {-1, 0, 1, -1}, {-1, 1, 1, -1}, {0, -1, -1, -1}, {0, 0, -1, -1}, {0, 0, -1, 1}, {0, 1, 1, -1}, {0, 1, 1, 1}},
	7: {{-1, -1, -1, -1}, {-1, 0, -1, -1}, {-1, 0, 1, -1}, {-1, 1, 1, -1}, {0, -2, -1, -1}, {0, -2, -1, 1}, {0, -1, -1, -1}, {0, -1, 1, -1}, {0, -1, 1, 1}, {0, 0, -1, -1}, {0, 1, 1, -1}, {1, -1, -1, 1}},
}

func pentagonFlavor(anchor lattice.Anchor) uint8 {
	f := 0
	if anchor.Flips[1] == lattice.YES {
		f += 2
	}
	q := anchor.Q
	flipSum := anchor.Flips[0] + anchor.Flips[1]
	if ((flipSum == -2 || flipSum == 2) && q > lattice.Quaternary1) || (flipSum == 0 && (q == lattice.Quaternary0 || q == lattice.Quaternary3)) {
		f += 1
	}
	if flipSum == -2 || flipSum == 2 {
		f += 4
	}
	return uint8(f)
}

func IsNeighbor(origin lattice.Anchor, candidate lattice.Anchor) bool {
	originFlavor := pentagonFlavor(origin)
	candidateFlavor := pentagonFlavor(candidate)
	if originFlavor == candidateFlavor {
		return false
	}
	relative := neighborPattern{
		int(candidate.Offset[0] - origin.Offset[0]),
		int(candidate.Offset[1] - origin.Offset[1]),
		int(candidate.Flips[0] * origin.Flips[0]),
		int(candidate.Flips[1] * origin.Flips[1]),
	}
	for _, n := range neighbors[originFlavor] {
		if relative == n {
			return true
		}
	}
	return false
}

func sortedUint64(values []uint64) []uint64 {
	sort.Slice(values, func(i, j int) bool { return values[i] < values[j] })
	return values
}
