package fulfillment

import (
	"container/list"
	"encoding/json"
	"io/ioutil"
	"path"
	"runtime"
	"time"

	"github.com/fguy/dispatch-simulation/config"
	"github.com/fguy/dispatch-simulation/entities"
	"github.com/fguy/dispatch-simulation/services/dispatch"
	"github.com/fguy/dispatch-simulation/utils"
	"go.uber.org/zap"
)

const (
	dispatchOrdesJSONPath = "../../dispatch_orders.json"
	tick                  = 1 * time.Second // every second the system should receive orders
	batch                 = 2               // 2 orders at a tick
)

type service struct {
	logger          *zap.Logger
	tickerDone      chan bool
	queue           *list.List
	dispatchService dispatch.Interface
	stat            *utils.Stat
	stopped         bool
	strategy        entities.Strategy
	orderQueue      *utils.OrderQueue
	courierQueue    *utils.CourierQueue
}

func (s *service) Start() error {
	strategyLoggingField := zap.Any("strategy", s.strategy)
	s.logger.Info("starting fulfillment", strategyLoggingField)
	s.tickerDone = make(chan bool, 1)
	go func() {
		ticker := time.NewTicker(tick)
		defer ticker.Stop()
		for {
			s.logger.Debug("fulfillment tick", zap.Duration("tick", tick), zap.Int("batch", batch), strategyLoggingField)
			select {
			case <-ticker.C:
				if s.queue.Len() < batch { // order ran out
					s.logger.Debug("not enough orders in the queue", zap.Int("queue_size", s.queue.Len()), strategyLoggingField)
					s.tickerDone <- true
					continue
				}
				for i := 0; i < batch; i++ { // pops out in batch size
					e := s.queue.Front()
					s.queue.Remove(e)
					order := e.Value.(*entities.Order)
					s.logger.Info("dispatching", zap.Any("order", order), strategyLoggingField)
					if err := s.dispatchService.Receive(s.strategy, order); err != nil {
						s.logger.Error("can not dispatch", zap.Any("order", order), zap.Error(err), strategyLoggingField)
					}
				}
			case <-s.tickerDone:
				s.logger.Debug("ticker done", strategyLoggingField)
				remaining := s.remainingDuration()
				s.logger.Debug("remaining", zap.Duration("time", remaining))
				time.Sleep(remaining + time.Second) // wait for completion
				s.Stop()
				return
			}
		}
	}()
	return nil
}

// Stop immediately
func (s *service) Stop() error {
	strategyLoggingField := zap.Any("strategy", s.strategy)
	if s.stopped {
		s.logger.Debug("already stopped", strategyLoggingField)
		return nil
	}
	s.logger.Info("stopping the fulfillment", strategyLoggingField)
	s.stat.PrintAvgAll(strategyLoggingField)
	s.stopped = true
	return nil
}

func (s *service) remainingDuration() time.Duration {
	couriers := s.courierQueue.All()
	orders := s.orderQueue.All()
	sum := time.Duration(0)
	for _, courier := range couriers {
		if courier.LeaveAt == nil {
			sum += courier.ArrivingIn
		} else if time.Now().Sub(*courier.LeaveAt) < courier.ArrivingIn {
			sum += time.Now().Sub(*courier.LeaveAt)
		}
	}

	for _, order := range orders {
		if order.ReadyAt == nil {
			if order.PrepStartedAt == nil {
				sum += order.PrepTime * time.Second
			} else {
				readyAt := order.PrepStartedAt.Add(order.PrepTime * time.Second)
				if readyAt.After(time.Now()) {
					sum += readyAt.Sub(time.Now())
				}
			}
		}
	}
	return sum
}

// New returns a service instance
func New(
	logger *zap.Logger,
	cfg *config.AppConfig,
	dispatchService dispatch.Interface,
	stat *utils.Stat,
	orderQueue *utils.OrderQueue,
	courierQueue *utils.CourierQueue,
) (Interface, error) {
	queue, err := loadOrders(logger)
	if err != nil {
		return nil, err
	}
	result := &service{
		logger:          logger,
		queue:           queue,
		dispatchService: dispatchService,
		stat:            stat,
		strategy:        cfg.Strategy,
		orderQueue:      orderQueue,
		courierQueue:    courierQueue,
	}

	return result, nil
}

func loadOrders(logger *zap.Logger) (*list.List, error) {
	result := list.New()
	// read dispatch_orders.json
	_, currentFilePath, _, _ := runtime.Caller(0)
	dir := path.Dir(currentFilePath)
	b, err := ioutil.ReadFile(path.Join(dir, dispatchOrdesJSONPath))
	if err != nil {
		panic(err)
	}

	var orders []*entities.Order
	if err := json.Unmarshal(b, &orders); err != nil {
		return nil, err
	}
	logger.Info("imported orders from json file", zap.Int("orders", len(orders)))
	// push entries in the queue
	for _, order := range orders {
		result.PushBack(order)
	}
	return result, nil
}
