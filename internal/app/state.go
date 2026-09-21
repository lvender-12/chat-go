package app

import (
	"chat-go/internal/config"
	"database/sql"

	amqp "github.com/rabbitmq/amqp091-go"
)

type State struct {
	DB     *sql.DB
	Config *config.Config
	Rabbit *amqp.Connection
}
