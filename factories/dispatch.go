package factories

import "github.com/fguy/dispatch-simulation/entities"

// Dispatch creates a default dispatch
func Dispatch() *entities.Dispatch {
	return &entities.Dispatch{
		Order:   Order(),
		Courier: Courier(),
	}
}
