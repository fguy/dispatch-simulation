package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStat(t *testing.T) {
	t.Parallel()

	instance := NewStat(logger)

	assert.EqualValues(t, 0, instance.avg(StatKindCourierWaitTime))
	assert.EqualValues(t, 0, instance.avg(StatKindFoodWaitTime))

	instance.Emit(StatKindCourierWaitTime, 1)
	instance.Emit(StatKindCourierWaitTime, 2)
	instance.Emit(StatKindCourierWaitTime, 3)

	instance.Emit(StatKindFoodWaitTime, 4)
	instance.Emit(StatKindFoodWaitTime, 5)
	instance.Emit(StatKindFoodWaitTime, 6)

	assert.EqualValues(t, 2, instance.avg(StatKindCourierWaitTime))
	assert.EqualValues(t, 5, instance.avg(StatKindFoodWaitTime))

	instance.PrintAvgAll()
}
