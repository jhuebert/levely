package service

import (
	"errors"
	"fmt"
	"github.com/jhuebert/levely/repository"
)

func (s *Service) FindAllPositions() ([]repository.Position, error) {
	return s.r.FindPositions()
}

func (s *Service) FindPosition(id int) (repository.Position, error) {
	return s.r.FindPosition(id)
}

func (s *Service) CreatePosition(position repository.Position) (repository.Position, error) {
	position.ID = 0 //Ensure ID is "blank" so we aren't overwriting existing data
	err := s.r.SavePosition(&position)
	return position, err
}

func (s *Service) UpdatePosition(id int, updated repository.Position) (repository.Position, error) {

	position, err := s.r.FindPosition(id)
	if err != nil {
		return updated, errors.New(fmt.Sprintf("position %v not found", id))
	}

	position.Name = updated.Name
	position.Roll = updated.Roll
	position.Pitch = updated.Pitch
	position.Favorite = updated.Favorite

	err = s.r.SavePosition(&position)

	return position, err

}

func (s *Service) DeletePosition(id int) error {
	return s.r.DeletePosition(id)
}
