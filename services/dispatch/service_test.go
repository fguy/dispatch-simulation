package dispatch

import (
	"errors"
	"testing"
	"time"

	"github.com/fguy/dispatch-simulation/entities"
	"github.com/fguy/dispatch-simulation/factories"
	mock_courier "github.com/fguy/dispatch-simulation/mocks/services/courier"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

var logger, _ = zap.NewDevelopment()

func TestReceive_Matched(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	courier := factories.Courier()
	order := factories.Order()

	mockCourierService := mock_courier.NewMockInterface(mockCtrl)
	mockCourierService.EXPECT().Dispatch().Return(courier, nil).Times(2)
	mockCourierService.EXPECT().Deliver(gomock.Any()).Return(nil).Times(1)
	mockCourierService.EXPECT().Deliver(gomock.Any()).Return(errors.New("")).Times(1)

	instance := New(
		logger,
		mockCourierService,
	)

	err := instance.Receive(entities.StrategyMatched, order)
	assert.NoError(t, err)

	err = instance.Receive(entities.StrategyMatched, order)
	assert.Error(t, err)
}

func TestReceive_FirstInFirstOut(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	dispatch := factories.Dispatch()

	mockCourierService := mock_courier.NewMockInterface(mockCtrl)
	//chan *entities.Dispatch, error
	ch := make(chan *entities.Dispatch, 1)
	go func() {
		ch <- dispatch
	}()
	mockCourierService.EXPECT().Wait(dispatch.Order).Return(ch, nil).Times(1)
	mockCourierService.EXPECT().Deliver(gomock.Any()).Return(nil).Times(1)

	instance := New(
		logger,
		mockCourierService,
	)

	err := instance.Receive(entities.StrategyFirstInFirstOut, dispatch.Order)
	time.Sleep(dispatch.Courier.ArrivingIn + dispatch.Order.PrepTime*time.Second)
	time.Sleep(1 * time.Second)
	assert.NoError(t, err)
}
