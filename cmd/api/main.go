package main

import (
	"context"
	"log"
	"time"

	"test.nahueldev23.com/internal/config"
	"test.nahueldev23.com/internal/database"
	"test.nahueldev23.com/internal/server"
	"test.nahueldev23.com/internal/session"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.New(cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	sessionRepository := session.NewRepository(db)

	cleanup := func() {
		cutoff := time.Now().AddDate(
			0,
			0,
			-cfg.SessionRetentionDays,
		)

		err := sessionRepository.DeleteOldSessions(
			context.Background(),
			cutoff,
		)

		if err != nil {
			log.Printf("session cleanup failed: %v", err)
		}
	}

	// Cleanup inicial
	cleanup()

	// Cleanup cada 24 horas
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	go func() {
		for {
			<-ticker.C
			cleanup()
		}
	}()

	router := server.New(db, cfg)

	log.Println("server running on port", cfg.Port)

	err = router.Run(":" + cfg.Port)
	if err != nil {
		log.Fatal(err)
	}
}
