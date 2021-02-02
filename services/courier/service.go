package courier

import (
	"container/list"
	"math/rand"
	"time"

	"github.com/fguy/dispatch-simulation/config"
	"github.com/fguy/dispatch-simulation/entities"
	"github.com/fguy/dispatch-simulation/errors"
	"github.com/fguy/dispatch-simulation/utils"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const nCouriers = 132

type service struct {
	logger       *zap.Logger
	orderQueue   *utils.OrderQueue
	couriers     *list.List
	courierQueue *utils.CourierQueue
	stat         *utils.Stat
}

func (s *service) Dispatch() (*entities.Courier, error) {
	if s.couriers.Len() == 0 {
		return nil, errors.ErrNoCourierAvailable
	}
	e := s.couriers.Front()
	s.couriers.Remove(e)

	courier := e.Value.(*entities.Courier)
	return courier, nil
}

func (s *service) Wait(order *entities.Order) (chan *entities.Dispatch, error) {
	s.orderQueue.Push(order)
	chDispatch := make(chan *entities.Dispatch)
	courier, err := s.Dispatch()
	if err != nil {
		return nil, err
	}
	s.courierQueue.PushAfterArrival(courier)
	chCourier := s.courierQueue.WaitForPop()
	ticker := time.NewTicker(1 * time.Second)
	go func(ticker *time.Ticker, chDispatch chan *entities.Dispatch, chCourier chan *entities.Courier) {
		for {
			<-ticker.C
			// If there are multiple orders available, pick up an arbitrary order.
			firstOrder := s.orderQueue.Pop()
			if firstOrder == nil {
				// If there are no available orders, couriers wait for the next available one.
				s.logger.Info("no order available for courier. waiting...")
				continue
			}
			ticker.Stop()
			// When there are multiple couriers waiting, the next available order is assigned to the earliest​ a​rrived courier.
			earliestArrivedCourier := <-chCourier
			chDispatch <- entities.NewDispatch(firstOrder, earliestArrivedCourier)
			return
		}
	}(ticker, chDispatch, chCourier)

	return chDispatch, nil
}

func (s *service) Deliver(dispatch *entities.Dispatch) error {
	pickupAt := time.Now()
	arrivedAt := dispatch.Courier.LeaveAt.Add(dispatch.Courier.ArrivingIn)

	s.stat.Emit(utils.StatKindCourierWaitTime, pickupAt.Sub(arrivedAt))
	s.stat.Emit(utils.StatKindFoodWaitTime, pickupAt.Sub(*dispatch.Order.ReadyAt))

	s.logger.Info("delivered", zap.Any("dispatch", dispatch))
	return nil
}

// New creates a service instance
func New(
	logger *zap.Logger,
	cfg *config.AppConfig,
	orderQueue *utils.OrderQueue,
	courierQueue *utils.CourierQueue,
	stat *utils.Stat,
) Interface {
	return &service{
		logger:       logger,
		orderQueue:   orderQueue,
		courierQueue: courierQueue,
		couriers:     loadCouriersInUniformDistribution(cfg.ArrivalTimes),
		stat:         stat,
	}
}

func loadCouriersInUniformDistribution(arrivalTimes *config.ArrivalTimes) *list.List {
	couriers := make([]*entities.Courier, nCouriers)
	for i := range couriers {
		interval := i % int(arrivalTimes.Interval())
		couriers[i] = &entities.Courier{
			ID:         uuid.New(),
			ArrivingIn: time.Duration(interval) + arrivalTimes.Min,
		}
	}
	// randomize
	rand.Shuffle(len(couriers), func(i, j int) {
		couriers[i], couriers[j] = couriers[j], couriers[i]
	})
	result := list.New()
	for _, courier := range couriers {
		result.PushBack(courier)
	}
	return result
}
