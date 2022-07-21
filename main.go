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
)

func main() {
	outfile, err := os.Create(time.Now().Format(nameFormat) + ".wav")
	if err != nil {
		panic(err)
	}
	defer outfile.Close()

	writer := wav.NewWriter(outfile, numSamples, numChannels, sampleRate, bitsPerSample)

	audio, err := rtaudio.Create(rtaudio.APIUnspecified)
	if err != nil {
		panic(err)
	}
	defer audio.Destroy()

	printAvailableApis()
	printAvailableDevices(audio)

	strings := []VibratingString{}
	fxs := []fx{}

	cb := func(out, in rtaudio.Buffer, dur time.Duration, status rtaudio.StreamStatus) int {
		samples := out.Float32()
		for i := 0; i < len(samples)/2; i++ {
			l, r := stringSamples(strings, fxs)
			samples[i*2], samples[i*2+1] = l, r

			s := toWavSample(l, r)
			err = writer.WriteSamples(s)
			if err != nil {
				panic(err)
			}
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

	decay := decayFactor
	overlap := true
	outOfPhase := false
	softDistortion := false
	stringTypes := []VibratingString{
		&GuitarString{},
		&RampAscString{},
		&RampDescString{},
		&SinString{},
		&SawString{},
		&SquareString{},
		&DoubleRampString{},
	}
	currentStringType := VibratingString(&GuitarString{})
	currentStringIndex := 0

	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		}
		note, ok := keysToNotes[char]
		frequency, name, octave := float32(0.0), "", 0
		if note != nil {
			frequency, name, octave = note.frequency, note.name, note.octave
		}
		switch key {
		case keyboard.KeyArrowLeft:
			currentStringIndex--
			if currentStringIndex < 0 {
				currentStringIndex = len(stringTypes) - 1
			}
			currentStringType = stringTypes[currentStringIndex]
		case keyboard.KeyArrowRight:
			currentStringIndex++
			if currentStringIndex == len(stringTypes) {
				currentStringIndex = 0
			}
			currentStringType = stringTypes[currentStringIndex]
		case keyboard.KeyHome:
			outOfPhase, fxs = switchFx(outOfPhase, outOfPhaseFx, fxs)
		case keyboard.KeyEnd:
			softDistortion, fxs = switchFx(softDistortion, softDistortionFx, fxs)
		case keyboard.KeySpace:
			overlap = !overlap
		case keyboard.KeyPgdn:
			decay -= 0.001
		case keyboard.KeyPgup:
			decay += 0.001
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
			strings = []VibratingString{}
		case keyboard.KeyEsc:
			return
		default:
			if ok {
				strings = removeDeadStrings(strings, defaultDuration)
				if overlap {
					strings = append(strings, currentStringType.Pluck(frequency, decay))
				} else {
					strings = []VibratingString{currentStringType.Pluck(frequency, decay)}
				}
			}
		}
		s := fmt.Sprintf("note %s%d, frequency %.3f, decay factor %.3f, overlap %v, type %v, fx %v, %d ringing strings, char %c",
			name, octave, frequency, decay, overlap, currentStringType, fxs, len(strings), char)
		fmt.Printf("\r%s", s)
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

func stringSamples(strings []VibratingString, fxs []fx) (float32, float32) {
	var sample float32
	for _, s := range strings {
		sample += s.Sample()
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
		{Values: [2]int{int(l * math.MaxInt32), int(r * math.MaxInt32)}},
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

func removeFx(fxs []fx, f fx) []fx {
	for i, v := range fxs {
		if f.name == v.name {
			return append(fxs[:i], fxs[i+1:]...)
		}
	}
	return fxs
}

func switchFx(flag bool, f fx, fxs []fx) (bool, []fx) {
	flag = !flag
	if flag {
		fxs = append(fxs, f)
	} else {
		fxs = removeFx(fxs, f)
	}
	return flag, fxs
}
