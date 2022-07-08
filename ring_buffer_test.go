package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRingBufferZeroCapacity(t *testing.T) {
	r := NewRingBuffer(0)
	assert.True(t, r.IsEmpty())
	assert.True(t, r.IsFull())
	assert.Equal(t, 0, r.Size())
}

func TestRingBuffer(t *testing.T) {
	r := NewRingBuffer(1)
	f := 0.1

	err := r.Enqueue(f)
	assert.NoError(t, err)
	assert.False(t, r.IsEmpty())
	assert.True(t, r.IsFull())
	assert.Equal(t, 1, r.Size())

	v, err := r.Peek()
	assert.NoError(t, err)
	assert.Equal(t, f, v)
	assert.False(t, r.IsEmpty())
	assert.True(t, r.IsFull())
	assert.Equal(t, 1, r.Size())

	v, err = r.Dequeue()
	assert.NoError(t, err)
	assert.Equal(t, f, v)
	assert.True(t, r.IsEmpty())
	assert.False(t, r.IsFull())
	assert.Equal(t, 0, r.Size())

	err = r.Enqueue(f)
	assert.NoError(t, err)
	err = r.Enqueue(f)
	assert.ErrorAs(t, err, &ErrFullQueue)

	_, err = r.Dequeue()
	assert.NoError(t, err)
	_, err = r.Dequeue()
	assert.ErrorAs(t, err, &ErrEmptyQueue)

	_, err = r.Peek()
	assert.ErrorAs(t, err, &ErrEmptyQueue)
}
