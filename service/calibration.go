package service

import (
	"github.com/asdine/storm/v3"
	"github.com/jhuebert/levely/repository"
	"github.com/sirupsen/logrus"
)

func (s *Service) GetCalibration() repository.Position {
	p, err := s.r.FindCalibration()
	if err != nil {
		logrus.Error(err)
		p = repository.Position{
			Pitch:       0,
			Roll:        0,
			Calibration: true,
		}
		if err == storm.ErrNotFound {
			err = s.r.SavePosition(&p)
			if err != nil {
				logrus.Error(err)
			}
		}
		return p
	}
	return p
}

func (s *Service) UpdateCalibration(updated repository.Position) (repository.Position, error) {

	current := s.GetCalibration()
	current.Calibration = true
	current.Roll = updated.Roll
	current.Pitch = updated.Pitch
	current.Favorite = false

	err := s.r.SavePosition(&current)

	return current, err
}
