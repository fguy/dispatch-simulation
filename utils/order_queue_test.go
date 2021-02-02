package utils

import (
	"testing"

	"github.com/fguy/dispatch-simulation/factories"
	"github.com/stretchr/testify/assert"
)

func TestOrderQueue_Push_Pop(t *testing.T) {
	t.Parallel()

	q := NewOrderQueue(logger)
	// pop when empty
	assert.Nil(t, q.Pop())

	order := factories.Order()
	q.Push(order)
	assert.Len(t, q.orders, 1)
	assert.NotNil(t, q.Pop())

	// pop when became empty
	assert.Nil(t, q.Pop())
}
