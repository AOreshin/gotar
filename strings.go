package main

import "math/rand"

type VibratingString interface {
	Pluck(frequency, decayFactor float32) VibratingString
	Tic()
	Sample() float32
	Time() int
}

type BaseString struct {
	decayFactor float32
	ringBuffer  *RingBuffer
	tics        int
}

func (b *BaseString) Sample() float32 {
	v, err := b.ringBuffer.Peek()
	if err != nil {
		panic(err)
	}
	return v
}

func (b *BaseString) Time() int {
	return b.tics
}

type GuitarString struct {
	BaseString
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
	v := g.decayFactor * (first + second)
	g.ringBuffer.Enqueue(v)
}

func (g *GuitarString) Pluck(frequency, decayFactor float32) VibratingString {
	capacity := int(float32(sampleRate) / frequency)
	r := NewRingBuffer(capacity)
	for i := 0; i < capacity; i++ {
		v := rand.Float32() - 0.5
		err := r.Enqueue(v)
		if err != nil {
			panic(err)
		}
	}
	return &GuitarString{
		BaseString: BaseString{
			decayFactor: decayFactor,
			ringBuffer:  r,
		},
	}
}

func (g *GuitarString) String() string {
	return "guitar"
}

type SynthString struct {
	BaseString
}

func (s *SynthString) Tic() {
	s.tics++
	first, err := s.ringBuffer.Dequeue()
	if err != nil {
		panic(err)
	}
	second, err := s.ringBuffer.Peek()
	if err != nil {
		panic(err)
	}
	v := s.decayFactor * (first + second)
	s.ringBuffer.Enqueue(v)
}

func (s *SynthString) Pluck(frequency, decayFactor float32) VibratingString {
	capacity := int(float32(sampleRate) / frequency)
	r := NewRingBuffer(capacity)
	step := 2.0 / float32(capacity)
	v := float32(-1) - step
	for i := 0; i < capacity; i++ {
		v += step
		err := r.Enqueue(v)
		if err != nil {
			panic(err)
		}
	}
	return &SynthString{
		BaseString: BaseString{
			decayFactor: decayFactor,
			ringBuffer:  r,
		},
	}
}

func (s *SynthString) String() string {
	return "synth"
}
