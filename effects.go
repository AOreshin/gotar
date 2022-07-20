package main

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
