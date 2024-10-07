package app

import (
	"context"
	"github.com/caarlos0/env/v6"
	"gitlab.rusklimat.ru/ecom/go-lib/errs"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	stdout "go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

type cfg struct {
	ServerName string `env:"SERVER_NAME" envDefault:"TestServer"`
	SentryDSN  string `env:"SENTRY_DSN" envDefault:""`
	ServiceEnv string `env:"SERVICE_ENV" envDefault:"dev"`
	LogLevel   string `env:"LOG_LEVEL" envDefault:"trace"`
}

type otlCfg struct {
	Endpoint    string `env:"OTL_ENDPOINT" envDefault:"localhost:4318"`
	ServiceName string `env:"OTL_SERVICE_NAME" envDefault:"serviceProducer"`
	ServiceEnv  string `env:"SERVICE_ENV" envDefault:"dev"`
}

func initTracer() *sdktrace.TracerProvider {
	var (
		exporter sdktrace.SpanExporter
		err      error
	)

	var cfg otlCfg
	e := env.Parse(&cfg)
	if e != nil {
		errs.Fatal(e)
	}

	if cfg.Endpoint == "" {
		exporter, err = stdout.New(stdout.WithPrettyPrint())
	} else {
		exporter, err = otlptracehttp.New(context.Background(),
			otlptracehttp.WithInsecure(),
			otlptracehttp.WithEndpoint(cfg.Endpoint),
		)
	}
	if err != nil {
		panic(err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String(cfg.ServiceName),
				attribute.String("environment", cfg.ServiceEnv),
			)),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp
}
