package main

import (
	"log"
	"net/http"

	"github.com/idesade/image-resizer/handlers"
)

func main() {
	http.Handle("/", &handlers.ResizeHandler{})

	err := http.ListenAndServe(":2000", nil)
	if err != nil {
		log.Fatal(err)
	}
}
