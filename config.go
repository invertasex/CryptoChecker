package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func GetBotToken() string {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env файл не найден, использую переменные окружения")
	}
	token := os.Getenv("TOKEN")
	if token == "" {
		log.Fatal("BOT TOKEN не задан в .env или переменных окружения")
	}
	return token
}
