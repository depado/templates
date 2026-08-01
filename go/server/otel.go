---
if: gin
---
package server

// This file bootstraps OpenTelemetry for all three signals (traces, metrics,
// logs). Exporters are configured via standard OTEL_* environment variables
// through the autoexport package - no manual protocol switching needed.
//
//	OTEL_EXPORTER_OTLP_ENDPOINT     OTLP collector address
//	OTEL_EXPORTER_OTLP_PROTOCOL     transport: grpc or http/protobuf
//	OTEL_RESOURCE_ATTRIBUTES        additional resource attributes (key=value,...)
//
// When the collector is unreachable, each signal is silently disabled.

{{ if .gin_otel }}
import (
	"context"
	"errors"
	"log/slog"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/contrib/exporters/autoexport"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	sdkresource "go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.41.0"

	"{{ .gitserver }}/{{ .owner }}/{{ .name }}/cmd"
)

// Telemetry holds the OpenTelemetry providers and an optional OTLP log
// handler. All signals are exported via OTLP, configured through standard
// OTEL_* environment variables.
type Telemetry struct {
	shutdown   func(context.Context) error
	logHandler slog.Handler
}

// NewTelemetry initialises OpenTelemetry for traces, metrics, and logs.
// When Instrument is disabled, a no-op instance is returned.
func NewTelemetry(c *cmd.Conf, l *slog.Logger) (*Telemetry, error) {
	t := &Telemetry{
		shutdown: func(_ context.Context) error { return nil },
	}

	if !c.Server.Instrument {
		return t, nil
	}

	ctx := context.Background()

	res, err := sdkresource.New(ctx,
		sdkresource.WithFromEnv(),
		sdkresource.WithTelemetrySDK(),
		sdkresource.WithAttributes(
			semconv.ServiceName("{{ .name }}"),
			semconv.ServiceVersion(cmd.Version),
		),
	)
	if err != nil {
		return nil, err
	}

	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	providers := make([]func(context.Context) error, 0, 3)

	// Traces - autoexport reads OTEL_TRACES_EXPORTER + OTEL_EXPORTER_OTLP_PROTOCOL
	{
		exp, err := autoexport.NewSpanExporter(ctx)
		if err != nil {
			l.Warn("OTLP trace exporter unavailable, traces disabled", "error", err)
		} else {
			tp := sdktrace.NewTracerProvider(
				sdktrace.WithBatcher(exp),
				sdktrace.WithResource(res),
			)
			otel.SetTracerProvider(tp)
			providers = append(providers, tp.Shutdown)
		}
	}

	// Metrics - autoexport reads OTEL_METRICS_EXPORTER + OTEL_EXPORTER_OTLP_PROTOCOL
	{
		reader, err := autoexport.NewMetricReader(ctx)
		if err != nil {
			l.Warn("OTLP metric exporter unavailable, metrics disabled", "error", err)
		} else {
			mp := sdkmetric.NewMeterProvider(
				sdkmetric.WithReader(reader),
				sdkmetric.WithResource(res),
			)
			otel.SetMeterProvider(mp)
			providers = append(providers, mp.Shutdown)

			if err := runtime.Start(runtime.WithMeterProvider(mp)); err != nil {
				l.Warn("runtime metrics failed to start", "error", err)
			}
		}
	}

	// Logs - autoexport reads OTEL_LOGS_EXPORTER + OTEL_EXPORTER_OTLP_PROTOCOL
	{
		exp, err := autoexport.NewLogExporter(ctx)
		if err != nil {
			l.Warn("OTLP log exporter unavailable, logs disabled", "error", err)
		} else {
			lp := sdklog.NewLoggerProvider(
				sdklog.WithProcessor(sdklog.NewBatchProcessor(exp)),
				sdklog.WithResource(res),
			)
			providers = append(providers, lp.Shutdown)
			t.logHandler = otelslog.NewHandler("{{ .name }}", otelslog.WithLoggerProvider(lp))
		}
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

// GinMiddleware returns nil - no instrumentation middleware.
func (t *Telemetry) GinMiddleware() gin.HandlerFunc { return nil }

// Shutdown is a no-op.
func (t *Telemetry) Shutdown(_ context.Context) error { return nil }

// Logger returns the base logger unchanged.
func (t *Telemetry) Logger(base *slog.Logger) *slog.Logger { return base }

{{ end }}
