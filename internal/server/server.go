package server

import (
	"log"
	"net/http"
	"time"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/internal/handlers"
)


type Server struct {
	logger *log.Logger
	server *http.Server
}


func NewServer(logger *log.Logger) *Server {
	// Создаем роутер
	mux := http.NewServeMux()

	
	mux.HandleFunc("/", handlers.HomeHandler)
	mux.HandleFunc("/upload", handlers.UploadHandler)

	
	httpServer := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ErrorLog:     logger,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	return &Server{
		logger: logger,
		server: httpServer,
	}
}


func (s *Server) Run() error {
	s.logger.Println("Сервер запущен на порту 8080")
	return s.server.ListenAndServe()
}
