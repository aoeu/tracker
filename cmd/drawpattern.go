package main

import (
	"flag"
	"fmt"
	"tracker"
)

func main() {
	args := struct {
		patternFile string
	}{}
	flag.StringVar(&args.patternFile, "pattern", "testpattern.trkr", "The pattern file to read.")
	flag.Parse()

	p, err := tracker.NewPattern(args.patternFile)
	if err != nil {
		panic(err)
	}

	fmt.Println(p)
}
