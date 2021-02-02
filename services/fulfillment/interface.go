package fulfillment

//go:generate mockgen -destination=../../mocks/services/fulfillment/interface.go github.com/fguy/dispatch-simulation/services/fulfillment Interface

// Interface is an interface of fulfillment service
type Interface interface {
	Start() error
	Stop() error
}

// Matched is an interface of fulfillment service with "matched" strategey
type Matched interface {
	Interface
}

// FirstInFirstOut is an interface of fulfillment service with "FIFO" strategey
type FirstInFirstOut interface {
	Interface
}
