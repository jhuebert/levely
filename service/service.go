package service

import (
	"github.com/jhuebert/levely/repository"
	"gobot.io/x/gobot/drivers/i2c"
)

type Service struct {
	r *repository.Repository
	d *i2c.MPU6050Driver
}

func New(r *repository.Repository, d *i2c.MPU6050Driver) *Service {
	return &Service{r, d}
}
