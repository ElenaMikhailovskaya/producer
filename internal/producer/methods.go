package producer

import (
	"github.com/caarlos0/env/v6"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"gitlab.rusklimat.ru/ecom/go-lib/errs"
)

func New() (*Client, *errs.Error) {
	var cfg cfg
	e := env.Parse(&cfg)
	if e != nil {
		return nil, errs.NewError(errs.FatalLevel, e.Error()).
			WrapWithSentry(errs.SentryCategoryFunc, nil)
	}

	conn, e := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9193",
		"sasl.mechanisms":   "SCRAM-SHA-256",
		"security.protocol": "SASL_PLAINTEXT",
		"sasl.username":     "root",
		"sasl.password":     "rootpassword",
	})
	if e != nil {
		return nil, errs.NewError(errs.ErrorLevel, e.Error())
	}

	return &Client{c: conn, cfg: cfg}, nil
}
