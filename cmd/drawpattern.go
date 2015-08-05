package main

import (
	"flag"
	"fmt"
	"github.com/nsf/termbox-go"
	"time"
	"tracker"
)

const (
	fg = termbox.ColorDefault
	bg = termbox.ColorDefault
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
	if err := termbox.Init(); err != nil {
		panic(err)
	}
	defer termbox.Close()
	termbox.SetCell(2, 3, 'A', fg, bg)

	//	Event = pattern.track[0].event[0]
	e := (*p)[0][0]
	newEventView(e).draw(1, 1)
	termbox.Flush()
	time.Sleep(2 * time.Second)
}

type eventView struct {
	*tracker.Event
	width, height int
	fg, bg        termbox.Attribute
}

func newEventView(e *tracker.Event) eventView {
	return eventView{height: 1, width: 0, fg: fg, bg: bg, Event: e}
}

func (ev eventView) draw(x, y int) {
	s := fmt.Sprintf("%v %v", ev.NoteNum, ev.Velocity)
	ev.width = len(s)
	for i, r := range s {
		termbox.SetCell(x+i, y, r, ev.fg, ev.bg)
	}
}

func drawString(x, y int, s string) {
	for i, r := range s {
		termbox.SetCell(x+i, y, r, fg, bg)
	}
}
/*
	for i, track := range *p {
		fmt.Println(i, track)
		for x, event := range track {
			fmt.Println(x, event)
		}
	}*/
