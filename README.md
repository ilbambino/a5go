# a5go

Go port of the `a5` spatial index from the TypeScript reference implementation in [`../a5`](../a5).

The port was done as a direct translation of the existing implementation and fixtures, not a redesign. The code follows the same geometric and indexing model, with Go package boundaries adjusted where needed to avoid import cycles.

## Status

- Full port of the current TypeScript implementation in this workspace
- Fixture-backed tests mirrored from the TypeScript repo
- Public API exposed from [`api.go`](./api.go)
- Verified with `go test ./...`

## Layout

- [`core`](./core): constants, coordinate transforms, serialization, hierarchy, compaction
- [`cells`](./cells): cell indexing, center conversion, boundary generation, point containment
- [`lattice`](./lattice): Hilbert and triangular lattice logic
- [`projections`](./projections): authalic, gnomonic, polyhedral, dodecahedron, CRS
- [`traversal`](./traversal): neighbors, grid disk, spherical cap
- [`geometry`](./geometry): planar and spherical polygon helpers
- [`internal`](./internal): wireframe/helper utilities mirrored from the TS repo

## Usage

```go
package main

import (
	"fmt"

	"a5go"
)

func main() {
	cell := a5go.LonLatToCell(a5go.LonLat{-3.7038, 40.4168}, 6)
	fmt.Println(a5go.U64ToHex(cell))
	fmt.Println(a5go.CellToLonLat(cell))
}
```

Run tests:

```bash
go test ./...
```

## Comparing Against The Official TypeScript Build

This repo includes a local comparison tool that runs the official TypeScript implementation from the sibling [`../a5`](../a5) repo and checks it against the Go port on canonical datasets.

### Prerequisite

The TypeScript repo must be built first:

```bash
cd ../a5
npm install
npm run build
```

The comparison tool loads `../a5/dist/a5.cjs`. It does not reimplement TS behavior in Go; it invokes the official build and compares outputs.

### Run

```bash
go run ./cmd/a5compare --ts-repo ../a5
```

Useful flags:

- `--points`: number of populated-place inputs to compare
- `--max-res`: highest resolution to compare for `lonLatToCell`

Example:

```bash
go run ./cmd/a5compare --ts-repo ../a5 --points 100 --max-res 10
```

The command compares:

- `lonLatToCell`
- `cellToLonLat`
- `cellToSpherical`
- `cellToBoundary` with `segments=1`
- `cellToBoundary` with `segments=auto`
- `sphericalCap` in compacted and uncompacted form

It exits non-zero on the first set of mismatches and prints the failing cases.
