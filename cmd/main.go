package main

import (
	"log"
	"tracker"
)

func main() {
	if err := tracker.New(); err != nil {
		log.Fatal(err)
	}
}
