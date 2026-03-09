package core

import (
	"a5go/internal/testutil"
	"fmt"
	"testing"
)

var resolutionMasks = []string{
	"0000001000000000000000000000000000000000000000000000000000000000",
	"0000000100000000000000000000000000000000000000000000000000000000",
	"0000000010000000000000000000000000000000000000000000000000000000",
	"0000000000100000000000000000000000000000000000000000000000000000",
	"0000000000001000000000000000000000000000000000000000000000000000",
	"0000000000000010000000000000000000000000000000000000000000000000",
	"0000000000000000100000000000000000000000000000000000000000000000",
	"0000000000000000001000000000000000000000000000000000000000000000",
	"0000000000000000000010000000000000000000000000000000000000000000",
	"0000000000000000000000100000000000000000000000000000000000000000",
	"0000000000000000000000001000000000000000000000000000000000000000",
	"0000000000000000000000000010000000000000000000000000000000000000",
	"0000000000000000000000000000100000000000000000000000000000000000",
	"0000000000000000000000000000001000000000000000000000000000000000",
	"0000000000000000000000000000000010000000000000000000000000000000",
	"0000000000000000000000000000000000100000000000000000000000000000",
	"0000000000000000000000000000000000001000000000000000000000000000",
	"0000000000000000000000000000000000000010000000000000000000000000",
	"0000000000000000000000000000000000000000100000000000000000000000",
	"0000000000000000000000000000000000000000001000000000000000000000",
	"0000000000000000000000000000000000000000000010000000000000000000",
	"0000000000000000000000000000000000000000000000100000000000000000",
	"0000000000000000000000000000000000000000000000001000000000000000",
	"0000000000000000000000000000000000000000000000000010000000000000",
	"0000000000000000000000000000000000000000000000000000100000000000",
	"0000000000000000000000000000000000000000000000000000001000000000",
	"0000000000000000000000000000000000000000000000000000000010000000",
	"0000000000000000000000000000000000000000000000000000000000100000",
	"0000000000000000000000000000000000000000000000000000000000001000",
	"0000000000000000000000000000000000000000000000000000000000000010",
}

func TestSerialization(t *testing.T) {
	if len(resolutionMasks) != MaxResolution {
		t.Fatalf("mask count mismatch")
	}
	expectedRemovalMask := uint64(0x03ffffffffffffff)
	if RemovalMask != expectedRemovalMask {
		t.Fatalf("removal mask mismatch")
	}

	origin0 := Origins[0]
	for i := range resolutionMasks {
		input := A5Cell{Origin: &origin0, Segment: 4, S: 0, Resolution: i}
		serialized, err := Serialize(input)
		if err != nil {
			t.Fatalf("serialize input: %v", err)
		}
		if fmt.Sprintf("%064b", serialized) != resolutionMasks[i] {
			t.Fatalf("resolution mask mismatch at %d", i)
		}
	}

	for i, binary := range resolutionMasks {
		var value uint64
		_, err := fmt.Sscanf(binary, "%b", &value)
		if err == nil {
			_ = value
		}
		if GetResolution(parseBinary(binary)) != i {
			t.Fatalf("resolution extraction mismatch at %d", i)
		}
	}

	serializedMax, err := Serialize(A5Cell{Origin: &origin0, Segment: 4, S: 0, Resolution: MaxResolution - 1})
	if err != nil {
		t.Fatalf("serialize max: %v", err)
	}
	if serializedMax != 0b10 {
		t.Fatalf("origin segment encoding mismatch")
	}

	if _, err := Serialize(A5Cell{Origin: &origin0, Segment: 0, S: 16, Resolution: 3}); err == nil {
		t.Fatalf("expected serialize error for oversized S")
	}
	if _, err := Serialize(A5Cell{Origin: &origin0, Segment: 0, S: 0, Resolution: 31}); err == nil {
		t.Fatalf("expected serialize error for oversized resolution")
	}

	var ids []string
	testutil.LoadJSON(t, "../../testdata/test-ids.json", &ids)
	for _, id := range ids {
		serialized := parseHex(id)
		deserialized := Deserialize(serialized)
		reserialized, err := Serialize(deserialized)
		if err != nil {
			t.Fatalf("reserialize: %v", err)
		}
		if reserialized != serialized {
			t.Fatalf("round trip mismatch for %s", id)
		}
	}

	for _, id := range ids {
		cell := parseHex(id)
		children, err := CellToChildren(cell)
		if err != nil {
			t.Fatalf("children: %v", err)
		}
		child := children[0]
		parent, err := CellToParent(child)
		if err != nil {
			t.Fatalf("parent: %v", err)
		}
		if parent != cell {
			t.Fatalf("parent/child round trip mismatch")
		}
		children, err = CellToChildren(cell)
		if err != nil {
			t.Fatalf("children: %v", err)
		}
		for _, c := range children {
			parent, err := CellToParent(c)
			if err != nil {
				t.Fatalf("parent: %v", err)
			}
			if parent != cell {
				t.Fatalf("child parent mismatch")
			}
		}
	}

	for _, id := range ids {
		cell := parseHex(id)
		currentResolution := GetResolution(cell)
		children, err := CellToChildren(cell, currentResolution)
		if err != nil {
			t.Fatalf("children at same resolution: %v", err)
		}
		if len(children) != 1 || children[0] != cell {
			t.Fatalf("same-resolution children mismatch")
		}
		parent, err := CellToParent(cell, currentResolution)
		if err != nil {
			t.Fatalf("parent at same resolution: %v", err)
		}
		if parent != cell {
			t.Fatalf("same-resolution parent mismatch")
		}
	}

	cell, err := Serialize(A5Cell{Origin: &origin0, Segment: 0, S: 0, Resolution: 0})
	if err != nil {
		t.Fatalf("serialize res0: %v", err)
	}
	children, err := CellToChildren(cell)
	if err != nil {
		t.Fatalf("children res0: %v", err)
	}
	if len(children) != 5 {
		t.Fatalf("res0 children count mismatch")
	}
	for _, child := range children {
		parent, err := CellToParent(child)
		if err != nil {
			t.Fatalf("parent res0 child: %v", err)
		}
		if parent != cell {
			t.Fatalf("non-hilbert parent mismatch")
		}
	}

	cell, err = Serialize(A5Cell{Origin: &origin0, Segment: 0, S: 0, Resolution: 1})
	if err != nil {
		t.Fatalf("serialize res1: %v", err)
	}
	children, err = CellToChildren(cell)
	if err != nil {
		t.Fatalf("children res1: %v", err)
	}
	if len(children) != 4 {
		t.Fatalf("res1 children count mismatch")
	}
	for _, child := range children {
		parent, err := CellToParent(child)
		if err != nil {
			t.Fatalf("parent res1 child: %v", err)
		}
		if parent != cell {
			t.Fatalf("non-hilbert to hilbert parent mismatch")
		}
	}

	cell, err = Serialize(A5Cell{Origin: &origin0, Segment: 0, S: 0, Resolution: 2})
	if err != nil {
		t.Fatalf("serialize res2: %v", err)
	}
	parent, err := CellToParent(cell, 1)
	if err != nil {
		t.Fatalf("parent res2: %v", err)
	}
	children, err = CellToChildren(parent)
	if err != nil {
		t.Fatalf("children parent res2: %v", err)
	}
	if len(children) != 4 {
		t.Fatalf("hilbert to non-hilbert child count mismatch")
	}
	found := false
	for _, child := range children {
		if child == cell {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected original child in parent children")
	}

	resolutions := []int{0, 1, 2, 3, 4}
	cells := make([]uint64, len(resolutions))
	for i, res := range resolutions {
		cells[i], err = Serialize(A5Cell{Origin: &origin0, Segment: 0, S: 0, Resolution: res})
		if err != nil {
			t.Fatalf("serialize resolution chain: %v", err)
		}
	}
	for i := 1; i < len(cells); i++ {
		parent, err := CellToParent(cells[i])
		if err != nil {
			t.Fatalf("chain parent: %v", err)
		}
		if parent != cells[i-1] {
			t.Fatalf("resolution chain parent mismatch")
		}
	}
	for i := 0; i < len(cells)-1; i++ {
		children, err := CellToChildren(cells[i])
		if err != nil {
			t.Fatalf("chain children: %v", err)
		}
		found := false
		for _, child := range children {
			if child == cells[i+1] {
				found = true
			}
		}
		if !found {
			t.Fatalf("resolution chain child mismatch")
		}
	}

	baseCell, err := Serialize(A5Cell{Origin: &origin0, Segment: 0, S: 0, Resolution: -1})
	if err != nil {
		t.Fatalf("serialize base cell: %v", err)
	}
	currentCells := []uint64{baseCell}
	expectedCounts := []int{12, 60, 240, 960}
	for resolution := 0; resolution < 4; resolution++ {
		allChildren := make([]uint64, 0)
		for _, cell := range currentCells {
			children, err := CellToChildren(cell)
			if err != nil {
				t.Fatalf("tree children: %v", err)
			}
			allChildren = append(allChildren, children...)
		}
		if len(allChildren) != expectedCounts[resolution] {
			t.Fatalf("expected %d cells at resolution %d, got %d", expectedCounts[resolution], resolution, len(allChildren))
		}
		currentCells = allChildren
	}

	res0Cells, err := GetRes0Cells()
	if err != nil {
		t.Fatalf("res0 cells: %v", err)
	}
	if len(res0Cells) != 12 {
		t.Fatalf("res0 cell count mismatch")
	}
	expectedHexValues := []string{"200000000000000", "600000000000000", "a00000000000000", "e00000000000000", "1200000000000000", "1600000000000000", "1a00000000000000", "1e00000000000000", "2200000000000000", "2600000000000000", "2a00000000000000", "2e00000000000000"}
	for i, cell := range res0Cells {
		if GetResolution(cell) != 0 {
			t.Fatalf("expected resolution 0")
		}
		if U64ToHex(cell) != expectedHexValues[i] {
			t.Fatalf("unexpected res0 hex at %d: %s", i, U64ToHex(cell))
		}
	}
}

func parseBinary(binary string) uint64 {
	var result uint64
	for _, c := range binary {
		result <<= 1
		if c == '1' {
			result |= 1
		}
	}
	return result
}

func parseHex(hex string) uint64 {
	value, err := HexToU64(hex)
	if err != nil {
		panic(err)
	}
	return value
}
