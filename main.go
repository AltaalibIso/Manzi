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
		// Handle panic
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

	config.LoadConfig()                            // Load configuration
	config.InitLogFile()                           // Initialize log file
	config.ReadLogJson()                           // Read logs in JSON format
	config.WriteLog("info", "Application started") // Log application start

	config.ConnectMongoDB()                                       // Connect to MongoDB
	config.WriteLog("info", "Connected to MongoDB successfully!") // Log successful connection

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM) // Handle termination signals

	router := gin.Default()       // Initialize router
	routes.InitUserRoutes(router) // Initialize user routes

	color.Green("Listening on port: %s", config.Config.App.Port) // Output port

	go func() {
		// Start HTTP server
		if err := http.ListenAndServe(config.Config.App.Port, router); err != nil {
			config.WriteLog("error", err.Error())
			log.Fatal(err)
		}
	}()

	<-signalChan                                                          // Wait for termination signal
	config.WriteLog("info", "Application stopped due to external signal") // Log application stop
	color.Blue("====================================")
	color.Red("  Stopping Manzi Microservice due to external signal...")
	color.Blue("====================================")
}
