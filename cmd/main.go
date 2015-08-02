package main

import (
	"log"
	"tracker"
)

func main() {
	if t, err := tracker.New(); err != nil {
		panic(err)
	} else {
		defer t.Exit()
		err := t.ApplySampler("/home/tasm/ir/src/tracker/cmd/config/waves.json")
		if err != nil {
			panic(err)
		}
		t.Run()
	}
	log.Println("Done.")
}
