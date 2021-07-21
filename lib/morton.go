package lib

func Encode3(xyz [3]uint16) uint64 {
	x64 := uint64(xyz[0])
	y64 := uint64(xyz[1])
	z64 := uint64(xyz[2])
	return spread16(x64) | spread16(y64)<<1 | spread16(z64)<<2
}

func EncodeInt32(x int32, y int32, z int32) uint64 {
	return Encode3([3]uint16{uint16(x), uint16(y), uint16(z)})
}

func Decode3(m uint64) [3]uint16 {
	xyz := [3]uint16{
		compact16(m),
		compact16(m >> 1),
		compact16(m >> 2),
	}
	return xyz
}

func Path3(code uint64) [16]uint8 {
	path := [16]uint8{}
	for i := 0; i < 16; i++ {
		path[15-i] = uint8(code & 0b111)
		code = code >> 3
	}

	return path
}

func Box3(bb [6]float64, path []uint8) [][6]float64 {
	boxes := [][6]float64{bb}
	origin := bb[:3]
	width := []float64{
		(bb[3] - bb[0]),
		(bb[4] - bb[1]),
		(bb[5] - bb[2]),
	}

	for i, curr := range path {
		div := 1 << (i + 1)
		divF64 := float64(div)

		// Level width
		lw := [3]float64{
			width[0] / divF64,
			width[1] / divF64,
			width[2] / divF64,
		}

		// Offset origin by level bits
		if curr&4 > 0 {
			origin[2] += lw[2]
		}
		if curr&2 > 0 {
			origin[1] += lw[1]
		}
		if curr&1 > 0 {
			origin[0] += lw[0]
		}

		// The level tree node
		boxes = append(boxes, [6]float64{
			origin[0],
			origin[1],
			origin[2],
			origin[0] + lw[0],
			origin[1] + lw[1],
			origin[2] + lw[2],
		})
	}

	return boxes
}

func spread16(val uint64) uint64 {
	val &= 0x00000000001fffff
	val = (val | val<<32) & 0x001f00000000ffff
	val = (val | val<<16) & 0x001f0000ff0000ff
	val = (val | val<<8) & 0x010f00f00f00f00f
	val = (val | val<<4) & 0x10c30c30c30c30c3
	val = (val | val<<2) & 0x1249249249249249
	return val
}

func compact16(val uint64) uint16 {
	val &= 0x1249249249249249
	val = (val ^ (val >> 2)) & 0x30c30c30c30c30c3
	val = (val ^ (val >> 4)) & 0xf00f00f00f00f00f
	val = (val ^ (val >> 8)) & 0x00ff0000ff0000ff
	val = (val ^ (val >> 16)) & 0x00ff00000000ffff
	val = (val ^ (val >> 32)) & 0x00000000001fffff
	return uint16(val)
}
