package service

import (
	"github.com/jhuebert/levely/config"
	"github.com/jhuebert/levely/filter"
	"github.com/jhuebert/levely/repository"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
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

	var p repository.Position

	if s.d == nil {
		logrus.Debug("returning zero data since driver is not set")
		return p
	}

	select {
	case p = <-s.c:
		logrus.Debug("received position data")
	case <-time.After(5 * time.Second):
		logrus.Warn("timed out waiting for position data")
	}

	return p
}

func (s *Service) captureData() {

	if s.d == nil {
		logrus.Warn("returning zero data since driver is not set")
		return
	}

	prefs := s.GetPreferences()

	fast := true
	fastPeriod := viper.GetDuration(config.AccelerometerUpdatePeriod)
	slowPeriod := viper.GetDuration(config.AccelerometerUpdateSleepPeriod)
	maxLastRequest := viper.GetDuration(config.AccelerometerUpdateSleepWait)
	ticker := time.NewTicker(fastPeriod)

	lastRequest := time.Now()
	rollFilter := filter.CreateFilter()
	pitchFilter := filter.CreateFilter()

	for ts := range ticker.C {

		p := s.calculatePosition(prefs)

		rollFilter.Add(p.Roll)
		pitchFilter.Add(p.Pitch)

		p.Roll = rollFilter.Value()
		p.Pitch = pitchFilter.Value()

		select {
		case s.c <- p:
			logrus.Debug("sent acceleration data")
			lastRequest = time.Now()
		default:
			logrus.Debug("no consumer waiting for acceleration data")
		}

		sinceLastRequest := ts.Sub(lastRequest)
		if fast && (sinceLastRequest > maxLastRequest) {
			logrus.Debug("decreasing update speed to save power")
			ticker.Reset(slowPeriod)
			fast = false
		} else if !fast && (sinceLastRequest < maxLastRequest) {
			logrus.Debug("increasing update speed due to demand")
			ticker.Reset(fastPeriod)
			fast = true
			logrus.Debug("updating preferences")
			prefs = s.GetPreferences()
		}
	}
}

type acceleration struct {
	x, y, z float64
}

func (s *Service) calculatePosition(prefs repository.Preferences) repository.Position {

	var p repository.Position

	err := s.d.GetData()
	if err != nil {
		logrus.Debug(err)
		return p
	}

	accel := acceleration{
		x: float64(s.d.Accelerometer.X),
		y: float64(s.d.Accelerometer.Y),
		z: float64(s.d.Accelerometer.Z),
	}

	switch prefs.OrientationRoll {
	case repository.AxisX:
		switch prefs.OrientationPitch {
		case repository.AxisY:
			p = calculatePosition(accel.x, accel.y, accel.z)
		case repository.AxisZ:
			p = calculatePosition(accel.x, accel.z, accel.y)
		default:
			logrus.Errorf("invalid orientation %v-%v", prefs.OrientationRoll, prefs.OrientationPitch)
		}
	case repository.AxisY:
		switch prefs.OrientationPitch {
		case repository.AxisX:
			p = calculatePosition(accel.y, accel.x, accel.z)
		case repository.AxisZ:
			p = calculatePosition(accel.y, accel.z, accel.x)
		default:
			logrus.Errorf("invalid orientation %v-%v", prefs.OrientationRoll, prefs.OrientationPitch)
		}
	case repository.AxisZ:
		switch prefs.OrientationPitch {
		case repository.AxisX:
			p = calculatePosition(accel.z, accel.x, accel.y)
		case repository.AxisY:
			p = calculatePosition(accel.z, accel.y, accel.x)
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

func (s *Service) getDuration(rate float64) time.Duration {
	limited := math.Min(math.Max(rate, 1), 1000)
	return time.Duration(1000/limited) * time.Millisecond
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
