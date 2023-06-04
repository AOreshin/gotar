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

	note := tview.NewTextView()
	note.SetTextAlign(tview.AlignCenter).
		SetBorder(true).
		SetTitle("Last note pressed")

	octave := tview.NewTextView()
	octave.SetText("0").
		SetTextAlign(tview.AlignCenter).
		SetBorder(true).
		SetTitle("Octave shift")

	header := tview.NewTextView()
	header.SetBorder(true).
		SetTitle("gotar").
		SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			handleKeys(event, note)
			handleOctaves(event, octave)
			return nil
		})

	flex := tview.NewFlex().
		AddItem(header, 0, 1, true).
		AddItem(note, 0, 1, false).
		AddItem(octave, 0, 1, false)

	err = tview.NewApplication().SetRoot(flex, true).Run()
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
