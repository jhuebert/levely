package filter

import (
	"embed"
	"encoding/json"
	"fmt"
	"github.com/jhuebert/levely/config"
	"github.com/spf13/viper"
	"io/ioutil"
	"strings"
)

//go:embed iir
var iirFiles embed.FS

func NewIir() (*Iir, error) {
	var data []byte
	var err error

	path := viper.GetString(config.AccelerometerFilterIirPath)
	if strings.TrimSpace(path) != "" {
		data, err = ioutil.ReadFile(path)
	} else {
		preset := viper.GetInt(config.AccelerometerFilterIirPreset)
		data, err = firFiles.ReadFile(fmt.Sprintf("iir/iir%v.json", preset))
	}

	f := &Iir{}
	if err != nil {
		return f, err
	}

	err = json.Unmarshal(data, f)
	return f, err
}

type Iir struct {
	Sections []*IirSection `json:"sections"`
	output   float64
}

func (f *Iir) Add(value float64) {
	f.output = f.add(0, value)
	for i := 1; i < len(f.Sections); i++ {
		f.output = f.add(i, f.output)
	}
}

func (f *Iir) add(sectionIndex int, input float64) float64 {
	section := f.Sections[sectionIndex]
	section.Add(input)
	return section.Value()
}

func (f *Iir) Value() float64 {
	return f.output
}

type IirSection struct {
	A0           float64 `json:"a0"`
	A1           float64 `json:"a1"`
	A2           float64 `json:"a2"`
	B0           float64 `json:"b0"`
	B1           float64 `json:"b1"`
	B2           float64 `json:"b2"`
	regX1, regX2 float64
	regY1, regY2 float64
	output       float64
}

func (f *IirSection) Add(input float64) {
	centerTap := input*f.B0 + f.B1*f.regX1 + f.B2*f.regX2
	f.output = f.A0*centerTap - f.A1*f.regY1 - f.A2*f.regY2

	f.regX2 = f.regX1
	f.regX1 = input
	f.regY2 = f.regY1
	f.regY1 = f.output
}

func (f *IirSection) Value() float64 {
	return f.output
}
