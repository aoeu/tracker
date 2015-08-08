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
	event
	line
	track
	pattern
}{
	// Embedding a view struct within a foo struct does cause
	// stuttering of the word "view" in config declaration, but turns
	// the foo constructor methods into one-liners.
	event{view: view{
		fg:        termbox.ColorBlue,
		bg:        termbox.ColorDefault,
		delimiter: " \u00B7",
	}},

	line{view: view{
		fg:        termbox.ColorRed,
		bg:        termbox.ColorDefault,
		delimiter: " | ",
	}},

	track{view: view{
		fg:        termbox.ColorGreen,
		bg:        termbox.ColorDefault,
		delimiter: " | ",
	}},

	pattern{view: view{
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
	newEvent(e).draw(1, 1)

	// Draw a tracker.Line - a series of tracker.Events that occur at the same moment in time.
	lines := (*p).GetLines()
	newLine(lines[0]).draw(5, 5)

	// Draw a tracker.Track - a series of tracker.Events played through time (by one instrument).
	track := (*p)[0]
	newTrack(track).draw(10, 7)

	// Draw a tracker.Pattern - (a series of tracker.Tracks drawn side by side).
	newPattern(p).draw(32, 10)

	// Redraw a tracker.Line over the tracker.Pattern
	newLine(lines[0]).draw(32, 10)

	// Redraw another tracker.Line over the tracker.Pattern to expose bugs.
	/*
		lineNum := 2
		newLine(lines[lineNum]).draw(32, 10 + lineNum)
	*/

	// Redraw a tracker.Track next to itself a few times.
	t := newTrack(track)
	t.draw(64, 32)
	t.draw(64 + t.width, 32)
	t.draw(64 + t.width * 2, 32)
	t.draw(64, 32 + t.height)

	termbox.Flush()
	time.Sleep(time.Duration(args.displayTime) * time.Second)
}

type view struct {
	width, height int
	fg, bg        termbox.Attribute
	delimiter     string
}

type pattern struct {
	*tracker.Pattern
	view
}

func newPattern(p *tracker.Pattern) *pattern {
	return &pattern{p, config.pattern.view}
}

func (pv *pattern) draw(x, y int) {
	pv.width, pv.height = 0, 0
	for _, t := range *pv.Pattern {
		tv := newTrack(t)
		tv.draw(x+pv.width, y)
		pv.width += tv.width
		if tv.height > pv.height {
			pv.height = tv.height
		}
	}
}

type track struct {
	tracker.Track
	view
}

func newTrack(t tracker.Track) *track {
	return &track{t, config.track.view}
}

func (tv *track) draw(x, y int) {
	tv.width, tv.height = 0, 0
	// TODO(aoeu): Reset width and height every call to draw?
	for _, e := range tv.Track {
		ev := newEvent(e)
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

type line struct {
	tracker.Line
	view
}

func newLine(l tracker.Line) *line {
	return &line{l, config.line.view}
}

func (lv *line) draw(x, y int) {
	lv.width, lv.height = 0, 0
	for _, e := range lv.Line {
		ev := newEvent(e)
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

type event struct {
	*tracker.Event
	view
}

func newEvent(e *tracker.Event) *event {
	return &event{e, config.event.view}
}

func (ev *event) draw(x, y int) {
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
