package service

import (
	"github.com/jhuebert/levely/repository"
	"github.com/sirupsen/logrus"
	"math"
	"time"
)

const (
	toDegrees float64 = 180 / math.Pi
)

func (s *Service) GetCurrentPosition() repository.Position {
	p := s.GetUncorrected()
	calibration := s.GetCalibration()
	p.Roll = p.Roll - calibration.Roll
	p.Pitch = p.Pitch - calibration.Pitch
	return p
}

func (s *Service) GetUncorrected() repository.Position {

	s.a.RLock()
	data := []float64{s.a.x, s.a.y, s.a.z}
	s.a.RUnlock()

	prefs := s.GetPreferences()

	var p repository.Position

	switch prefs.OrientationRoll {
	case repository.AxisX:
		switch prefs.OrientationPitch {
		case repository.AxisY:
			p = calculatePosition(data[0], data[1], data[2])
		case repository.AxisZ:
			p = calculatePosition(data[0], data[2], data[1])
		default:
			logrus.Errorf("invalid orientation %v-%v", prefs.OrientationRoll, prefs.OrientationPitch)
		}
	case repository.AxisY:
		switch prefs.OrientationPitch {
		case repository.AxisX:
			p = calculatePosition(data[1], data[0], data[2])
		case repository.AxisZ:
			p = calculatePosition(data[1], data[2], data[0])
		default:
			logrus.Errorf("invalid orientation %v-%v", prefs.OrientationRoll, prefs.OrientationPitch)
		}
	case repository.AxisZ:
		switch prefs.OrientationPitch {
		case repository.AxisX:
			p = calculatePosition(data[2], data[0], data[1])
		case repository.AxisY:
			p = calculatePosition(data[2], data[1], data[0])
		default:
			logrus.Errorf("invalid orientation %v-%v", prefs.OrientationRoll, prefs.OrientationPitch)
		}
	default:
		logrus.Errorf("invalid orientation %v-%v", prefs.OrientationRoll, prefs.OrientationPitch)
	}

	if prefs.OrientationInvertRoll {
		p.Roll = -p.Roll
	}

	if prefs.OrientationInvertPitch {
		p.Pitch = -p.Pitch
	}

	return p
}

func (s *Service) captureData() {

	if s.d == nil {
		logrus.Warn("returning zero data since driver is not set")
		return
	}

	p := s.GetPreferences()
	rate := math.Min(math.Max(p.AccelerometerRate, 1), 1000)
	deviceDuration := time.Duration(1000/rate) * time.Millisecond
	ticker := time.NewTicker(deviceDuration)
	smoothing := p.AccelerometerSmoothing

	values := []float64{0, 0, 0}

	for tick := range ticker.C {
		s.getRawData(values)

		s.a.Lock()
		s.a.x = updateSmooth(s.a.timestamp, s.a.x, tick, values[0], smoothing)
		s.a.y = updateSmooth(s.a.timestamp, s.a.y, tick, values[1], smoothing)
		s.a.z = updateSmooth(s.a.timestamp, s.a.z, tick, values[2], smoothing)
		s.a.timestamp = tick
		s.a.Unlock()

		//TODO Need to support updating rate and smoothing while running
	}

	//TODO Need to be able to exit on command
}

func (s *Service) getRawData(values []float64) {
	err := s.d.GetData()
	if err != nil {
		logrus.Debug(err)
		return
	}
	data := s.d.Accelerometer
	values[0] = float64(data.X)
	values[1] = float64(data.Y)
	values[2] = float64(data.Z)
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

func updateSmooth(previousTime time.Time, previousValue float64, currentTime time.Time, currentValue float64, smoothing float64) float64 {
	elapsedTime := currentTime.Sub(previousTime)
	valueDifference := currentValue - previousValue
	ratio := float64(elapsedTime.Milliseconds()) / smoothing
	if ratio > 1 {
		ratio = 1
		logrus.Debug("ratio is too large")
	}
	return previousValue + (ratio * valueDifference)
}
