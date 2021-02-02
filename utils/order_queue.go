package utils

import (
	"sync"

	"github.com/fguy/dispatch-simulation/entities"
	"go.uber.org/zap"
)

// OrderQueue is a blocking queue for First-in-first-out strategy
type OrderQueue struct {
	sync.RWMutex
	logger *zap.Logger
	orders []*entities.Order
}

// Push add an order item on the queue
func (q *OrderQueue) Push(order *entities.Order) {
	q.Lock()
	defer q.Unlock()

	q.logger.Debug("adding an order to the queue", zap.Any("order", order))
	q.orders = append(q.orders, order)
}

// Pop return and removes an order item from the queue
func (q *OrderQueue) Pop() *entities.Order {
	q.Lock()
	defer q.Unlock()
	if len(q.orders) == 0 {
		return nil
	}
	order := q.orders[0]
	q.logger.Debug("popping an order from the queue", zap.Any("order", order))
	result := q.orders[0]
	//q.orders[0] = nil // prevent memory leaks
	q.orders = q.orders[1:]
	return result
}

// All returns all order items in the queue
func (q *OrderQueue) All() []*entities.Order {
	q.RLock()
	defer q.RUnlock()
	return q.orders
}

// NewOrderQueue returns a new instance of OrderQueue
func NewOrderQueue(logger *zap.Logger) *OrderQueue {
	return &OrderQueue{
		logger: logger,
	}
}
