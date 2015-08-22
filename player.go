package tracker

import (
	"time"
)

type Player struct {
	BPM int
	PatternTable
	Clock     chan time.Time
	isPlaying bool
	stop      chan bool
}

func InitPlayer(p PatternTable) *Player { //TODO(aoeu): Rename to NewPlayer
	return &Player{
		BPM:          120,
		PatternTable: p,
		Clock:        make(chan time.Time, 1),
		stop:         make(chan bool),
	}
}

func (p *Player) TogglePlayback() {
	if p.isPlaying {
		p.Stop()
	} else {
		go p.Play()
	}
}

func (p *Player) Stop() {
	p.stop <- true
}

func (p *Player) Play() {
	p.isPlaying = true
	defer func() { p.isPlaying = false }()
	beatLen := p.nsPerBeat()
	for _, pattern := range p.PatternTable {
		for _, line := range pattern.GetLines() {
			p.Clock <- time.Now()
			for _, e := range line {
				if e.Generator != nil {
					// TODO(aoeu): Reconsider ownership of Events and Generators.
					go e.Generator.Play(*e)
				}
			}
			select {
			case <-time.After(time.Duration(beatLen) * time.Nanosecond):
			case <-p.stop:
				return
			}
		}
	}
}

func (p *Player) nsPerBeat() int {
	return 60000000000 / p.BPM
}
