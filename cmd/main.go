package main

import (
	"github.com/kiricle/api-homework/internal/handlers"
	"github.com/kiricle/api-homework/internal/router"
	"github.com/kiricle/api-homework/internal/services"
	"github.com/kiricle/api-homework/internal/storage/cache"
	"github.com/kiricle/api-homework/internal/storage/postgres"
	"log/slog"
	"os"
)

// @title API Server
// @host localhost:8080

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout,
		&slog.HandlerOptions{Level: slog.LevelInfo}),
	)

	postgresStorage, err := postgres.NewStorage(log)
	if err != nil {
		panic(err)
	}
	cacheStorage := cache.NewCache()

	bookService := services.NewBookService(postgresStorage, cacheStorage, log)
	bookHandler := handlers.NewBookHandler(log, bookService)

	r := router.SetupRouter(bookHandler)

	log.Info("Starting server on port 8080")
	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}
