package view

import (
	"fmt"
	"github.com/nsf/termbox-go"
	"tracker"
)

var Config = NewDefaultConfig()

type ViewConfig struct {
	Screen
	NoteNum
	Velocity
	Event
	Line
	Track
	Pattern
}

func NewDefaultConfig() *ViewConfig {
	return &ViewConfig{
		TermScreen{},

		NoteNum{View: View{
			Fg: termbox.ColorBlue,
			Bg: termbox.ColorDefault,
		}},

		Velocity{View: View{
			Fg: termbox.ColorCyan,
			Bg: termbox.ColorDefault,
		}},

		Event{View: View{
			Fg:        termbox.ColorBlue,
			Bg:        termbox.ColorDefault,
			delimiter: " \u00B7",
		}},

		Line{View: View{
			Fg:        termbox.ColorRed,
			Bg:        termbox.ColorDefault,
			delimiter: " | ",
		}},

		Track{View: View{
			Fg:        termbox.ColorGreen,
			Bg:        termbox.ColorDefault,
			delimiter: " | ",
		}},

		Pattern{View: View{
			Fg: termbox.ColorYellow,
			Bg: termbox.ColorDefault,
		}},
	}
}

type View struct {
	Width, Height int
	Fg, Bg        termbox.Attribute
	delimiter     string
}

type Pattern struct {
	*tracker.Pattern
	View
}

func NewPattern(p *tracker.Pattern) *Pattern {
	return &Pattern{p, Config.Pattern.View}
}

func (pv *Pattern) Draw(x, y int) {
	pv.Width, pv.Height = 0, 0
	for _, t := range *pv.Pattern {
		tv := NewTrack(t)
		tv.Draw(x+pv.Width, y)
		pv.Width += tv.Width
		if tv.Height > pv.Height {
			pv.Height = tv.Height
		}
	}
}

type Track struct {
	tracker.Track
	View
}

func NewTrack(t tracker.Track) *Track {
	return &Track{t, Config.Track.View}
}

func (tv *Track) Draw(x, y int) {
	tv.Width, tv.Height = 0, 0
	for _, e := range tv.Track {
		ev := NewEvent(e)
		ev.Draw(x, y+tv.Height)
		if ev.Width > tv.Width {
			tv.Width = ev.Width
		}
		for i, r := range tv.delimiter {
			Config.Screen.SetCell(x+tv.Width+i, y+tv.Height, r, tv.Fg, tv.Bg)
		}
		tv.Height += ev.Height
	}
	tv.Width += len(tv.delimiter)
}

type Line struct {
	tracker.Line
	View
}

func NewLine(l tracker.Line) *Line {
	return &Line{l, Config.Line.View}
}

func (lv *Line) Draw(x, y int) {
	lv.Width, lv.Height = 0, 0
	for _, e := range lv.Line {
		ev := NewEvent(e)
		ev.Draw(x+lv.Width, y)
		lv.Width += ev.Width
		if ev.Height > lv.Height {
			lv.Height = ev.Height
		}
		for i, r := range lv.delimiter {
			termbox.SetCell(x+lv.Width+i, y, r, lv.Fg, lv.Bg)
		}
		lv.Width += len(lv.delimiter)
	}
}

type Event struct {
	*tracker.Event
	View
}

func NewEvent(e *tracker.Event) *Event {
	return &Event{e, Config.Event.View}
}

func (ev *Event) Draw(x, y int) {
	ev.Width, ev.Height = 0, 0
	n := NewNoteNum(ev.NoteNum)
	n.Draw(x+ev.Width, y)
	ev.Width += n.Width
	if n.Height > ev.Height {
		ev.Height = n.Height
	}
	for i, r := range ev.delimiter {
		termbox.SetCell(x+ev.Width+i, y, r, ev.Fg, ev.Bg)
		ev.Width += i
	}
	v := NewVelocity(ev.Velocity)
	v.Draw(x+ev.Width, y)
	ev.Width += v.Width
	if v.Height > ev.Height {
		v.Height = ev.Height
	}
}

type NoteNum struct {
	tracker.NoteNum
	View
}

func NewNoteNum(n tracker.NoteNum) *NoteNum {
	return &NoteNum{n, Config.NoteNum.View}
}

func (n *NoteNum) Draw(x, y int) {
	n.Width, n.Height = 0, 0
	s := fmt.Sprintf("%v", n.NoteNum)
	n.Width = len(s)
	n.Height = 1
	for i, r := range s {
		Config.Screen.SetCell(x+i, y, r, n.Fg, n.Bg)
	}
}

type Velocity struct {
	tracker.Velocity
	View
}

func NewVelocity(v tracker.Velocity) *Velocity {
	return &Velocity{v, Config.Velocity.View}
}

func (v *Velocity) Draw(x, y int) {
	v.Width, v.Height = 0, 0
	s := fmt.Sprintf("%v", v.Velocity)
	v.Width = len(s)
	v.Height = 1
	for i, r := range s {
		Config.Screen.SetCell(x+i, y, r, v.Fg, v.Bg)
	}
}

type Screen interface {
	Init() error
	Close()
	SetCell(x, y int, r rune, Fg, Bg termbox.Attribute)
	Flush()
}

type TermScreen struct{}

func (t TermScreen) Init() error { return termbox.Init() }
func (t TermScreen) Close()      { termbox.Close() }
func (t TermScreen) Flush()      { termbox.Flush() }
func (t TermScreen) SetCell(x, y int, r rune, Fg, Bg termbox.Attribute) {
	termbox.SetCell(x, y, r, Fg, Bg)
}

type MockScreen struct {
	cells [][]rune
}

func NewMockScreen(Width, Height int) *MockScreen {
	m := &MockScreen{}
	m.cells = make([][]rune, Height)
	for i, _ := range m.cells {
		m.cells[i] = make([]rune, Width)
	}
	return m
}

func (m *MockScreen) Init() error { return nil }

func (m *MockScreen) Close() {}

func (m *MockScreen) SetCell(x, y int, r rune, Fg, Bg termbox.Attribute) {
	m.cells[y][x] = r
}

func (m *MockScreen) Flush() {
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

func (m *MockScreen) clear() {
	for x := 0; x < len(m.cells); x++ {
		for y := 0; y < len(m.cells[x]); y++ {
			m.cells[x][y] = 0
		}
	}
}
