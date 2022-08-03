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
			g.currentStringTypes =
				append(g.currentStringTypes, g.stringTypes[g.currentStringIndex])
		case '\'':
			g.currentStringTypes =
				[]VibratingString{g.stringTypes[g.currentStringIndex]}
		case '[':
			g.currentFxIndex--
			if g.currentFxIndex < 0 {
				g.currentFxIndex = len(g.fxTypes) - 1
			}
		case ']':
			g.currentFxIndex++
			if g.currentFxIndex == len(g.fxTypes) {
				g.currentFxIndex = 0
			}
		case '-':
			g.activeFx = append(g.activeFx, g.fxTypes[g.currentFxIndex])
		case '=':
			g.activeFx = []fx{}
		case '\\':
			g.record = !g.record
			if g.record {
				name := time.Now().Format(nameFormat) + ".wav"
				outfile, err := os.Create(name)
				if err != nil {
					panic(err)
				}
				g.file = outfile
				defer outfile.Close()
				g.writer = wav.NewWriter(outfile, numSamples, numChannels, sampleRate, bitsPerSample)
			} else {
				g.file.Close()
			}
		case '@':
			g.volume += 0.01
		case '#':
			g.volume -= 0.01
		}
		switch key {
		case keyboard.KeyTab:
			g.playLoop = !g.playLoop
		case keyboard.KeyHome:
			g.recordLoop = !g.recordLoop
			if g.recordLoop {
				g.loop = [2]*PeekBuffer{{}, {}}
			} else {
				g.loops = append(g.loops, g.loop)
			}
		case keyboard.KeyEnd:
			if len(g.loops) > 0 {
				g.loops = g.loops[:len(g.loops)-1]
			}
		case keyboard.KeyArrowLeft:
			g.currentStringIndex--
			if g.currentStringIndex < 0 {
				g.currentStringIndex = len(g.stringTypes) - 1
			}
		case keyboard.KeyArrowRight:
			g.currentStringIndex++
			if g.currentStringIndex == len(g.stringTypes) {
				g.currentStringIndex = 0
			}
		case keyboard.KeySpace:
			g.overlap = !g.overlap
		case keyboard.KeyPgdn:
			g.decay -= 0.001
		case keyboard.KeyPgup:
			g.decay += 0.001
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
			g.ringingStrings = []VibratingString{}
		case keyboard.KeyEsc:
			return
		default:
			if ok {
				g.ringingStrings = removeDeadStrings(g.ringingStrings, defaultDuration)
				if g.overlap {
					for _, strType := range g.currentStringTypes {
						g.ringingStrings = append(g.ringingStrings,
							strType.Pluck(note.frequency, g.decay))
					}
				} else {
					newStrings := []VibratingString{}
					for _, strType := range g.currentStringTypes {
						newStrings = append(newStrings,
							strType.Pluck(note.frequency, g.decay))
					}
					g.ringingStrings = newStrings
				}
			}
		}
		printState(char, note, g)
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
