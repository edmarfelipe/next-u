package tracer

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"go.opentelemetry.io/otel/trace"
)

const (
	service     = "identity"
	environment = "production"
	id          = 1
)

func StartSpan(ctx context.Context, componetName string, spanName string) (context.Context, trace.Span) {
	tr := otel.Tracer(componetName)
	return tr.Start(ctx, spanName)
}

func NewTracer() error {
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint())
	if err != nil {
		return err
	}

	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(service),
			attribute.String("environment", environment),
		)),
	)

	otel.SetTracerProvider(tp)
	return nil
}
