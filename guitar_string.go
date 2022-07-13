package main

import "math/rand"

const (
	SAMPLING_RATE = 44100
	DECAY_FACTOR  = 0.994 * 0.5
)

type GuitarString struct {
	ringBuffer *RingBuffer
	tics       int
}

func (g *GuitarString) Tic() {
	g.tics++
	first, err := g.ringBuffer.Dequeue()
	if err != nil {
		panic(err)
	}
	second, err := g.ringBuffer.Peek()
	if err != nil {
		panic(err)
	}
	v := DECAY_FACTOR * (first + second)
	g.ringBuffer.Enqueue(v)
}

func (g *GuitarString) Sample() float32 {
	v, err := g.ringBuffer.Peek()
	if err != nil {
		panic(err)
	}
	return v
}

func (g *GuitarString) Time() int {
	return g.tics
}

func NewGuitarString(frequency float32) *GuitarString {
	capacity := int(SAMPLING_RATE / frequency)
	r := NewRingBuffer(capacity)
	for i := 0; i < capacity; i++ {
		v := rand.Float32() - 0.5
		err := r.Enqueue(v)
		if err != nil {
			panic(err)
		}
	}
	return &GuitarString{
		ringBuffer: r,
	}
}
