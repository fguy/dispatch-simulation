package utils

import (
	"sync"
	"time"

	"go.uber.org/zap"
)

// StatKind is the kind of statistic
type StatKind int

const (
	// StatKindCourierWaitTime -
	StatKindCourierWaitTime = iota
	// StatKindFoodWaitTime -
	StatKindFoodWaitTime
)

func (s StatKind) String() string {
	return [...]string{
		"courier wait time between order ready and pickup",
		"food wait time between arrival and order pickup",
	}[s]
}

// Stat is to collect time series
type Stat struct {
	logger *zap.Logger
	data   *sync.Map
}

// Emit emits a metric point
func (s *Stat) Emit(kind StatKind, duration time.Duration) {
	v, _ := s.data.LoadOrStore(kind, []time.Duration{})
	timeSeries := v.([]time.Duration)
	s.data.Store(kind, append(timeSeries, duration))
}

// PrintAvgAll prints average statictics for all kinds
func (s *Stat) PrintAvgAll(extraFields ...zap.Field) {
	s.data.Range(func(k, v interface{}) bool {
		s.printAvg(k.(StatKind), extraFields...)
		return true
	})
}

func (s *Stat) printAvg(kind StatKind, fields ...zap.Field) {
	s.logger.Info("Average time taken (milliseconds)", append(fields, zap.String("kind", kind.String()), zap.Int64("time_ms", s.avg(kind).Milliseconds()))...)
}

func (s *Stat) avg(kind StatKind) time.Duration {
	sum := time.Duration(0)
	v, ok := s.data.Load(kind)
	if !ok {
		s.logger.Debug("empty data", zap.String("kind", kind.String()))
		return 0
	}
	timeSeries := v.([]time.Duration)
	for _, item := range timeSeries {
		sum += item
	}
	if sum == 0 {
		return 0
	}
	return sum / time.Duration(len(timeSeries))
}

// NewStat returns a new instance of Stat
func NewStat(logger *zap.Logger) *Stat {
	return &Stat{
		logger: logger,
		data:   &sync.Map{},
	}
}
