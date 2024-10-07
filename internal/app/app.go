/*
Package app содержит функцию Start для инициализации всех интерфейсов и запуска приложения на указанном
порту
*/
package app

import (
	"context"
	"errors"
	"github.com/caarlos0/env/v6"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/getsentry/sentry-go"
	"github.com/sirupsen/logrus"
	"gitlab.rusklimat.ru/ecom/go-lib/errs"
	"net/http"
	"service/internal/producer"
	"service/internal/service"
	"service/internal/transport/http_transport"
	"sync"
	"time"
)

func Start() {
	var cfg cfg
	e := env.Parse(&cfg)
	if e != nil {
		panic(e)
	}

	errs.SetupLogger(cfg.ServerName, cfg.SentryDSN, cfg.ServiceEnv, cfg.LogLevel)
	defer sentry.Flush(2 * time.Second)

	tp := initTracer()

	defer func() {
		if err := tp.ForceFlush(context.Background()); err != nil {
			errs.Err(err)
		}
		if err := tp.Shutdown(context.Background()); err != nil {
			errs.Err(err)
		}
	}()

	logChan := make(errs.LogChan, 1000)

	producerClient, err := producer.New()
	if err != nil {
		errs.Fatal(err)
	}

	defer producerClient.Close()
	defer producerClient.Flush()

	logic, err := service.New(
		service.WithProducer(producerClient))
	if err != nil {
		errs.Fatal(err)
	}

	transport, err := http_transport.New(logic, logChan)
	if err != nil {
		errs.Fatal(err)
	}

	wg := new(sync.WaitGroup)
	wg.Add(1)
	go errs.LogWatcher(logChan, wg)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-ctx.Done():
				return

			default:
				event := producerClient.GetEvent()
				switch e := event.(type) {
				case kafka.Error:
					errs.NewError(errs.ErrorLevel, e.Error()).Wrap("producer.GetEvent")
				}
			}

		}
	}()

	connectionsClosed := make(chan struct{})
	go transport.GracefulShutdown(connectionsClosed)

	if err := transport.Listen(); errors.Is(err, http.ErrServerClosed) {
		logChan <- errs.NewError(logrus.FatalLevel, err.Error()) // wrap not required
	} else {
		<-connectionsClosed // wait for success close connections
	}

	// close db connections

	close(logChan)
	cancel()
	wg.Wait()
}
