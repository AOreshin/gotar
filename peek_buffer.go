package main

type PeekBuffer struct {
	buffer    []float32
	peekIndex int
}

func (p *PeekBuffer) Append(f float32) {
	p.buffer = append(p.buffer, f)
}

func (p *PeekBuffer) Tic() int {
	i := p.peekIndex
	p.peekIndex++
	if p.peekIndex == len(p.buffer) {
		p.peekIndex = 0
	}
	return i
}

func (p *PeekBuffer) Get(i int) float32 {
	if i >= len(p.buffer) {
		return 0.0
	}
	return p.buffer[i]
}
