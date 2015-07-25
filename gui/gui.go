package gui

import ("github.com/nsf/termbox-go"
			"time"
			"fmt"
			. "tracker")

func New() {
	termboxInit()
	defer termbox.Close()

	drawPattern(2, 2, testPattern)
	refresh()

	time.Sleep(2*time.Second)

}

func refresh() {
	termbox.Flush()
}

func termboxInit() {
	termbox.Init()
	termbox.SetInputMode(termbox.InputEsc)
	termbox.SetOutputMode(termbox.Output256)
}

//	Prints text to the screen
func prints(x, y int, n interface{}) {
	drawString(x, y, fmt.Sprint(n))
}

func drawString(x, y int, s string) {
	for i, r := range s {
		drawChar(x + i, y, r)
	}
}

func drawChar(x, y int, r rune) {
	termbox.SetCell(x, y, r, termbox.ColorDefault, termbox.ColorDefault)
}
/*
func drawEvent(x, y int, e event) {
	prints(
}
*/
func drawPattern(x, y int, p Pattern) {
	l := p.GetLines()
	for i, v := range l {
		for z, e := range v {
			prints(x + (z * 10), y + 3 + i, e.NoteNum)
			prints(x + (z * 10), y + 3 + i, e.Velocity)
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



var gen1 = Generator{}
var gen2 = Generator{}
var gen3 = Generator{}

var testPattern = Pattern{
        []Track{
                Event{1, 1, gen1},
                Event{2, 1, gen1},
                Event{3, 1, gen1},
                Event{4, 1,  gen1},
                Event{5, 1, gen1},
        },
        []Track {
                Event{64, 127, gen2},
                Event{60, 127, gen2},
               Event{67, 127, gen2},
        },
        []Track {
                Event{127, 127, gen3},
        },
}
