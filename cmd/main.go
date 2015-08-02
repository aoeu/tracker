package main

import (
	"log"
	"tracker"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	if t, err := tracker.New(); err != nil {
		panic(err)
	} else {
		defer t.Exit()
		t.ApplySampler("/home/tasm/ir/src/tracker/cmd/config/waves.json")
		t.Run()
	}
	log.Println("Done.")
}
