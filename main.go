package main

import (
	"fmt"
	"time"

	"github.com/eiannone/keyboard"
	"github.com/yanel/go-rtaudio/src/contrib/go/rtaudio"
)

const (
	AUDIO_BUFFER          = 16
	DEFAULT_DURATION_TICS = 600000
)

func main() {
	audio, err := rtaudio.Create(rtaudio.APIUnspecified)
	if err != nil {
		panic(err)
	}
	defer audio.Destroy()

	printAvailableApis()
	printAvailableDevices(audio)

	params := rtaudio.StreamParams{
		DeviceID:     uint(audio.DefaultOutputDevice()),
		NumChannels:  2,
		FirstChannel: 0,
	}
	options := rtaudio.StreamOptions{
		Flags: rtaudio.FlagsMinimizeLatency,
	}

	duration := DEFAULT_DURATION_TICS
	strings := []*GuitarString{}

	cb := func(out, in rtaudio.Buffer, dur time.Duration, status rtaudio.StreamStatus) int {
		samples := out.Float32()
		for i := 0; i < len(samples)/2; i++ {
			samples[i*2], samples[i*2+1] = stringSamples(strings, duration)
		}
		return 0
	}
	err = audio.Open(&params, nil, rtaudio.FormatFloat32, SAMPLING_RATE, AUDIO_BUFFER, cb, &options)
	if err != nil {
		panic(err)
	}
	defer audio.Close()
	audio.Start()
	defer audio.Stop()

	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()

	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		}
		frequency, ok := notes[char]
		switch key {
		case keyboard.KeyPgdn:
			if duration > 1 {
				duration /= 2
			}
		case keyboard.KeyPgup:
			duration *= 2
		case keyboard.KeyArrowUp:
			for k := range notes {
				notes[k] *= 2
			}
		case keyboard.KeyArrowDown:
			for k := range notes {
				notes[k] /= 2
			}
		case keyboard.KeyEsc:
			return
		default:
			if ok {
				strings = removeDeadStrings(strings, duration)
				strings = append(strings, NewGuitarString(frequency))
			}
		}
		s := fmt.Sprintf("note %s, frequency %.3f, duration = %d tics, %d ringing strings", keysToNotes[char], frequency, duration, len(strings))
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

func stringSamples(strings []*GuitarString, duration int) (float32, float32) {
	var sample float32
	for _, s := range strings {
		if s.Time() > duration {
			continue
		}
		sample += s.Sample()
		s.Tic()
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
