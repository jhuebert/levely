package repository

import (
	"github.com/asdine/storm/v3"
	"github.com/asdine/storm/v3/q"
)

type Position struct {
	ID          int     `storm:"id,increment" json:"id,omitempty"`
	Name        string  `json:"name,omitempty"`
	Pitch       float64 `json:"pitch"`
	Roll        float64 `json:"roll"`
	Favorite    bool    `storm:"index" json:"favorite,omitempty"`
	Calibration bool    `storm:"index" json:"calibration,omitempty"`
}

func (r *Repository) FindCalibration() (Position, error) {
	var entity Position
	err := r.db.One("Calibration", true, &entity)
	return entity, err
}

func (r *Repository) FindPositions() ([]Position, error) {
	var entities []Position
	s := q.Eq("Calibration", false)
	err := r.db.Select(s).Find(&entities)
	if err == storm.ErrNotFound {
		return []Position{}, nil
	}
	return entities, err
}

func (r *Repository) FindPosition(id int) (Position, error) {
	var entity Position
	err := r.db.One("ID", id, &entity)
	return entity, err
}

func (r *Repository) SavePosition(entity *Position) error {
	return r.db.Save(entity)
}

func (r *Repository) DeletePosition(id int) error {
	entity, err := r.FindPosition(id)
	if err != nil {
		return err
	}
	return r.db.DeleteStruct(&entity)
}
