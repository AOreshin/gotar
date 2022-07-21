package main

import (
	"math"
	"math/rand"
)

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

type RampAscString struct {
	BaseString
}

func (s *RampAscString) Tic() {
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

func (s *RampAscString) Pluck(frequency, decayFactor float32) VibratingString {
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
	return &RampAscString{
		BaseString: BaseString{
			decayFactor: decayFactor,
			ringBuffer:  r,
		},
	}
}

func (s *RampAscString) String() string {
	return "ramp asc"
}

type RampDescString struct {
	BaseString
}

func (s *RampDescString) Tic() {
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

func (s *RampDescString) Pluck(frequency, decayFactor float32) VibratingString {
	capacity := int(float32(sampleRate) / frequency)
	r := NewRingBuffer(capacity)
	step := 2.0 / float32(capacity)
	v := float32(1) + step
	for i := 0; i < capacity; i++ {
		v -= step
		err := r.Enqueue(v)
		if err != nil {
			panic(err)
		}
	}
	return &RampAscString{
		BaseString: BaseString{
			decayFactor: decayFactor,
			ringBuffer:  r,
		},
	}
}

func (s *RampDescString) String() string {
	return "ramp desc"
}

type SinString struct {
	BaseString
}

func (s *SinString) Tic() {
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

func (s *SinString) Pluck(frequency, decayFactor float32) VibratingString {
	capacity := int(float32(sampleRate) / frequency)
	r := NewRingBuffer(capacity)
	angle := 2 * math.Pi / float32(capacity)
	for i := 0; i < capacity; i++ {
		v := float32(math.Sin(float64(angle) * float64(i)))
		err := r.Enqueue(v)
		if err != nil {
			panic(err)
		}
	}
	return &RampAscString{
		BaseString: BaseString{
			decayFactor: decayFactor,
			ringBuffer:  r,
		},
	}
}

func (s *SinString) String() string {
	return "sin"
}

type SawString struct {
	BaseString
}

func (s *SawString) Tic() {
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

func (s *SawString) Pluck(frequency, decayFactor float32) VibratingString {
	capacity := int(float32(sampleRate) / frequency)
	r := NewRingBuffer(capacity)
	limit := float32(1)
	step := 2 * 4 / float32(capacity)
	v := -limit
	for i := 0; i < capacity; i++ {
		v += step
		if v >= limit {
			v = -limit
		}
		err := r.Enqueue(v)
		if err != nil {
			panic(err)
		}
	}
	return &SawString{
		BaseString: BaseString{
			decayFactor: decayFactor,
			ringBuffer:  r,
		},
	}
}

func (s *SawString) String() string {
	return "saw"
}

type SquareString struct {
	BaseString
}

func (s *SquareString) Tic() {
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

func (s *SquareString) Pluck(frequency, decayFactor float32) VibratingString {
	capacity := int(float32(sampleRate) / frequency)
	r := NewRingBuffer(capacity)
	n := capacity / 2
	v := float32(1)
	for i := 0; i < capacity; i++ {
		if i%n == 0 {
			v = -v
		}
		err := r.Enqueue(v)
		if err != nil {
			panic(err)
		}
	}
	return &SawString{
		BaseString: BaseString{
			decayFactor: decayFactor,
			ringBuffer:  r,
		},
	}
}

func (s *SquareString) String() string {
	return "square"
}

type DoubleRampString struct {
	BaseString
}

func (s *DoubleRampString) Tic() {
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

func (s *DoubleRampString) Pluck(frequency, decayFactor float32) VibratingString {
	capacity := int(float32(sampleRate) / frequency)
	r := NewRingBuffer(capacity)
	limit := float32(1)
	step := 2 * 2 / float32(capacity)
	v := -limit
	for i := 0; i < capacity; i++ {
		v += step
		if v >= limit || v <= -limit {
			step = -step
		}
		err := r.Enqueue(v)
		if err != nil {
			panic(err)
		}
	}
	return &SawString{
		BaseString: BaseString{
			decayFactor: decayFactor,
			ringBuffer:  r,
		},
	}
}

func (s *DoubleRampString) String() string {
	return "triangle"
}
