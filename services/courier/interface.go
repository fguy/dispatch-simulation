package courier

import (
	"github.com/fguy/dispatch-simulation/entities"
)

//go:generate mockgen -destination=../../mocks/services/courier/interface.go github.com/fguy/dispatch-simulation/services/courier Interface

// Interface is an interface of courier service
type Interface interface {
	Dispatch() (*entities.Courier, error)
	Wait(*entities.Order) (chan *entities.Dispatch, error)
	Deliver(*entities.Dispatch) error
}
