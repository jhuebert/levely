package filter

import (
	"embed"
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/jhuebert/levely/config"
	"github.com/spf13/viper"
	"io/ioutil"
	"strings"
)

//go:embed fir
var firFiles embed.FS

type Fir struct {

	// Array of filter coefficients
	Coefficients []float64 `json:"coefficients"`

	// Array of sample history
	samples []float64

	// Index to place the next sample
	index int
}

func NewFir() (*Fir, error) {
	var data []byte
	var err error

	path := viper.GetString(config.AccelerometerFilterFirPath)
	if strings.TrimSpace(path) != "" {
		data, err = ioutil.ReadFile(path)
	} else {
		preset := viper.GetInt(config.AccelerometerFilterFirPreset)
		data, err = firFiles.ReadFile(fmt.Sprintf("fir/fir%v.json", preset))
	}

	f := &Fir{}
	if err != nil {
		return f, err
	}

	err = json.Unmarshal(data, f)

	if err == nil {
		f.samples = make([]float64, len(f.Coefficients))
	}

	return f, err
}

// Add Adds a sample to the filter. If the filter already is at maximum size, the oldest sample is removed.
func (f *Fir) Add(value float64) {

	// Add the sample to the array
	f.samples[f.index] = value
	f.index++

	// Reset the index for the next sample if it is currently beyond the end of the array
	if f.index == len(f.samples) {
		f.index = 0
	}
}

func (f *Fir) Value() float64 {
	output := 0.0
	for i := 0; i < len(f.Coefficients); i++ {
		output += f.Coefficients[i] * f.samples[i]
	}
	return output
}
