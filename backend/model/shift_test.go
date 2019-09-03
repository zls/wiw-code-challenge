package model

import (
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	now := time.Now()
	later := now.Add(time.Minute * 10)
	_, err := NewShift(1, 1, now, later)
	if err != nil {
		t.Error("Expected shift got ", err)
	}
}
