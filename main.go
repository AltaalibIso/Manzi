package main

import (
	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
	"log"
	"manzi/config"
	"manzi/routes"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			config.WriteLog("info", "Application stopped due to panic")
			color.Blue("====================================")
			color.Red("  Stopping Manzi Microservice due to panic...")
			color.Blue("====================================")
			log.Fatal(r)
		}
	}()

	color.Blue("====================================")
	color.Cyan("  Starting Manzi Microservice...")
	color.Blue("====================================")

	config.LoadConfig()
	config.InitLogFile()
	config.ReadLogJson()
	config.WriteLog("info", "Application started")

	config.ConnectMongoDB()
	config.WriteLog("info", "Connected to MongoDB successfully!")

	// Настройка канала для обработки сигналов
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	router := gin.Default()
	routes.InitUserRoutes(router)

	color.Green("Listening on port: %s", config.Config.App.Port)

	go func() {
		if err := http.ListenAndServe(config.Config.App.Port, router); err != nil {
			config.WriteLog("error", err.Error())
			log.Fatal(err)
		}
	}()

	// Ожидание сигнала завершения
	<-signalChan
	config.WriteLog("info", "Application stopped due to external signal")
	color.Blue("====================================")
	color.Red("  Stopping Manzi Microservice due to external signal...")
	color.Blue("====================================")
}
