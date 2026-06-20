package tracing

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
)

type ExporterType string

const (
	ExporterOTLP    ExporterType = "otlp"
	ExporterConsole ExporterType = "console"
)

type Tracing struct {
	serviceName    string
	serviceVersion string
	environment    string

	endpoint string
	insecure bool

	exporterType ExporterType
}

type OptFunc func(*Tracing) error

func WithServiceName(v string) OptFunc {
	return func(t *Tracing) error {
		t.serviceName = v
		return nil
	}
}

func WithServiceVersion(v string) OptFunc {
	return func(t *Tracing) error {
		t.serviceVersion = v
		return nil
	}
}

func WithEnvironment(v string) OptFunc {
	return func(t *Tracing) error {
		t.environment = v
		return nil
	}
}

func WithEndpoint(v string) OptFunc {
	return func(t *Tracing) error {
		t.endpoint = v
		return nil
	}
}

func WithInsecure() OptFunc {
	return func(t *Tracing) error {
		t.insecure = true
		return nil
	}
}

func WithConsoleExporter() OptFunc {
	return func(t *Tracing) error {
		t.exporterType = ExporterConsole
		return nil
	}
}

func WithOTLPExporter() OptFunc {
	return func(t *Tracing) error {
		t.exporterType = ExporterOTLP
		return nil
	}
}

func (t *Tracing) buildResource(
	ctx context.Context,
) (*resource.Resource, error) {

	return resource.New(
		ctx,

		resource.WithFromEnv(),

		resource.WithTelemetrySDK(),

		resource.WithAttributes(
			semconv.ServiceName(t.serviceName),

			semconv.ServiceVersion(t.serviceVersion),

			semconv.DeploymentEnvironmentName(
				t.environment,
			),
		),
	)
}

func (t *Tracing) newConsoleExporter() (sdktrace.SpanExporter, error) {

	return stdouttrace.New(
		stdouttrace.WithPrettyPrint(),
	)
}

func (t *Tracing) newOTLPExporter(
	ctx context.Context,
) (sdktrace.SpanExporter, error) {

	opts := []otlptracehttp.Option{
		otlptracehttp.WithEndpoint(
			t.endpoint,
		),
	}

	if t.insecure {
		opts = append(
			opts,
			otlptracehttp.WithInsecure(),
		)
	}

	return otlptracehttp.New(
		ctx,
		opts...,
	)
}

func (t *Tracing) buildExporter(
	ctx context.Context,
) (sdktrace.SpanExporter, error) {

	switch t.exporterType {

	case ExporterConsole:
		return t.newConsoleExporter()

	case ExporterOTLP:
		return t.newOTLPExporter(ctx)

	default:
		return t.newOTLPExporter(ctx)
	}
}

func (t *Tracing) buildProvider(
	ctx context.Context,
	exp sdktrace.SpanExporter,
) (*sdktrace.TracerProvider, error) {

	res, err := t.buildResource(ctx)
	if err != nil {
		return nil, err
	}

	return sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(res),
	), nil
}

func New(
	ctx context.Context,
	opts ...OptFunc,
) (*sdktrace.TracerProvider, error) {

	cfg := &Tracing{
		serviceVersion: "1.0.0",
		environment:    "development",

		endpoint: "localhost:4318",

		insecure: true,

		exporterType: ExporterOTLP,
	}

	for _, opt := range opts {
		if err := opt(cfg); err != nil {
			return nil, err
		}
	}

	if cfg.serviceName == "" {
		return nil, fmt.Errorf(
			"service name is required",
		)
	}

	exp, err := cfg.buildExporter(ctx)
	if err != nil {
		return nil, fmt.Errorf(
			"create exporter: %w",
			err,
		)
	}

	tp, err := cfg.buildProvider(
		ctx,
		exp,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"create tracer provider: %w",
			err,
		)
	}

	otel.SetTracerProvider(tp)

	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	return tp, nil
}

func Shutdown(
	ctx context.Context,
	tp *sdktrace.TracerProvider,
) error {

	if tp == nil {
		return nil
	}

	return tp.Shutdown(ctx)
}
