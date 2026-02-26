package main

import (
	"log"

	"red_packet/backend/internal/config"
	"red_packet/backend/internal/database"
	"red_packet/backend/internal/http/router"
	"red_packet/backend/internal/service"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config failed: %v", err)
	}

	db, err := database.Init(cfg)
	if err != nil {
		log.Fatalf("init db failed: %v", err)
	}

	services := service.NewContainer(db, cfg)
	e := router.New(services, cfg)
	addr := ":" + cfg.Server.Port
	if err := e.Start(addr); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
