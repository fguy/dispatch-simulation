package entities

import "time"

// Dispatch entity
type Dispatch struct {
	Order     *Order
	Courier   *Courier
	CreatedAt time.Time
}

// NewDispatch creates a new Dispatch entity
func NewDispatch(
	order *Order,
	courier *Courier,
) *Dispatch {
	return &Dispatch{
		Order:     order,
		Courier:   courier,
		CreatedAt: time.Now(),
	}
}
