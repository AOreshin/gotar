package main

import (
	"fmt"
	"os"
	"time"

	"github.com/eiannone/keyboard"
	"github.com/youpy/go-wav"
)

func handleInput() {
	fmt.Print("\n\033[s")

	for {
		r, k, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		}

		if k == keyboard.KeyEsc {
			return
		}

		handleStrings(r, k)
		changeOctaves(k)
		handleFx(r)
		handleRecord(r)
		handleVolume(r)
		handleLoops(k)

		note, ok := runesToNotes[r]
		if ok {
			pluckString(note)
		}

		printState(r, note)
	}
}

func handleVolume(r rune) {
	switch r {
	case '@':
		vState.volume += 0.01
	case '#':
		vState.volume -= 0.01
	}
}

func handleFx(r rune) {
	switch r {
	case '[':
		fState.currentFxIndex--
		if fState.currentFxIndex < 0 {
			fState.currentFxIndex = len(fState.fxTypes) - 1
		}
	case ']':
		fState.currentFxIndex++
		if fState.currentFxIndex == len(fState.fxTypes) {
			fState.currentFxIndex = 0
		}
	case '-':
		fState.activeFx = append(fState.activeFx, fState.fxTypes[fState.currentFxIndex])
	case '=':
		fState.activeFx = []fx{}
	}
}

func handleRecord(r rune) {
	switch r {
	case '\\':
		rState.record = !rState.record
		if rState.record {
			name := time.Now().Format(nameFormat) + ".wav"
			outfile, err := os.Create(name)
			if err != nil {
				panic(err)
			}
			rState.file = outfile
			defer outfile.Close()
			rState.writer = wav.NewWriter(outfile, numSamples, numChannels, sampleRate, bitsPerSample)
		} else {
			rState.file.Close()
		}
	}
}

func handleLoops(k keyboard.Key) {
	switch k {
	case keyboard.KeyTab:
		lState.playLoop = !lState.playLoop
	case keyboard.KeyHome:
		lState.recordLoop = !lState.recordLoop
		if lState.recordLoop {
			lState.loop = [2]*PeekBuffer{{}, {}}
		} else {
			lState.loops = append(lState.loops, lState.loop)
		}
	case keyboard.KeyEnd:
		if len(lState.loops) > 0 {
			lState.loops = lState.loops[:len(lState.loops)-1]
		}
	}
}

func changeOctaves(k keyboard.Key) {
	switch k {
	case keyboard.KeyArrowUp:
		for k := range runesToNotes {
			note := runesToNotes[k]
			note.frequency *= 2
			note.octave++
		}
	case keyboard.KeyArrowDown:
		for k := range runesToNotes {
			note := runesToNotes[k]
			note.frequency /= 2
			note.octave--
		}
	}
}

func handleStrings(r rune, k keyboard.Key) {
	switch r {
	case ';':
		sState.currentStringTypes =
			append(sState.currentStringTypes, sState.stringTypes[sState.currentStringIndex])
	case '\'':
		sState.currentStringTypes =
			[]VibratingString{sState.stringTypes[sState.currentStringIndex]}
	}

	switch k {
	case keyboard.KeyArrowLeft:
		sState.currentStringIndex--
		if sState.currentStringIndex < 0 {
			sState.currentStringIndex = len(sState.stringTypes) - 1
		}
	case keyboard.KeyArrowRight:
		sState.currentStringIndex++
		if sState.currentStringIndex == len(sState.stringTypes) {
			sState.currentStringIndex = 0
		}
	case keyboard.KeySpace:
		sState.overlap = !sState.overlap
	case keyboard.KeyPgdn:
		sState.decay -= 0.001
	case keyboard.KeyPgup:
		sState.decay += 0.001
	case keyboard.KeyEnter:
		sState.ringingStrings = []VibratingString{}
	}
}

func pluckString(n *note) {
	sState.ringingStrings = removeDeadStrings(sState.ringingStrings, defaultDuration)
	if sState.overlap {
		for _, strType := range sState.currentStringTypes {
			sState.ringingStrings = append(sState.ringingStrings,
				strType.Pluck(n.frequency, sState.decay))
		}
	} else {
		newStrings := []VibratingString{}
		for _, strType := range sState.currentStringTypes {
			newStrings = append(newStrings,
				strType.Pluck(n.frequency, sState.decay))
		}
		sState.ringingStrings = newStrings
	}
}

func removeDeadStrings(strings []VibratingString, duration int) []VibratingString {
	ringingStrings := []VibratingString{}
	for _, s := range strings {
		if s.Time() > duration {
			continue
		}
		ringingStrings = append(ringingStrings, s)
	}
	return ringingStrings
}
