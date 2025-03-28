// main.go
package main

import (
	"context"
	"log"

	"github.com/ricardomussett/gotest/config"
	"github.com/ricardomussett/gotest/handlers"
	"github.com/ricardomussett/gotest/services"

	"github.com/gofiber/fiber/v2"
)

func main() {
	cfg := config.LoadConfig()

	// Inicializar MongoDB
	mongoSvc, err := services.NewMongoService(cfg.MongoDBURI, cfg.DBName, cfg.Collection)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongoSvc.Client.Disconnect(context.Background())

	// Iniciar servidor TCP
	tcpServer := services.NewTCPServer(cfg.TCPPort, mongoSvc)
	go tcpServer.Start()

	// Iniciar servidor HTTP con Fiber
	app := fiber.New()

	app.Get("/status", handlers.StatusHandler)

	log.Fatal(app.Listen(":3000"))
}
