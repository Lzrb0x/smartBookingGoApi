package main

import (
	"fmt"
	"log"

	"github.com/Lzrb0x/smartBookingGoApi/internal/config"
	"github.com/Lzrb0x/smartBookingGoApi/internal/database"
	"github.com/Lzrb0x/smartBookingGoApi/internal/server"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables from system")
	}

	cfg, err := config.Load()
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %s", err))
	}

	db, err := database.New(cfg)
	if err != nil {
		panic(fmt.Sprintf("failed to connect to database: %s", err))
	}

	server := server.NewServer(cfg, db)
	err = server.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("http server error: %s", err))
	}
}
