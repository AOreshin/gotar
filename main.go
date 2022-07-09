package main

import (
	"github.com/eiannone/keyboard"
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
)

const (
	AUDIO_BUFFER = 2048
	MAX_TIME     = 600000
)

var notes = map[rune]float64{
	'q': 110,
	'2': 116.54,
	'w': 123.47,
	'e': 130.81,
	'4': 138.59,
	'r': 146.83,
	'5': 155.56,
	't': 164.81,
	'y': 174.61,
	'7': 185.00,
	'u': 196.00,
	'8': 207.65,
	'i': 220.00,
	'9': 233.08,
	'o': 246.94,
	'p': 261.63,
	'-': 277.18,
	'[': 293.66,
	'=': 311.13,

	'z': 27.50,
	'x': 29.14,
	'd': 30.87,
	'c': 32.70,
	'f': 34.65,
	'v': 36.71,
	'g': 38.89,
	'b': 41.20,
	'n': 43.65,
	'j': 46.25,
	'm': 49.00,
	'k': 51.91,
	',': 55.00,
	'.': 58.27,
	';': 61.74,
	'/': 65.41,
}

func main() {
	speaker.Init(SAMPLING_RATE, AUDIO_BUFFER)
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
			str := NewGuitarString(frequency)
			speaker.Play(stringStreamer(str))
		}
		if key == keyboard.KeyEsc {
			break
		}
	}
}

func stringStreamer(str *GuitarString) beep.Streamer {
	return beep.StreamerFunc(func(samples [][2]float64) (n int, ok bool) {
		for i := 0; i < len(samples); i++ {
			s := str.Sample()
			s = softDistortion(s)
			samples[i][0] = s
			samples[i][1] = s
			str.Tic()
		}
		if str.Time() > MAX_TIME {
			return 0, false
		}
		return len(samples), true
	})
}

func heavyDistortion(s float64) float64 {
	if s > 0.01 {
		s = 0.2
	}
	if s < -0.01 {
		s = -0.2
	}
	return s
}

func softDistortion(s float64) float64 {
	if s > 0.01 {
		s = 0.1
	}
	if s < -0.01 {
		s = -0.1
	}
	return s
}
