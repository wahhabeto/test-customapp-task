package app

import (
	"errors"
	"math/rand"
	"time"
)

type Engine interface {
	GenerateMultiplier() (float64, error)
}

type engine struct {
	rtp float64
}

func NewEngine(rtp float64) (Engine, error) {
	if rtp <= 0 || rtp > 1.0 {
		return nil, errors.New("rtp must be between 0 and 1.0")
	}
	return &engine{
		rtp: rtp,
	}, nil
}

func (e *engine) GenerateMultiplier() (float64, error) {
	rand.Seed(time.Now().UnixNano())
	mean := 1.0

	stddev := 0.5 + (1.0-e.rtp)*0.5

	multiplier := mean + stddev*rand.NormFloat64()

	if multiplier < 1.0 {
		multiplier = 1.0
	} else if multiplier > 10000.0 {
		multiplier = 10000.0
	}

	return multiplier, nil
}
