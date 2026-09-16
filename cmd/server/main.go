package main

import (
	"chat-go/internal/app"
	"chat-go/internal/config"
	"chat-go/internal/database"
	"chat-go/internal/router"
	"log/slog"

	"github.com/gofiber/fiber/v3"
)

func main() {
	logger := slog.Default()
	Routeapp := fiber.New()

	conf, err := config.LoadConfig("config/config.json")
	if err != nil {
		slog.Error("failed to load config", "error", err)
		panic(err)
	}

	db, err := database.LoadDB(conf.Database.Driver, conf.Database.User, conf.Database.Password, conf.Database.Name, logger)
	if err != nil {
		slog.Error("failed to load db", "error", err)
		panic(err)
	}

	err = database.MigrateDB(db, logger)
	if err != nil {
		slog.Error("failed to migrate db", "error", err)
		panic(err)
	}

	state := &app.State{
		DB:     db,
		Config: conf,
	}

	router.SetupRouter(Routeapp, state, logger)

	slog.Error("failed to start server", "error", Routeapp.Listen(":3000"))
}
