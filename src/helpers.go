package src

import (
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"strconv"
)

const (
	ImgFmtPNG  = "png"
	ImgFmtJPEG = "jpeg"
)

func stringToBin(s string) []uint8 {
	var binString string
	for _, c := range s {
		// Convert character to ASCII value
		asciiVal := int(c)

		// Convert ASCII value to binary and pad with leading zeros
		binaryRep := fmt.Sprintf("%08b", asciiVal)

		// Append binary representation to result
		binString += binaryRep
	}

	return listBinToListNumber(binString)
}

func listBinToListNumber(binString string) []uint8 {
	binArray := make([]uint8, len(binString))
	for i, c := range binString {
		if c == '1' {
			binArray[i] = 1
		} else {
			binArray[i] = 0
		}
	}
	return binArray
}

func toNumber(r []uint8, offset int, seek int) (uint32, int) {
	// retrieve first 'seek' digits (binary) from array and convert to uint8
	values := make([]uint8, 0)
	var i int
	for i = offset; i < len(r) && i < offset+seek; i++ {
		values = append(values, r[i]<<(seek+offset-i-1))
	}

	var value uint8
	for _, v := range values {
		value |= v
	}

	return uint32(value), i
}

func binaryToHumanReadable(binaryString string) (string, error) {
	var result string

	// Split the binary string into chunks of 8 bits (1 byte)
	for i := 0; i < len(binaryString); i += 8 {
		// Get the next 8 bits
		end := i + 8
		if end > len(binaryString) {
			end = len(binaryString)
		}
		byteString := binaryString[i:end]

		// Convert the byte string to an integer
		num, err := strconv.ParseInt(byteString, 2, 64)
		if err != nil {
			return "", err
		}

		// Convert the integer to a character and append it to the result
		result += string(rune(num))
	}

	return result, nil
}

func getMaxStringSize(imgCap int) int {
	return (imgCap*9 + 64) / 8
}

func parseImage(imgFmt string, reader io.Reader) (image.Image, error) {
	switch imgFmt {
	case ImgFmtPNG:
		return png.Decode(reader)
	case ImgFmtJPEG:
		return jpeg.Decode(reader)
	default:
		return nil, errors.New("the image format is not supported")
	}
}
