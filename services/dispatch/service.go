package dispatch

import (
	"github.com/fguy/dispatch-simulation/entities"
	"github.com/fguy/dispatch-simulation/services/courier"
	"go.uber.org/zap"
)

type service struct {
	logger         *zap.Logger
	courierService courier.Interface
}

func (s *service) Receive(strategy entities.Strategy, order *entities.Order) error {
	s.logger.Info("received an order", zap.Any("strategy", strategy), zap.Any("order", order))
	switch strategy {
	// a courier is dispatched for a specific order and may only pick up that order
	case entities.StrategyMatched:
		return s.receiveMatched(order)
	// a courier picks up the next available order upon arrival.
	case entities.StrategyFirstInFirstOut:
		return s.receiveFirstInFirstOut(order)
	}
	return nil
}

func (s *service) receiveMatched(order *entities.Order) error {
	courier, err := s.courierService.Dispatch()
	if err != nil {
		return err
	}
	dispatch := entities.NewDispatch(order, courier)
	order.StartPrep()
	courier.WaitForArrival()
	return s.postAction(dispatch)
}

func (s *service) receiveFirstInFirstOut(order *entities.Order) error {
	ch, err := s.courierService.Wait(order)
	if err != nil {
		return err
	}
	order.StartPrep()
	go func(ch chan *entities.Dispatch) {
		dispatch := <-ch
		s.postAction(dispatch)
	}(ch)
	return nil
}

func (s *service) postAction(dispatch *entities.Dispatch) error {
	if err := dispatch.Order.EndPrep(); err != nil {
		return err
	}
	err := s.courierService.Deliver(dispatch)
	return err
}

// New creates a service instance
func New(
	logger *zap.Logger,
	courierService courier.Interface,
) Interface {
	return &service{
		logger:         logger,
		courierService: courierService,
	}
}
