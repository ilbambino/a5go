package projections

var geodeticToAuthalic = []float64{
	-2.2392098386786394e-03,
	2.1308606513250217e-06,
	-2.5592576864212742e-09,
	3.3701965267802837e-12,
	-4.6675453126112487e-15,
	6.6749287038481596e-18,
}

var authalicToGeodetic = []float64{
	2.2392089963541657e-03,
	2.8831978048607556e-06,
	5.0862207399726603e-09,
	1.0201812377816100e-11,
	2.1912872306767718e-14,
	4.9284235482523806e-17,
}

type AuthalicProjection struct{}

func (AuthalicProjection) applyCoefficients(phi float64, coefficients []float64) float64 {
	sinPhi := sin(phi)
	cosPhi := cos(phi)
	x := 2 * (cosPhi - sinPhi) * (cosPhi + sinPhi)

	u0 := x*coefficients[5] + coefficients[4]
	u1 := x*u0 + coefficients[3]
	u0 = x*u1 - u0 + coefficients[2]
	u1 = x*u0 - u1 + coefficients[1]
	u0 = x*u1 - u0 + coefficients[0]

	return phi + 2*sinPhi*cosPhi*u0
}

func (a AuthalicProjection) Forward(phi float64) float64 {
	return a.applyCoefficients(phi, geodeticToAuthalic)
}

func (a AuthalicProjection) Inverse(phi float64) float64 {
	return a.applyCoefficients(phi, authalicToGeodetic)
}
