package service

import (
	"github.com/jhuebert/levely/repository"
	"gobot.io/x/gobot/drivers/i2c"
	"time"
)

type Service struct {
	r *repository.Repository
	d *i2c.MPU6050Driver
	c chan Acceleration
}

type Acceleration struct {
	x, y, z   float64
	timestamp time.Time
}

func New(r *repository.Repository, d *i2c.MPU6050Driver) *Service {
	s := &Service{r, d, make(chan Acceleration)}
	go s.captureData()
	return s
}
