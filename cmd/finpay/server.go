package main

import (
	"finpay/internal/config"
	routes "finpay/internal/routes"
	"finpay/pkg/db"

	"context"
	"log"
)

func main() {
	cfg := config.Load()

	mongo := db.Connect(cfg.MONGODB_URI, cfg.DBName)

	defer mongo.Client.Disconnect(context.Background())

	r := routes.Routes(mongo)

	r.Run(":" + cfg.Port)

	log.Printf("Finpay running on localhost:%s", cfg.Port)
}
