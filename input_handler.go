package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/eiannone/keyboard"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/youpy/go-wav"
)

func handleKeys(event *tcell.EventKey, view *tview.TextView) {
	handleStrings(event)

	r := event.Rune()
	note, ok := runesToNotes[r]
	if !ok {
		return
	}
	pluckString(note)

	view.SetText(runesToNotes[event.Rune()].String() + "\nRune: " + string(event.Rune()))
}

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

		// handleStrings(k)
		// handleOctave / s(k)
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
	case '0':
		fState.bypassFx = !fState.bypassFx
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

func handleOctaves(event *tcell.EventKey, view *tview.TextView) {
	currentOctaveStr := view.GetText(true)
	currentOctave, err := strconv.Atoi(currentOctaveStr)
	if err != nil {
		panic(err)
	}

	switch event.Key() {
	case tcell.KeyUp:
		for k := range runesToNotes {
			note := runesToNotes[k]
			note.frequency *= 2
			note.octave++
		}
		currentOctave++
	case tcell.KeyDown:
		for k := range runesToNotes {
			note := runesToNotes[k]
			note.frequency /= 2
			note.octave--
		}
		currentOctave--
	}

	view.SetText(strconv.Itoa(currentOctave))
}

func handlePolyphonic(event *tcell.EventKey, view *tview.TextView) {
	if event.Key() == tcell.KeyEnter {
		sState.polyphonic = !sState.polyphonic
	}

	view.SetText(strconv.FormatBool(sState.polyphonic))
}

func handleStrings(key *tcell.EventKey) {
	switch key.Rune() {
	case ';':
		sState.currentStringTypes =
			append(sState.currentStringTypes, sState.stringTypes[sState.currentStringIndex])
	case '\'':
		sState.currentStringTypes =
			[]VibratingString{sState.stringTypes[sState.currentStringIndex]}
	}

	switch key.Key() {
	case tcell.KeyLeft:
		sState.currentStringIndex--
		if sState.currentStringIndex < 0 {
			sState.currentStringIndex = len(sState.stringTypes) - 1
		}
	case tcell.KeyRight:
		sState.currentStringIndex++
		if sState.currentStringIndex == len(sState.stringTypes) {
			sState.currentStringIndex = 0
		}
	case tcell.KeyPgDn:
		sState.decay -= 0.001
	case tcell.KeyPgUp:
		sState.decay += 0.001
		// case tcell.KeyEnter:
		// 	sState.ringingStrings = []VibratingString{}
	}
}

func pluckString(n *note) {
	sState.ringingStrings = removeDeadStrings(sState.ringingStrings, defaultDuration)
	if sState.polyphonic {
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
