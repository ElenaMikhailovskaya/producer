/*
Package http_transport содержит в себе описание хендлеров сервиса и
валидацию запросов по типу значений и их формату.
*/
package http_transport

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gitlab.rusklimat.ru/ecom/go-lib/errs"
)

type Producer interface {
	Produce(ctx context.Context) *errs.Error
}

type Server struct {
	a         *fiber.App
	validator *validator.Validate
	cfg       cfg
	logs      errs.LogChan

	service Producer
}

type cfg struct {
	Host   string `env:"HOST" envDefault:":1235"`
	TLSKey string `env:"TLS_KEY" envDefault:""`
	TLSPem string `env:"TLS_PEM" envDefault:""`
}
