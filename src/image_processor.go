package src

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"io"
	"strconv"
)

func newSecretImage(inputString string, imgFmt string, reader io.Reader) (image.Image, error) {
	img, err := parseImage(imgFmt, reader)
	if err != nil {
		return nil, err
	}

	// Get the image dimensions
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	if !isValidImageCap(width*height, inputString) {
		return nil, errors.New("image dimensions are not enough")
	}

	rgbaImg := image.NewRGBA(img.Bounds())
	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			rgbaImg.Set(x, y, img.At(x, y))
		}
	}

	bits := stringToBin(inputString)
	sl := fmt.Sprintf("%064b", len(bits))
	bits = append(listBinToListNumber(sl), bits...)

	var (
		offset int
		val    uint32
	)

outer:
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			if len(bits) <= offset {
				break outer
			}

			// Get the pixel at position (i, j)
			r, g, b, a := img.At(i, j).RGBA()

			val, offset = toNumber(bits, offset, 3)
			r = r>>8>>3<<3 | val

			val, offset = toNumber(bits, offset, 3)
			g = g>>8>>3<<3 | val

			val, offset = toNumber(bits, offset, 3)
			b = b>>8>>3<<3 | val

			rgbaImg.Set(i, j, color.RGBA{
				R: uint8(r),
				G: uint8(g),
				B: uint8(b),
				A: uint8(a),
			})
		}
	}

	return rgbaImg, err
}

func readSecretImage(imgFmt string, reader io.Reader) (string, error) {
	img, err := parseImage(imgFmt, reader)
	if err != nil {
		fmt.Println("Error decoding image:", err)
		return "", err
	}

	width, height := img.Bounds().Max.X, img.Bounds().Max.Y

	var (
		bits            []rune
		msgBinLength    int64
		totalBitCounted int64
	)

outer:
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			r, g, b, _ := img.At(i, j).RGBA()

			if msgBinLength > 0 && msgBinLength <= totalBitCounted {
				break outer
			}

			bits = append(bits, []rune(fmt.Sprintf("%08b", r>>8))[5:]...)
			bits = append(bits, []rune(fmt.Sprintf("%08b", g>>8))[5:]...)
			bits = append(bits, []rune(fmt.Sprintf("%08b", b>>8))[5:]...)

			totalBitCounted += 9

			if len(bits) >= 64 {
				if msgBinLength > 0 {
					continue
				}

				bits = bits[64:]
				msgBinLength, err = strconv.ParseInt(string(bits[:64]), 2, 64)
				totalBitCounted -= 64
			}
		}
	}

	if totalBitCounted < msgBinLength {
		return "", errors.New("image data is not enough")
	}

	bits = bits[:msgBinLength]

	return binaryToHumanReadable(string(bits))
}

func getImgMaxCap(imgFmt string, reader io.Reader) (int, error) {
	img, err := parseImage(imgFmt, reader)
	if err != nil {
		return 0, err
	}

	if img == nil {
		return 0, errors.New("image is nil")
	}

	return img.Bounds().Dx() * img.Bounds().Dy(), nil
}
