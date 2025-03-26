package src

const MinImageCapacity = 12

func isValidImageFormat(imgFmt string) bool {
	return imgFmt == ImgFmtPNG || imgFmt == ImgFmtJPEG
}

// 9 bits for each pixel with R,G,and B properties
// minimum 64 bits for message binary length
func isValidImageCap(imgCap int, msg string) bool {
	return imgCap >= MinImageCapacity && len(stringToBin(msg))+64 < imgCap*9
}
