package main

import (
	"fmt"

	"github.com/yanel/go-rtaudio/src/contrib/go/rtaudio"
)

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

func printState(r rune, n *note, s *state) {
	if n == nil {
		n = &note{}
	}
	fmt.Print("\033[u\033[0J")
	fmt.Printf("note \033[1;32m%s%d\033[0m frequency \033[1;32m%.3f\033[0m\r\n",
		n.name, n.octave, n.frequency)
	fmt.Printf("decay factor \033[1;32m%.3f\033[0m\r\n", s.decay)
	fmt.Printf("volume \033[1;32m%.3f\033[0m\r\n", s.volume)
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
	fmt.Printf("recording loop \033[1;32m%v\033[0m\r\n", s.recordLoop)
	fmt.Printf("play recorded loops \033[1;32m%v\033[0m\r\n", s.playLoop)
	fmt.Printf("\033[1;32m%d\033[0m loops playing\r\n", len(s.loops))
	fmt.Printf("char \033[1;32m%c\033[0m\r\n", r)
}
