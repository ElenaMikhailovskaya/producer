package http_transport

import (
	"github.com/caarlos0/env/v6"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/contrib/otelfiber/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"gitlab.rusklimat.ru/ecom/go-lib/errs"
	"gitlab.rusklimat.ru/ecom/go-lib/tools/den"
)

func New(service Producer, logs errs.LogChan) (*Server, *errs.Error) {
	var cfg cfg
	e := env.Parse(&cfg)
	if e != nil {
		return nil, errs.NewError(errs.FatalLevel, e.Error())
	}

	s := new(Server)
	s.cfg = cfg

	s.service = service
	if service == nil {
		return nil, errs.NewError(errs.FatalLevel, "param `service` is required")
	}

	s.logs = logs
	if logs == nil {
		return nil, errs.NewError(errs.FatalLevel, "param `logs` is required")
	}

	s.validator = validator.New(validator.WithRequiredStructEnabled())
	s.a = newFiberApp(errs.NewFiberLogger(logs))

	return s.setupHandlers(), nil
}

func (s *Server) setupHandlers() *Server {
	s.a.Get("/api/ping", s.ping) // ping не зависит от версии

	v1 := s.a.Group("/api/v1")

	v1.Get("/", s.ProducerHandler)

	return s
}

func newFiberApp(lg errs.Logger4Fiber) *fiber.App {
	a := fiber.New(fiber.Config{
		JSONEncoder: func(v any) ([]byte, error) {
			b, e := den.EncodeJson(v)
			if e != nil {
				return nil, e.GetErr()
			}
			return b.Bytes(), nil
		},
	})

	// обертка над запросами для OTEL
	// теперь в каждом контексте есть span и его нужно(!) использовать
	a.Use(otelfiber.Middleware())

	a.Use(logger.New(logger.Config{
		Format: "{\"status\": ${status}, \"duration\": \"${latency}\", \"method\": \"${method}\", \"path\": \"${path}\", \"resp\": \"${resBody}\"}\n",
		Output: &lg,
	}))
	a.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))
	a.Use(recover.New())

	return a
}
