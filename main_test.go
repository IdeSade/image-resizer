package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/idesade/image-resizer/handlers"
	"github.com/pkg/errors"
)

func TestServer(t *testing.T) {
	handlers.Testing = true

	server := httptest.NewServer(handlers.NewResizeHandler(time.Second))
	defer server.Close()

	u, err := makeQuery(server.URL)
	if err != nil {
		t.Fatal(err)
	}

	err = checkRequest(u, "FromCache", "false")
	if err != nil {
		t.Fatal(err)
	}

	err = checkRequest(u, "FromCache", "true")
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(2 * time.Second)

	err = checkRequest(u, "FromCache", "false")
	if err != nil {
		t.Fatal(err)
	}
}

func makeQuery(serverUrl string) (string, error) {
	u, err := url.Parse(serverUrl)
	if err != nil {
		return "", errors.Wrap(err, "parse serverUrl")
	}

	q := u.Query()
	q.Add("width", "200")
	q.Add("height", "200")
	q.Add("url", "https://www.datapipe.com/blog/wp-content/uploads/2015/12/big-data-will-drive-the-next-phase-of-innovation-in-mobile-computing.jpg")
	u.RawQuery = q.Encode()

	return u.String(), nil
}

func checkRequest(url, key, value string) error {
	resp, err := http.Get(url)
	if err != nil {
		return errors.Wrap(err, "get")
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.Errorf("Received non-200 response: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "readAll")
	}

	if len(body) == 0 {
		return errors.New("Body is empty")
	}

	if resp.Header.Get(key) != value {
		return errors.Errorf("Wrong header: %s != %s", resp.Header.Get(key), value)
	}

	return nil
}
