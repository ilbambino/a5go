package core

import "fmt"

func ChildrenAt(index uint64, childResolution int) ([]uint64, error) {
	cell := Deserialize(index)
	currentResolution := cell.Resolution

	if childResolution < currentResolution {
		return nil, fmt.Errorf("target resolution %d must be >= current resolution %d", childResolution, currentResolution)
	}
	if childResolution > MaxResolution {
		return nil, fmt.Errorf("target resolution %d exceeds maximum resolution %d", childResolution, MaxResolution)
	}
	if childResolution == currentResolution {
		return []uint64{index}, nil
	}

	newOrigins := []*Origin{cell.Origin}
	newSegments := []int{cell.Segment}
	if currentResolution == -1 {
		newOrigins = make([]*Origin, len(Origins))
		for i := range Origins {
			newOrigins[i] = &Origins[i]
		}
	}
	if (currentResolution == -1 && childResolution > 0) || currentResolution == 0 {
		newSegments = []int{0, 1, 2, 3, 4}
	}

	startResolution := currentResolution
	if startResolution < FirstHilbertResolution-1 {
		startResolution = FirstHilbertResolution - 1
	}
	resolutionDiff := childResolution - startResolution
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
				child, err := Serialize(A5Cell{Origin: newOrigin, Segment: newSegment, S: newS, Resolution: childResolution})
				if err != nil {
					return nil, err
				}
				children = append(children, child)
			}
		}
	}
	return children, nil
}

func ParentAt(index uint64, parentResolution int) (uint64, error) {
	cell := Deserialize(index)
	currentResolution := cell.Resolution

	if parentResolution == -1 {
		return WorldCell, nil
	}
	if parentResolution < 0 {
		return 0, fmt.Errorf("target resolution %d cannot be negative", parentResolution)
	}
	if parentResolution > currentResolution {
		return 0, fmt.Errorf("target resolution %d must be <= current resolution %d", parentResolution, currentResolution)
	}
	if parentResolution == currentResolution {
		return index, nil
	}

	resolutionDiff := currentResolution - parentResolution
	shiftedS := cell.S >> uint(2*resolutionDiff)
	return Serialize(A5Cell{Origin: cell.Origin, Segment: cell.Segment, S: shiftedS, Resolution: parentResolution})
}

func Resolution(index uint64) int {
	return GetResolution(index)
}

func Res0Cells() ([]uint64, error) {
	return ChildrenAt(WorldCell, 0)
}
