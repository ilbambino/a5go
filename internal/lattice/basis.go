package lattice

func IJToKJ(ij IJ) KJ {
	return KJ{ij[0] + ij[1], ij[1]}
}

func KJToIJ(kj KJ) IJ {
	return IJ{kj[0] - kj[1], kj[1]}
}
