package main

import (
	"log"

	"cms/server/internal/config"
	"cms/server/internal/db"
	"cms/server/internal/transport/http"
)

func main() {
	cfg := config.Load()
	dbConn := db.MustOpen(cfg)

	r := http.NewRouter(cfg, dbConn)
	addr := ":" + cfg.AppPort
	log.Printf("server listening on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}
}
