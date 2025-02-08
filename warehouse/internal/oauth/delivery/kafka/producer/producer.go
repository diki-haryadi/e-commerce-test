package oauthKafkaProducer

import (
	"context"

	"github.com/segmentio/kafka-go"

	oauthDomain "github.com/diki-haryadi/go-micro-template/internal/oauth/domain"
	kafkaProducer "github.com/diki-haryadi/ztools/kafka/producer"
)

type producer struct {
	createWriter *kafkaProducer.Writer
}

func NewProducer(w *kafkaProducer.Writer) oauthDomain.KafkaProducer {
	return &producer{createWriter: w}
}

func (p *producer) PublishCreateEvent(ctx context.Context, messages ...kafka.Message) error {
	return p.createWriter.Client.WriteMessages(ctx, messages...)
}
