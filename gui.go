package tracker

import ("github.com/nsf/termbox-go"
			"fmt"
			)

type dir int

const(UP dir = iota
		DOWN
		LEFT
		RIGHT)

func New() (*Tracker, error) {
	if err := termboxInit(); err != nil {
		return &Tracker{}, err
	}
	// TODO(aoeu): Don't hardcode the pattern name, use flags.
	t, err := NewTracker("cmd/testpattern.trkr")
	if err != nil {
		return &Tracker{}, err
	}
	t.screen.printThings()
	return t, nil 
}


type screen struct {
	fg, bg		termbox.Attribute
	editMode		bool
	cX, cY		int
	redraw chan bool
	currentPattern *Pattern // TODO(aoeu): Do we need this?
	lineOffset int
}

func NewScreen() *screen {
	return  &screen{
		fg: termbox.ColorDefault, 
		bg: termbox.ColorDefault, 
		editMode: false,
		redraw: make(chan bool),
		lineOffset: -1,
	}
}

func (t *Tracker) Run() {
	t.screen.refresh()
	t.UserIn()
}

func (t *Tracker) Exit() {
	termbox.Close()
}

func (s *screen) printThings() {
	s.printEditMode()
	s.drawTable()
	//	s.drawCursor()
}

func (s *screen) drawTable() {
	switch {
	case s.editMode:
		s.drawPattern(5, 5, s.cY, s.cX, *s.currentPattern)
	default:
		s.drawPattern(5, 5, s.cY, -1, *s.currentPattern)
	}
}
/*
func (s *screen) drawCursor() {
	if s.editMode {
		s.bg = termbox.ColorBlue
		s.drawChar(s.cX, s.cY, ' ')
		s.bg = termbox.ColorDefault
	}
}
*/
func (s *screen) moveC(d dir) {
	maxX, maxY :=  len(*s.currentPattern)-1, s.currentPattern.maxTrackLen()-1
	if s.editMode {
		switch d {
		case UP:
			if s.cY > 0 {
				s.cY--
			}
		case DOWN:
			if s.cY < maxY {
				s.cY++
			}
		case RIGHT:
			if s.cX < maxX {
				s.cX++
			}
		case LEFT:
			if s.cX > 0 {
			s.cX--
			}
		}
	}
}

const (
	editMode = "EDIT MODE"
	playbackMode = "Press 'e' to edit."
)

func (s *screen) printEditMode() {
	if s.editMode {
		s.prints(1, 1, editMode)
	} else {
		s.prints(1, 1, playbackMode)
	}
}

func termboxInit() error {
	if err := termbox.Init(); err != nil {
		return err
	}
	termbox.SetInputMode(termbox.InputEsc)
//	termbox.SetOutputMode(termbox.Output256)
	return nil
}

//	Prints text to the screen
func (s screen) prints(x, y int, n interface{}) {
	switch n.(type){
	case int:
		s.drawString(x, y, fmt.Sprintf("%3d", n))
	default:
		s.drawString(x, y, fmt.Sprint(n))
	}
}

func (s *screen) refresh() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	s.printThings()
	termbox.Flush()
}

func (t *Tracker) UserIn() {
	for {
		keyEvents := make(chan termbox.Event)
		go func() {
			for {	
				e := termbox.PollEvent()
				switch e.Type {
				case termbox.EventKey:
					keyEvents <- e
				}
			}
		}()
		select {
		case e := <-keyEvents:
			switch e.Key {
			case termbox.KeyEsc:
				return
			case termbox.KeyArrowUp:
				t.screen.moveC(UP)
			case termbox.KeyArrowDown:
				t.screen.moveC(DOWN)
			case termbox.KeyArrowRight:
				t.screen.moveC(RIGHT)
			case termbox.KeyArrowLeft:
				t.screen.moveC(LEFT)
			}
			switch e.Ch {
			case 'e':
				if t.isPlaying {
					t.TogglePlayback()
				}
				if t.screen.editMode {
					t.screen.editMode = false
					t.screen.refresh()
				} else {
					t.screen.editMode = true
					// TODO(aoeu): t.Stop() ? 
				}
			case 'p':
				t.screen.editMode = false
				t.screen.refresh()
				t.TogglePlayback()
			}
			t.screen.refresh()
		case <-t.screen.redraw:
			t.screen.refresh()
		}
	}
}

func (s screen) drawString(x, y int, str string) {
	for i, r := range str {
		s.drawChar(x + i, y, r)
	}
}

func (s screen) drawChar(x, y int, r rune) {
	termbox.SetCell(x, y, r, s.fg, s.bg)
}

func (s *screen) drawPattern(x, y, hr, hc int, p Pattern) {
	defer func() {
		s.fg = termbox.ColorBlue
		s.bg = termbox.ColorDefault
	}()
	for i, l := range p.GetLines() {
		s.fg = termbox.ColorBlue
		s.bg = termbox.ColorDefault
		s.prints(x, y + 3 + i, i)
		for z, e := range l {
			if i == hr && z == hc || i == s.lineOffset {
				s.bg = termbox.ColorRed
			} else {
				s.bg = termbox.ColorDefault
			}
			s.fg = termbox.ColorDefault
			s.prints(x + (z * 7) + 3, y + 3 + i, e.NoteNum)
			s.fg = termbox.ColorGreen
			s.prints(3 + x + (z * 7) + 3, y + 3 + i, e.Velocity)
		}
	}
	//	Function that calls function drawEvent to draw an event next to the previous event (from left to write)
	//	Note, generator, effect, parameter
}


//	Windows
//	-Draw to windows
//	-Dynamically adds coordinate offsets based on window position
//	-Record table consisting of window IDs
//	-Draw in window like normally drawing in termbox buffer

//	NOTESTSEOTNSET
//	Struct variable indicating "theme:
//	Functions use current "Theme" when printing/drawing

