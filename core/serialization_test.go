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

	origin0 := *Origins[0]
	for i := range resolutionMasks {
		input := A5Cell{Origin: &origin0, Segment: 4, S: 0, Resolution: i}
		serialized := Serialize(input)
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

	if Serialize(A5Cell{Origin: &origin0, Segment: 4, S: 0, Resolution: MaxResolution - 1}) != 0b10 {
		t.Fatalf("origin segment encoding mismatch")
	}

	mustPanicWith(t, "S (16) is too large for resolution level 3", func() {
		Serialize(A5Cell{Origin: &origin0, Segment: 0, S: 16, Resolution: 3})
	})
	mustPanicWith(t, "Resolution (31) is too large", func() {
		Serialize(A5Cell{Origin: &origin0, Segment: 0, S: 0, Resolution: 31})
	})

	var ids []string
	testutil.LoadJSON(t, "../testdata/test-ids.json", &ids)
	for _, id := range ids {
		serialized := parseHex(id)
		deserialized := Deserialize(serialized)
		reserialized := Serialize(deserialized)
		if reserialized != serialized {
			t.Fatalf("round trip mismatch for %s", id)
		}
	}

	for _, id := range ids {
		cell := parseHex(id)
		child := CellToChildren(cell)[0]
		parent := CellToParent(child)
		if parent != cell {
			t.Fatalf("parent/child round trip mismatch")
		}
		children := CellToChildren(cell)
		for _, c := range children {
			if CellToParent(c) != cell {
				t.Fatalf("child parent mismatch")
			}
		}
	}

	for _, id := range ids {
		cell := parseHex(id)
		currentResolution := GetResolution(cell)
		children := CellToChildren(cell, currentResolution)
		if len(children) != 1 || children[0] != cell {
			t.Fatalf("same-resolution children mismatch")
		}
		if CellToParent(cell, currentResolution) != cell {
			t.Fatalf("same-resolution parent mismatch")
		}
	}

	cell := Serialize(A5Cell{Origin: &origin0, Segment: 0, S: 0, Resolution: 0})
	children := CellToChildren(cell)
	if len(children) != 5 {
		t.Fatalf("res0 children count mismatch")
	}
	for _, child := range children {
		if CellToParent(child) != cell {
			t.Fatalf("non-hilbert parent mismatch")
		}
	}

	cell = Serialize(A5Cell{Origin: &origin0, Segment: 0, S: 0, Resolution: 1})
	children = CellToChildren(cell)
	if len(children) != 4 {
		t.Fatalf("res1 children count mismatch")
	}
	for _, child := range children {
		if CellToParent(child) != cell {
			t.Fatalf("non-hilbert to hilbert parent mismatch")
		}
	}

	cell = Serialize(A5Cell{Origin: &origin0, Segment: 0, S: 0, Resolution: 2})
	parent := CellToParent(cell, 1)
	children = CellToChildren(parent)
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
		cells[i] = Serialize(A5Cell{Origin: &origin0, Segment: 0, S: 0, Resolution: res})
	}
	for i := 1; i < len(cells); i++ {
		if CellToParent(cells[i]) != cells[i-1] {
			t.Fatalf("resolution chain parent mismatch")
		}
	}
	for i := 0; i < len(cells)-1; i++ {
		children := CellToChildren(cells[i])
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

	baseCell := Serialize(A5Cell{Origin: &origin0, Segment: 0, S: 0, Resolution: -1})
	currentCells := []uint64{baseCell}
	expectedCounts := []int{12, 60, 240, 960}
	for resolution := 0; resolution < 4; resolution++ {
		allChildren := make([]uint64, 0)
		for _, cell := range currentCells {
			allChildren = append(allChildren, CellToChildren(cell)...)
		}
		if len(allChildren) != expectedCounts[resolution] {
			t.Fatalf("expected %d cells at resolution %d, got %d", expectedCounts[resolution], resolution, len(allChildren))
		}
		currentCells = allChildren
	}

	res0Cells := GetRes0Cells()
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

func mustPanicWith(t *testing.T, message string, fn func()) {
	t.Helper()
	defer func() {
		value := recover()
		if value == nil {
			t.Fatalf("expected panic %q", message)
		}
		if got, ok := value.(string); !ok || got != message {
			t.Fatalf("panic = %v want %q", value, message)
		}
	}()
	fn()
}
