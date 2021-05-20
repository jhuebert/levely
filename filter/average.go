package filter

import (
	"github.com/jhuebert/levely/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type MovingAverage struct {

	// Array of sample history
	samples []float64

	// Current size of the moving average history. This will never be larger than the samples array.
	currentSize int

	// Index to place the next sample
	index int

	// Sum of the samples
	sum float64
}

func NewMovingAverage() *MovingAverage {
	size := viper.GetInt(config.AccelerometerFilterAverageSize)
	logrus.Debugf("creating moving average filter with size=%v", size)
	return &MovingAverage{
		samples:     make([]float64, size),
		currentSize: 0,
		index:       0,
		sum:         0,
	}
}

// Add Adds a sample to the moving average. If the moving average already is at maximum size, the oldest sample is removed.
func (f *MovingAverage) Add(value float64) {

	// Remove the oldest sample's value from the sum if the array is full
	if f.currentSize == len(f.samples) {
		f.sum -= f.samples[f.index]
	} else {
		f.currentSize++
	}

	// Add the sample to the array and update the sum
	f.samples[f.index] = value
	f.index++
	f.sum += value

	// Reset the index for the next sample if it is currently beyond the end of the array
	if f.index == len(f.samples) {
		f.index = 0
	}
}

func (f *MovingAverage) Value() float64 {
	if f.currentSize > 0 {
		return f.sum / float64(f.currentSize)
	}
	return 0
}
