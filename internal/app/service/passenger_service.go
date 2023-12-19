// internal/app/service/passenger_service.go
package service

import (
	"context"

	"github.com/shindesatish/titanic-service/pkg/model"
)

type Repository interface {
	GetAllPassengers() ([]model.Passenger, error)
	GetPassengerByID(passengerID uint) (*model.Passenger, error)
	GetPassengerAttributes(passengerID uint, attributes []string) (*model.Passenger, error)
	GetFareHistogram() (map[string]int, error)
}

type PassengerService struct {
	Repository Repository
}

func NewPassengerService(repository Repository) *PassengerService {
	return &PassengerService{Repository: repository}
}

func (s *PassengerService) GetAllPassengers(ctx context.Context) ([]model.Passenger, error) {
	return s.Repository.GetAllPassengers()
}

func (s *PassengerService) GetPassengerByID(ctx context.Context, passengerID uint) (*model.Passenger, error) {
	return s.Repository.GetPassengerByID(passengerID)
}

func (s *PassengerService) GetPassengerAttributes(ctx context.Context, passengerID uint, attributes []string) (*model.Passenger, error) {
	return s.Repository.GetPassengerAttributes(passengerID, attributes)
}

func (s *PassengerService) GetFareHistogram(ctx context.Context) (map[string]int, error) {
	return s.Repository.GetFareHistogram()
}
