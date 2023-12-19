// internal/app/repository/csv_repository.go
package repository

import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"strconv"

	"github.com/shindesatish/titanic-service/pkg/model"
)

type CSVRepository struct {
	Path string
}

func NewCSVRepository(path string) *CSVRepository {
	return &CSVRepository{Path: path}
}

func (r *CSVRepository) GetAllPassengers() ([]model.Passenger, error) {
	file, err := os.Open(r.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to open CSV file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV records: %v", err)
	}

	var passengers []model.Passenger
	for i := 1; i < len(records); i++ {
		passenger, err := convertCSVRecordToPassenger(records[i])
		if err != nil {
			return nil, fmt.Errorf("failed to convert CSV record to Passenger: %v", err)
		}
		passengers = append(passengers, *passenger)
	}

	return passengers, nil
}

func (r *CSVRepository) GetPassengerByID(passengerID uint) (*model.Passenger, error) {
	file, err := os.Open(r.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to open CSV file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV records: %v", err)
	}

	for i := 1; i < len(records); i++ {
		record := records[i]

		id, err := strconv.Atoi(record[0])
		if err != nil {
			return nil, fmt.Errorf("failed to convert ID to integer: %v", err)
		}

		if uint(id) == passengerID {
			passenger, err := convertCSVRecordToPassenger(record)
			if err != nil {
				return nil, fmt.Errorf("failed to convert CSV record to Passenger: %v", err)
			}
			return passenger, nil
		}
	}

	return nil, fmt.Errorf("passenger not found with ID %d", passengerID)
}

func (r *CSVRepository) GetPassengerAttributes(passengerID uint, attributes []string) (*model.Passenger, error) {
	file, err := os.Open(r.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to open CSV file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV records: %v", err)
	}

	for i := 1; i < len(records); i++ {
		record := records[i]
		id, err := strconv.Atoi(record[0])
		if err != nil {
			return nil, fmt.Errorf("failed to convert ID to integer: %v", err)
		}

		if uint(id) == passengerID {
			passenger, err := convertCSVRecordToPassengerWithAttributes(record, attributes)
			if err != nil {
				return nil, fmt.Errorf("failed to convert CSV record to Passenger: %v", err)
			}
			return passenger, nil
		}
	}

	return nil, fmt.Errorf("passenger not found with ID %d", passengerID)
}

func (r *CSVRepository) GetFareHistogram() (map[string]int, error) {
	file, err := os.Open(r.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to open CSV file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV records: %v", err)
	}

	var fares []float64
	for i := 1; i < len(records); i++ {
		record := records[i]
		fare, err := strconv.ParseFloat(record[9], 64)
		if err != nil {
			return nil, fmt.Errorf("failed to convert Fare to float64: %v", err)
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

// Helper function to convert CSV record to Passenger
func convertCSVRecordToPassenger(record []string) (*model.Passenger, error) {
	// Ensure that the CSV record has the expected number of fields
	if len(record) != 12 {
		return nil, fmt.Errorf("invalid CSV record format")
	}

	// Parse fields from the CSV record
	passengerID, err := strconv.Atoi(record[0])
	if err != nil {
		return nil, fmt.Errorf("failed to convert ID to integer: %v", err)
	}

	survived, err := strconv.Atoi(record[1])
	if err != nil {
		return nil, fmt.Errorf("failed to convert Survived to integer: %v", err)
	}

	pclass, err := strconv.Atoi(record[2])
	if err != nil {
		return nil, fmt.Errorf("failed to convert Pclass to integer: %v", err)
	}

	var age float64
	fmt.Println(record)
	if record[5] != "" {
		age, err = strconv.ParseFloat(record[5], 64)
		if err != nil {
			return nil, fmt.Errorf("failed to convert Age to float64: %v", err)
		}
	}

	sibSp, err := strconv.Atoi(record[6])
	if err != nil {
		return nil, fmt.Errorf("failed to convert SibSp to integer: %v", err)
	}

	parch, err := strconv.Atoi(record[7])
	if err != nil {
		return nil, fmt.Errorf("failed to convert Parch to integer: %v", err)
	}

	fare, err := strconv.ParseFloat(record[9], 64)
	if err != nil {
		return nil, fmt.Errorf("failed to convert Fare to float64: %v", err)
	}

	// Create a model.Passenger instance
	passenger := &model.Passenger{
		PassengerID: passengerID,
		Survived:    survived,
		Pclass:      pclass,
		Name:        record[3],
		Sex:         record[4],
		Age:         age,
		SibSp:       sibSp,
		Parch:       parch,
		Ticket:      record[8],
		Fare:        fare,
		Cabin:       record[10],
		Embarked:    record[11],
	}

	return passenger, nil
}

// Helper function to convert CSV record to Passenger with specific attributes
func convertCSVRecordToPassengerWithAttributes(record []string, attributes []string) (*model.Passenger, error) {
	// Ensure that the CSV record has the expected number of fields
	if len(record) != 12 {
		return nil, fmt.Errorf("invalid CSV record format")
	}

	// Create a map to store the attribute values
	attributeValues := make(map[string]interface{})

	// Parse fields from the CSV record and store in the map
	for _, attribute := range attributes {
		switch attribute {
		case "PassengerId":
			value, err := strconv.Atoi(record[0])
			if err != nil {
				return nil, fmt.Errorf("failed to convert ID to integer: %v", err)
			}
			attributeValues["PassengerId"] = uint(value)
		case "Survived":
			value, err := strconv.Atoi(record[1])
			if err != nil {
				return nil, fmt.Errorf("failed to convert Survived to integer: %v", err)
			}
			attributeValues["Survived"] = value
		case "Pclass":
			value, err := strconv.Atoi(record[2])
			if err != nil {
				return nil, fmt.Errorf("failed to convert Pclass to integer: %v", err)
			}
			attributeValues["Pclass"] = value
		case "Name":
			attributeValues["Name"] = record[3]
		case "Sex":
			attributeValues["Sex"] = record[4]
		case "Age":
			value, err := strconv.ParseFloat(record[5], 64)
			if err != nil {
				return nil, fmt.Errorf("failed to convert Age to float64: %v", err)
			}
			attributeValues["Age"] = value
		case "SibSp":
			value, err := strconv.Atoi(record[6])
			if err != nil {
				return nil, fmt.Errorf("failed to convert SibSp to integer: %v", err)
			}
			attributeValues["SibSp"] = value
		case "Parch":
			value, err := strconv.Atoi(record[7])
			if err != nil {
				return nil, fmt.Errorf("failed to convert Parch to integer: %v", err)
			}
			attributeValues["Parch"] = value
		case "Ticket":
			attributeValues["Ticket"] = record[8]
		case "Fare":
			value, err := strconv.ParseFloat(record[9], 64)
			if err != nil {
				return nil, fmt.Errorf("failed to convert Fare to float64: %v", err)
			}
			attributeValues["Fare"] = value
		case "Cabin":
			attributeValues["Cabin"] = record[10]
		case "Embarked":
			attributeValues["Embarked"] = record[11]
		// Add cases for other attributes as needed
		default:
			return nil, fmt.Errorf("unknown attribute: %s", attribute)
		}
	}

	id := 0
	if val, ok := attributeValues["PassengerId"]; ok && val != nil {
		// Attempt to type assert only if the value is not nil
		id = val.(int)
	}

	// Create a model.Passenger instance with specific attributes
	passenger := &model.Passenger{
		PassengerID: id,
		Survived:    attributeValues["Survived"].(int),
		Pclass:      attributeValues["Pclass"].(int),
		Name:        attributeValues["Name"].(string),
		Sex:         attributeValues["Sex"].(string),
		Age:         attributeValues["Age"].(float64),
		SibSp:       attributeValues["SibSp"].(int),
		Parch:       attributeValues["Parch"].(int),
		Ticket:      attributeValues["Ticket"].(string),
		Fare:        attributeValues["Fare"].(float64),
		Cabin:       attributeValues["Cabin"].(string),
		Embarked:    attributeValues["Embarked"].(string),
		// Add assignments for other attributes as needed
	}

	return passenger, nil
}
