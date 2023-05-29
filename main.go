package main

import (
	"math"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/yanel/go-rtaudio/src/contrib/go/rtaudio"
)

const (
	buffer                 = 16
	defaultDuration        = 600000
	numSamples      uint32 = math.MaxUint32
	numChannels     uint16 = 2
	firstChannel    uint   = 0
	sampleRate      uint32 = 48000
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

	app := tview.NewApplication()

	note := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetChangedFunc(func() {
			app.Draw()
		})
	note.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		handleKey(event)
		if runesToNotes[event.Rune()] != nil {
			note.SetText(runesToNotes[event.Rune()].String() + "\nRune: " + string(event.Rune()))
		}
		return nil
	}).
		SetBorder(true).
		SetTitle("Last note pressed")

	// flex := tview.NewFlex().
	// 	AddItem(note, 0, 1, false).
	// 	SetBorder(true).
	// 	SetTitle("gotar").
	// 	SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
	// 		handleKey(event)
	// 		note.SetText(event.Name())
	// 		return nil
	// 	})

	err = app.SetRoot(note, true).Run()
	if err != nil {
		panic(err)
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
