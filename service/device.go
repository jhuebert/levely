package service

import (
	"errors"
	"github.com/jhuebert/levely/repository"
	"github.com/sirupsen/logrus"
	"math"
)

const (
	toDegrees float64 = 180 / math.Pi
)

func (s *Service) GetCurrentPosition() (repository.Position, error) {

	p, err := s.GetUncorrected()
	if err != nil {
		return p, err
	}

	calibration := s.GetCalibration()
	p.Roll = p.Roll - calibration.Roll
	p.Pitch = p.Pitch - calibration.Pitch

	return p, nil
}

func (s *Service) GetUncorrected() (repository.Position, error) {

	data, err := s.getRawData()
	if err != nil {
		return repository.Position{}, err
	}

	prefs := s.GetPreferences()

	var p repository.Position

	switch prefs.Orientation.Width {
	case repository.AxisX:
		switch prefs.Orientation.Length {
		case repository.AxisY:
			p = calculatePosition(data[0], data[1], data[2])
		case repository.AxisZ:
			p = calculatePosition(data[0], data[2], data[1])
		default:
			return p, errors.New("invalid orientation")
		}
	case repository.AxisY:
		switch prefs.Orientation.Length {
		case repository.AxisX:
			p = calculatePosition(data[1], data[0], data[2])
		case repository.AxisZ:
			p = calculatePosition(data[1], data[2], data[0])
		default:
			return p, errors.New("invalid orientation")
		}
	case repository.AxisZ:
		switch prefs.Orientation.Length {
		case repository.AxisX:
			p = calculatePosition(data[2], data[0], data[1])
		case repository.AxisY:
			p = calculatePosition(data[2], data[1], data[0])
		default:
			return p, errors.New("invalid orientation")
		}
	default:
		return p, errors.New("invalid orientation")
	}

	if prefs.Orientation.InvertRoll {
		p.Roll = -p.Roll
	}

	if prefs.Orientation.InvertPitch {
		p.Pitch = -p.Pitch
	}

	return p, nil
}

func (s *Service) getRawData() ([]float64, error) {
	if s.d == nil {
		logrus.Warn("returning zero data since driver is not set")
		return []float64{0, 0, 0}, nil
	}

	err := s.d.GetData()
	if err != nil {
		return nil, err
	}
	data := s.d.Accelerometer
	return []float64{float64(data.X), float64(data.Y), float64(data.Z)}, nil
}

func calculatePosition(x, y, z float64) repository.Position {
	return repository.Position{
		Roll:  adjustedAtan2(-x, math.Sqrt(y*y+z*z)),
		Pitch: adjustedAtan2(y, z),
	}
}

func adjustedAtan2(y, x float64) float64 {
	r := math.Atan2(y, x) * toDegrees
	for r > 90 {
		r -= 180
	}
	for r < -90 {
		r += 180
	}
	return r
}
