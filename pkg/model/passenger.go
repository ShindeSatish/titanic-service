// model/passenger.go
package model

import "database/sql"

type Passenger struct {
	PassengerID int     `json:"PassengerId"`
	Survived    int     `json:"Survived"`
	Pclass      int     `json:"Pclass"`
	Name        string  `json:"Name"`
	Sex         string  `json:"Sex"`
	Age         float64 `json:"Age"`
	SibSp       int     `json:"SibSp"`
	Parch       int     `json:"Parch"`
	Ticket      string  `json:"Ticket"`
	Fare        float64 `json:"Fare"`
	Cabin       string  `json:"Cabin"`
	Embarked    string  `json:"Embarked"`
}

// NullFloat64 represents a float64 that may be null.
type NullFloat64 struct {
	sql.NullFloat64
}
