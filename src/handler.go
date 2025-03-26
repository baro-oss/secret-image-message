package src

import (
	"encoding/json"
	"image/jpeg"
	"image/png"
	"net/http"
	"strings"
)

func HandleEncodeImage(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("img")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	imgFileName := strings.Split(header.Filename, ".")
	imgFmt := imgFileName[len(imgFileName)-1]

	if !isValidImageFormat(imgFmt) {
		http.Error(w, "Invalid image format. Only PNG and JPEG are supported", http.StatusBadRequest)
		return
	}

	messages := r.URL.Query()["message"]
	if len(messages) == 0 {
		http.Error(w, "message parameter is required", http.StatusBadRequest)
		return
	}

	img, err := newSecretImage(messages[0], imgFmt, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	switch imgFmt {
	case ImgFmtPNG:
		err = png.Encode(w, img)
	case ImgFmtJPEG:
		err = jpeg.Encode(w, img, nil)
	}

	w.Header().Set("Content-Type", "image/"+imgFmt)
	w.WriteHeader(http.StatusOK)
}

func HandleDecodeImage(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("img")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	imgFileName := strings.Split(header.Filename, ".")
	imgFmt := imgFileName[len(imgFileName)-1]

	if !isValidImageFormat(imgFmt) {
		http.Error(w, "Invalid image format. Only PNG and JPEG are supported", http.StatusBadRequest)
		return
	}

	msg, err := readSecretImage(imgFmt, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(msg))
	w.WriteHeader(http.StatusOK)
}

func HandleGetMaxCapacity(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("img")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	imgFileName := strings.Split(header.Filename, ".")
	imgFmt := imgFileName[len(imgFileName)-1]

	if !isValidImageFormat(imgFmt) {
		http.Error(w, "Invalid image format. Only PNG and JPEG are supported", http.StatusBadRequest)
		return
	}

	maxCap, err := getImgMaxCap(imgFmt, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]int{"max_capacity": getMaxStringSize(maxCap)})
}
