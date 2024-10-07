package producer

import (
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"go.opentelemetry.io/otel"
	"time"
)

type Client struct {
	c   *kafka.Producer
	cfg cfg
}

// пример конфига подключения к кафке
type cfg struct {
	RetriesForProduce time.Duration `env:"KAFKA_RETRIES_FOR_PRODUCE" endDefault:"5s"`
}

var Tracer = otel.Tracer("default")
