package main

import (
	"os"

	"github.com/youpy/go-wav"
)

type stringsState struct {
	stringTypes        []VibratingString
	currentStringTypes []VibratingString
	ringingStrings     []VibratingString
	decay              float32
	overlap            bool
	currentStringIndex int
}

type fxState struct {
	fxTypes        []fx
	activeFx       []fx
	currentFxIndex int
}

type recordState struct {
	record bool
	writer *wav.Writer
	file   *os.File
}

type loopState struct {
	recordLoop bool
	playLoop   bool
	loop       [2]*PeekBuffer
	loops      [][2]*PeekBuffer
}

type volumeState struct {
	volume float32
}

var (
	sState = &stringsState{
		stringTypes: []VibratingString{
			&GuitarString{},
			&RampAscString{},
			&RampDescString{},
			&SinString{},
			&SawString{},
			&SquareString{},
			&DoubleRampString{},
			&DrumString{},
		},
		currentStringTypes: []VibratingString{&GuitarString{}},
		currentStringIndex: 0,
		overlap:            true,
		decay:              decayFactor,
	}
	fState = &fxState{
		fxTypes: []fx{
			outOfPhaseFx,
			vibrato,
		},
		activeFx:       []fx{},
		currentFxIndex: 0,
	}
	rState = &recordState{
		record: false,
	}
	vState = &volumeState{
		volume: 1.0,
	}
	lState = &loopState{
		playLoop: true,
	}
)
