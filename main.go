package main

import (
	"fmt"
	"time"

	"github.com/eiannone/keyboard"
	"github.com/yanel/go-rtaudio/src/contrib/go/rtaudio"
)

const (
	AUDIO_BUFFER = 16
	MAX_TIME     = 600000
)

var ()

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
	strings := []*GuitarString{}
	cb := func(out, in rtaudio.Buffer, dur time.Duration, status rtaudio.StreamStatus) int {
		samples := out.Float32()
		for i := 0; i < len(samples)/2; i++ {
			s := stringSamples(strings)
			samples[i*2] = s
			samples[i*2+1] = s
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
		if ok {
			strings = removeDeadStrings(strings)
			strings = append(strings, NewGuitarString(frequency))
			s := fmt.Sprintf("%d polyphonic strings", len(strings))
			fmt.Printf("\r%s", s)
		}
		if key == keyboard.KeyEsc {
			break
		}
	}
}

func removeDeadStrings(strings []*GuitarString) []*GuitarString {
	ringingStrings := []*GuitarString{}
	for _, s := range strings {
		if s.Time() > MAX_TIME {
			continue
		}
		ringingStrings = append(ringingStrings, s)
	}
	return ringingStrings
}

func stringSamples(strings []*GuitarString) float32 {
	sample := 0.
	for _, str := range strings {
		s := str.Sample()
		s = softDistortion(s)
		sample += s
		str.Tic()
	}
	return float32(sample)
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
