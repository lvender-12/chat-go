package rabbit

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	amqp "github.com/rabbitmq/amqp091-go"
)

const ChatExchange = "chat.fanout"

type Event struct {
	ConversationID uint64          `json:"conversation_id"`
	Message        json.RawMessage `json:"message"`
}

type Bus struct {
	ch       *amqp.Channel
	exchange string
	queue    string
	logger   *slog.Logger
}

func NewBus(conn *amqp.Connection, logger *slog.Logger) (*Bus, error) {
	ch, err := OpenChannel(conn)
	if err != nil {
		return nil, err
	}

	bus := &Bus{
		ch:       ch,
		exchange: ChatExchange,
		logger:   logger,
	}

	if err := bus.setup(); err != nil {
		_ = ch.Close()
		return nil, err
	}

	return bus, nil
}

func (b *Bus) setup() error {
	if err := b.ch.ExchangeDeclare(
		b.exchange,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return fmt.Errorf("declare exchange: %w", err)
	}

	q, err := b.ch.QueueDeclare(
		"",
		false,
		true,
		true,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("declare queue: %w", err)
	}
	b.queue = q.Name

	if err := b.ch.QueueBind(
		b.queue,
		"",
		b.exchange,
		false,
		nil,
	); err != nil {
		return fmt.Errorf("bind queue: %w", err)
	}

	b.logger.Info(
		"rabbitmq chat bus ready",
		"exchange", b.exchange,
		"queue", b.queue,
	)

	return nil
}

func (b *Bus) Publish(ctx context.Context, conversationID uint64, message []byte) error {
	body, err := json.Marshal(Event{
		ConversationID: conversationID,
		Message:        message,
	})
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}

	if err := b.ch.PublishWithContext(
		ctx,
		b.exchange,
		"",
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         body,
		},
	); err != nil {
		return fmt.Errorf("publish: %w", err)
	}

	b.logger.Debug(
		"rabbitmq published",
		"conversation_id", conversationID,
		"bytes", len(message),
	)

	return nil
}

func (b *Bus) Consume(handler func(Event)) error {
	deliveries, err := b.ch.Consume(
		b.queue,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("consume: %w", err)
	}

	go func() {
		for d := range deliveries {
			var event Event
			if err := json.Unmarshal(d.Body, &event); err != nil {
				b.logger.Error("rabbitmq invalid event", "error", err)
				continue
			}

			handler(event)
		}

		b.logger.Info("rabbitmq consumer stopped")
	}()

	b.logger.Info("rabbitmq consumer started", "queue", b.queue)

	return nil
}

func (b *Bus) Close() error {
	if b.ch == nil {
		return nil
	}
	return b.ch.Close()
}
