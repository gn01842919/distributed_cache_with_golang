package main

import (
	"./http"

	"./cache"
)

func main() {
	c := cache.New("inmemory")
	http.New(c).Listen()
}
