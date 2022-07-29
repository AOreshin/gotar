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

	cb := func(out, in rtaudio.Buffer, dur time.Duration, status rtaudio.StreamStatus) int {
		samples := out.Float32()
		for i := 0; i < len(samples)/2; i++ {
			l, r := stringSamples(initialState.ringingStrings, initialState.activeFx)
			samples[i*2], samples[i*2+1] = l, r

			if initialState.record && initialState.writer != nil {
				s := toWavSample(l, r)
				err = initialState.writer.WriteSamples(s)
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

	err = audio.Open(rtAudioParams(audio), nil, rtaudio.FormatFloat32,
		uint(sampleRate), buffer, cb, rtAudioOptions())
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

	fmt.Print("\n\033[s")

	inputHandler(initialState)
}

func inputHandler(s *state) {
	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		}
		note, ok := keysToNotes[char]
		switch char {
		case ';':
			s.currentStringTypes =
				append(s.currentStringTypes, s.stringTypes[s.currentStringIndex])
		case '\'':
			s.currentStringTypes =
				[]VibratingString{s.stringTypes[s.currentStringIndex]}
		case '[':
			s.currentFxIndex--
			if s.currentFxIndex < 0 {
				s.currentFxIndex = len(s.fxTypes) - 1
			}
		case ']':
			s.currentFxIndex++
			if s.currentFxIndex == len(s.fxTypes) {
				s.currentFxIndex = 0
			}
		case '-':
			s.activeFx = append(s.activeFx, s.fxTypes[s.currentFxIndex])
		case '=':
			s.activeFx = []fx{}
		case '\\':
			s.record = !s.record
			if s.record {
				name := time.Now().Format(nameFormat) + ".wav"
				outfile, err := os.Create(name)
				if err != nil {
					panic(err)
				}
				s.file = outfile
				defer outfile.Close()
				s.writer = wav.NewWriter(outfile, numSamples, numChannels, sampleRate, bitsPerSample)
			} else {
				s.file.Close()
			}
		}
		switch key {
		case keyboard.KeyArrowLeft:
			s.currentStringIndex--
			if s.currentStringIndex < 0 {
				s.currentStringIndex = len(s.stringTypes) - 1
			}
		case keyboard.KeyArrowRight:
			s.currentStringIndex++
			if s.currentStringIndex == len(s.stringTypes) {
				s.currentStringIndex = 0
			}
		case keyboard.KeySpace:
			s.overlap = !s.overlap
		case keyboard.KeyPgdn:
			s.decay -= 0.001
		case keyboard.KeyPgup:
			s.decay += 0.001
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
			s.ringingStrings = []VibratingString{}
		case keyboard.KeyEsc:
			return
		default:
			if ok {
				s.ringingStrings = removeDeadStrings(s.ringingStrings, defaultDuration)
				if s.overlap {
					for _, strType := range s.currentStringTypes {
						s.ringingStrings = append(s.ringingStrings,
							strType.Pluck(note.frequency, s.decay))
					}
				} else {
					newStrings := []VibratingString{}
					for _, strType := range s.currentStringTypes {
						newStrings = append(newStrings,
							strType.Pluck(note.frequency, s.decay))
					}
					s.ringingStrings = newStrings
				}
			}
		}
		printState(char, note, s)
	}
}

func printState(r rune, n *note, s *state) {
	if n == nil {
		n = &note{}
	}
	fmt.Print("\033[u")
	fmt.Printf("note \033[1;32m%s%d\033[0m frequency \033[1;32m%.3f\033[0m\r\n",
		n.name, n.octave, n.frequency)
	fmt.Printf("decay factor \033[1;32m%.3f\033[0m\r\n", s.decay)
	fmt.Printf("record \033[1;32m%v\033[0m\r\n", s.record)
	if s.record {
		fmt.Printf("writing to \033[1;32m%v\033[0m\r\n", s.file.Name())
	}
	fmt.Printf("overlap \033[1;32m%v\033[0m\r\n", s.overlap)
	fmt.Printf("types \033[1;32m%v\033[0m\r\n", s.currentStringTypes)
	fmt.Printf("selected type \033[1;32m%v\033[0m\r\n", s.stringTypes[s.currentStringIndex])
	fmt.Printf("fx type \033[1;32m%v\033[0m\r\n", s.activeFx)
	fmt.Printf("selected fx \033[1;32m%v\033[0m\r\n", s.fxTypes[s.currentFxIndex])
	fmt.Printf("\033[1;32m%d\033[0m ringing strings\r\n", len(s.ringingStrings))
	fmt.Printf("char \033[1;32m%c\033[0m\r\n", r)
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

func printAvailableApis() {
	fmt.Println("RtAudio version: ", rtaudio.Version())
	for _, api := range rtaudio.CompiledAPI() {
		fmt.Println("Compiled API: ", api)
	}
}

func printAvailableDevices(audio rtaudio.RtAudio) {
	devices, err := audio.Devices()
	if err != nil {
		panic(err)
	}
	for _, d := range devices {
		fmt.Printf("Audio device: %#v\n", d)
	}
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
