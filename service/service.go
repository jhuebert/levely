package service

import (
	"github.com/jhuebert/levely/repository"
	"gobot.io/x/gobot/drivers/i2c"
	"sync"
	"time"
)

type Service struct {
	r *repository.Repository
	d *i2c.MPU6050Driver
	a *Acceleration
}

type Acceleration struct {
	sync.RWMutex
	timestamp time.Time
	x, y, z   float64
}

func New(r *repository.Repository, d *i2c.MPU6050Driver) *Service {
	s := &Service{r, d, &Acceleration{
		timestamp: time.Now(),
		x:         0,
		y:         0,
		z:         0,
	}}
	go s.captureData()
	return s
}
