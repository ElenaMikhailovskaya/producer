package producer

import (
	"context"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"gitlab.rusklimat.ru/ecom/go-lib/errs"
	"service/internal/tools"
	"time"
)

// Produce отправляем сообщение в топик кафки
func (c *Client) Produce(ctx context.Context, key string, val []byte, topic string) *errs.Error {
	ctx, span := Tracer.Start(ctx, "Produce")
	defer span.End()
	span.AddEvent("Send message to topic")

	fmt.Println(span.SpanContext().TraceID())

	msg := kafka.Message{
		Key:   []byte(key),
		Value: val,
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
	}

	// пример добавления данных телеметрии (TraceID и SpanID в заголовок сообщения кафки)
	tools.SetTraceToMessageHeader(ctx, &msg)

	e := c.c.Produce(&msg, nil)
	if e != nil {
		return errs.NewError(errs.ErrorLevel, e.Error())
	}

	return nil
}

// Retries returns the number of retries for produce.
func (c *Client) Retries() time.Duration { return c.cfg.RetriesForProduce }

// GetEvent returns the next event from the client's event channel.
// For listen errors
func (c *Client) GetEvent() kafka.Event { return <-c.c.Events() }

func (c *Client) Flush() {
	c.c.Flush(15 * 1000)
}

func (c *Client) Close() {
	c.c.Close()
}
