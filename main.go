package main

import (
	"log"
	"net/http"
	"time"

	"github.com/idesade/image-resizer/handlers"
)

func main() {
	http.Handle("/", handlers.NewResizeHandler(time.Hour))

	err := http.ListenAndServe(":2000", nil)
	if err != nil {
		log.Fatal(err)
	}
}
