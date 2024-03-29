package filter

import (
	"github.com/jhuebert/levely/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Filter interface {
	Add(value float64)
	Value() float64
}

func CreateFilter() Filter {
	selected := viper.GetString(config.AccelerometerFilterSelected)
	logrus.Debugf("creating filter: %v", selected)
	switch selected {
	case config.FilterSmoother:
		return NewSmoother()
	case config.FilterAverage:
		return NewMovingAverage()
	case config.FilterPassthrough:
		return NewPassthrough()
	default:
		logrus.Warnf("unknown filter type: %v", selected)
		return NewSmoother()
	}
}
