package mq

import (
	"context"

	"github.com/google/wire"
	"github.com/webmin7761/go-school/homework/final/internal/mq/kafka"
)

// ProviderSet is mq providers.
var ProviderSet = wire.NewSet(kafka.NewKafkaClient)

type MessageQueue interface {
	Produce(msg string) error
	Consume(consume func(ctx context.Context, msg string, err error) error) error
}
