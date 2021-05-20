package filter

import (
	"github.com/sirupsen/logrus"
)

type Passthrough struct {
	output float64
}

func NewPassthrough() *Passthrough {
	logrus.Debug("creating passthrough filter")
	return &Passthrough{}
}

func (f *Passthrough) Add(value float64) {
	f.output = value
}

func (f *Passthrough) Value() float64 {
	return f.output
}
