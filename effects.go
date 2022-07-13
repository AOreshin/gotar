package main

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
