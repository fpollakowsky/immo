package main

import (
	"immo/pkg/read"
	"time"
)

func main() {
	for {
		read.Landeseigen()

		time.Sleep(30 * time.Second)
	}
}
