package main

import (
	"log"
	"os"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/internal/server" 
)

func main() {
	// Создаем логгер
	logger := log.New(os.Stdout, "MORSE-CONVERTER: ", log.LstdFlags|log.Lshortfile)

	// Создаем сервер
	srv := server.NewServer(logger)

	// Запускаем сервер
	if err := srv.Run(); err != nil {
		logger.Fatal("Ошибка при запуске сервера:", err)
	}
}
