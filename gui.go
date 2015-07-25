package main

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

func drawString(x, y int, s string) {
	for i, r := range s {
		termbox.SetCell(x + i, y, r, termbox.ColorDefault, termbox.ColorDefault)
	}
}

func prints(x, y int, n interface{}) {
	drawString(x, y, fmt.Sprint(n))
}
