package http_transport

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"gitlab.rusklimat.ru/ecom/go-lib/errs"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"net/http"
	"service/internal/models"
)

func (s *Server) ping(c *fiber.Ctx) error {
	_, span := models.Tracer.Start(c.UserContext(), "test1")
	span.AddEvent("New event")
	defer span.End()
	return c.SendStatus(http.StatusOK)
}

func (s *Server) ProducerHandler(c *fiber.Ctx) error {
	fmt.Println("test")
	ctx, span := models.Tracer.Start(c.UserContext(), "ProducerHandler",
		trace.WithSpanKind(trace.SpanKindProducer),
		trace.WithAttributes(
			attribute.String("peer.service", "kafka"),
			attribute.String("connection_type", "messaging_system"),
			attribute.String("virtual_node", "client")))

	span.AddEvent("ProducerHandler start")
	defer span.End()

	err := s.service.Produce(ctx)
	if err != nil {
		errs.NewError(errs.ErrorLevel, err.Error())
	}

	return c.SendStatus(http.StatusOK)
}
