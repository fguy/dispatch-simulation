package utils

import (
	"testing"
	"time"

	"github.com/fguy/dispatch-simulation/factories"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

var logger, _ = zap.NewDevelopment()

func TestCourierQueue_PushAfterArriaval(t *testing.T) {
	t.Parallel()

	q := NewCourierQueue(logger)
	ts := time.Now()
	courier := factories.Courier()
	q.PushAfterArrival(courier)
	time.Sleep(courier.ArrivingIn)
	tsOffset := time.Now().Sub(ts).Seconds()
	assert.Equal(t, courier, q.couriers[0])
	assert.GreaterOrEqual(t, tsOffset, courier.ArrivingIn.Seconds())
}

func TestCourierQueue_WaitForPop(t *testing.T) {
	t.Parallel()

	q := NewCourierQueue(logger)
	ts := time.Now()
	courier := factories.Courier()
	ch := q.WaitForPop()
	go func() {
		received := <-ch
		tsOffset := time.Now().Sub(ts).Seconds()
		assert.Equal(t, courier, received)
		// assert includes sleep times
		assert.GreaterOrEqual(t, tsOffset, float64(2))
	}()
	time.Sleep(2 * time.Second)
	q.couriers = append(q.couriers, courier)
	time.Sleep(2 * time.Second) // allow to wait for the next tick
}
