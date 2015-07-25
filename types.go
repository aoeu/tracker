package tracker

import (
	"fmt"
)

// An Event represents a musical event to occur at a point in time.
type Event struct {
	NoteNum  int
	Velocity int
	*Generator
}

func (e Event) String() string {
	return fmt.Sprintf("%v %v", e.NoteNum, e.Velocity)
}

// A Generator maps to a musical device that can play an Event.
type Generator interface {
	Play(e Event)
	String()
}

// A Track is a series of Events meant to be played sequentially
// through time on one or more Generators.
type Track []Event

// A Line is a set of Events meant to be played concurrently
// at a single moment in time on one or more Generators.
type Line []Event

// A Pattern is a set of Tracks meant to be played
// concurrently through time.
type Pattern []Track

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
	for i, track := range p {
		if len(track) < minLen {
			minLen = len(track)
		}
	}
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
		case len(track) < offset:
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
