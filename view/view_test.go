package view

import (
	"bytes"
	"os"
	"testing"
	"tracker"
)

func testDrawBufferedTrack(t *testing.T) {

	var gen1 = tracker.MockGenerator{}
	var gen2 = tracker.MockGenerator{}
	var gen3 = tracker.MockGenerator{}

	p := tracker.Pattern{
		tracker.Track{
			&tracker.Event{1, 127, gen1},
			&tracker.Event{4, 127, gen1},
			&tracker.Event{4, 127, gen1},
			&tracker.Event{1, 127, gen1},
			&tracker.Event{2, 127, gen1},
			&tracker.Event{4, 127, gen1},
			&tracker.Event{4, 127, gen1},
			&tracker.Event{1, 127, gen1},
		},
		tracker.Track{
			&tracker.Event{0, 127, gen2},
			&tracker.Event{2, 127, gen2},
			&tracker.Event{3, 127, gen2},
		},
		tracker.Track{
			&tracker.Event{127, 127, gen3},
		},
	}
	var b bytes.Buffer
	Config.Screen = NewMockScreen(&b, 36, 8)

	pv := NewPattern(&p)
	pv.DrawBuffered(0, 0)
	Config.Screen.Flush()

	b.WriteTo(os.Stdout)
}

func TestDrawNoteNum(t *testing.T) {
	var n tracker.NoteNum = 1
	var b bytes.Buffer
	nv := NewNoteNum(n)
	Config.Screen = NewMockScreen(&b, nv.maxwidth, 1)
	nv.Draw(0, 0)
	Config.Screen.Flush()
	if s := b.String(); len(s)-len("\n") != nv.maxwidth {
		l := len(s) - len("\n")
		t.Errorf("Expected string of length %v but actual was %v: '%v'", nv.maxwidth, l, s[:l])
	}
}
