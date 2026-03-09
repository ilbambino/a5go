package a5go_test

import (
	"a5go"
	"testing"

	a5ext "github.com/a5geo/a5-go"
)

var (
	benchCellSink       uint64
	benchCellsSink      []uint64
	benchLonLatSink     a5go.LonLat
	benchLonLatExtSink  a5ext.LonLat
	benchResolutionSink int
)

func BenchmarkCompareImplementations(b *testing.B) {
	rootCells, err := a5go.GetRes0Cells()
	if err != nil {
		b.Fatalf("a5go GetRes0Cells: %v", err)
	}
	extRootCells, err := a5ext.GetRes0Cells()
	if err != nil {
		b.Fatalf("a5go-ext GetRes0Cells: %v", err)
	}
	if len(extRootCells) == 0 {
		b.Fatal("a5go-ext returned no res0 cells")
	}

	points := []struct {
		goPoint  a5go.LonLat
		extPoint a5ext.LonLat
		res      int
	}{
		{goPoint: a5go.LonLat{-3.7038, 40.4168}, extPoint: a5ext.LonLat{-3.7038, 40.4168}, res: 6},
		{goPoint: a5go.LonLat{-73.9857, 40.7484}, extPoint: a5ext.LonLat{-73.9857, 40.7484}, res: 10},
		{goPoint: a5go.LonLat{151.2093, -33.8688}, extPoint: a5ext.LonLat{151.2093, -33.8688}, res: 12},
		{goPoint: a5go.LonLat{139.6917, 35.6895}, extPoint: a5ext.LonLat{139.6917, 35.6895}, res: 15},
	}

	childSource := rootCells[3]
	childTargetRes := 8
	childrenGo, err := a5go.CellToChildren(childSource, childTargetRes)
	if err != nil {
		b.Fatalf("a5go CellToChildren setup: %v", err)
	}
	childrenExt, err := a5ext.CellToChildren(childSource, childTargetRes)
	if err != nil {
		b.Fatalf("a5go-ext CellToChildren setup: %v", err)
	}
	if len(childrenExt) == 0 {
		b.Fatal("a5go-ext produced no children")
	}

	parentSource := childrenGo[len(childrenGo)/2]
	parentTargetRes := 3
	compactSource := childrenGo[:256]

	b.Run("GetResolution", func(b *testing.B) {
		cell := parentSource
		b.Run("a5go", func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				benchResolutionSink = a5go.GetResolution(cell)
			}
		})
		b.Run("a5go-ext", func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				benchResolutionSink = a5ext.GetResolution(cell)
			}
		})
	})

	b.Run("LonLatToCell", func(b *testing.B) {
		b.Run("a5go", func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				sample := points[i%len(points)]
				cell, err := a5go.LonLatToCell(sample.goPoint, sample.res)
				if err != nil {
					b.Fatalf("a5go LonLatToCell: %v", err)
				}
				benchCellSink = cell
			}
		})
		b.Run("a5go-ext", func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				sample := points[i%len(points)]
				cell, err := a5ext.LonLatToCell(sample.extPoint, sample.res)
				if err != nil {
					b.Fatalf("a5go-ext LonLatToCell: %v", err)
				}
				benchCellSink = cell
			}
		})
	})

	b.Run("CellToLonLat", func(b *testing.B) {
		cell := parentSource
		b.Run("a5go", func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				benchLonLatSink = a5go.CellToLonLat(cell)
			}
		})
		b.Run("a5go-ext", func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				benchLonLatExtSink = a5ext.CellToLonLat(cell)
			}
		})
	})

	b.Run("CellToChildren", func(b *testing.B) {
		cell := childSource
		resolution := childTargetRes
		b.Run("a5go", func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				cells, err := a5go.CellToChildren(cell, resolution)
				if err != nil {
					b.Fatalf("a5go CellToChildren: %v", err)
				}
				benchCellsSink = cells
			}
		})
		b.Run("a5go-ext", func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				cells, err := a5ext.CellToChildren(cell, resolution)
				if err != nil {
					b.Fatalf("a5go-ext CellToChildren: %v", err)
				}
				benchCellsSink = cells
			}
		})
	})

	b.Run("CellToParent", func(b *testing.B) {
		cell := parentSource
		resolution := parentTargetRes
		b.Run("a5go", func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				parent, err := a5go.CellToParent(cell, resolution)
				if err != nil {
					b.Fatalf("a5go CellToParent: %v", err)
				}
				benchCellSink = parent
			}
		})
		b.Run("a5go-ext", func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				parent, err := a5ext.CellToParent(cell, resolution)
				if err != nil {
					b.Fatalf("a5go-ext CellToParent: %v", err)
				}
				benchCellSink = parent
			}
		})
	})

	b.Run("Compact", func(b *testing.B) {
		source := append([]uint64(nil), compactSource...)
		b.Run("a5go", func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				benchCellsSink = a5go.Compact(source)
			}
		})
		b.Run("a5go-ext", func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				benchCellsSink = a5ext.Compact(source)
			}
		})
	})

	b.Run("Uncompact", func(b *testing.B) {
		compactedGo := a5go.Compact(compactSource)
		compactedExt := a5ext.Compact(compactSource)
		targetResolution := childTargetRes

		b.Run("a5go", func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				cells, err := a5go.Uncompact(compactedGo, targetResolution)
				if err != nil {
					b.Fatalf("a5go Uncompact: %v", err)
				}
				benchCellsSink = cells
			}
		})
		b.Run("a5go-ext", func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				cells, err := a5ext.Uncompact(compactedExt, targetResolution)
				if err != nil {
					b.Fatalf("a5go-ext Uncompact: %v", err)
				}
				benchCellsSink = cells
			}
		})
	})
}
