package service

import (
	"fmt"
	"github.com/asdine/storm/v3"
	"github.com/jhuebert/levely/repository"
	"github.com/sirupsen/logrus"
)

func (s *Service) getDefaultCalibration() repository.Position {
	return repository.Position{
		Name:        "Calibration",
		Calibration: true,
		Favorite:    false,
		Pitch:       0,
		Roll:        0,
	}
}

func (s *Service) GetCalibration() repository.Position {

	p, err := s.r.FindCalibration()
	if err == nil {
		return p
	}

	logrus.Errorf("calibration does not exist, using default: %v", err)
	p = s.getDefaultCalibration()

	if err == storm.ErrNotFound {
		logrus.Debug("saving default calibration")
		err = s.r.SavePosition(&p)
		if err != nil {
			logrus.Errorf("could not save default calibration: %v", err)
		}
	}

	return p
}

func (s *Service) UpdateCalibration(updated repository.Position) (repository.Position, error) {

	c := s.GetCalibration()

	c.Roll = updated.Roll
	c.Pitch = updated.Pitch
	c.Calibration = true
	c.Favorite = false

	err := s.r.SavePosition(&c)
	if err != nil {
		err = fmt.Errorf("could not save calibration: %v", err)
	}

	return c, err
}
