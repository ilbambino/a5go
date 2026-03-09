package core

import (
	"errors"
	"math/bits"
)

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
	shifted := index >> 1
	if shifted == 0 {
		return -1
	}

	tz := bits.TrailingZeros64(shifted)
	if tz >= 55 {
		return 56 - tz
	}
	return (58 - tz) / 2
}

func Deserialize(index uint64) A5Cell {
	resolution := GetResolution(index)
	if resolution == -1 {
		return A5Cell{Origin: &Origins[0], Segment: 0, S: 0, Resolution: resolution}
	}

	top6Bits := int(index >> 58)
	var origin *Origin
	var segment int

	if resolution == 0 {
		origin = &Origins[top6Bits]
		segment = 0
	} else {
		origin = &Origins[top6Bits/5]
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

func Serialize(cell A5Cell) (uint64, error) {
	origin, segment, s, resolution := cell.Origin, cell.Segment, cell.S, cell.Resolution
	if resolution > MaxResolution {
		return 0, errors.New("resolution exceeds maximum")
	}
	if resolution == -1 {
		return WorldCell, nil
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
			return 0, errors.New("S value too large for resolution")
		}
		index += s << (HilbertStartBit - hilbertBits)
	}

	index |= uint64(1) << (HilbertStartBit - r)
	return index, nil
}

func CellToChildren(index uint64, childResolution ...int) ([]uint64, error) {
	cell := Deserialize(index)
	newResolution := cell.Resolution + 1
	if len(childResolution) > 0 {
		newResolution = childResolution[0]
	}
	return ChildrenAt(index, newResolution)
}

func CellToParent(index uint64, parentResolution ...int) (uint64, error) {
	cell := Deserialize(index)
	newResolution := cell.Resolution - 1
	if len(parentResolution) > 0 {
		newResolution = parentResolution[0]
	}
	return ParentAt(index, newResolution)
}

func GetRes0Cells() ([]uint64, error) {
	return Res0Cells()
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
