package errors

import "errors"

// ErrOrderPrepNotStarted stands for the order hasn't been started for prep
var ErrOrderPrepNotStarted = errors.New("order prep not started")
