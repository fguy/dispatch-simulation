package factories

import (
	"math/rand"
	"time"

	"github.com/fguy/dispatch-simulation/entities"
	"github.com/google/uuid"
)

// Courier creates a default courier
func Courier() *entities.Courier {
	return &entities.Courier{
		ID:         uuid.New(),
		ArrivingIn: time.Duration(rand.Int63n(int64(12))),
	}
}
