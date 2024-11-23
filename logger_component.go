package main

import (
	"sync"
	"time"

	"github.com/fatih/color"
)

// AdjustString adjusts a string to a specified length, considering only printable characters.
// ANSI escape codes (e.g., for colors) are ignored when calculating the length.
// func adjust(s string, adjustValue string, length int) string {
// 	// Regular expression to match ANSI escape codes
// 	ansiRegex := regexp.MustCompile(`\x1b\[[0-9;]*m`)

// 	// Remove ANSI codes to calculate the printable length
// 	printable := ansiRegex.ReplaceAllString(s, "")

// 	// Calculate the length of the printable string
// 	printableLength := len(printable)

// 	if printableLength > length {
// 		// Trim the string to the desired length, considering ANSI codes
// 		trimIndex := 0
// 		count := 0
// 		for i, r := range s {
// 			if ansiRegex.MatchString(string(r)) {
// 				continue
// 			}
// 			count++
// 			if count == length {
// 				trimIndex = i + len(string(r))
// 				break
// 			}
// 		}
// 		return s[:trimIndex]
// 	}

// 	// Pad with spaces if necessary
// 	padding := strings.Repeat(adjustValue, length-printableLength)
// 	return s + padding
// }

type Dirty interface {
	Dirty()
}

type LoggerComponent struct {
	lock   sync.RWMutex
	dirty  Dirty
	buffer *Buffer
	lines  []string
}

func NewLoggerComponent(dirty Dirty, width int) *LoggerComponent {
	return &LoggerComponent{dirty: dirty, buffer: NewBuffer(width-2, 7), lines: []string{"", "", "", "", ""}}
}

func (l *LoggerComponent) Log(v string) {
	l.LogColor(v, color.FgWhite)
}

func (l *LoggerComponent) Warn(v string) {
	l.LogColor(v, color.FgHiYellow)
}

func (l *LoggerComponent) Error(v string) {
	l.LogColor(v, color.FgRed)
}

func (l *LoggerComponent) LogColor(v string, value ...color.Attribute) {
	faint := color.New(color.Faint).SprintFunc()
	cv := color.New(value...).SprintFunc()

	now := time.Now()

	// Format the timestamp
	timestamp := now.Format("15:04:05 ")

	l.lock.Lock()
	defer l.lock.Unlock()

	if len(l.lines) >= 5 {
		l.lines = l.lines[1:]
	}
	l.lines = append(l.lines, faint(timestamp)+cv(v))

	l.dirty.Dirty()
}

func (l *LoggerComponent) Render() *Buffer {
	l.lock.RLock()
	defer l.lock.RUnlock()

	l.buffer.Clear()
	l.buffer.DrawBoxWithTitle(0, 0, l.buffer.width-1, l.buffer.height-1, "Logs")
	for y, line := range l.lines {
		l.buffer.WriteString(2, y+1, line)
	}
	return l.buffer
}

func (l *LoggerComponent) View() string {
	l.lock.RLock()
	defer l.lock.RUnlock()

	// faint := color.New(color.Faint).SprintFunc()

	l.buffer.Clear()
	l.buffer.DrawBox(0, 0, l.buffer.width-1, l.buffer.height-1)
	for y, line := range l.lines {
		l.buffer.WriteString(2, y+1, line)
	}
	// l.buffer.WriteString(2, 2, faint("Hello"))

	return l.buffer.String()
}
