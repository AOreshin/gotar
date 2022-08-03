package main

import (
	"errors"
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

func callback(out, in rtaudio.Buffer, dur time.Duration, status rtaudio.StreamStatus) int {
	samples := out.Float32()
	for i := 0; i < len(samples)/2; i++ {
		l, r := stringSamples(
			g.ringingStrings,
			g.activeFx,
		)

		v := g.volume
		l, r = l*v, r*v

		if g.recordLoop {
			g.loop[0].Append(l)
			g.loop[1].Append(r)
		}

		if len(g.loops) > 0 {
			baseLoopIndex := g.loops[0][0].Tic()
			g.loops[0][1].Tic()
			for i := 0; i < len(g.loops); i++ {
				loop := g.loops[i]
				l += loop[0].Get(baseLoopIndex)
				r += loop[1].Get(baseLoopIndex)
			}
		}

		samples[i*2], samples[i*2+1] = l, r

		if g.record && g.writer != nil {
			s := toWavSample(l, r)
			err := g.writer.WriteSamples(s)
			if err != nil {
				if errors.Is(err, os.ErrClosed) {
					return 0
				}
				panic(err)
			}
		}
	}
	return 0
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

func stringSamples(strings []VibratingString, fxs []fx) (float32, float32) {
	var sample float32
	for _, s := range strings {
		sample += s.Sample() * 0.25
		s.Tic()
	}
	if sample > 1 {
		sample = 1
	}
	if sample < -1 {
		sample = -1
	}
	l, r := sample, sample
	for _, f := range fxs {
		l, r = f.apply(l, r)
	}
	return l, r
}

func toWavSample(l, r float32) []wav.Sample {
	return []wav.Sample{
		{Values: [2]int{int(l * floatToInt), int(r * floatToInt)}},
	}
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
