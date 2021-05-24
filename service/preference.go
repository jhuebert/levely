package service

import (
	"errors"
	"net/http"

	"github.com/asdine/storm/v3"
	"github.com/jhuebert/levely/repository"
	"github.com/sirupsen/logrus"
)

func (s *Service) getDefaultPreferences() repository.Preferences {
	return repository.Preferences{
		ID:                     repository.PreferencesId,
		Version:                1,
		DimensionLength:        240,
		DimensionWidth:         96,
		DimensionUnits:         repository.UnitInches,
		OrientationPitch:       repository.AxisY,
		OrientationRoll:        repository.AxisX,
		OrientationInvertPitch: false,
		OrientationInvertRoll:  false,
	}
}

func (s *Service) GetPreferences() repository.Preferences {

	p, err := s.r.GetPreferences()
	if err == nil {
		return p
	}

	logrus.Errorf("no saved preferences, using default preferences: %v", err)
	p = s.getDefaultPreferences()

	if err == storm.ErrNotFound {
		logrus.Debug("saving default preferences")
		p, err = s.r.UpdatePreferences(p)
		if err != nil {
			logrus.Errorf("could not save default preferences: %v", err)
		}
	}

	return p
}

func (s *Service) ExportPreferences(w http.ResponseWriter) error {
	return s.r.ExportPreferences(w)
}

func (s *Service) UpdatePreferences(updated repository.Preferences) (repository.Preferences, error) {

	p := s.GetPreferences()

	err := isValid(updated)
	if err != nil {
		return p, err
	}

	p.DimensionWidth = updated.DimensionWidth
	p.DimensionLength = updated.DimensionLength
	p.DimensionUnits = updated.DimensionUnits
	p.OrientationPitch = updated.OrientationPitch
	p.OrientationRoll = updated.OrientationRoll
	p.OrientationInvertPitch = updated.OrientationInvertPitch
	p.OrientationInvertRoll = updated.OrientationInvertRoll

	return s.r.UpdatePreferences(p)
}

func isValid(p repository.Preferences) error {

	if p.DimensionWidth <= 0.0 {
		return errors.New("dimension width must be greater than zero")
	}

	if p.DimensionLength <= 0.0 {
		return errors.New("dimension length must be greater than zero")
	}

	if (p.DimensionUnits != repository.UnitInches) && (p.DimensionUnits != repository.UnitCentimeters) {
		return errors.New("dimension unit is not valid")
	}

	if (p.OrientationRoll != repository.AxisX) && (p.OrientationRoll != repository.AxisY) && (p.OrientationRoll != repository.AxisZ) {
		return errors.New("roll axis is not set")
	}

	if (p.OrientationPitch != repository.AxisX) && (p.OrientationPitch != repository.AxisY) && (p.OrientationPitch != repository.AxisZ) {
		return errors.New("pitch axis is not set")
	}

	if p.OrientationRoll == p.OrientationPitch {
		return errors.New("pitch and roll axis must not be the same value")
	}

	return nil
}
