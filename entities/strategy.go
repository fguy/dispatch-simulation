package entities

// Strategy is a courier dispatch strategy
type Strategy string

const (
	// StrategyMatched means that a courier is dispatched for a specific order and may only pick up that order
	StrategyMatched Strategy = "matched"
	// StrategyFirstInFirstOut means that a courier picks up the next available order upon arrival.
	StrategyFirstInFirstOut Strategy = "fifo"
)
