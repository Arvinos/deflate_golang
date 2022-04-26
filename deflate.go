package go_deflate

import "math"

type Compressor interface {
	Compress(src []byte, dst []byte) (in uint32, out uint32, res uint32)
	Decompress(src []byte, dst []byte) (in uint32, out uint32, res uint32)
}

type Deflate struct {
}

func writeStoredBlock(src []byte, dst []byte) (in uint32, out uint32, res uint32) {
	if len(src) > math.MaxInt16 {
		return 0, 0, 1
	}

	if len(dst) < len(src)+4 {
		return 0, 0, 1
	}

	lenV := uint16(len(src))
	nlenV := ^lenV

	dst[0] = byte(lenV >> 8)
	dst[1] = byte(lenV & 0xFF)

	dst[2] = byte(nlenV >> 8)
	dst[3] = byte(nlenV & 0xFF)

	outBlock := dst[4:]

	copy(outBlock, src)

	return uint32(lenV), uint32(4 + lenV), 0
}

func readStoredBlock(src []byte, dst []byte) (in uint32, out uint32, res uint32) {
	if 4 > len(src) {
		return 0, 0, 1
	}

	lenV := (uint16(src[0]) << 8) | uint16(src[1])
	nlenV := (uint16(src[2]) << 8) | uint16(src[3])

	if lenV != ^nlenV {
		return 0, 0, 1
	}

	if int(lenV)+4 > len(src) {
		return 0, 0, 1
	}

	if int(lenV) > len(dst) {
		return 0, 0, 1
	}

	inByteBlock := src[4:]
	outByteBlock := dst[:lenV]

	copy(outByteBlock, inByteBlock)

	return uint32(4 + lenV), uint32(lenV), 0
}

func (state *Deflate) Compress(src []byte, dst []byte) (in uint32, out uint32, res uint32) {
	in = 0
	out = 0
	res = 1

	if res == 1 {
		in, out, res = writeStoredBlock(src, dst)
	}

	return in, out, res
}

func (state *Deflate) Decompress(src []byte, dst []byte) (in uint32, out uint32, res uint32) {
	in = 0
	out = 0
	res = 1

	if res == 1 {
		in, out, res = readStoredBlock(src, dst)
	}

	return in, out, res
}
