package dto

// Validate the accepted attributes
var AllowedAttributes = []string{
	"PassengerID",
	"Survived",
	"Pclass",
	"Name",
	"Sex",
	"Age",
	"SibSp",
	"Parch",
	"Ticket",
	"Fare",
	"Cabin",
	"Embarked",
}

func ValidAttribute(key string) bool {
	return Contains(AllowedAttributes, key)
}

// contains checks if a string is present in a slice
func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
