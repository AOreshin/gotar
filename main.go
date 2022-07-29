package main

import (
	"fmt"
	"math"
	// "os"
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
	// outfile, err := os.Create(time.Now().Format(nameFormat) + ".wav")
	// if err != nil {
	// 	panic(err)
	// }
	// defer outfile.Close()

	// writer := wav.NewWriter(outfile, numSamples, numChannels, sampleRate, bitsPerSample)

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

			// s := toWavSample(l, r)
			// err = writer.WriteSamples(s)
			// if err != nil {
			// 	panic(err)
			// }
		}
		return 0
	}

	err = audio.Open(rtAudioParams(audio), nil, rtaudio.FormatFloat32, uint(sampleRate), buffer, cb, rtAudioOptions())
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

	fmt.Print("\n\n\n\n\n\n\n\n\n")

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
		case '+':
			s.activeFx = []fx{}
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
	for i := 0; i < 9; i++ {
		fmt.Print("\033[A")
		fmt.Print("\033[2K")
	}
	fmt.Printf(`
note %s%d, frequency %.3f
decay factor %.3f
overlap %v
types %v
selected type %v
fx %v
selected fx %v
%d ringing strings
char %c
`,
		n.name,
		n.octave,
		n.frequency,
		s.decay,
		s.overlap,
		s.currentStringTypes,
		s.stringTypes[s.currentStringIndex],
		s.activeFx,
		s.fxTypes[s.currentFxIndex],
		len(s.ringingStrings),
		r,
	)
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
