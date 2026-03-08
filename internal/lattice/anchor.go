package lattice

func isGroup2Orientation(orientation Orientation) bool {
	return orientation == OrientationUW || orientation == OrientationWU
}

func ComputeQ(offset IJ, flips [2]Flip, orientation Orientation) Quaternary {
	i, j := int(offset[0]), int(offset[1])
	flip0, flip1 := flips[0], flips[1]

	imod2 := i & 1
	jmod2 := j & 1
	f0idx := int((flip0 + 1) >> 1)
	f1idx := int((flip1 + 1) >> 1)

	if isGroup2Orientation(orientation) {
		group2Lookup := [2][2][2][2]Quaternary{
			{{{0, 3}, {3, 0}}, {{3, 2}, {2, 3}}},
			{{{2, 1}, {1, 2}}, {{1, 0}, {0, 1}}},
		}
		return group2Lookup[imod2][jmod2][f0idx][f1idx]
	}
	if imod2 == 0 {
		if jmod2 == 0 {
			return Quaternary0
		}
		return Quaternary2
	}
	oddILookup := [2][2][2]Quaternary{
		{{3, 1}, {1, 3}},
		{{1, 3}, {3, 1}},
	}
	return oddILookup[jmod2][f0idx][f1idx]
}

func OffsetFlipsToAnchor(offset IJ, flips [2]Flip, orientation Orientation) Anchor {
	q := ComputeQ(offset, flips, orientation)
	return Anchor{Q: q, Offset: offset, Flips: flips}
}
