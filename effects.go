package main

import "math"

type fx struct {
	name  string
	apply func(l, r float32) (float32, float32)
}

func (f fx) String() string {
	return f.name
}

var (
	outOfPhaseFx = fx{
		name: "out of phase",
		apply: func(l, r float32) (float32, float32) {
			return l, -r

		},
	}
	softDistortionFx = fx{
		name: "soft distortion",
		apply: func(l, r float32) (float32, float32) {
			l = clip(l, 0.01, 0.1)
			r = clip(r, 0.01, 0.1)
			return l, r
		},
	}
	heavyDistortionFx = fx{
		name: "heavy distortion",
		apply: func(l, r float32) (float32, float32) {
			l = clip(l, 0.01, 0.2)
			r = clip(r, 0.01, 0.2)
			return l, r
		},
	}
	vibratoCount = 0
	vibrato      = fx{
		name: "vibrato",
		apply: func(l, r float32) (float32, float32) {
			m := float32(math.Sin(2 * math.Pi * 2 * float64(vibratoCount) / float64(sampleRate)))
			vibratoCount++
			return l * m, r * m
		},
	}
)

func clip(v, threshold, level float32) float32 {
	if v > threshold {
		v = level
	}
	if v < -threshold {
		v = level
	}
	return v
}
