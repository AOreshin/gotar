package main

import (
	"fmt"
	"math"
	"os"

	"time"

	"github.com/eiannone/keyboard"
	"github.com/yanel/go-rtaudio/src/contrib/go/rtaudio"
	"github.com/youpy/go-wav"
)

const (
	buffer                 = 16
	defaultDuration        = 600000
	numSamples      uint32 = math.MaxUint32
	numChannels     uint16 = 2
	firstChannel    uint   = 0
	sampleRate      uint32 = 44100
	bitsPerSample   uint16 = 32
	nameFormat             = "2006-02-01 15-04-05"
	decayFactor            = float32(0.994 * 0.5)
	floatToInt             = math.MaxInt32 / 4
)

func main() {
	audio, err := rtaudio.Create(rtaudio.APIUnspecified)
	if err != nil {
		panic(err)
	}
	defer audio.Destroy()

	printAvailableApis()
	printAvailableDevices(audio)

	err = audio.Open(rtAudioParams(audio), nil, rtaudio.FormatFloat32,
		uint(sampleRate), buffer, callback, rtAudioOptions())
	if err != nil {
		panic(err)
	}
	defer audio.Close()

	audio.Start()
	defer audio.Stop()

	err = keyboard.Open()
	if err != nil {
		panic(err)
	}
	defer keyboard.Close()

	handleInput()
}

func handleInput() {
	fmt.Print("\n\033[s")

	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		}
		note, ok := keysToNotes[char]
		switch char {
		case ';':
			sState.currentStringTypes =
				append(sState.currentStringTypes, sState.stringTypes[sState.currentStringIndex])
		case '\'':
			sState.currentStringTypes =
				[]VibratingString{sState.stringTypes[sState.currentStringIndex]}
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
		case '@':
			vState.volume += 0.01
		case '#':
			vState.volume -= 0.01
		}
		switch key {
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
		case keyboard.KeyArrowUp:
			for k := range keysToNotes {
				note := keysToNotes[k]
				note.frequency *= 2
				note.octave++
			}
		case keyboard.KeyArrowDown:
			for k := range keysToNotes {
				note := keysToNotes[k]
				note.frequency /= 2
				note.octave--
			}
		case keyboard.KeyEnter:
			sState.ringingStrings = []VibratingString{}
		case keyboard.KeyEsc:
			return
		default:
			if ok {
				sState.ringingStrings = removeDeadStrings(sState.ringingStrings, defaultDuration)
				if sState.overlap {
					for _, strType := range sState.currentStringTypes {
						sState.ringingStrings = append(sState.ringingStrings,
							strType.Pluck(note.frequency, sState.decay))
					}
				} else {
					newStrings := []VibratingString{}
					for _, strType := range sState.currentStringTypes {
						newStrings = append(newStrings,
							strType.Pluck(note.frequency, sState.decay))
					}
					sState.ringingStrings = newStrings
				}
			}
		}
		printState(char, note)
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

func rtAudioParams(audio rtaudio.RtAudio) *rtaudio.StreamParams {
	return &rtaudio.StreamParams{
		DeviceID:     uint(audio.DefaultOutputDevice()),
		NumChannels:  uint(numChannels),
		FirstChannel: uint(firstChannel),
	}
}

func rtAudioOptions() *rtaudio.StreamOptions {
	return &rtaudio.StreamOptions{
		Flags: rtaudio.FlagsMinimizeLatency,
	}
}
