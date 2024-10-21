package config

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"path/filepath"
)

type Struct struct {
	App struct {
		Port string `json:"port"`
	} `json:"app"`
	Database struct {
		URI  string `json:"uri"`
		Name string `json:"name"`
	} `json:"database"`
}

var Config Struct
var WD string
var Client *mongo.Client

func ConnectMongoDB() {
	var err error
	clientOptions := options.Client().ApplyURI(Config.Database.URI)
	Client, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = Client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func LoadConfig() {
	var err error

	WD, err = os.Getwd()
	if err != nil {
		fmt.Println(err)
		return
	}

	data, err := os.ReadFile(filepath.Join(WD, "config", "config.json"))
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(data, &Config)
	if err != nil {
		log.Fatal(err)
	}
}
