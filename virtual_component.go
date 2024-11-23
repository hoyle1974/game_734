package main

type VirtualComponent struct {
	x, y   int
	buffer *Buffer
	dirty  Dirty
}

func NewVirtualComponent(width int, height int, dirty Dirty) *VirtualComponent {
	v := &VirtualComponent{
		buffer: NewBuffer(width, height),
		dirty:  dirty,
	}
	v.buffer.DrawBox(5, 5, 20, 20)
	v.buffer.DrawBox(15, 25, 35, 45)
	v.buffer.DrawBox(45, 5, 60, 20)
	v.buffer.DrawBox(65, 5, 70, 20)
	return v
}

func (v *VirtualComponent) Move(dx, dy int) {
	v.x += dx
	v.y += dy
	v.dirty.Dirty()
}

func (v *VirtualComponent) Render(x1, y1, x2, y2 int) *Buffer {
	b := NewBuffer(x2-x1, y2-y1)

	b.CopyFromBuffer(v.x, v.y, v.buffer)

	return b
}
