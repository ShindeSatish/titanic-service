// internal/app/repository/sqlite_repository.go
package repository

import (
	"database/sql"
	"fmt"
	"sort"

	"github.com/shindesatish/titanic-service/pkg/model"
)

type SQLiteRepository struct {
	DB *sql.DB
}

func NewSQLiteRepository(db *sql.DB) *SQLiteRepository {
	return &SQLiteRepository{DB: db}
}

func (r *SQLiteRepository) GetAllPassengers() ([]model.Passenger, error) {
	rows, err := r.DB.Query("SELECT * FROM titanic")
	if err != nil {
		return nil, fmt.Errorf("failed to query passengers: %v", err)
	}
	defer rows.Close()

	var passengers []model.Passenger
	for rows.Next() {
		var passenger model.Passenger
		var cabin, embarked sql.NullString
		var age sql.NullFloat64

		err := rows.Scan(
			&passenger.PassengerID, &passenger.Survived, &passenger.Pclass, &passenger.Name,
			&passenger.Sex, &age, &passenger.SibSp, &passenger.Parch,
			&passenger.Ticket, &passenger.Fare, &cabin, &embarked,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan passenger row: %v", err)
		}
		passenger.Age = age.Float64
		passenger.Cabin = cabin.String
		passenger.Embarked = embarked.String
		passengers = append(passengers, passenger)
	}

	return passengers, nil
}

func (r *SQLiteRepository) GetPassengerByID(passengerID uint) (*model.Passenger, error) {
	row := r.DB.QueryRow("SELECT * FROM titanic WHERE PassengerID = ?", passengerID)
	var passenger model.Passenger
	var cabin, embarked sql.NullString
	var age sql.NullFloat64
	err := row.Scan(
		&passenger.PassengerID, &passenger.Survived, &passenger.Pclass, &passenger.Name,
		&passenger.Sex, &age, &passenger.SibSp, &passenger.Parch,
		&passenger.Ticket, &passenger.Fare, &cabin, &embarked,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get passenger by ID: %v", err)
	}
	passenger.Age = age.Float64
	passenger.Cabin = cabin.String
	passenger.Embarked = embarked.String
	return &passenger, nil
}

func (r *SQLiteRepository) GetPassengerAttributes(passengerID uint, attributes []string) (*model.Passenger, error) {
	// Ensure attributes are not empty to prevent SQL injection
	if len(attributes) == 0 {
		return nil, fmt.Errorf("attributes cannot be empty")
	}

	// Construct SQL query with prepared statement
	query := "SELECT " + JoinAttributes(attributes) + " FROM titanic WHERE PassengerID = ?"
	stmt, err := r.DB.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("error preparing SQL statement: %v", err)
	}
	defer stmt.Close()

	// Execute the prepared statement
	row := stmt.QueryRow(passengerID)

	// Create a map to store attribute values
	attributeValues := make(map[string]interface{})

	// Scan the row and populate the attributeValues map
	for _, attr := range attributes {
		var value interface{}
		err := row.Scan(&value)
		if err != nil {
			return nil, fmt.Errorf("error scanning attribute %s: %v", attr, err)
		}
		attributeValues[attr] = value
	}

	// Convert the attributeValues map to a Passenger instance
	passenger, err := ConvertAttributesToPassenger(attributeValues)
	if err != nil {
		return nil, fmt.Errorf("error converting attributes to Passenger: %v", err)
	}

	return passenger, nil
}

func (r *SQLiteRepository) GetFareHistogram() (map[string]int, error) {
	rows, err := r.DB.Query("SELECT Fare FROM titanic")
	if err != nil {
		return nil, fmt.Errorf("failed to query fares: %v", err)
	}
	defer rows.Close()

	var fares []float64
	for rows.Next() {
		var fare float64
		err := rows.Scan(&fare)
		if err != nil {
			return nil, fmt.Errorf("failed to scan fare row: %v", err)
		}
		fares = append(fares, fare)
	}

	// Sort the fares in ascending order
	sort.Float64s(fares)

	// Calculate percentiles
	percentiles := []float64{25, 50, 75, 90, 95, 99}
	percentileValues := make(map[string]float64)
	for _, p := range percentiles {
		idx := int(float64(len(fares)-1) * (p / 100.0))
		percentileValues[fmt.Sprintf("%.2f%%", p)] = fares[idx]
	}

	// Create histogram buckets and count occurrences
	fareHistogram := make(map[string]int)
	for _, fare := range fares {
		for label, pValue := range percentileValues {
			if fare <= pValue {
				fareHistogram[label]++
				break
			}
		}
	}

	return fareHistogram, nil
}

// ConvertAttributesToPassenger converts a map of attribute values to a Passenger instance
func ConvertAttributesToPassenger(attributeValues map[string]interface{}) (*model.Passenger, error) {
    passenger := &model.Passenger{}

    for key, value := range attributeValues {
        switch key {
        case "PassengerID":
            passenger.PassengerID = value.(int)
        case "Survived":
            passenger.Survived = value.(int)
        case "Pclass":
            passenger.Pclass = value.(int)
        case "Name":
            passenger.Name = value.(string)
        case "Sex":
            passenger.Sex = value.(string)
        case "Age":
            if value != nil {
                passenger.Age = value.(float64)
            }
        case "SibSp":
            passenger.SibSp = value.(int)
        case "Parch":
            passenger.Parch = value.(int)
        case "Ticket":
            passenger.Ticket = value.(string)
        case "Fare":
            if value != nil {
                passenger.Fare =  value.(float64) 
            }
        case "Cabin":
            if value != nil {
                passenger.Cabin = value.(string)
            }
        case "Embarked":
            if value != nil {
                passenger.Embarked = value.(string)
            }
        // Add more cases for other fields
        default:
            return nil, fmt.Errorf("unknown attribute: %s", key)
        }
    }

    return passenger, nil
}