package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/disintegration/imaging"
)

type ResizeHandler struct {
}

func (r *ResizeHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	widthA := req.URL.Query().Get("width")
	heightA := req.URL.Query().Get("height")
	u := req.URL.Query().Get("url")

	log.Printf("Receive params: width - %q, height - %q, url - %q", widthA, heightA, u)

	width, err := strconv.Atoi(widthA)
	if err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}

	height, err := strconv.Atoi(heightA)
	if err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}

	resp, err := http.Get(u)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	log.Print("Get image")

	img, err := imaging.Decode(resp.Body)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}
	log.Print("Decode image")

	dstImg := imaging.Resize(img, width, height, imaging.Lanczos)
	log.Print("Resize image")

	err = imaging.Encode(w, dstImg, imaging.JPEG)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}
	log.Print("Encode image")
}

func writeError(w http.ResponseWriter, err error, status int) {
	log.Print(err)
	w.WriteHeader(status)
	w.Write([]byte(err.Error()))
}
