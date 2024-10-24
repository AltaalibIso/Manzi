package repositories

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"manzi/config"
)

func CheckUserExists(username string) (bool, error) {
	collection := config.Client.Database(config.Config.Database.Name).Collection("users")

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

	_, err := collection.InsertOne(context.TODO(), bson.M{"username": username, "password": password})
	if err != nil {
		return err
	}

	return nil
}
