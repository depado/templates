module {{ .gitserver }}/{{ .owner }}/{{ .name }}

go 1.26

require (
	{{ if .gin -}}
	github.com/gin-contrib/cors v1.7.8
	github.com/gin-gonic/gin v1.12.0
	{{- if .gin_otel }}
	go.opentelemetry.io/contrib/bridges/otelslog v0.20.1
	go.opentelemetry.io/contrib/exporters/autoexport v0.71.0
	go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin v0.71.0
	go.opentelemetry.io/contrib/instrumentation/runtime v0.71.0
	go.opentelemetry.io/otel v1.46.0
	go.opentelemetry.io/otel/sdk v1.46.0
	go.opentelemetry.io/otel/sdk/log v0.22.0
	go.opentelemetry.io/otel/sdk/metric v1.46.0
	{{ end -}}
	{{ end -}}
	github.com/lmittmann/tint v1.2.0
	github.com/mattn/go-isatty v0.0.24
	github.com/spf13/cobra v1.10.2
	github.com/spf13/viper v1.21.0
)
