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

type LogDataStruct struct {
	Log []LogItem `json:"log"`
}

var (
	Config  Struct
	WD      string
	Client  *mongo.Client
	LogData LogDataStruct
)

func HandleError(err error, suppressLog ...bool) {
	if err != nil {
		// Проверяем, передан ли параметр suppressLog, и его значение
		if len(suppressLog) == 0 || !suppressLog[0] {
			WriteLog("info", "Application stopped")
		}

		color.Blue("====================================")
		color.Red("  Stopping Manzi Microservice...")
		color.Blue("====================================")
		log.Fatal(err)
	}
}

func ConnectMongoDB() {
	clientOptions := options.Client().
		ApplyURI(Config.Database.URI).
		SetConnectTimeout(5 * time.Second) // Устанавливаем тайм-аут для подключения

	var err error
	Client, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		WriteLog("error", fmt.Sprintf("Failed to connect to MongoDB: %v", err))
		HandleError(err)
	}

	// Устанавливаем тайм-аут для пинга
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	err = Client.Ping(ctx, nil)
	if err != nil {
		WriteLog("error", fmt.Sprintf("Failed to ping MongoDB: %v", err))
		HandleError(err)
	}
}

func LoadConfig() {
	var err error

	WD, err = os.Getwd()
	if err != nil {
		HandleError(err, true)
	}

	data, err := os.ReadFile(filepath.Join(WD, "config", "config.json"))
	if err != nil {
		HandleError(err, true)
	}

	err = json.Unmarshal(data, &Config)
	if err != nil {
		HandleError(err, true)
	}
}

func InitLogFile() {
	filePath := filepath.Join(WD, Config.Log.PathFile)
	dir := filepath.Dir(filePath)
	// Проверка и создание директории, если она не существует
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.Mkdir(dir, 0755); err != nil {
			HandleError(err, true)
		}
	}

	// Проверка и создание файла log.json, если его нет
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		initialData := map[string][]struct{}{"log": {}}

		// Создание и запись шаблонного JSON в файл
		file, err := os.Create(filePath)
		if err != nil {
			HandleError(err, true)
		}
		defer file.Close()

		if err := json.NewEncoder(file).Encode(initialData); err != nil {
			HandleError(err, true)
		}
	}
}

func ReadLogJson() {
	file, err := os.Open(filepath.Join(WD, Config.Log.PathFile))
	if err != nil {
		HandleError(err, true)
	}

	defer func() {
		if cerr := file.Close(); cerr != nil {
			HandleError(err, true)
		}
	}()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&LogData); err != nil {
		HandleError(err, true)
	}
}

func WriteLog(logType, message string) {
	newLog := LogItem{
		Type:    logType,
		Time:    time.Now(),
		Message: message,
	}

	LogData.Log = append(LogData.Log, newLog)

	filePath := filepath.Join(WD, Config.Log.PathFile)
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		HandleError(err, true)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(LogData); err != nil {
		HandleError(err, true)
	}
}
