package filter

import (
	"github.com/jhuebert/levely/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"sort"
)

type MovingMedian struct {

	// Array of sample history
	samples []float64

	// Array of sorted samples
	sorted []float64

	// Current size of the moving average history. This will never be larger than the samples array.
	currentSize int

	// Index to place the next sample
	index int

	// Indicates if the length of samples is even
	even bool
}

func NewMovingMedian() *MovingMedian {
	size := viper.GetInt(config.AccelerometerFilterMedianSize)
	logrus.Debugf("creating moving median filter with size=%v", size)
	even := size%2 == 0
	return &MovingMedian{
		samples:     make([]float64, size),
		sorted:      make([]float64, size),
		currentSize: 0,
		index:       0,
		even:        even,
	}
}

// Add Adds a sample to the moving average. If the moving average already is at maximum size, the oldest sample is removed.
func (f *MovingMedian) Add(value float64) {

	// Remove the oldest sample's value from the sum if the array is full
	if f.currentSize == len(f.samples) {
		//f.sum -= f.samples[f.index]
	} else {
		f.currentSize++
	}

	// Add the sample to the array and update the sum
	f.samples[f.index] = value
	f.index++
	//f.sum += value

	// Reset the index for the next sample if it is currently beyond the end of the array
	if f.index == len(f.samples) {
		f.index = 0
	}
}

func (f *MovingMedian) Value() float64 {
	copy(f.sorted, f.samples)
	sort.Float64s(f.sorted)
	if f.even {
		//TODO
		return f.sorted[len(f.sorted)/2]
	}
	return f.sorted[len(f.sorted)/2]
}
