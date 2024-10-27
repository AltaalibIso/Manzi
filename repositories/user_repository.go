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

func ValidateInput(input string) error {
	if DangerousCharsRegex.MatchString(input) {
		return ErrInvalidInput
	}
	return nil
}

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

func CreateUser(username, password string) error {
	collection := config.Client.Database(config.Config.Database.Name).Collection("users")
	if err := ValidateInput(username); err != nil {
		return err
	}
	if err := ValidateInput(password); err != nil {
		print(err.Error())
		return err
	}
	_, err := collection.InsertOne(context.TODO(), bson.M{"username": username, "password": password})
	if err != nil {
		return err
	}

	return nil
}
