package handlers

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/disintegration/imaging"
	"github.com/idesade/image-resizer/utils"
	"github.com/pkg/errors"
)

var Testing = false

type ResizeHandler struct {
	cacheImages *utils.Cache
}

func NewResizeHandler(duration time.Duration) *ResizeHandler {
	return &ResizeHandler{cacheImages: utils.NewCache(duration)}
}

func (h *ResizeHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	width, height, u, err := getParams(req)
	if err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}

	key := makeKey(width, height, u)

	data, ok := h.cacheImages.GetItem(key)
	if ok {
		log.Print("Image cache hit")
		if Testing {
			w.Header().Add("FromCache", "true")
		}
		w.Write(data)
		return
	}

	resp, err := http.Get(u)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	log.Print("Get image")

	data, err = resizeImage(resp.Body, width, height)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	h.cacheImages.AddItem(key, data)

	if Testing {
		w.Header().Add("FromCache", "false")
	}

	w.Write(data)
}

func getParams(req *http.Request) (int, int, string, error) {
	width := req.URL.Query().Get("width")
	height := req.URL.Query().Get("height")
	u := req.URL.Query().Get("url")

	w, err := strconv.Atoi(width)
	if err != nil {
		return 0, 0, "", errors.Wrap(err, "parse width")
	}

	h, err := strconv.Atoi(height)
	if err != nil {
		return 0, 0, "", errors.Wrap(err, "parse height")
	}

	return w, h, u, nil
}

func makeKey(w, h int, u string) string {
	return fmt.Sprintf("%d+%d+%s", w, h, u)
}

func resizeImage(r io.Reader, w, h int) ([]byte, error) {
	img, err := imaging.Decode(r)
	if err != nil {
		return nil, errors.Wrap(err, "decode image")
	}
	log.Print("Decode image")

	dstImg := imaging.Resize(img, w, h, imaging.Lanczos)
	log.Print("Resize image")

	buf := bytes.NewBuffer(make([]byte, 0, bytes.MinRead))

	err = imaging.Encode(buf, dstImg, imaging.JPEG)
	if err != nil {
		return nil, errors.Wrap(err, "encode image")
	}
	log.Print("Encode image")

	return buf.Bytes(), nil
}

func writeError(w http.ResponseWriter, err error, status int) {
	log.Print(err)
	w.WriteHeader(status)
	w.Write([]byte(err.Error()))
}
