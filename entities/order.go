package entities

import (
	"time"

	"github.com/fguy/dispatch-simulation/errors"
	"github.com/google/uuid"
)

// Order entity
type Order struct {
	ID            uuid.UUID
	Name          string
	PrepTime      time.Duration
	PrepStartedAt *time.Time
	ReadyAt       *time.Time
}

// StartPrep marks the order is preparing
func (o *Order) StartPrep() {
	now := time.Now()
	o.PrepStartedAt = &now
}

// EndPrep waits the order to be ready
func (o *Order) EndPrep() error {
	now := time.Now()
	if o.PrepStartedAt == nil {
		return errors.ErrOrderPrepNotStarted
	}
	left := o.PrepTime - now.Sub(*o.PrepStartedAt)
	if left > 0 { // time left to be ready
		time.Sleep(left)
		o.ReadyAt = &now
	} else { // ready before dispatch
		readyAt := o.PrepStartedAt.Add(o.PrepTime)
		o.ReadyAt = &readyAt
	}
	return nil
}
