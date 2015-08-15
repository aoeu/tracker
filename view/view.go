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
			Fg:       termbox.ColorBlue,
			Bg:       termbox.ColorDefault,
			maxwidth: 4,
		}},

		Velocity{View: View{
			Fg: termbox.ColorCyan,
			Bg: termbox.ColorDefault,
			maxwidth: 4,
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
	width, height int
	Fg, Bg        termbox.Attribute
	delimiter     string
	maxwidth      int
}

func (v View) Width() int {
	return v.width
}

func (v View) Height() int {
	return v.height
}
type Pattern struct {
	*tracker.Pattern
	View
}

func NewPattern(p *tracker.Pattern) *Pattern {
	return &Pattern{p, Config.Pattern.View}
}

func (pv *Pattern) Draw(x, y int) {
	pv.width, pv.height = 0, 0
	for _, t := range *pv.Pattern {
		tv := NewTrack(t)
		tv.Draw(x+pv.width, y)
		pv.width += tv.width
		if tv.height > pv.height {
			pv.height = tv.height
		}
	}
}

func (pv *Pattern) DrawBuffered(x, y int) {
	maxTrackLen := 0
	for _, t := range *pv.Pattern {
		if len(t) > maxTrackLen {
			maxTrackLen = len(t)
		}
	}
	for _, t := range *pv.Pattern {
		tv := NewTrack(t)
		tv.DrawBuffered(x+pv.width, y, maxTrackLen)
		pv.width += tv.width
		if tv.height > pv.height {
			pv.height = tv.height
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
	tv.width, tv.height = 0, 0
	for _, e := range tv.Track {
		ev := NewEvent(e)
		ev.Draw(x, y+tv.height)
		if ev.width > tv.width {
			tv.width = ev.width
		}
		for i, r := range tv.delimiter {
			Config.Screen.SetCell(x+tv.width+i, y+tv.height, r, tv.Fg, tv.Bg)
		}
		tv.height += ev.height
	}
	tv.width += len(tv.delimiter)
}

func (tv *Track) DrawBuffered(x, y, maxTrackLen int) {
	tv.Draw(x, y)
	switch l := len(tv.Track); {
	case l > maxTrackLen:
		panic("Track was drawn and has a length longer than the provided maximum.")
	case l == maxTrackLen:
		return
	}
	tv.width -= len(tv.delimiter)
	for i := 0; i < maxTrackLen - len(tv.Track); i++ {
		ev := NewEvent(&tracker.Event{NoteNum: 77, Velocity: 777})
		ev.Draw(x, y+tv.height)
		if ev.width > tv.width {
			tv.width = ev.width
		}
		for i, r := range tv.delimiter {
			Config.Screen.SetCell(x+tv.width+i, y+tv.height, r, tv.Fg, tv.Bg)
		}
		tv.height += ev.height
	}
	tv.width += len(tv.delimiter)
}



type Line struct {
	tracker.Line
	View
}

func NewLine(l tracker.Line) *Line {
	return &Line{l, Config.Line.View}
}

func (lv *Line) Draw(x, y int) {
	lv.width, lv.height = 0, 0
	for _, e := range lv.Line {
		ev := NewEvent(e)
		ev.Draw(x+lv.width, y)
		lv.width += ev.width
		if ev.height > lv.height {
			lv.height = ev.height
		}
		for i, r := range lv.delimiter {
			termbox.SetCell(x+lv.width+i, y, r, lv.Fg, lv.Bg)
		}
		lv.width += len(lv.delimiter)
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
	ev.width, ev.height = 0, 0
	n := NewNoteNum(ev.NoteNum)
	n.Draw(x+ev.width, y)
	ev.width += n.width
	if n.height > ev.height {
		ev.height = n.height
	}
	for i, r := range ev.delimiter {
		termbox.SetCell(x+ev.width+i, y, r, ev.Fg, ev.Bg)
		ev.width += i
	}
	v := NewVelocity(ev.Velocity)
	v.Draw(x+ev.width, y)
	ev.width += v.width
	if v.height > ev.height {
		v.height = ev.height
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
	n.width, n.height = 0, 0
	s := fmt.Sprintf("%v", n.NoteNum)
	for i := 0; i < n.maxwidth-len(s); i++ {
		s = " " + s
	}
	n.width = len(s)
	n.height = 1
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
	v.width, v.height = 0, 0
	s := fmt.Sprintf("%v", v.Velocity)
	for i := 0; i < v.maxwidth-len(s); i++ {
		s = " " + s
	}
	v.width = len(s)
	v.height = 1
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

func NewMockScreen(width, height int) *MockScreen {
	m := &MockScreen{}
	m.cells = make([][]rune, height)
	for i, _ := range m.cells {
		m.cells[i] = make([]rune, width)
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
