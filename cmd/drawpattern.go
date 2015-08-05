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

	if err := termbox.Init(); err != nil {
		panic(err)
	}
	defer termbox.Close()
	//termbox.SetCell(2, 3, 'A', fg, bg)

	//	Event = pattern.track[0].event[0]
	//e := (*p)[0][0]
	// newEventView(e).draw(1, 1)
	lines := (*p).GetLines()
	newLineView(lines[0]).draw(5, 5)
	termbox.Flush()
	time.Sleep(2 * time.Second)
}

type lineView struct {
	tracker.Line
	width, height int
	fg, bg        termbox.Attribute
	delimiter     string
}

func newLineView(l tracker.Line) *lineView {
	return &lineView{Line: l,
		fg:        termbox.ColorGreen,
		bg:        termbox.ColorDefault,
		delimiter: " | ",
	}

}

func (lv *lineView) draw(x, y int) {
	for _, e := range lv.Line {
		ev := newEventView(e)
		ev.draw(x+lv.width, y)
		lv.width += ev.width
		if ev.height > lv.height {
			lv.height = ev.height
		}
		for i, r := range lv.delimiter {
			termbox.SetCell(x+lv.width+i, y, r, lv.fg, lv.bg)
		}
		lv.width += len(lv.delimiter)
	}
}

type eventView struct {
	*tracker.Event
	width, height int
	fg, bg        termbox.Attribute
}

func newEventView(e *tracker.Event) *eventView {
	return &eventView{height: 1, width: 0, fg: fg, bg: bg, Event: e}
}

func (ev *eventView) draw(x, y int) {
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
