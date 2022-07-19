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

	strings := []*GuitarString{}

	cb := func(out, in rtaudio.Buffer, dur time.Duration, status rtaudio.StreamStatus) int {
		samples := out.Float32()
		for i := 0; i < len(samples)/2; i++ {
			l, r := stringSamples(strings)
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

	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		}
		frequency, ok := notes[char]
		switch key {
		case keyboard.KeySpace:
			overlap = !overlap
		case keyboard.KeyPgdn:
			decay -= 0.001
		case keyboard.KeyPgup:
			decay += 0.001
		case keyboard.KeyArrowUp:
			for k := range notes {
				notes[k] *= 2
			}
		case keyboard.KeyArrowDown:
			for k := range notes {
				notes[k] /= 2
			}
		case keyboard.KeyEnter:
			strings = []*GuitarString{}
		case keyboard.KeyEsc:
			return
		default:
			if ok {
				strings = removeDeadStrings(strings, defaultDuration)
				if overlap {
					strings = append(strings, NewGuitarString(frequency, decay))
				} else {
					strings = []*GuitarString{NewGuitarString(frequency, decay)}
				}
			}
		}
		s := fmt.Sprintf("note %s, frequency %.3f, decay factor %.3f, overlap %v, %d ringing strings, char %c",
			keysToNotes[char], frequency, decay, overlap, len(strings), char)
		fmt.Printf("\r%s", s)
	}
}

func removeDeadStrings(strings []*GuitarString, duration int) []*GuitarString {
	ringingStrings := []*GuitarString{}
	for _, s := range strings {
		if s.Time() > duration {
			continue
		}
		ringingStrings = append(ringingStrings, s)
	}
	return ringingStrings
}

func stringSamples(strings []*GuitarString) (float32, float32) {
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
	return sample, sample
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
