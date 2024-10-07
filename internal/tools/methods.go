package tools

import (
	"context"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"gitlab.rusklimat.ru/ecom/go-lib/errs"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

// SetTraceToMessageHeader получает span из контекста и записывает TraceID и SpanID в заголовки кафки, добавляет заголовки сообщения кафки в контекст
// (эту функцию вызываем в producer Produce)
func SetTraceToMessageHeader(ctx context.Context, msg *kafka.Message) {
	propagators := propagation.TraceContext{}
	spanContext := trace.SpanFromContext(ctx).SpanContext()
	setValue(msg, TraceID, []byte(spanContext.TraceID().String()))
	setValue(msg, SpanID, []byte(spanContext.SpanID().String()))

	carrier := newMessageCarrier(msg)
	propagators.Inject(ctx, carrier)
}

// UpdateContext получает TraceID и SpanID из заголовка сообщения кафки из контекста
func UpdateContext(ctx context.Context, msg *kafka.Message) (context.Context, *errs.Error) {

	propagators := propagation.TraceContext{}
	carrier := newMessageCarrier(msg)
	ctxNew := propagators.Extract(ctx, carrier)

	return ctxNew, nil
}
