package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/idesade/image-resizer/handlers"
)

func TestServer(t *testing.T) {
	server := httptest.NewServer(&handlers.ResizeHandler{})
	defer server.Close()

	u, err := url.Parse(server.URL)
	if err != nil {
		t.Fatal(err)
	}

	q := u.Query()
	q.Add("width", "200")
	q.Add("height", "200")
	q.Add("url", "https://www.datapipe.com/blog/wp-content/uploads/2015/12/big-data-will-drive-the-next-phase-of-innovation-in-mobile-computing.jpg")
	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Fatalf("Received non-200 response: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	if len(body) == 0 {
		t.Fatal("Body is empty")
	}
}
