package main

import (
	"fmt"
	"sort"
)

type PlayerStatsComponent struct {
	dirty  Dirty
	buffer *Buffer
	stats  map[string]string
}

func NewPlayerStatsComponent(dirty Dirty, height int) *PlayerStatsComponent {
	return &PlayerStatsComponent{
		dirty:  dirty,
		buffer: NewBuffer(32, height),
		stats:  map[string]string{},
	}
}

func (p *PlayerStatsComponent) Render() *Buffer {
	p.buffer.Clear()
	p.buffer.DrawBoxWithTitle(0, 0, p.buffer.width-1, p.buffer.height-1, "Stats")
	idx := 0

	// Extract keys from the map
	keys := make([]string, 0, len(p.stats))
	for key := range p.stats {
		keys = append(keys, key)
	}

	// Sort the keys
	sort.Strings(keys)

	// Print the map in alphabetical key order
	for _, key := range keys {
		p.buffer.WriteString(2, idx+1, fmt.Sprintf("%s : %s", key, p.stats[key]))
		idx++
	}

	return p.buffer
}

func (p *PlayerStatsComponent) SetStat(key string, value string) {
	// p.dirty.Dirty()
	p.stats[key] = value
}

func (p *PlayerStatsComponent) DelStats(key string) {
	// p.dirty.Dirty()
	delete(p.stats, key)
}
