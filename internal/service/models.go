/*
Package service хранит в себе бизнес логику приложения и взаимодействует с остальными интеграциями
*/
package service

import (
	"context"
	"gitlab.rusklimat.ru/ecom/go-lib/errs"
	"sync"
)

type Server struct {
	userMu *sync.Map

	prod Producer
}

type Producer interface {
	Produce(ctx context.Context, key string, val []byte, topic string) *errs.Error
}
