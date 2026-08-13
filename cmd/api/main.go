package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"test.nahueldev23.com/internal/config"
	"test.nahueldev23.com/internal/database"
	"test.nahueldev23.com/internal/logging"
	"test.nahueldev23.com/internal/server"
	"test.nahueldev23.com/internal/session"
)

func main() {
	logger := logging.New(
		slog.New(
			slog.NewJSONHandler(os.Stdout, nil),
		),
	)

	cfg, err := config.Load()
	if err != nil {
		logger.Error(context.Background(), "failed to load config", "error", err)
		os.Exit(1)
	}

	db, err := database.New(cfg.DatabaseURL)
	if err != nil {
		logger.Error(context.Background(), "failed to connect to database", "error", err)
		os.Exit(1)
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
			logger.Error(context.Background(), "session cleanup failed", "error", err)
			return
		}

		logger.Info(
			context.Background(),
			"session cleanup completed",
			"retention_days", cfg.SessionRetentionDays,
		)
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

	router := server.New(db, cfg, logger)

	logger.Info(context.Background(), "server starting", "port", cfg.Port)

	err = router.Run(":" + cfg.Port)
	if err != nil {
		logger.Error(context.Background(), "server stopped", "error", err)
		os.Exit(1)
	}
}
