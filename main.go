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
			initialState.currentStringTypes =
				append(initialState.currentStringTypes, initialState.stringTypes[initialState.currentStringIndex])
		case '\'':
			initialState.currentStringTypes =
				[]VibratingString{initialState.stringTypes[initialState.currentStringIndex]}
		case '[':
			initialState.currentFxIndex--
			if initialState.currentFxIndex < 0 {
				initialState.currentFxIndex = len(initialState.fxTypes) - 1
			}
		case ']':
			initialState.currentFxIndex++
			if initialState.currentFxIndex == len(initialState.fxTypes) {
				initialState.currentFxIndex = 0
			}
		case '-':
			initialState.activeFx = append(initialState.activeFx, initialState.fxTypes[initialState.currentFxIndex])
		case '=':
			initialState.activeFx = []fx{}
		case '\\':
			initialState.record = !initialState.record
			if initialState.record {
				name := time.Now().Format(nameFormat) + ".wav"
				outfile, err := os.Create(name)
				if err != nil {
					panic(err)
				}
				initialState.file = outfile
				defer outfile.Close()
				initialState.writer = wav.NewWriter(outfile, numSamples, numChannels, sampleRate, bitsPerSample)
			} else {
				initialState.file.Close()
			}
		case '@':
			initialState.volume += 0.01
		case '#':
			initialState.volume -= 0.01
		}
		switch key {
		case keyboard.KeyHome:
			initialState.recordLoop = !initialState.recordLoop
			if initialState.recordLoop {
				initialState.loop = [2]*PeekBuffer{{}, {}}
			} else {
				initialState.recordLoopStarted = false
				if len(initialState.loops) != 0 {
					baseLoopLen := len(initialState.loops[0])
					newLoopLen := len(initialState.loop)
					if newLoopLen > baseLoopLen {
						initialState.loop[0].Cut(baseLoopLen)
						initialState.loop[1].Cut(baseLoopLen)
						initialState.loops = append(initialState.loops, initialState.loop)
					} else {
						for i := 0; i < baseLoopLen-newLoopLen; i++ {
							initialState.loop[0].Append(0.0)
							initialState.loop[1].Append(0.0)
						}
					}
				}
				initialState.loops = append(initialState.loops, initialState.loop)
			}
		case keyboard.KeyEnd:
			if len(initialState.loops) > 0 {
				initialState.loops = initialState.loops[:len(initialState.loops)-1]
			}
		case keyboard.KeyArrowLeft:
			initialState.currentStringIndex--
			if initialState.currentStringIndex < 0 {
				initialState.currentStringIndex = len(initialState.stringTypes) - 1
			}
		case keyboard.KeyArrowRight:
			initialState.currentStringIndex++
			if initialState.currentStringIndex == len(initialState.stringTypes) {
				initialState.currentStringIndex = 0
			}
		case keyboard.KeySpace:
			initialState.overlap = !initialState.overlap
		case keyboard.KeyPgdn:
			initialState.decay -= 0.001
		case keyboard.KeyPgup:
			initialState.decay += 0.001
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
			initialState.ringingStrings = []VibratingString{}
		case keyboard.KeyEsc:
			return
		default:
			if ok {
				initialState.ringingStrings = removeDeadStrings(initialState.ringingStrings, defaultDuration)
				if initialState.overlap {
					for _, strType := range initialState.currentStringTypes {
						initialState.ringingStrings = append(initialState.ringingStrings,
							strType.Pluck(note.frequency, initialState.decay))
					}
				} else {
					newStrings := []VibratingString{}
					for _, strType := range initialState.currentStringTypes {
						newStrings = append(newStrings,
							strType.Pluck(note.frequency, initialState.decay))
					}
					initialState.ringingStrings = newStrings
				}
			}
		}
		printState(char, note, initialState)
	}
}

func callback(out, in rtaudio.Buffer, dur time.Duration, status rtaudio.StreamStatus) int {
	samples := out.Float32()
	for i := 0; i < len(samples)/2; i++ {
		l, r := stringSamples(
			initialState.ringingStrings,
			initialState.activeFx,
		)

		v := initialState.volume
		l, r = l*v, r*v

		if initialState.recordLoop {
			if len(initialState.loops) == 0 {
				initialState.recordLoopStarted = true
			} else {
				if initialState.loops[0][0].AtStart() {
					initialState.recordLoopStarted = true
				}
			}
		}

		if initialState.recordLoopStarted {
			initialState.loop[0].Append(l)
			initialState.loop[1].Append(r)
		}

		for i := 0; i < len(initialState.loops); i++ {
			loop := initialState.loops[i]
			l += loop[0].Peek()
			r += loop[1].Peek()
		}

		samples[i*2], samples[i*2+1] = l, r

		if initialState.record && initialState.writer != nil {
			s := toWavSample(l, r)
			err := initialState.writer.WriteSamples(s)
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
