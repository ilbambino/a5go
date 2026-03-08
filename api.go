package a5go

import (
	"a5go/cells"
	"a5go/core"
	"a5go/traversal"
)

type CellBoundaryOptions = cells.CellBoundaryOptions
type Degrees = core.Degrees
type Radians = core.Radians
type Spherical = core.Spherical
type LonLat = core.LonLat

const (
	WorldCell              = core.WorldCell
	FirstHilbertResolution = core.FirstHilbertResolution
)

func CellToBoundary(cellID uint64, options ...CellBoundaryOptions) []LonLat {
	return cells.CellToBoundary(cellID, options...)
}

func CellToLonLat(cellID uint64) LonLat {
	return cells.CellToLonLat(cellID)
}

func CellToSpherical(cellID uint64) Spherical {
	return cells.CellToSpherical(cellID)
}

func LonLatToCell(lonLat LonLat, resolution int) uint64 {
	return cells.LonLatToCell(lonLat, resolution)
}

func HexToU64(hex string) (uint64, error) {
	return core.HexToU64(hex)
}

func U64ToHex(index uint64) string {
	return core.U64ToHex(index)
}

func CellToParent(index uint64, parentResolution ...int) uint64 {
	return core.CellToParent(index, parentResolution...)
}

func CellToChildren(index uint64, childResolution ...int) []uint64 {
	return core.CellToChildren(index, childResolution...)
}

func GetResolution(index uint64) int {
	return core.GetResolution(index)
}

func GetRes0Cells() []uint64 {
	return core.GetRes0Cells()
}

func GetNumCells(resolution int) float64 {
	return core.GetNumCells(resolution)
}

func GetNumChildren(parentResolution, childResolution int) float64 {
	return core.GetNumChildren(parentResolution, childResolution)
}

func CellArea(resolution int) float64 {
	return core.CellArea(resolution)
}

func Compact(cells []uint64) []uint64 {
	return core.Compact(cells)
}

func Uncompact(cells []uint64, targetResolution int) []uint64 {
	return core.Uncompact(cells, targetResolution)
}

func GridDisk(cellID uint64, k int) []uint64 {
	return traversal.GridDisk(cellID, k)
}

func GridDiskVertex(cellID uint64, k int) []uint64 {
	return traversal.GridDiskVertex(cellID, k)
}

func SphericalCap(cellID uint64, radius float64) []uint64 {
	return traversal.SphericalCap(cellID, radius)
}
