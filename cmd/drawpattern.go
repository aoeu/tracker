package main

import (
	"flag"
	"fmt"
	"github.com/nsf/termbox-go"
	"strings"
	"time"
	"tracker"
)

const (
	fg = termbox.ColorDefault
	bg = termbox.ColorDefault
)

var config = struct {
	eventView
	lineView
	trackView
	patternView
}{
	// Embedding a view struct within a fooView struct does cause
	// stuttering of the word "view" in config declaration, but turns
	// the fooView constructor methods into one-liners.
	eventView{view: view{
		fg:        termbox.ColorBlue,
		bg:        termbox.ColorDefault,
		delimiter: " \u00B7 ",
	}},

	lineView{view: view{
		fg:        termbox.ColorRed,
		bg:        termbox.ColorDefault,
		delimiter: " | ",
	}},

	trackView{view: view{
		fg:        termbox.ColorGreen,
		bg:        termbox.ColorDefault,
		delimiter: " | ",
	}},

	patternView{view: view{
		fg: termbox.ColorYellow,
		bg: termbox.ColorDefault,
	}},
}

func main() {
	args := struct {
		patternFile string
		displayTime int
	}{}
	flag.StringVar(&args.patternFile, "pattern", "testpattern.trkr", "The pattern file to read.")
	flag.IntVar(&args.displayTime, "sec", 5, "The number of seconds to show the pattern for.")
	flag.Parse()

	p, err := tracker.NewPattern(args.patternFile)
	if err != nil {
		panic(err)
	}

	if err := termbox.Init(); err != nil {
		panic(err)
	}
	defer termbox.Close()

	// Draw a character with termbox.
	termbox.SetCell(2, 3, 'A', fg, bg)

	// Draw an tracker.Event, with helper functions that draw sequences of characters.
	// Event = pattern.track[0].event[0]
	e := (*p)[0][0]
	newEventView(e).draw(1, 1)

	// Draw a tracker.Line - a series of tracker.Events that occur at the same moment in time.
	lines := (*p).GetLines()
	newLineView(lines[0]).draw(5, 5)

	// Draw a tracker.Track - a series of tracker.Events played through time (by one instrument).
	track := (*p)[0]
	newTrackView(track).draw(10, 7)

	// Draw a tracker.Pattern - (a series of tracker.Tracks drawn side by side).
	newPatternView(p).draw(32, 10)

	// Redraw a tracker.Line over the tracker.Pattern
	newLineView(lines[0]).draw(32, 10)

	// Redraw another tracker.Line over the tracker.Pattern to expose bugs.
	/*
		lineNum := 2
		newLineView(lines[lineNum]).draw(32, 10 + lineNum)
	*/

	termbox.Flush()
	time.Sleep(time.Duration(args.displayTime) * time.Second)
}

type view struct {
	width, height int
	fg, bg        termbox.Attribute
	delimiter     string
}

type patternView struct {
	*tracker.Pattern
	view
}

func newPatternView(p *tracker.Pattern) *patternView {
	return &patternView{p, config.patternView.view}
}

func (pv *patternView) draw(x, y int) {
	for _, t := range *pv.Pattern {
		tv := newTrackView(t)
		tv.draw(x+pv.width, y)
		pv.width += tv.width
		if tv.height > pv.height {
			pv.height = tv.height
		}
	}
}

type trackView struct {
	tracker.Track
	view
}

func newTrackView(t tracker.Track) *trackView {
	return &trackView{t, config.trackView.view}
}

func (tv *trackView) draw(x, y int) {
	// TODO(aoeu): Reset width and height every call to draw?
	for _, e := range tv.Track {
		ev := newEventView(e)
		ev.draw(x, y+tv.height)
		if ev.width > tv.width {
			tv.width = ev.width
		}
		for i, r := range tv.delimiter {
			termbox.SetCell(x+tv.width+i, y+tv.height, r, tv.fg, tv.bg)
		}
		tv.height += ev.height
	}
	tv.width += len(tv.delimiter)
}

type lineView struct {
	tracker.Line
	view
}

func newLineView(l tracker.Line) *lineView {
	return &lineView{l, config.lineView.view}
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
	view
}

func newEventView(e *tracker.Event) *eventView {
	return &eventView{e, config.eventView.view}
}

func (ev *eventView) draw(x, y int) {
	s := fmt.Sprintf("%v%v%v", ev.NoteNum, ev.delimiter, ev.Velocity)
	ev.width = len(s)
	ev.height = 1 + (1 * strings.Count(s, "\n"))
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
