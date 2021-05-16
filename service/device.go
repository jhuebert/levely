package service

import (
	"github.com/jhuebert/levely/config"
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
	var accel Acceleration

	if s.d == nil {
		logrus.Debug("returning zero data since driver is not set")
		return p
	}

	select {
	case accel = <-s.c:
		logrus.Debug("received acceleration data")
	case <-time.After(5 * time.Second):
		logrus.Warn("timed out waiting for acceleration data")
		return p
	}

	prefs := s.GetPreferences()
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

func (s *Service) captureData() {

	if s.d == nil {
		logrus.Warn("returning zero data since driver is not set")
		return
	}

	lastPrefsUpdate := time.Now()
	prefsUpdatePeriod := viper.GetDuration(config.AccelerometerPreferencesUpdatePeriod)
	prefs := s.GetPreferences()

	fast := true
	fastPeriod := s.getDuration(prefs.AccelerometerRate)
	slowPeriod := viper.GetDuration(config.AccelerometerUpdateSleepPeriod)
	maxLastRequest := viper.GetDuration(config.AccelerometerUpdateSleepWait)
	ticker := time.NewTicker(fastPeriod)

	lastRequest := time.Now()
	accel := Acceleration{
		timestamp: time.Now(),
	}

	for ts := range ticker.C {

		err := s.d.GetData()
		if err != nil {
			logrus.Debug(err)
			continue
		}

		accel.x = updateSmooth(accel.timestamp, accel.x, ts, s.d.Accelerometer.X, prefs.AccelerometerSmoothing)
		accel.y = updateSmooth(accel.timestamp, accel.y, ts, s.d.Accelerometer.Y, prefs.AccelerometerSmoothing)
		accel.z = updateSmooth(accel.timestamp, accel.z, ts, s.d.Accelerometer.Z, prefs.AccelerometerSmoothing)
		accel.timestamp = ts

		select {
		case s.c <- accel:
			logrus.Debug("sent acceleration data")
			lastRequest = time.Now()
		default:
			logrus.Debug("no consumer waiting for acceleration data")
		}

		sinceLastRequest := ts.Sub(lastRequest)
		if fast && (sinceLastRequest > maxLastRequest) {
			logrus.Info("decreasing update speed")
			ticker.Reset(slowPeriod)
			fast = false
		} else if !fast && (sinceLastRequest < maxLastRequest) {
			logrus.Info("increasing update speed")
			ticker.Reset(fastPeriod)
			fast = true
		} else if fast && ts.Sub(lastPrefsUpdate) > prefsUpdatePeriod {
			logrus.Infof("updating preferences")
			prefs = s.GetPreferences()
			ticker.Reset(s.getDuration(prefs.AccelerometerRate))
			lastPrefsUpdate = ts
		}
	}
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

func updateSmooth(previousTime time.Time, previousValue float64, currentTime time.Time, currentValue int16, smoothing float64) float64 {
	valueDifference := float64(currentValue) - previousValue
	ratio := float64(currentTime.Sub(previousTime).Milliseconds()) / smoothing
	if ratio > 1 {
		logrus.Debug("ratio is too large")
		ratio = 1
	}
	return previousValue + (ratio * valueDifference)
}
