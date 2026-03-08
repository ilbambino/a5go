package projections

import "a5go/core"

type GnomonicProjection struct{}

func (GnomonicProjection) Forward(spherical core.Spherical) core.Polar {
	return core.Polar{tan(spherical[1]), spherical[0]}
}

func (GnomonicProjection) Inverse(polar core.Polar) core.Spherical {
	return core.Spherical{polar[1], atan(polar[0])}
}
