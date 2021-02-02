package errors

import "errors"

// ErrNoCourierAvailable stands for no available courier in the queue
var ErrNoCourierAvailable = errors.New("no courier available")
