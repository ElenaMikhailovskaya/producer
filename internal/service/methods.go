package service

import (
	"context"
	"gitlab.rusklimat.ru/ecom/go-lib/errs"
	"service/internal/models"
	"sync"
)

func WithProducer(c Producer) func(*Server) {
	return func(s *Server) {
		s.prod = c
	}
}

// New ждет на вход массив опций, поскольку может быть несколько интеграций и для тестов нужны будут не все.
// В случае с инициализацией транспорта по-другому, поскольку там всегда будет только вызов слоя логики
func New(opts ...func(server *Server)) (*Server, *errs.Error) {
	c := new(Server)
	c.userMu = new(sync.Map)

	for _, o := range opts {
		o(c)
	}
	return c, nil
}

func (s *Server) Produce(ctx context.Context) *errs.Error {
	ctx, span := models.Tracer.Start(ctx, "Produce()")
	span.AddEvent("start produce")
	defer span.End()

	topic := "testTopic"

	err := s.prod.Produce(ctx, "testKey", []byte("testValue"), topic)
	if err != nil {
		return errs.NewError(errs.ErrorLevel, err.Error())
	}

	return nil

}
