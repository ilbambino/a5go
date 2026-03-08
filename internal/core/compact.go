package core

import (
	"sort"
)

func Uncompact(cells []uint64, targetResolution int) []uint64 {
	n := 0
	resolutions := make([]int, len(cells))
	for i, cell := range cells {
		resolution := GetResolution(cell)
		resolutionDiff := targetResolution - resolution
		if resolutionDiff < 0 {
			panicString("Cannot uncompact cell at resolution " + itoa(resolution) + " to lower resolution " + itoa(targetResolution))
		}
		resolutions[i] = resolution
		n += int(GetNumChildren(resolution, targetResolution))
	}

	result := make([]uint64, n)
	offset := 0
	for i, cell := range cells {
		resolution := resolutions[i]
		numChildren := int(GetNumChildren(resolution, targetResolution))
		if numChildren == 1 {
			result[offset] = cell
		} else {
			children := CellToChildren(cell, targetResolution)
			copy(result[offset:], children)
		}
		offset += numChildren
	}
	return result
}

func Compact(cells []uint64) []uint64 {
	if len(cells) == 0 {
		return []uint64{}
	}

	seen := make(map[uint64]struct{}, len(cells))
	currentCells := make([]uint64, 0, len(cells))
	for _, cell := range cells {
		if _, ok := seen[cell]; ok {
			continue
		}
		seen[cell] = struct{}{}
		currentCells = append(currentCells, cell)
	}
	sort.Slice(currentCells, func(i, j int) bool {
		return currentCells[i] < currentCells[j]
	})

	changed := true
	for changed {
		changed = false
		result := make([]uint64, 0, len(currentCells))
		i := 0
		for i < len(currentCells) {
			cell := currentCells[i]
			resolution := GetResolution(cell)
			if resolution < 0 {
				result = append(result, cell)
				i++
				continue
			}

			expectedChildren := 4
			if resolution < FirstHilbertResolution {
				if resolution == 0 {
					expectedChildren = 12
				} else {
					expectedChildren = 5
				}
			}

			if i+expectedChildren <= len(currentCells) {
				hasAllSiblings := true
				if IsFirstChild(cell, resolution) {
					stride := GetStride(resolution)
					for j := 1; j < expectedChildren; j++ {
						expectedCell := cell + uint64(j)*stride
						if currentCells[i+j] != expectedCell {
							hasAllSiblings = false
							break
						}
					}
				} else {
					hasAllSiblings = false
				}

				if hasAllSiblings {
					parent := CellToParent(cell)
					result = append(result, parent)
					i += expectedChildren
					changed = true
					continue
				}
			}

			result = append(result, cell)
			i++
		}
		currentCells = result
	}

	return currentCells
}
