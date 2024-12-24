package main

import (
	"log"
	"os"

	"github.com/alfanzain/custom-agent-allocation/config"
	"github.com/alfanzain/custom-agent-allocation/handlers"
	"github.com/alfanzain/custom-agent-allocation/services"
	"github.com/labstack/echo/v4"
)

func main() {
	config.Init()

	e := echo.New()

	// Logging purpose
	// redisPolling := pollings.NewRedisPolling(config.RedisClient)
	// go redisPolling.StartRedisPolling()

	qiscusService := services.NewQiscusService(os.Getenv("QISCUS_BASE_URL"), os.Getenv("QISCUS_APP_ID"), os.Getenv("QISCUS_SECRET_KEY"))
	queueService := services.NewQueueService(config.RedisClient, config.Ctx)
	agentService := services.NewAgentService(config.DB)

	allocateAgentHandler := handlers.NewAllocateAgentHandler(qiscusService, queueService, agentService)
	markAsSolvedHandler := handlers.NewMarkAsSolvedHandler(qiscusService, queueService, agentService)

	e.POST("/allocate-agent/webhook", allocateAgentHandler.AllocateAgentWebhook)
	e.POST("/mark-as-solved/webhook", markAsSolvedHandler.MarkAsSolvedWebhook)
	e.GET("/alive", handlers.AliveCheck)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s...\n", port)
	if err := e.Start(":" + port); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
