package main

import "fmt"

type PlayerStatsComponent struct {
	dirty  Dirty
	buffer *Buffer
	stats  map[string]string
}

func NewPlayerStatsComponent(dirty Dirty) *PlayerStatsComponent {
	return &PlayerStatsComponent{
		dirty:  dirty,
		buffer: NewBuffer(32, 16),
		stats:  map[string]string{},
	}
}

func (p *PlayerStatsComponent) Render() *Buffer {
	p.buffer.Clear()
	p.buffer.DrawBoxWithTitle(0, 0, p.buffer.width-1, p.buffer.height-1, "Stats")
	idx := 0
	for k, v := range p.stats {
		p.buffer.WriteString(2, idx+1, fmt.Sprintf("%s : %s", k, v))
		idx++
	}

	return p.buffer
}

func (p *PlayerStatsComponent) SetStat(key string, value string) {
	p.dirty.Dirty()
	p.stats[key] = value
}

func (p *PlayerStatsComponent) DelStats(key string) {
	p.dirty.Dirty()
	delete(p.stats, key)
}
