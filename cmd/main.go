package main

import (
    "database/sql"
    "fmt"
    "go.uber.org/zap"
    "log"
    "music-library/internal/api"
    "music-library/internal/config"
    "music-library/internal/repository"
    "music-library/internal/service"

    _ "github.com/lib/pq"
    "github.com/golang-migrate/migrate/v4"
    "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/file"
    _ "music-library/docs" // This line is important for swagger
)

// @title           Music Library API
// @version         1.0
// @description     A REST API for managing a music library
// @host            localhost:8080
// @BasePath        /api/v1

func main() {
    // Load configuration
    cfg, err := config.LoadConfig()
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

    // Initialize logger
    logger, err := zap.NewProduction()
    if err != nil {
        log.Fatalf("Failed to create logger: %v", err)
    }
    defer logger.Sync()

    // Connect to database
    db, err := sql.Open("postgres", cfg.GetDBConnString())
    if err != nil {
        logger.Fatal("Failed to connect to database", zap.Error(err))
    }
    defer db.Close()

    // Run migrations
    driver, err := postgres.WithInstance(db, &postgres.Config{})
    if err != nil {
        logger.Fatal("Failed to create migration driver", zap.Error(err))
    }

    m, err := migrate.NewWithDatabaseInstance(
        "file://D:/zadanie/migrations",
        "postgres",
        driver,
    )
    if err != nil {
        logger.Fatal("Failed to create migration instance", zap.Error(err))
    }

    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
        logger.Fatal("Failed to run migrations", zap.Error(err))
    }

    // Initialize components
    songRepo := repository.NewSongRepository(db)
    musicAPIClient := service.NewMusicAPIClient(cfg.MusicAPIURL)
    songService := service.NewSongService(songRepo, musicAPIClient, logger)
    handler := api.NewHandler(songService, logger)
    router := api.SetupRouter(handler)

    // Start server
    addr := fmt.Sprintf(":%s", cfg.ServerPort)
    logger.Info("Starting server", zap.String("addr", addr))
    if err := router.Run(addr); err != nil {
        logger.Fatal("Failed to start server", zap.Error(err))
    }
}
