package factories

import (
	"time"

	"github.com/fguy/dispatch-simulation/entities"
	"github.com/google/uuid"
)

// Order creates a default order
func Order() *entities.Order {
	return &entities.Order{
		ID:       uuid.MustParse("0ff534a7-a7c4-48ad-b6ec-7632e36af950"),
		Name:     "Cheese Pizza",
		PrepTime: time.Duration(13),
	}
}
