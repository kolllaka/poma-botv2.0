package main

import (
	"log"

	"github.com/KoLLlaka/poma-botv2.0/internal/db"
)

func main() {
	if err := db.DownDB("static/db/music.db"); err != nil {
		log.Println(err)
	}
}
