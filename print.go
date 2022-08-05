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

func printState(r rune, n *note) {
	if n == nil {
		n = &note{}
	}
	fmt.Print("\033[u\033[0J")
	fmt.Printf("note \033[1;32m%s%d\033[0m frequency \033[1;32m%.3f\033[0m\r\n",
		n.name, n.octave, n.frequency)
	fmt.Printf("overlap \033[1;32m%v\033[0m\r\n", sState.overlap)
	fmt.Printf("decay factor \033[1;32m%.3f\033[0m\r\n", sState.decay)
	fmt.Printf("\033[1;32m%d\033[0m ringing strings\r\n", len(sState.ringingStrings))
	fmt.Printf("types \033[1;32m%v\033[0m\r\n", sState.currentStringTypes)
	fmt.Printf("selected type \033[1;32m%v\033[0m\r\n", sState.stringTypes[sState.currentStringIndex])
	fmt.Printf("bypass fx \033[1;32m%v\033[0m\r\n", fState.bypassFx)
	fmt.Printf("fx type \033[1;32m%v\033[0m\r\n", fState.activeFx)
	fmt.Printf("selected fx \033[1;32m%v\033[0m\r\n", fState.fxTypes[fState.currentFxIndex])
	fmt.Printf("volume \033[1;32m%.3f\033[0m\r\n", vState.volume)
	fmt.Printf("record \033[1;32m%v\033[0m\r\n", rState.record)
	if rState.record {
		fmt.Printf("writing to \033[1;32m%v\033[0m\r\n", rState.file.Name())
	}
	fmt.Printf("recording loop \033[1;32m%v\033[0m\r\n", lState.recordLoop)
	fmt.Printf("play loops \033[1;32m%v\033[0m\r\n", lState.playLoop)
	fmt.Printf("\033[1;32m%d\033[0m loops recorded\r\n", len(lState.loops))
	fmt.Printf("char \033[1;32m%c\033[0m\r\n", r)
}
