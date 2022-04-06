package main

import (
	"immo/pkg/read"
	"time"
)

func main() {
	for {
		//read.Gewobag()
		read.Wbm()

		time.Sleep(30 * time.Second)
	}
}
