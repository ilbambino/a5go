package core

const (
	FirstHilbertResolution        = 2
	MaxResolution                 = 30
	HilbertStartBit               = uint(58)
	RemovalMask            uint64 = 0x03ffffffffffffff
	WorldCell              uint64 = 0
)

func GetResolution(index uint64) int {
	if index == 0 {
		return -1
	}

	resolution := MaxResolution - 1
	shifted := index >> 1
	if shifted == 0 {
		return -1
	}

	low32 := uint32(shifted & 0xFFFFFFFF)
	var remaining uint32

	if low32 == 0 {
		shifted >>= 32
		resolution -= 16
		remaining = uint32(shifted)
	} else {
		remaining = low32
	}

	if (remaining & 0xFFFF) == 0 {
		remaining >>= 16
		resolution -= 8
	}
	if resolution >= 6 && (remaining&0xFF) == 0 {
		remaining >>= 8
		resolution -= 4
	}
	if resolution >= 4 && (remaining&0xF) == 0 {
		remaining >>= 4
		resolution -= 2
	}
	for resolution > -1 && (remaining&0b1) == 0 {
		resolution--
		if resolution < FirstHilbertResolution {
			remaining >>= 1
		} else {
			remaining >>= 2
		}
	}
	return resolution
}

func Deserialize(index uint64) A5Cell {
	resolution := GetResolution(index)
	if resolution == -1 {
		return A5Cell{Origin: Origins[0], Segment: 0, S: 0, Resolution: resolution}
	}

	top6Bits := int(index >> 58)
	var origin *Origin
	var segment int

	if resolution == 0 {
		origin = Origins[top6Bits]
		segment = 0
	} else {
		origin = Origins[top6Bits/5]
		segment = (top6Bits + origin.FirstQuintant) % 5
	}

	if resolution < FirstHilbertResolution {
		return A5Cell{Origin: origin, Segment: segment, S: 0, Resolution: resolution}
	}

	hilbertLevels := resolution - FirstHilbertResolution + 1
	hilbertBits := uint(2 * hilbertLevels)
	shift := HilbertStartBit - hilbertBits
	s := (index & RemovalMask) >> shift
	return A5Cell{Origin: origin, Segment: segment, S: s, Resolution: resolution}
}

func Serialize(cell A5Cell) uint64 {
	origin, segment, s, resolution := cell.Origin, cell.Segment, cell.S, cell.Resolution
	if resolution > MaxResolution {
		panicString("Resolution (" + itoa(resolution) + ") is too large")
	}
	if resolution == -1 {
		return WorldCell
	}

	var r uint
	if resolution < FirstHilbertResolution {
		r = uint(resolution + 1)
	} else {
		hilbertResolution := 1 + resolution - FirstHilbertResolution
		r = uint(2*hilbertResolution + 1)
	}

	segmentN := (segment - origin.FirstQuintant + 5) % 5
	var index uint64
	if resolution == 0 {
		index = uint64(origin.ID) << 58
	} else {
		index = uint64(5*origin.ID+segmentN) << 58
	}

	if resolution >= FirstHilbertResolution {
		hilbertLevels := resolution - FirstHilbertResolution + 1
		hilbertBits := uint(2 * hilbertLevels)
		if s >= (uint64(1) << hilbertBits) {
			panicString("S (" + uitoa(s) + ") is too large for resolution level " + itoa(resolution))
		}
		index += s << (HilbertStartBit - hilbertBits)
	}

	index |= uint64(1) << (HilbertStartBit - r)
	return index
}

func CellToChildren(index uint64, childResolution ...int) []uint64 {
	cell := Deserialize(index)
	newResolution := cell.Resolution + 1
	if len(childResolution) > 0 {
		newResolution = childResolution[0]
	}
	currentResolution := cell.Resolution

	if newResolution < currentResolution {
		panicString("Target resolution (" + itoa(newResolution) + ") must be equal to or greater than current resolution (" + itoa(currentResolution) + ")")
	}
	if newResolution > MaxResolution {
		panicString("Target resolution (" + itoa(newResolution) + ") exceeds maximum resolution (" + itoa(MaxResolution) + ")")
	}
	if newResolution == currentResolution {
		return []uint64{index}
	}

	newOrigins := []*Origin{cell.Origin}
	newSegments := []int{cell.Segment}
	if currentResolution == -1 {
		newOrigins = Origins
	}
	if (currentResolution == -1 && newResolution > 0) || currentResolution == 0 {
		newSegments = []int{0, 1, 2, 3, 4}
	}

	startResolution := currentResolution
	if startResolution < FirstHilbertResolution-1 {
		startResolution = FirstHilbertResolution - 1
	}
	resolutionDiff := newResolution - startResolution
	childrenCount := 1
	for i := 0; i < resolutionDiff; i++ {
		childrenCount *= 4
	}
	shiftedS := cell.S << uint(2*resolutionDiff)
	children := make([]uint64, 0, len(newOrigins)*len(newSegments)*childrenCount)
	for _, newOrigin := range newOrigins {
		for _, newSegment := range newSegments {
			for i := 0; i < childrenCount; i++ {
				newS := shiftedS + uint64(i)
				children = append(children, Serialize(A5Cell{Origin: newOrigin, Segment: newSegment, S: newS, Resolution: newResolution}))
			}
		}
	}
	return children
}

func CellToParent(index uint64, parentResolution ...int) uint64 {
	cell := Deserialize(index)
	newResolution := cell.Resolution - 1
	if len(parentResolution) > 0 {
		newResolution = parentResolution[0]
	}
	currentResolution := cell.Resolution

	if newResolution == -1 {
		return WorldCell
	}
	if newResolution < 0 {
		panicString("Target resolution (" + itoa(newResolution) + ") cannot be negative")
	}
	if newResolution > currentResolution {
		panicString("Target resolution (" + itoa(newResolution) + ") must be equal to or less than current resolution (" + itoa(currentResolution) + ")")
	}
	if newResolution == currentResolution {
		return index
	}

	resolutionDiff := currentResolution - newResolution
	shiftedS := cell.S >> uint(2*resolutionDiff)
	return Serialize(A5Cell{Origin: cell.Origin, Segment: cell.Segment, S: shiftedS, Resolution: newResolution})
}

func GetRes0Cells() []uint64 {
	return CellToChildren(WorldCell, 0)
}

func IsFirstChild(index uint64, resolution ...int) bool {
	res := GetResolution(index)
	if len(resolution) > 0 {
		res = resolution[0]
	}
	if res < 2 {
		top6Bits := int(index >> HilbertStartBit)
		childCount := 12
		if res != 0 {
			childCount = 5
		}
		return top6Bits%childCount == 0
	}
	sPosition := uint(2 * (MaxResolution - res))
	sMask := uint64(3) << sPosition
	return (index & sMask) == 0
}

func GetStride(resolution int) uint64 {
	if resolution < 2 {
		return uint64(1) << HilbertStartBit
	}
	sPosition := uint(2 * (MaxResolution - resolution))
	return uint64(1) << sPosition
}
