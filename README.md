# a5go

Go port of the [`a5`](https://github.com/felixpalmer/a5) spatial index from the upstream TypeScript reference implementation.

The port was done as a direct translation of the existing implementation and fixtures, then refactored toward a more idiomatic Go API. The root `a5go` package is the intended public surface; implementation packages now live under `internal/`.

## Status

- Full port of the current TypeScript implementation in this workspace
- Fixture-backed tests mirrored from the TypeScript repo
- Public API exposed from [`api.go`](./api.go)
- Verified with `go test ./...`

## Layout

- [`api.go`](./api.go), [`public_types.go`](./public_types.go): public API, including `Cell` and `Point`
- [`internal/core`](./internal/core): constants, coordinate transforms, serialization, hierarchy, compaction
- [`internal/cells`](./internal/cells): cell indexing, center conversion, boundary generation, point containment
- [`internal/lattice`](./internal/lattice): Hilbert and triangular lattice logic
- [`internal/projections`](./internal/projections): authalic, gnomonic, polyhedral, dodecahedron, CRS
- [`internal/traversal`](./internal/traversal): neighbors, grid disk, spherical cap
- [`internal/geometry`](./internal/geometry): planar and spherical polygon helpers
- [`internal`](./internal): wireframe/helper utilities mirrored from the TS repo

## Usage

```go
package main

import (
	"fmt"

	"a5go"
)

func main() {
	point := a5go.Point{Lon: -3.7038, Lat: 40.4168}
	cell, err := point.Cell(6)
	if err != nil {
		panic(err)
	}
	fmt.Println(cell.Hex())
	fmt.Println(cell.Center())
}
```

Run tests:

```bash
go test ./...
```

Run the local cross-implementation benchmark against `a5go-ext`:

```bash
go test -run ^$ -bench BenchmarkCompareImplementations -benchmem
```

## Comparing Against The Official TypeScript Build

This repo includes a local comparison tool that runs the official TypeScript implementation from a local checkout of the upstream repo, [`felixpalmer/a5`](https://github.com/felixpalmer/a5), and checks it against the Go port on canonical datasets.

### Prerequisite

First, clone and build the upstream TypeScript repo somewhere on your machine:

```bash
git clone https://github.com/felixpalmer/a5.git
cd a5
npm install
npm run build
```

The comparison tool loads `<your-a5-checkout>/dist/a5.cjs`. It does not reimplement TS behavior in Go; it invokes the official build and compares outputs.

### Run

```bash
go run ./cmd/a5compare --ts-repo /path/to/a5
```

Useful flags:

- `--ts-repo`: path to a local checkout of [`felixpalmer/a5`](https://github.com/felixpalmer/a5)
- `--points`: number of populated-place inputs to compare
- `--max-res`: highest resolution to compare for `lonLatToCell`

You can also set `A5_TS_REPO=/path/to/a5` instead of passing `--ts-repo`.

Example:

```bash
A5_TS_REPO=/path/to/a5 go run ./cmd/a5compare --points 100 --max-res 10
```

The command compares:

- `lonLatToCell`
- `cellToLonLat`
- `cellToSpherical`
- `cellToBoundary` with `segments=1`
- `cellToBoundary` with `segments=auto`
- `sphericalCap` in compacted and uncompacted form

It exits non-zero on the first set of mismatches and prints the failing cases.

## API Notes

The public indexing and hierarchy entry points return errors for invalid input instead of panicking. In particular:

- `Point.Cell`
- `LonLatToCell`
- `CellToParent`
- `CellToChildren`
- `GetRes0Cells`
- `Res0Cells`
- `Uncompact`
