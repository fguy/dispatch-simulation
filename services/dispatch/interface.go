package dispatch

import (
	"github.com/fguy/dispatch-simulation/entities"
)

//go:generate mockgen -destination=../../mocks/services/dispatch/interface.go github.com/fguy/dispatch-simulation/services/dispatch Interface

// Interface is a interface of dispatch service
type Interface interface {
	Receive(entities.Strategy, *entities.Order) error
}
