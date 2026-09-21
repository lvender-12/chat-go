package rabbit

import (
	"fmt"
	"log/slog"
	"net/url"

	"chat-go/internal/config"

	amqp "github.com/rabbitmq/amqp091-go"
)

func Connect(cfg config.RabbitConfig, logger *slog.Logger) (*amqp.Connection, error) {
	host := cfg.Host
	if host == "" {
		host = "127.0.0.1"
	}

	port := cfg.Port
	if port == "" {
		port = "5672"
	}

	u := url.URL{
		Scheme: "amqp",
		User:   url.UserPassword(cfg.Name, cfg.Password),
		Host:   fmt.Sprintf("%s:%s", host, port),
	}

	conn, err := amqp.Dial(u.String())
	if err != nil {
		return nil, fmt.Errorf("rabbitmq dial: %w", err)
	}

	logger.Info("rabbitmq connected", "host", host, "port", port)

	return conn, nil
}

func OpenChannel(conn *amqp.Connection) (*amqp.Channel, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("rabbitmq channel: %w", err)
	}

	return ch, nil
}
