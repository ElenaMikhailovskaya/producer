package models

import "go.opentelemetry.io/otel"

var Tracer = otel.Tracer("default")
