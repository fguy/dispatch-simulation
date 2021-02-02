package utils

import (
	"sync"
	"time"

	"github.com/fguy/dispatch-simulation/entities"
	"go.uber.org/zap"
)

// CourierQueue is a blocking queue for First-in-first-out strategy
type CourierQueue struct {
	sync.RWMutex
	logger   *zap.Logger
	couriers []*entities.Courier
}

// PushAfterArrival add a courier on the queue after their ArrivingIn time
func (q *CourierQueue) PushAfterArrival(courier *entities.Courier) {
	q.logger.Debug("a courier to be pushed to the queue", zap.Any("courier", courier))
	courier.WaitForArrival()
	q.logger.Debug("a courier arrived thus enqueue", zap.Any("courier", courier))
	q.Lock()
	defer q.Unlock()
	q.couriers = append(q.couriers, courier)
}

// WaitForPop gets and removes a courier from the queue when courier show up
func (q *CourierQueue) WaitForPop() chan *entities.Courier {
	ticker := time.NewTicker(1 * time.Second)
	ch := make(chan *entities.Courier, 1)
	// tick until there's available courier
	go func(ticker *time.Ticker) {
		for {
			<-ticker.C
			q.Lock()
			defer q.Unlock()
			if len(q.couriers) == 0 {
				q.logger.Debug("no courier in the queue")
				continue
			}
			q.logger.Debug("pop a courier", zap.Any("courier", q.couriers[0]))
			ch <- q.couriers[0]
			//q.couriers[0] = nil // prevent memory leaks
			q.couriers = q.couriers[1:]
			ticker.Stop()
			return
		}
	}(ticker)
	return ch
}

// All returns all courier items in the queue
func (q *CourierQueue) All() []*entities.Courier {
	q.RLock()
	defer q.RUnlock()
	return q.couriers
}

// NewCourierQueue returns a new instance of CourierQueue
func NewCourierQueue(logger *zap.Logger) *CourierQueue {
	return &CourierQueue{
		logger: logger,
	}
}
