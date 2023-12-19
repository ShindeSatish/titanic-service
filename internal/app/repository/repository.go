package repository

import (
	"strings"

	"github.com/shindesatish/titanic-service/pkg/model"
)

type Repository interface {
	GetAllPassengers() ([]model.Passenger, error)
	GetPassengerByID(passengerID uint) (*model.Passenger, error)
	GetPassengerAttributes(passengerID uint, attributes []string) (*model.Passenger, error)
	GetFareHistogram() (map[string]int, error)
}

// JoinAttributes joins a list of attributes into a comma-separated string
func JoinAttributes(attributes []string) string {
	return strings.Join(attributes, ", ")
}
