package view

import (
	"fmt"
	"github.com/nsf/termbox-go"
	"tracker"
)

var config = newDefaultConfig()

type viewConfig struct {
	screen
	noteNum
	velocity
	event
	line
	track
	pattern
}

func newDefaultConfig() *viewConfig {
	return &viewConfig{
		termScreen{},

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
