package main

import "fmt"

type note struct {
	frequency float32
	name      string
	octave    int
}

func (n *note) String() string {
	if n == nil {
		return ""
	}
	return fmt.Sprintf("Name: %s\nOctave: %d\nFrequency: %.3f", n.name, n.octave, n.frequency)
}

var runesToNotes = map[rune]*note{
	'q': {110, "A", 2},
	'2': {116.54, "Bb", 2},
	'w': {123.47, "B", 2},
	'e': {130.81, "С", 3},
	'4': {138.59, "С#", 3},
	'r': {146.83, "D", 3},
	'5': {155.56, "D#", 3},
	't': {164.81, "E", 3},
	'y': {174.61, "F", 3},
	'7': {185.00, "F#", 3},
	'u': {196.00, "G", 3},
	'8': {207.65, "G#", 3},
	'i': {220.00, "A", 3},
	'9': {233.08, "Bb", 3},
	'o': {246.94, "B", 3},
	'p': {261.63, "С", 4},

	'z': {27.50, "A", 0},
	's': {29.14, "Bb", 0},
	'x': {30.87, "B", 0},
	'c': {32.70, "С", 1},
	'f': {34.65, "С#", 1},
	'v': {36.71, "D", 1},
	'g': {38.89, "D#", 1},
	'b': {41.20, "E", 1},
	'n': {43.65, "F", 1},
	'j': {46.25, "F#", 1},
	'm': {49.00, "G", 1},
	'k': {51.91, "G#", 1},
	',': {55.00, "A", 1},
	'l': {58.27, "Bb", 1},
	'.': {61.74, "B", 1},
	'/': {65.41, "С", 2},
}
