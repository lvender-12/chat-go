package main

import (
	"chat-go/internal/app"
	"chat-go/internal/config"
	"chat-go/internal/database"
	"chat-go/internal/rabbit"
	"chat-go/internal/router"
	"log/slog"
	"os"

	"github.com/gofiber/contrib/v3/swaggerui"
	"github.com/gofiber/fiber/v3"
)

func main() {
	var logger *slog.Logger

	Routeapp := fiber.New(fiber.Config{
		ErrorHandler: app.ErrorHandler,
	})

	conf, err := config.LoadConfig("config/config.json")
	if err != nil {
		slog.Error("failed to load config", "error", err)
		panic(err)
	}

	if conf.App.LogLevel == "debug" {
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	} else {
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
	}
	db, err := database.LoadDB(conf.Database.Driver, conf.Database.User, conf.Database.Password, conf.Database.Name, logger)
	if err != nil {
		slog.Error("failed to load db", "error", err)
		panic(err)
	}

	if conf.Migration.Enabled {
		err = database.MigrateDB(db, logger, conf.Migration.Path)
		if err != nil {
			slog.Error("failed to migrate db", "error", err)
			panic(err)
		}
	}

	rabbitConn, err := rabbit.Connect(conf.Rabbit, logger)
	if err != nil {
		slog.Error("failed to connect rabbitmq", "error", err)
		panic(err)
	}
	defer rabbitConn.Close()

	state := &app.State{
		DB:     db,
		Config: conf,
		Rabbit: rabbitConn,
	}

	Routeapp.Use(swaggerui.New(swaggerui.Config{
		BasePath: "/",
		FilePath: "./docs/swagger.json",
		Path:     "docs",
		Title:    "chat go API",
		CacheAge: 0,
	}))
	router.SetupRouter(Routeapp, state, logger)

	slog.Error("failed to start server", "error", Routeapp.Listen(":3000"))
}
