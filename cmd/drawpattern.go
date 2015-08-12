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

var config = struct {
	screen
	noteNum
	velocity
	event
	line
	track
	pattern
}{
	termScreen{},

	// Embedding a view struct within a foo struct does cause
	// stuttering of the word "view" in config declaration, but turns
	// the foo constructor methods into one-liners.
	noteNum{view: view{
		fg: termbox.ColorBlue,
		bg: termbox.ColorDefault,
	}},

	velocity{view: view{
		fg: termbox.ColorCyan,
		bg: termbox.ColorDefault,
	}},

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
		patternFile   string
		displayTime   int
		useMockScreen bool
	}{}
	flag.StringVar(&args.patternFile, "pattern", "testpattern.trkr", "The pattern file to read.")
	flag.IntVar(&args.displayTime, "sec", 5, "The number of seconds to show the pattern for.")
	flag.BoolVar(&args.useMockScreen, "mock", false, "Print a mock screen instead of drawing with termbox")
	flag.Parse()

	p, err := tracker.NewPattern(args.patternFile)
	if err != nil {
		panic(err)
	}

	if args.useMockScreen {
		config.screen = newMockScreen(200, 58)
		args.displayTime = 0
	}

	if err := config.screen.Init(); err != nil {
		panic(err)
	}
	defer config.screen.Close()

	// Draw a character with termbox.
	config.screen.SetCell(2, 3, 'A', fg, bg)

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
	t.draw(64+t.width, 32)
	t.draw(64+t.width*2, 32)
	t.draw(64, 32+t.height)

	config.screen.Flush()
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
	for _, e := range tv.Track {
		ev := newEvent(e)
		ev.draw(x, y+tv.height)
		if ev.width > tv.width {
			tv.width = ev.width
		}
		for i, r := range tv.delimiter {
			config.screen.SetCell(x+tv.width+i, y+tv.height, r, tv.fg, tv.bg)
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
	ev.width, ev.height = 0, 0
	n := newNoteNum(ev.NoteNum)
	n.draw(x+ev.width, y)
	ev.width += n.width
	if n.height > ev.height {
		ev.height = n.height
	}
	for i, r := range ev.delimiter {
		termbox.SetCell(x+ev.width+i, y, r, ev.fg, ev.bg)
		ev.width += i
	}
	v := newVelocity(ev.Velocity)
	v.draw(x+ev.width, y)
	ev.width += v.width
	if v.height > ev.height {
		v.height = ev.height
	}
}

type noteNum struct {
	tracker.NoteNum
	view
}

func newNoteNum(n tracker.NoteNum) *noteNum {
	return &noteNum{n, config.noteNum.view}
}

func (n *noteNum) draw(x, y int) {
	n.width, n.height = 0, 0
	s := fmt.Sprintf("%v", n.NoteNum)
	n.width = len(s)
	n.height = 1
	for i, r := range s {
		config.screen.SetCell(x+i, y, r, n.fg, n.bg)
	}
}

type velocity struct {
	tracker.Velocity
	view
}

func newVelocity(v tracker.Velocity) *velocity {
	return &velocity{v, config.velocity.view}
}

func (v *velocity) draw(x, y int) {
	v.width, v.height = 0, 0
	s := fmt.Sprintf("%v", v.Velocity)
	v.width = len(s)
	v.height = 1
	for i, r := range s {
		config.screen.SetCell(x+i, y, r, v.fg, v.bg)
	}
}

func drawString(x, y int, s string) {
	for i, r := range s {
		config.screen.SetCell(x+i, y, r, fg, bg)
		// termbox.SetCell(x+i, y, r, fg, bg)
	}
}

type screen interface {
	Init() error
	Close()
	SetCell(x, y int, r rune, fg, bg termbox.Attribute)
	Flush()
}

type termScreen struct{}

func (t termScreen) Init() error { return termbox.Init() }
func (t termScreen) Close()      { termbox.Close() }
func (t termScreen) Flush()      { termbox.Flush() }
func (t termScreen) SetCell(x, y int, r rune, fg, bg termbox.Attribute) {
	termbox.SetCell(x, y, r, fg, bg)
}

type mockScreen struct {
	cells [][]rune
}

func newMockScreen(width, height int) *mockScreen {
	m := &mockScreen{}
	m.cells = make([][]rune, height)
	for i, _ := range m.cells {
		m.cells[i] = make([]rune, width)
	}
	return m
}

func (m *mockScreen) Init() error { return nil }

func (m *mockScreen) Close() {}

func (m *mockScreen) SetCell(x, y int, r rune, fg, bg termbox.Attribute) {
	m.cells[y][x] = r
}

func (m *mockScreen) Flush() {
	for x, row := range m.cells {
		for y, r := range row {
			if r == 0 {
				m.cells[x][y] = ' '
			}
		}
	}
	for _, row := range m.cells {
		fmt.Println(string(row))
	}
	m.clear()
}

func (m *mockScreen) clear() {
	for x := 0; x < len(m.cells); x++ {
		for y := 0; y < len(m.cells[x]); y++ {
			m.cells[x][y] = 0
		}
	}
}

/*
	for i, track := range *p {
		fmt.Println(i, track)
		for x, event := range track {
			fmt.Println(x, event)
		}
	}*/
