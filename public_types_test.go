package a5go_test

import (
	"a5go"
	"testing"
)

func TestPublicCellAndPointAPI(t *testing.T) {
	point := a5go.Point{Lon: -3.7038, Lat: 40.4168}
	cell, err := point.Cell(6)
	if err != nil {
		t.Fatalf("point cell: %v", err)
	}
	if cell.Resolution() != 6 {
		t.Fatalf("resolution mismatch: got %d want 6", cell.Resolution())
	}

	parsed, err := a5go.Parse(cell.Hex())
	if err != nil {
		t.Fatalf("parse cell hex: %v", err)
	}
	if parsed != cell {
		t.Fatalf("parsed cell mismatch: got %v want %v", parsed, cell)
	}

	center := cell.Center()
	if center == (a5go.Point{}) {
		t.Fatalf("zero center returned for non-world cell")
	}

	boundary := cell.Boundary(a5go.CellBoundaryOptions{ClosedRing: true, Segments: 1})
	if len(boundary) == 0 {
		t.Fatal("expected non-empty boundary")
	}

	children, err := cell.Children()
	if err != nil {
		t.Fatalf("children: %v", err)
	}
	if len(children) == 0 {
		t.Fatal("expected children")
	}

	parent, err := cell.Parent()
	if err != nil {
		t.Fatalf("parent: %v", err)
	}
	if parent.Resolution() != 5 {
		t.Fatalf("parent resolution mismatch: got %d want 5", parent.Resolution())
	}
}
