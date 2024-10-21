package main

import (
	"log"
	"manzi/config"
	"manzi/routes"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	log.Println("Starting application...")

	config.LoadConfig()
	log.Println("Loaded config:", config.Config)

	config.ConnectMongoDB()
	log.Println("Connected to MongoDB!")

	router := gin.Default()
	routes.InitUserRoutes(router)

	log.Println("Listening on", config.Config.App.Port)
	if err := http.ListenAndServe(config.Config.App.Port, router); err != nil {
		log.Fatal(err)
	}
}
