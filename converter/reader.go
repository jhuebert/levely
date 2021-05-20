package converter

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/jhuebert/levely/filter"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

//go:embed iir.txt
var iirFilter string

//go:embed fir.txt
var firFilter string

func WriteIirFilter() {
	lines := strings.Split(iirFilter, "\n")
	sections := []*filter.IirSection{}
	for i := 0; i < len(lines); i++ {
		if strings.HasPrefix(lines[i], "Sect") {
			sections = append(sections, createSection(lines[i+1:i+7]))
		}
	}
	f := &filter.Iir{Sections: sections}
	printJson(f)
}

func WriteFirFilter() {
	lines := strings.Split(firFilter, "\n")
	cs := []float64{}
	for _, line := range lines {
		cs = append(cs, parseFirFloat(strings.TrimSpace(line)))
	}
	f := &filter.Fir{Coefficients: cs}
	printJson(f)
}

func printJson(f filter.Filter) {
	j, err := json.Marshal(f)
	if err != nil {
		logrus.Error(err)
		return
	}
	fmt.Println(string(j))
}

func createSection(lines []string) *filter.IirSection {
	return &filter.IirSection{
		A0: parseIirFloat(lines[0]),
		A1: parseIirFloat(lines[1]),
		A2: parseIirFloat(lines[2]),
		B0: parseIirFloat(lines[3]),
		B1: parseIirFloat(lines[4]),
		B2: parseIirFloat(lines[5]),
	}
}

func parseIirFloat(text string) float64 {
	v, err := strconv.ParseFloat(strings.Fields(text)[1], 64)
	if err != nil {
		logrus.Error(err)
	}
	return v
}

func parseFirFloat(text string) float64 {
	v, err := strconv.ParseFloat(strings.TrimSpace(text), 64)
	if err != nil {
		logrus.Error(err)
	}
	return v
}
