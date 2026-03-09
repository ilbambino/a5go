package a5go

import "a5go/internal/core"

type Cell uint64

type Point struct {
	Lon float64
	Lat float64
}

func World() Cell {
	return Cell(WorldCell)
}

func ParseCell(hex string) (Cell, error) {
	value, err := HexToU64(hex)
	if err != nil {
		return 0, err
	}
	return Cell(value), nil
}

func (c Cell) Uint64() uint64 {
	return uint64(c)
}

func (c Cell) Hex() string {
	return U64ToHex(uint64(c))
}

func (c Cell) Resolution() int {
	return Resolution(uint64(c))
}

func (c Cell) Parent() (Cell, error) {
	return c.ParentAt(c.Resolution() - 1)
}

func (c Cell) ParentAt(resolution int) (Cell, error) {
	parent, err := ParentAt(uint64(c), resolution)
	return Cell(parent), err
}

func (c Cell) Children() ([]Cell, error) {
	return c.ChildrenAt(c.Resolution() + 1)
}

func (c Cell) ChildrenAt(resolution int) ([]Cell, error) {
	children, err := ChildrenAt(uint64(c), resolution)
	if err != nil {
		return nil, err
	}
	return toCells(children), nil
}

func (c Cell) Boundary(opts CellBoundaryOptions) []Point {
	return toPoints(Boundary(uint64(c), opts))
}

func (c Cell) Center() Point {
	return fromCoreLonLat(CellToLonLat(uint64(c)))
}

func (c Cell) GridDisk(k int) []Cell {
	return toCells(GridDisk(uint64(c), k))
}

func (c Cell) GridDiskVertices(k int) []Cell {
	return toCells(GridDiskVertex(uint64(c), k))
}

func (c Cell) SphericalCap(radius float64) []Cell {
	return toCells(SphericalCap(uint64(c), radius))
}

func (p Point) LonLat() LonLat {
	return core.LonLat{p.Lon, p.Lat}
}

func (p Point) Cell(resolution int) (Cell, error) {
	cell, err := LonLatToCell(p.LonLat(), resolution)
	if err != nil {
		return 0, err
	}
	return Cell(cell), nil
}

func fromCoreLonLat(value core.LonLat) Point {
	return Point{Lon: value[0], Lat: value[1]}
}

func toPoints(values []core.LonLat) []Point {
	points := make([]Point, len(values))
	for i, value := range values {
		points[i] = fromCoreLonLat(value)
	}
	return points
}

func toCells(values []uint64) []Cell {
	cells := make([]Cell, len(values))
	for i, value := range values {
		cells[i] = Cell(value)
	}
	return cells
}
