package main

import (
	"flag"
	"os"
	"time"
	"tracker"
	"tracker/view"
)

func main() {
	args := struct {
		patternFile   string
		displayTime   int
		useMockScreen bool
	}{}
	flag.StringVar(&args.patternFile, "pattern", "testpattern.trkr", "The pattern file to read.")
	flag.IntVar(&args.displayTime, "sec", 5, "The number of seconds to show the pattern for.")
	flag.BoolVar(&args.useMockScreen, "mock", false, "Print a mock Screen instead of Drawing with termbox")
	flag.Parse()

	p, err := tracker.NewPattern(args.patternFile)
	if err != nil {
		panic(err)
	}

	if args.useMockScreen {
		view.Config.Screen = view.NewMockScreen(os.Stdout, 200, 58)
		args.displayTime = 0
	}

	screen := view.Config.Screen
	if err := screen.Init(); err != nil {
		panic(err)
	}
	defer screen.Close()
	defer time.Sleep(time.Duration(args.displayTime) * time.Second)

	pt := make(tracker.PatternTable, 1)
	pt[0] = p
	ui := view.NewUI(&pt)
	for i, _ := range p.GetLines() {
		ui.Selection.LineOffset = i
		ui.Draw()
		screen.Flush()
		time.Sleep(1 * time.Second)
	}

}
