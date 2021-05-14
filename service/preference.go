package service

import (
	"github.com/asdine/storm/v3"
	"github.com/jhuebert/levely/repository"
	"github.com/sirupsen/logrus"
)

var (
	preferencesDefault = repository.Preferences{
		Dimensions: repository.DimensionPreferences{
			Length: 240,
			Width:  96,
			Units:  repository.UnitInches,
		},
		Orientation: repository.OrientationPreferences{
			Length:      repository.AxisX,
			Width:       repository.AxisY,
			InvertPitch: false,
			InvertRoll:  false,
		},
		Tolerance:  0.1,
		UpdateRate: 1,
	}
)

func (s *Service) GetPreferences() repository.Preferences {
	p, err := s.r.GetPreferences()
	if err != nil {
		logrus.Error(err)
		p = preferencesDefault
		if err == storm.ErrNotFound {
			p, err = s.r.UpdatePreferences(p)
			if err != nil {
				logrus.Error(err)
			}
		}
		return p
	}
	return p
}

func (s *Service) UpdatePreferences(updated repository.Preferences) (repository.Preferences, error) {
	return s.r.UpdatePreferences(updated)
}
