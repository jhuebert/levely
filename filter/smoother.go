package filter

import (
	"github.com/jhuebert/levely/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"time"
)

type Smoother struct {
	smoothing  float64
	lastUpdate time.Time
	smoothed   float64
}

func NewSmoother() *Smoother {
	smoothing := viper.GetFloat64(config.AccelerometerFilterSmootherSmoothing)
	logrus.Debugf("creating smoothing filter with smoothing=%v", smoothing)
	return &Smoother{
		smoothing:  smoothing,
		lastUpdate: time.Now(),
		smoothed:   0,
	}
}

func (f *Smoother) Add(value float64) {
	now := time.Now()
	valueDifference := value - f.smoothed
	ratio := float64(now.Sub(f.lastUpdate).Milliseconds()) / f.smoothing
	if ratio > 1 {
		logrus.Debug("ratio is too large")
		ratio = 1
	}
	f.smoothed += ratio * valueDifference
	f.lastUpdate = now
}

func (f *Smoother) Value() float64 {
	return f.smoothed
}
