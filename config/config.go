package config

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"path/filepath"
	"time"
)

type Struct struct {
	App struct {
		Port string `json:"port"`
	} `json:"app"`
	Database struct {
		URI  string `json:"uri"`
		Name string `json:"name"`
	} `json:"database"`
	Log struct {
		PathFile string `json:"path_file"`
	} `json:"log"`
}

type LogItem struct {
	Type    string    `json:"type"`
	Time    time.Time `json:"time"`
	Message string    `json:"message"`
}

var LogData []LogItem

var (
	Config Struct
	WD     string
	Client *mongo.Client
)

// HandleError handles errors and logs them.
func HandleError(err error, suppressLog ...bool) {
	if err != nil {
		if len(suppressLog) == 0 || !suppressLog[0] {
			WriteLog("info", "Application stopped")
		}

		color.Blue("====================================")
		color.Red("  Stopping Manzi Microservice...")
		color.Blue("====================================")
		log.Fatal(err)
	}
}

// ConnectMongoDB connects to the MongoDB database.
func ConnectMongoDB() {
	clientOptions := options.Client().
		ApplyURI(Config.Database.URI).
		SetConnectTimeout(5 * time.Second)

	var err error
	Client, err = mongo.Connect(context.TODO(), clientOptions)
	HandleError(err)

	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	if err = Client.Ping(ctx, nil); err != nil {
		WriteLog("error", fmt.Sprintf("Failed to ping MongoDB: %v", err))
		HandleError(err)
	}
}

// LoadConfig loads the application configuration from a JSON file.
func LoadConfig() {
	var err error
	WD, err = os.Getwd()
	HandleError(err)

	data, err := os.ReadFile(filepath.Join(WD, "config", "config.json"))
	HandleError(err)

	if err = json.Unmarshal(data, &Config); err != nil {
		HandleError(err)
	}
}

// InitLogFile initializes the log file.
func InitLogFile() {
	filePath := filepath.Join(WD, Config.Log.PathFile)
	dir := filepath.Dir(filePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.Mkdir(dir, 0755); err != nil {
			HandleError(err)
		}
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		file, err := os.Create(filePath)
		HandleError(err)
		defer file.Close()

		if err := json.NewEncoder(file).Encode([]LogItem{}); err != nil {
			HandleError(err)
		}
	}
}

// ReadLogJson reads the log data from the JSON file.
func ReadLogJson() {
	file, err := os.Open(filepath.Join(WD, Config.Log.PathFile))
	HandleError(err)
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&LogData); err != nil {
		HandleError(err)
	}
}

// WriteLog writes a new log entry to the log file.
func WriteLog(logType, message string) {
	newLog := LogItem{
		Type:    logType,
		Time:    time.Now(),
		Message: message,
	}

	LogData = append(LogData, newLog)

	filePath := filepath.Join(WD, Config.Log.PathFile)
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC, 0644)
	HandleError(err)
	defer file.Close()

	encoder := json.NewEncoder(file)
	HandleError(encoder.Encode(LogData))
}
