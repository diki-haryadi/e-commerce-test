package productKafkaProducer

import (
	"context"

	"github.com/segmentio/kafka-go"

	productDomain "github.com/diki-haryadi/go-micro-template/internal/product/domain"
	kafkaProducer "github.com/diki-haryadi/ztools/kafka/producer"
)

type producer struct {
	createWriter *kafkaProducer.Writer
}

func NewProducer(w *kafkaProducer.Writer) productDomain.KafkaProducer {
	return &producer{createWriter: w}
}

func (p *producer) PublishCreateEvent(ctx context.Context, messages ...kafka.Message) error {
	return p.createWriter.Client.WriteMessages(ctx, messages...)
}
