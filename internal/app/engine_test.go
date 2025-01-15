package app

import (
	"math/rand"
	"testing"
	"time"
)

func withinRange(value, min, max float64) bool {
	return value >= min && value <= max
}

func TestNewEngineValidRTP(t *testing.T) {
	rtp := 0.5
	engine, err := NewEngine(rtp)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if engine == nil {
		t.Error("expected non-nil engine")
	}
}

func TestNewEngineInvalidRTP(t *testing.T) {
	rtps := []float64{-0.1, 0, 1.1}
	for _, rtp := range rtps {
		_, err := NewEngine(rtp)
		if err == nil {
			t.Errorf("expected error for rtp %v", rtp)
		}
	}
}

func TestGenerateMultiplierValid(t *testing.T) {
	rand.Seed(0)

	rtp := 0.5
	engine, err := NewEngine(rtp)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	for i := 0; i < 100; i++ {
		multiplier, err := engine.GenerateMultiplier()
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if !withinRange(multiplier, 1.0, 10000.0) {
			t.Errorf("expected multiplier in range [1.0, 10000.0], got %v", multiplier)
		}
	}
}

func TestGenerateMultiplierWithDifferentRTPs(t *testing.T) {
	rtps := []float64{0.1, 0.25, 0.5, 0.75, 1.0}
	for _, rtp := range rtps {
		engine, err := NewEngine(rtp)
		if err != nil {
			t.Errorf("expected no error for rtp %v, got %v", rtp, err)
		}

		for i := 0; i < 100; i++ {
			multiplier, err := engine.GenerateMultiplier()
			if err != nil {
				t.Errorf("expected no error for rtp %v, got %v", rtp, err)
			}
			if !withinRange(multiplier, 1.0, 10000.0) {
				t.Errorf("expected multiplier in range [1.0, 10000.0] for rtp %v, got %v", rtp, multiplier)
			}
		}
	}
}

func TestGenerateMultiplierRange(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	rtp := 0.5
	engine, err := NewEngine(rtp)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	for i := 0; i < 10000; i++ {
		multiplier, err := engine.GenerateMultiplier()
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if !withinRange(multiplier, 1.0, 10000.0) {
			t.Errorf("expected multiplier in range [1.0, 10000.0], got %v", multiplier)
		}
	}
}
