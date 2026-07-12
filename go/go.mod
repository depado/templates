module {{ .gitserver }}/{{ .owner }}/{{ .name }}

go 1.26

require (
	{{ if .gin -}}
	github.com/gin-contrib/cors v1.7.7
	github.com/gin-gonic/gin v1.12.0
	{{- if .gin_otel }}
	go.opentelemetry.io/contrib/bridges/otelslog v0.19.0
	go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin v0.69.0
	go.opentelemetry.io/contrib/instrumentation/runtime v0.69.0
	go.opentelemetry.io/otel v1.44.0
	go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc v0.20.0
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc v1.44.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.44.0
	go.opentelemetry.io/otel/sdk v1.44.0
	go.opentelemetry.io/otel/sdk/log v0.20.0
	go.opentelemetry.io/otel/sdk/metric v1.44.0
	{{ end -}}
	{{ end -}}
	github.com/lmittmann/tint v1.2.0
	github.com/mattn/go-isatty v0.0.22
	github.com/spf13/cobra v1.10.2
	github.com/spf13/viper v1.21.0
)
