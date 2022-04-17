package go_deflate

import "math"

type Compressor interface {
	Compress(src []byte, dst []byte) (in uint32, out uint32, res uint32)
	Decompress(src []byte, dst []byte) (in uint32, out uint32, res uint32)
}

type Deflate struct {
	storedBlock bool
	blockFinal  bool
	eosFound    bool
	output      []byte
	input       []byte
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

func (state *Deflate) writeStoredDeflateHeader(isFinalBlock bool) (res uint32) {
	if isFinalBlock {
		state.output[0] = 0b10000000
	} else {
		state.output[0] = 0
	}

	state.output = state.output[1:]

	return 0
}

func (state *Deflate) readHeader() (res uint32) {
	typeHeader := state.input[0] >> 5

	if (typeHeader & 0b100) != 0 {
		state.blockFinal = true
	}

	typeHeader &= 0b011

	switch typeHeader {
	case 0:
		state.storedBlock = true
		state.input = state.input[1:]
	case 1:
		fallthrough
	case 2:
		fallthrough
	case 3:
		return 1
	}

	return 0
}

func (state *Deflate) Compress(src []byte, dst []byte) (in uint32, out uint32, res uint32) {
	in = 0
	out = 0
	res = 1

	state.input = src[:]
	state.output = dst[:]

	if res == 1 {
		result := state.writeStoredDeflateHeader(true)

		if result != 0 {
			return in, uint32(len(dst) - len(state.output)), result
		}

		bytesIn, bytesOut, result := writeStoredBlock(state.input, state.output)

		state.input = state.input[bytesIn:]
		state.output = state.output[bytesOut:]

		res = result
	}

	return uint32(len(src) - len(state.input)), uint32(len(dst) - len(state.output)), res
}

func (state *Deflate) Decompress(src []byte, dst []byte) (in uint32, out uint32, res uint32) {
	in = 0
	out = 0
	res = 1

	state.input = src[:]
	state.output = dst[:]

	for int(in) < len(src) && !state.eosFound {
		res = state.readHeader()

		if res != 0 {
			return in, out, res
		}

		if state.storedBlock {
			bytesIn, bytesOut, result := readStoredBlock(state.input, state.output)
			state.eosFound = state.blockFinal

			if result != 0 {
				return in, out, result
			}

			state.input = state.input[bytesIn:]
			state.output = state.output[bytesOut:]

		} else {
			// to implement
		}

		if res != 0 {
			return in, out, res
		}
	}

	return uint32(len(src) - len(state.input)), uint32(len(dst) - len(state.output)), res
}
