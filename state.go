package main

type state struct {
	stringTypes        []VibratingString
	currentStringTypes []VibratingString
	ringingStrings     []VibratingString
	currentStringIndex int
	fxTypes            []fx
	activeFx           []fx
	currentFxIndex     int
	overlap            bool
	decay              float32
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
	decay:          decayFactor,
}
