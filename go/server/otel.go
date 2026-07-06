---
if: gin
---
package server

// This file uses OTLP over gRPC (port 4317). To switch to HTTP/protobuf
// (port 4318), replace the three exporter imports and constructors:
//
//	otlptracegrpc   → otlptracehttp
//	otlpmetricgrpc  → otlpmetrichttp
//	otlploggrpc     → otlploghttp
//
// The API surface (New, WithInsecure, etc.) is identical. You'll also
// need the corresponding go.mod entries:
//
//	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp
//	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp
//	go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp

{{ if .gin_otel }}
import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"

	"{{ .gitserver }}/{{ .owner }}/{{ .name }}/cmd"
)

// Telemetry holds the OpenTelemetry providers, shutdown hook, and an
// optional OTLP log handler. All signals are exported via OTLP (gRPC),
// configured through standard OTEL_* environment variables.
//
//	OTEL_EXPORTER_OTLP_ENDPOINT     OTLP collector address (default localhost:4317)
//	OTEL_SERVICE_NAME               service name (defaults to the binary name)
//	OTEL_RESOURCE_ATTRIBUTES        additional resource attributes (key=value,...)
//
// When the collector is unreachable, each signal is silently disabled so
// the application can still start.
type Telemetry struct {
	shutdown   func(context.Context) error
	logHandler slog.Handler // nil when log export is disabled
}

// NewTelemetry initialises OpenTelemetry with OTLP (gRPC) for traces,
// metrics, and logs. When Instrument is disabled, a no-op instance is
// returned so that callers can treat it uniformly.
func NewTelemetry(c *cmd.Conf, l *slog.Logger) (*Telemetry, error) {
	t := &Telemetry{
		shutdown: func(_ context.Context) error { return nil },
	}

	if !c.Server.Instrument {
		return t, nil
	}

	ctx := context.Background()

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName("{{ .name }}"),
			semconv.ServiceVersion(cmd.Version),
		),
		resource.WithFromEnv(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// W3C TraceContext + Baggage propagation
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	providers := make([]func(context.Context) error, 0, 3)

	// Traces — OTLP gRPC
	traceExp, err := otlptracegrpc.New(ctx)
	if err != nil {
		l.Warn("OTLP trace exporter unavailable, traces disabled", "error", err)
	} else {
		tp := sdktrace.NewTracerProvider(
			sdktrace.WithBatcher(traceExp),
			sdktrace.WithResource(res),
		)
		otel.SetTracerProvider(tp)
		providers = append(providers, tp.Shutdown)
	}

	// Metrics — OTLP gRPC
	metricExp, err := otlpmetricgrpc.New(ctx)
	if err != nil {
		l.Warn("OTLP metric exporter unavailable, metrics disabled", "error", err)
	} else {
		mp := sdkmetric.NewMeterProvider(
			sdkmetric.WithReader(sdkmetric.NewPeriodicReader(metricExp)),
			sdkmetric.WithResource(res),
		)
		otel.SetMeterProvider(mp)
		providers = append(providers, mp.Shutdown)

		if err := runtime.Start(runtime.WithMeterProvider(mp)); err != nil {
			l.Warn("runtime metrics failed to start", "error", err)
		}
	}

	// Logs — OTLP gRPC
	logExp, err := otlploggrpc.New(ctx)
	if err != nil {
		l.Warn("OTLP log exporter unavailable, logs disabled", "error", err)
	} else {
		lp := sdklog.NewLoggerProvider(
			sdklog.WithProcessor(sdklog.NewBatchProcessor(logExp)),
			sdklog.WithResource(res),
		)
		providers = append(providers, lp.Shutdown)
		t.logHandler = otelslog.NewHandler("{{ .name }}", otelslog.WithLoggerProvider(lp))
	}

	t.shutdown = func(ctx context.Context) error {
		var errs []error
		for _, s := range providers {
			errs = append(errs, s(ctx))
		}
		return errors.Join(errs...)
	}

	return t, nil
}

// GinMiddleware returns the OpenTelemetry Gin middleware that creates a span
// and records HTTP metrics for every request.
func (t *Telemetry) GinMiddleware() gin.HandlerFunc {
	return otelgin.Middleware("{{ .name }}")
}

// Shutdown flushes and stops all OTEL providers.
func (t *Telemetry) Shutdown(ctx context.Context) error {
	return t.shutdown(ctx)
}

// Logger returns a logger that fans out to both the base handler (stderr)
// and the OTLP log handler (when instrumented). When OTLP logging is not
// configured, the base logger is returned unchanged.
func (t *Telemetry) Logger(base *slog.Logger) *slog.Logger {
	if t.logHandler == nil {
		return base
	}
	return slog.New(multiHandler{[]slog.Handler{base.Handler(), t.logHandler}})
}

{{ else }}

import (
	"context"
	"log/slog"

	"github.com/gin-gonic/gin"

	"{{ .gitserver }}/{{ .owner }}/{{ .name }}/cmd"
)

// Telemetry is a no-op when OpenTelemetry is not configured.
type Telemetry struct{}

// NewTelemetry returns a no-op Telemetry instance.
func NewTelemetry(_ *cmd.Conf, _ *slog.Logger) (*Telemetry, error) {
	return &Telemetry{}, nil
}

// GinMiddleware returns nil — no instrumentation middleware.
func (t *Telemetry) GinMiddleware() gin.HandlerFunc { return nil }

// Shutdown is a no-op.
func (t *Telemetry) Shutdown(_ context.Context) error { return nil }

// Logger returns the base logger unchanged.
func (t *Telemetry) Logger(base *slog.Logger) *slog.Logger { return base }

{{ end }}
