package tracker

import (
	"fmt"
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
