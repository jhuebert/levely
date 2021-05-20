package service

import (
	"github.com/jhuebert/levely/repository"
	"gobot.io/x/gobot/drivers/i2c"
)

type Service struct {
	r *repository.Repository
	d *i2c.MPU6050Driver
	c chan repository.Position
}

func New(r *repository.Repository, d *i2c.MPU6050Driver) *Service {
	s := &Service{r, d, make(chan repository.Position)}
	go s.captureData()
	return s
}
