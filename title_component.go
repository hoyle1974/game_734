package main

type TitleComponent struct {
	buffer *Buffer
	title  string
}

func NewTitleComponent(width int, title string) *TitleComponent {
	return &TitleComponent{
		title:  title,
		buffer: NewBuffer(width-2, 1),
	}
}

func (t *TitleComponent) Render() *Buffer {
	t.buffer.Clear()
	t.buffer.WriteString(t.buffer.width/2-len(t.title)/2, 0, t.title)

	return t.buffer
}
