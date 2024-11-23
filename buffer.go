package main

import (
	"strings"
	"unicode/utf8"
)

func CountColumns(s string) int {
	columnCount := 0
	for len(s) > 0 {
		// Decode the next rune
		r, size := utf8.DecodeRuneInString(s)
		if r == utf8.RuneError {
			// If an invalid rune is encountered, break out
			break
		}

		// Determine if the rune is a wide character (typically emoji or CJK)
		if isWideRune(r) {
			// Wide characters take 2 columns
			columnCount += 2
		} else {
			// Regular characters take 1 column
			columnCount++
		}

		// Move past the decoded rune
		s = s[size:]
	}
	return columnCount
}

// isWideRune checks if a rune is considered "wide" (takes up 2 columns visually).
func isWideRune(r rune) bool {
	// Wide characters typically fall in certain ranges like CJK or emojis
	return (r >= 0x4E00 && r <= 0x9FFF) || // CJK characters
		(r >= 0x3040 && r <= 0x30FF) || // Hiragana/Katakana
		(r >= 0x1F600 && r <= 0x1F64F) || // Emoticons
		(r >= 0x1F300 && r <= 0x1F5FF) || // Miscellaneous Symbols and Pictographs
		(r >= 0x1F680 && r <= 0x1F6FF) || // Transport and Map Symbols
		(r >= 0x1F700 && r <= 0x1F77F) || // Alchemical Symbols
		(r >= 0x1F780 && r <= 0x1F7FF) || // Geometric Shapes Extended
		(r >= 0x1F800 && r <= 0x1F8FF) || // Supplemental Arrows-C
		(r >= 0x1F900 && r <= 0x1F9FF) || // Supplemental Symbols and Pictographs
		(r >= 0x1FA00 && r <= 0x1FA6F) || // Chess Symbols
		(r >= 0x1FA70 && r <= 0x1FAFF) // Symbols and Pictographs Extended-A
}

type Display interface {
	String() string
	Columns() int
}

type EscapeSequence struct {
	value string
}

func (e EscapeSequence) String() string { return e.value }
func (e EscapeSequence) Columns() int   { return 0 }

type Printable struct {
	value   string
	columns int
}

func (e Printable) String() string { return e.value }
func (e Printable) Columns() int   { return e.columns }

func splitByPrintables(s string) []Display {
	var result []Display
	var current string
	var inAnsiCode bool

	for _, r := range s {
		if r == '\x1b' {
			// If we're already in an ANSI code, append the current color code to the result
			if inAnsiCode {
				result = append(result, EscapeSequence{value: "\x1b" + current})
				current = ""
			}
			inAnsiCode = true
		} else if inAnsiCode {
			// Append the current rune to the color code
			current += string(r)
			if r == 'm' {
				// End of ANSI code, append it to the result
				result = append(result, EscapeSequence{value: "\x1b" + current})
				current = ""
				inAnsiCode = false
			}
		} else {
			// Not in an ANSI code, append the rune to the result directly
			result = append(result, Printable{value: string(r), columns: CountColumns(string(r))})
		}
	}

	// Handle any remaining characters
	if current != "" {
		result = append(result, Printable{value: current, columns: CountColumns(current)})
	}

	return result
}

type BufferCell struct {
	escapeCode string
	printable  Printable
}

type Buffer struct {
	width  int
	height int
	data   [][]BufferCell
}

// NewBuffer creates a new buffer with given width and height, initially cleared.
func NewBuffer(width, height int) *Buffer {
	b := &Buffer{
		width:  width,
		height: height,
	}
	b.Clear()

	return b
}

func (b *Buffer) set(x, y int, v string) {
	// Ensure the position is within bounds and the string fits within the buffer
	// Find where to set it

	if x >= 0 && x < b.width && y >= 0 && y < b.height {
		printables := splitByPrintables(v)

		// Handle and keep the escape codes we are copying over
		idx := 0
		escapeCodes := ""
		for _, p := range printables {
			switch p.(type) {
			case EscapeSequence:
			case Printable:
				if x+idx >= 0 && x+idx < b.width {
					escapeCodes += b.data[y][x+idx].escapeCode
					b.data[y][x+idx].escapeCode = ""
					idx++
				}
			default:
				panic("Nope")
			}
		}
		if x+idx >= 0 && x+idx < b.width {
			b.data[y][x+idx].escapeCode = escapeCodes + b.data[y][x+idx].escapeCode
		}

		idx = 0
		for _, p := range printables {
			switch v := p.(type) {
			case EscapeSequence:
				if x+idx >= 0 && x+idx < b.width {
					b.data[y][x+idx].escapeCode += v.value
				}
			case Printable:
				if x+idx >= 0 && x+idx < b.width {
					b.data[y][x+idx].printable = v
					idx++
				}
			default:
				panic("Nope")
			}
		}
	}
}

// String returns the buffer as a string, preserving emojis and ANSI codes.
func (b *Buffer) String() string {
	var builder strings.Builder
	for _, row := range b.data {
		for x := 0; x < len(row); {
			cell := row[x]
			builder.WriteString(cell.escapeCode + cell.printable.String())
			x += cell.printable.Columns()
		}
		builder.WriteString("\n")
	}
	return builder.String()
}

// Lines returns the buffer as an array of strings, one per line.
func (b *Buffer) Lines() []string {
	lines := make([]string, b.height)
	// for i := 0; i < b.height; i++ {
	// 	lines[i] = strings.Join(b.data[i], "")
	// }
	return lines
}

// WriteString writes a string to the buffer starting at position (x, y), trimming as needed.
func (b *Buffer) WriteString(x, y int, s string) {
	b.set(x, y, s)
}

// WriteBuffer writes the contents of another buffer onto this buffer starting at (x, y), clipping as needed.
func (b *Buffer) WriteBuffer(x, y int, other *Buffer) {
	// fmt.Println(other.height, other.width)
	for i := 0; i < other.height; i++ {
		for j := 0; j < other.width; j++ {
			if x+j < b.width && y+i < b.height {
				b.set(x+j, y+i, other.data[i][j].escapeCode+other.data[i][j].printable.String())
			}
		}
	}
}

func (b *Buffer) CopyFromBuffer(x, y int, other *Buffer) {
	for i := y; i < y+b.height; i++ {
		for j := x; j < x+b.width; j++ {
			xx := j - x
			yy := i - y
			if xx >= 0 && yy >= 0 && xx <= b.width-1 && yy <= b.height-1 {
				b.set(xx, yy, other.data[i][j].escapeCode+other.data[i][j].printable.String())
			}
		}
	}
}

// Clear clears the buffer, filling it with spaces.
func (b *Buffer) Clear() {
	data := make([][]BufferCell, b.height)
	for y := range data {
		temp := make([]BufferCell, b.width)
		for x := range temp {
			temp[x].printable = Printable{value: " ", columns: 1}
		}
		data[y] = temp
	}
	b.data = data
}

func (b *Buffer) DrawBoxWithTitle(x1, y1, x2, y2 int, title string) {
	b.DrawBox(x1, y1, x2, y2)
	b.WriteString(x1+1, y1, "┤"+title+"├")
}

// DrawBox draws a box using the Unicode box-drawing characters.
func (b *Buffer) DrawBox(x1, y1, x2, y2 int) {
	if x1 >= 0 && x1 < b.width && y1 >= 0 && y1 < b.height {
		b.set(x1, y1, "┌")
		b.set(x2, y1, "┐")
		b.set(x1, y2, "└")
		b.set(x2, y2, "┘")

		for y := y1 + 1; y < y2; y++ {
			if x1 >= 0 && x1 < b.width {
				b.set(x1, y, "│")
			}
			if x2 >= 0 && x2 < b.width {
				b.set(x2, y, "│")
			}
		}

		for x := x1 + 1; x < x2; x++ {
			if y1 >= 0 && y1 < b.height {
				b.set(x, y1, "─")
			}
			if y2 >= 0 && y2 < b.height {
				b.set(x, y2, "─")
			}
		}
	}
}
