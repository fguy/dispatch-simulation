package fulfillment

import (
	"testing"
	"time"

	"github.com/fguy/dispatch-simulation/config"
	"github.com/fguy/dispatch-simulation/entities"
	mock_dispatch "github.com/fguy/dispatch-simulation/mocks/services/dispatch"
	"github.com/fguy/dispatch-simulation/utils"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

var logger, _ = zap.NewDevelopment()

func Test_Success(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)

	mockDispatchService := mock_dispatch.NewMockInterface(mockCtrl)
	mockDispatchService.EXPECT().Receive(gomock.Any(), gomock.Any()).Return(nil).Times(132)
	instance, err := New(
		logger,
		&config.AppConfig{
			Strategy: entities.StrategyMatched,
		},
		mockDispatchService,
		utils.NewStat(logger),
		utils.NewOrderQueue(logger),
		utils.NewCourierQueue(logger),
	)
	assert.NoError(t, err)
	timer := time.NewTimer(3 * time.Second)
	go func() {
		<-timer.C
		defer instance.Stop()
	}()
	instance.Start()
	time.Sleep(3 * tick)
}
