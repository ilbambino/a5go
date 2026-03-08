package internal

import (
	"a5go/internal/core"
	"fmt"
)

func GeneratePentagonSVG() string {
	const (
		width  = 64
		height = 64
		cx     = 7.0
		cy     = 7.0
		scale  = 32.0
	)

	vertices := core.PentagonShapeDef.GetVertices()
	points := ""
	for i, vertex := range vertices {
		x := vertex[0]*scale + cx
		y := vertex[1]*scale + cy
		if i > 0 {
			points += " "
		}
		points += fmt.Sprintf("%g,%g", x, y)
	}

	return fmt.Sprintf(
		"<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"%d\" height=\"%d\" viewBox=\"0 0 %d %d\">\n  <polygon points=\"%s\" fill=\"none\" stroke=\"black\" stroke-width=\"2\"/>\n</svg>",
		width, height, width, height, points,
	)
}
