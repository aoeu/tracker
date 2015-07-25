package tracker

import ("github.com/nsf/termbox-go"
			"time"
			"fmt")

func main() {
	termboxInit()
	defer termbox.Close()
	prints(5, 5, "Hello World!")
	prints(5, 6, 123)
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

func drawEvent(x, y int, e event) {
}

func drawPatterns(x, y int, currentLine int) {
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
