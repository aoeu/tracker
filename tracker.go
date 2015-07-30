package tracker

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"io/ioutil"
	"time"
)

// A mockGenerator is only intended for testing or debugging.
type MockGenerator struct{}

func (m MockGenerator) Play(e Event)   {}
func (m MockGenerator) String() string { return "Mock generator." }
func (e Event) String() string {
	return fmt.Sprintf("%v %v", e.NoteNum, e.Velocity)
}

func NewTrack(g Generator, velocity int, notes ...int) Track {
	t := make(Track, len(notes))
	for i := 0; i < len(t); i++ {
		t[i] = Event{Generator: g, NoteNum: notes[i], Velocity: velocity}
	}
	return t
}

func NewPattern(filePath string) (*Pattern, error) {
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return &Pattern{}, err
	}
	// TODO(aoeu): Is there a programmatic way of registering the types?
	gob.Register(Pattern{})
	gob.Register(MockGenerator{})
	dec := gob.NewDecoder(bytes.NewReader(b))
	pat := Pattern{}
	err = dec.Decode(&pat)
	if err != nil && err != io.EOF {
		return &Pattern{}, err
	}
	return &pat, nil
}

func (p Pattern) maxTrackLen() int {
	maxLen := 0
	for _, track := range p {
		if len(track) > maxLen {
			maxLen = len(track)
		}
	}
	return maxLen
}

func (p Pattern) minTrackLen() int {
	// TODO(aoeu): Is this actually needed?
	minLen := int(^uint(0) >> 1)
	for _, track := range p {
		if len(track) < minLen {
			minLen = len(track)
		}
	}
	return minLen
}

// GetLine returns a Line containing the Events
// associated with all of the Tracks at a given
// offest in a Pattern.
//
// If a Pattern is thought of as a "table", and
// Tracks are thought of as "columns", GetLine
// returns "row" of the table.
func (p Pattern) GetLine(offset int) Line {
	l := make(Line, len(p))
	for i, track := range p {
		switch {
		case len(track) > offset:
			l[i] = track[offset]
		default:
			l[i] = Event{}
		}
	}
	return l
}

// GetLines returns a series of Line types
// containing the Events associated with all of
// the Tracks in a Pattern.
//
// Any Track that is shorter in length than others
// in the pattern is still still represented in a
// respective Line with empty Event values (as padding).
func (p Pattern) GetLines() []Line {
	maxTrackLen := p.maxTrackLen()
	l := make([]Line, maxTrackLen)
	for i := range l {
		l[i] = p.GetLine(i)
	}
	return l
}

func NewPlayer(filepath string) (*Player, error) {
	p := Player{}
	p.PatternTable = make(PatternTable, 1)
	if pattern, err := NewPattern(filepath); err != nil {
		return &Player{}, err
	} else {
		p.PatternTable[0] = pattern
	}
	p.BPM = 120 // TODO(aoeu): Don't hardcode the BPM.
	return &p, nil

}

func (t *Tracker) TogglePlayback() {
	if t.isPlaying {
		t.Stop()
	} else {
		go t.Play()
	}
}

func (t *Tracker) Stop() {
	t.stop <- true
}

func (t *Tracker) Play() {
	t.isPlaying = true
	defer func() { 
		t.isPlaying = false 
		t.screen.lineOffset= -1
		t.screen.redraw <- true
	}()
	nsPerBeat := 60000000000 / t.Player.BPM
	for _, pattern := range t.Player.PatternTable {
		for _, line := range pattern.GetLines() {
			t.screen.lineOffset += 1
			t.screen.redraw <- true
			for _, e := range line {
				if e.Generator != nil {
					go e.Generator.Play(e) // TODO(aoeu): Reconsider ownership of Events and Generators.
				}
			}
			select {
			case <-time.After(time.Duration(nsPerBeat) * time.Nanosecond):
			case <-t.stop:
				return
			}
		}
	}

}
