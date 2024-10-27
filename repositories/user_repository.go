package repositories

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"manzi/config"
	"regexp"
)

var (
	ErrInvalidInput     = errors.New("input contains dangerous characters")
	DangerousCharsRegex = regexp.MustCompile(`[{}\[\]$:"';\\/*&|%<>"()+=?~` + "`" + `#]`)
)

// ValidateInput checks for dangerous characters in the input.
func ValidateInput(input string) error {
	if DangerousCharsRegex.MatchString(input) {
		return ErrInvalidInput
	}
	return nil
}

// CheckUserExists checks if a user exists in the database.
func CheckUserExists(username string) (bool, error) {
	collection := config.Client.Database(config.Config.Database.Name).Collection("users")
	if err := ValidateInput(username); err != nil {
		return false, err
	}

	err := collection.FindOne(context.TODO(), bson.M{"username": username}).Err()
	if errors.Is(err, mongo.ErrNoDocuments) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

// CreateUser inserts a new user into the database.
func CreateUser(username, password string) error {
	collection := config.Client.Database(config.Config.Database.Name).Collection("users")
	if err := ValidateInput(username); err != nil {
		return err
	}
	if err := ValidateInput(password); err != nil {
		return err
	}
	_, err := collection.InsertOne(context.TODO(), bson.M{"username": username, "password": password})
	if err != nil {
		return err
	}

	return nil
}
