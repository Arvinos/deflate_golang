package deflate_golang

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"testing"
)

const MaxOffset int = 32768
const MaxMatch int = 258
const MinMatch int = 3

func lzGenerator(size int, matchProbability float32) []byte {
	vector := make([]byte, size)
	actualLength := 0

	for actualLength < size {
		value := rand.Intn(math.MaxUint8)
		matchToken := rand.Float32()

		if matchToken < matchProbability && actualLength > MinMatch {
			var offset = rand.Intn(MaxOffset)
			// Offset is not greater than the actualLength
			offset = int(math.Min(float64(offset), float64(actualLength)))

			// match is in range [MinMatch : MaxMatch]
			var match = rand.Intn(MaxMatch-MinMatch) + MinMatch
			// match is not greater than the rest size
			match = int(math.Min(float64(match), float64(size-actualLength)))
			// match is not greater than a current offset
			match = int(math.Min(float64(match), float64(offset)))

			matchSource := vector[actualLength-offset : actualLength-offset+match]
			matchDestination := vector[actualLength : actualLength+match]

			copy(matchDestination, matchSource)

			actualLength += match
		} else {
			literal := byte(value)
			vector[actualLength] = literal

			actualLength++
		}
	}

	return vector
}

func TestDeflateStoredWR(t *testing.T) {
	testData := lzGenerator(96, 0.25)
	compressedData := make([]byte, len(testData)*2)
	outputData := make([]byte, len(testData))

	in, out, ret := writeStoredBlock(testData, compressedData)

	if ret != 0 {
		panic("Compression error")
	}

	fmt.Println("Compressed bytes = ", in, ", Output bytes = ", out, ", Ratio = ", float64(in)/float64(out))

	in, out, ret = readStoredBlock(compressedData, outputData)

	if ret != 0 {
		panic("Decompression error")
	}

	fmt.Println("Input bytes = ", in, ", Decompressed bytes = ", out)

	if !bytes.Equal(testData, outputData) {
		fmt.Println(testData)
		fmt.Println(outputData)
		panic("Data is not equal")
	}
}

func TestDeflateFixed(t *testing.T) {
	testData := lzGenerator(96, 0.25)
	compressedData := make([]byte, len(testData)*2)
	outputData := make([]byte, len(testData))

	fmt.Println(testData)

	deflate := new(Deflate)

	in, out, ret := deflate.Compress(testData, compressedData)

	if ret != 0 {
		panic("Compression error")
	}

	fmt.Println("Compressed bytes = ", in, ", Output bytes = ", out, ", Ratio = ", float64(in)/float64(out))

	in, out, ret = deflate.Decompress(compressedData, outputData)

	if ret != 0 {
		panic("Decompression error")
	}

	fmt.Println("Input bytes = ", in, ", Decompressed bytes = ", out)

	if !bytes.Equal(testData, outputData) {
		fmt.Println(testData)
		fmt.Println(outputData)
		panic("Data is not equal")
	}
}

func TestDeflateDynamic(t *testing.T) {
	testData := lzGenerator(96, 0.25)
	compressedData := make([]byte, len(testData)*2)
	outputData := make([]byte, len(testData))

	fmt.Println(testData)

	deflate := new(Deflate)

	in, out, ret := deflate.Compress(testData, compressedData)

	if ret != 0 {
		panic("Compression error")
	}

	fmt.Println("Compressed bytes = ", in, ", Output bytes = ", out, ", Ratio = ", float64(in)/float64(out))

	in, out, ret = deflate.Decompress(compressedData, outputData)

	if ret != 0 {
		panic("Decompression error")
	}

	fmt.Println("Input bytes = ", in, ", Decompressed bytes = ", out)
}

func TestDeflate(t *testing.T) {
	testData := lzGenerator(96, 0.25)
	compressedData := make([]byte, len(testData)*2)
	outputData := make([]byte, len(testData))

	fmt.Println(testData)

	deflate := new(Deflate)

	in, out, ret := deflate.Compress(testData, compressedData)

	if ret != 0 {
		panic("Compression error")
	}

	fmt.Println("Compressed bytes = ", in, ", Output bytes = ", out, ", Ratio = ", float64(in)/float64(out))

	in, out, ret = deflate.Decompress(compressedData, outputData)

	if ret != 0 {
		panic("Decompression error")
	}

	fmt.Println("Input bytes = ", in, ", Decompressed bytes = ", out)

	if !bytes.Equal(testData, outputData) {
		fmt.Println(testData)
		fmt.Println(outputData)
		panic("Data is not equal")
	}
}
