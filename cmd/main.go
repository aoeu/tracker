package main

import (
	"log"
	"tracker"
)

func main() {
	if t, err := tracker.New(); err != nil {
		log.Fatal(err)
	} else {
		t.Run()
		defer t.Exit()
	}
	log.Println("Done.")
}
