package main

import (
	"log"

	"test.nahueldev23.com/internal/config"
	"test.nahueldev23.com/internal/database"
	"test.nahueldev23.com/internal/server"
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

	router := server.New(db, cfg)

	log.Println("server running on port", cfg.Port)

	err = router.Run(":" + cfg.Port)
	if err != nil {
		log.Fatal(err)
	}

}
