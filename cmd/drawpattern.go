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

	// Draw a character.

	// Draw an tracker.Event, with helper functions that Draw sequences of characters.
	// Event = pattern.track[0].event[0]
	e := (*p)[0][0]
	view.NewEvent(e).Draw(1, 1)

	// Draw a tracker.Line - a series of tracker.Events that occur at the same moment in time.
	lines := (*p).GetLines()
	view.NewLine(lines[0]).Draw(5, 5)

	// Draw a tracker.Track - a series of tracker.Events played through time (by one instrument).
	track := (*p)[0]
	view.NewTrack(track).Draw(10, 7)

	// Draw a tracker.Pattern - (a series of tracker.Tracks Drawn side by side).
	view.NewPattern(p).DrawBuffered(32, 10)

	// Redraw a tracker.Line over the tracker.Pattern
	view.NewLine(lines[0]).Draw(32, 10)

	// Draw another tracker.Line, but highlighted, over the tracker.Pattern.
	lineNum := 2
	l := view.NewLine(lines[lineNum])
	l.Highlight()
	l.Draw(32, 10+lineNum)

	// Redraw a tracker.Track next to itself a few times.
	t := view.NewTrack(track)
	t.Draw(64, 32)
	t.Draw(64+t.Width(), 32)
	t.Draw(64+t.Width()*2, 32)
	t.Draw(64, 32+t.Height())

	screen.Flush()
	time.Sleep(time.Duration(args.displayTime) * time.Second)
}
