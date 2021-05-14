package repository

import "github.com/asdine/storm/v3"

type Repository struct {
	db *storm.DB
}

func New(path string) (*Repository, error) {
	db, err := storm.Open(path)
	return &Repository{db}, err
}

func (r *Repository) Get(bucketName string, key interface{}, to interface{}) error {
	return r.db.Get(bucketName, key, to)
}

//func init(db *storm.DB) {
//	db.
//}
