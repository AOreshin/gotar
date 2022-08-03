package main

type PeekBuffer struct {
	buffer    []float32
	peekIndex int
}

func (p *PeekBuffer) Append(f float32) {
	p.buffer = append(p.buffer, f)
}

func (p *PeekBuffer) Cut(i int) {
	p.buffer = p.buffer[:i]
}

func (p *PeekBuffer) AtStart() bool {
	return p.peekIndex == 0
}

func (p *PeekBuffer) Peek() float32 {
	f := p.buffer[p.peekIndex]
	p.peekIndex++
	if p.peekIndex == len(p.buffer) {
		p.peekIndex = 0
	}
	return f
}
