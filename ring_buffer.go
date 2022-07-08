package main

import "errors"

var (
	ErrFullQueue  = errors.New("full queue")
	ErrEmptyQueue = errors.New("empty queue")
)

type RingBuffer struct {
	first, last, capacity, size int
	buffer                      []float64
}

func (r *RingBuffer) Size() int {
	return r.size
}

func (r *RingBuffer) IsEmpty() bool {
	return r.size == 0
}

func (r *RingBuffer) IsFull() bool {
	return r.size == r.capacity
}

func (r *RingBuffer) Enqueue(f float64) error {
	if r.IsFull() {
		return ErrFullQueue
	}
	r.size++
	r.buffer[r.last] = f
	r.last++
	if r.last == r.capacity {
		r.last = 0
	}
	return nil
}

func (r *RingBuffer) Dequeue() (float64, error) {
	if r.IsEmpty() {
		return 0, ErrEmptyQueue
	}
	r.size--
	f := r.buffer[r.first]
	r.first++
	if r.first == r.capacity {
		r.first = 0
	}
	return f, nil
}

func (r *RingBuffer) Peek() (float64, error) {
	if r.IsEmpty() {
		return 0, ErrEmptyQueue
	}
	return r.buffer[r.first], nil
}

func NewRingBuffer(capacity int) *RingBuffer {
	return &RingBuffer{
		capacity: capacity,
		buffer:   make([]float64, capacity),
	}
}
