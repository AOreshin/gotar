package main

import (
	"os"

	"github.com/youpy/go-wav"
)

type state struct {
	stringTypes        []VibratingString
	currentStringTypes []VibratingString
	ringingStrings     []VibratingString
	currentStringIndex int
	fxTypes            []fx
	activeFx           []fx
	currentFxIndex     int
	overlap            bool
	record             bool
	writer             *wav.Writer
	file               *os.File
	decay              float32
	recordLoop         bool
	recordLoopStarted  bool
	loop               [2]*PeekBuffer
	loops              [][2]*PeekBuffer
	volume             float32
}

var initialState = &state{
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
	fxTypes: []fx{
		outOfPhaseFx,
		vibrato,
	},
	activeFx:       []fx{},
	currentFxIndex: 0,
	overlap:        true,
	record:         false,
	decay:          decayFactor,
	volume:         1.0,
}
