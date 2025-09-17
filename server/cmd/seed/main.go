package main

import (
	"log"

	"cms/server/internal/config"
	"cms/server/internal/db"
	"cms/server/internal/seed"
)

func main() {
	cfg := config.Load()
	dbConn := db.MustOpen(cfg)
	if err := seed.Run(dbConn); err != nil {
		log.Fatal(err)
	}
	log.Println("seed completed")
}
