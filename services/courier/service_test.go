package courier

import (
	"math"
	"testing"
	"time"

	"github.com/fguy/dispatch-simulation/config"
	"github.com/fguy/dispatch-simulation/entities"
	"github.com/fguy/dispatch-simulation/errors"
	"github.com/fguy/dispatch-simulation/factories"
	"github.com/fguy/dispatch-simulation/utils"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

var (
	logger, _ = zap.NewDevelopment()
	cfg       = &config.AppConfig{
		ArrivalTimes: &config.ArrivalTimes{
			Min: 3,
			Max: 15,
		},
	}
)

func TestUniformDistribution(t *testing.T) {
	t.Parallel()

	instance := New(
		logger,
		cfg,
		nil,
		nil,
		nil,
	)
	couriers := instance.(*service).couriers
	counts := map[time.Duration]float64{}
	for couriers.Len() > 0 {
		e := couriers.Front()
		courier := e.Value.(*entities.Courier)
		couriers.Remove(e)
		if _, ok := counts[courier.ArrivingIn]; !ok {
			counts[courier.ArrivingIn] = 1
		} else {
			counts[courier.ArrivingIn] = counts[courier.ArrivingIn] + 1
		}
	}
	min := math.MaxFloat64
	max := float64(0)
	for k := range counts {
		min = math.Min(min, counts[k])
		max = math.Max(max, counts[k])
	}
	assert.LessOrEqual(t, max-min, float64(1))
}

func TestDispatch_Success(t *testing.T) {
	t.Parallel()

	instance := New(
		logger,
		cfg,
		nil,
		nil,
		nil,
	)

	_, err := instance.Dispatch()
	assert.NoError(t, err)
}

func TestDispatch_ErrNoCourierAvailable(t *testing.T) {
	t.Parallel()

	instance := New(
		logger,
		cfg,
		nil,
		nil,
		nil,
	)
	for i := 0; i < nCouriers; i++ {
		_, err := instance.Dispatch()
		assert.NoError(t, err)
	}
	_, err := instance.Dispatch()
	assert.Error(t, err)
	assert.Equal(t, err, errors.ErrNoCourierAvailable)
}

func TestWait_Success(t *testing.T) {
	t.Parallel()

	instance := New(
		logger,
		cfg,
		utils.NewOrderQueue(zap.L()),
		utils.NewCourierQueue(zap.L()),
		utils.NewStat(logger),
	)
	order := factories.Order()
	ch, err := instance.Wait(order)
	assert.NoError(t, err)
	dispatch := <-ch
	assert.Equal(t, order, dispatch.Order)
}
