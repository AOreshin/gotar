package main

type note struct {
	frequency float32
	name      string
	octave    int
}

var notes = map[rune]float32{
	'q': 110,
	'2': 116.54,
	'w': 123.47,
	'e': 130.81,
	'4': 138.59,
	'r': 146.83,
	'5': 155.56,
	't': 164.81,
	'y': 174.61,
	'7': 185.00,
	'u': 196.00,
	'8': 207.65,
	'i': 220.00,
	'9': 233.08,
	'o': 246.94,
	'p': 261.63,

	'z': 27.50,
	's': 29.14,
	'x': 30.87,
	'c': 32.70,
	'f': 34.65,
	'v': 36.71,
	'g': 38.89,
	'b': 41.20,
	'n': 43.65,
	'j': 46.25,
	'm': 49.00,
	'k': 51.91,
	',': 55.00,
	'l': 58.27,
	'.': 61.74,
	'/': 65.41,
}

var keysToNotes = map[rune]string{
	'q': "A",
	'2': "Bb",
	'w': "B",
	'e': "C",
	'4': "C#",
	'r': "D",
	'5': "D#",
	't': "E",
	'y': "F",
	'7': "F#",
	'u': "G",
	'8': "G#",
	'i': "A",
	'9': "Bb",
	'o': "B",
	'p': "C",

	'z': "A",
	's': "Bb",
	'x': "B",
	'c': "C",
	'f': "C#",
	'v': "D",
	'g': "D#",
	'b': "E",
	'n': "F",
	'j': "F#",
	'm': "G",
	'k': "G#",
	',': "A",
	'l': "Bb",
	'.': "B",
	'/': "C",
}
