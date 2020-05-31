package main

import "strconv"

type BeaufortScale struct {
	Start float32
	End   float32
	Bft   int
}

var SCALES = [...]BeaufortScale{
	{
		Start: 0,
		End:   0.4,
		Bft:   0,
	}, {
		Start: 0.5,
		End:   1.5,
		Bft:   1,
	}, {
		Start: 1.6,
		End:   3.3,
		Bft:   2,
	}, {
		Start: 3.4,
		End:   5.5,
		Bft:   3,
	}, {
		Start: 5.6,
		End:   7.9,
		Bft:   4,
	}, {
		Start: 8,
		End:   10.7,
		Bft:   5,
	}, {
		Start: 10.8,
		End:   13.8,
		Bft:   6,
	}, {
		Start: 13.9,
		End:   17.1,
		Bft:   7,
	}, {
		Start: 17.2,
		End:   20.7,
		Bft:   8,
	}, {
		Start: 20.8,
		End:   24.4,
		Bft:   9,
	}, {
		Start: 24.5,
		End:   28.4,
		Bft:   10,
	}, {
		Start: 28.5,
		End:   32.6,
		Bft:   11,
	}, {
		Start: 32.7,
		End:   99,
		Bft:   12,
	},
}

func MeterPerSecondToBeaufortScale(meterPerSecond float32) int {
	for _, scale := range SCALES {
		if isBetween(scale, meterPerSecond) {
			return scale.Bft
		}
	}

	return -1
}

func DisplayBeaufort(beaufort int) string {
	return strconv.Itoa(beaufort) + " bft"
}

func isBetween(scale BeaufortScale, number float32) bool {
	return scale.Start <= number && number <= scale.End
}
