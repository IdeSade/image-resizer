package utils

import (
	"testing"
	"time"
)

const key = "1"
const data = "123456789"

func TestCache(t *testing.T) {
	cache := NewCache(time.Second)

	cache.AddItem(key, []byte(data))

	c, ok := cache.GetItem(key)
	if !ok {
		t.Fatal("Item not found")
	}

	if string(c) != data {
		t.Fatalf("Cache data not equals const data: %s != %s", string(c), data)
	}

	time.Sleep(2 * time.Second)

	_, ok = cache.GetItem(key)
	if ok {
		t.Fatal("Item found but must removed")
	}
}
