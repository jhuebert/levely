package repository

import (
	"github.com/asdine/storm/v3"
	"github.com/sirupsen/logrus"
)

type Repository struct {
	db *storm.DB
}

func New(path string) (*Repository, error) {
	logrus.Infof("opening databse: %v", path)
	db, err := storm.Open(path)
	return &Repository{db}, err
}

func (r *Repository) Get(bucketName string, key interface{}, to interface{}) error {
	return r.db.Get(bucketName, key, to)
}

func (r *Repository) Close() {
	err := r.db.Close()
	if err != nil {
		logrus.Warnf("error encountered when closing database: %v", err)
	}
}
