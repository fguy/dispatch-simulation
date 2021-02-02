package entities

import (
	"time"

	"github.com/google/uuid"
)

// Courier entity
type Courier struct {
	ID         uuid.UUID
	ArrivingIn time.Duration
	LeaveAt    *time.Time
}

// WaitForArrival wait for ArrivingIn duration
func (c *Courier) WaitForArrival() {
	now := time.Now()
	c.LeaveAt = &now
	time.Sleep(c.ArrivingIn)
}
