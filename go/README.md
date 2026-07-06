<h1 align="center">{{ .name }}</h1>

<h2 align="center">

  ![Go Version](https://img.shields.io/badge/Go%20Version-latest-brightgreen.svg)
  [![Go Report Card](https://goreportcard.com/badge/{{ .gitserver }}/{{ .owner }}/{{ .name }})](https://goreportcard.com/report/{{ .gitserver }}/{{ .owner }}/{{ .name }})
  {{ if .license }}[![License](https://img.shields.io/badge/license-{{ .license }}-blue.svg)](https://{{ .gitserver }}/{{ .owner }}/{{ .name }}/blob/main/LICENSE){{ end }}

</h2>

<h2 align="center">{{ .description }}</h2>

## Getting Started

After generating this project, run the following commands to initialize dependencies:

```bash
go mod tidy
```

Then build and run:

```bash
make build
./{{ .name }}
```
{{ if .conf }}
## Configuration

A sample configuration file (`conf.yml`) is provided with documented defaults.
Configuration is loaded with the following precedence:

1. **CLI flags** (e.g. `--server.port 9090`)
2. **Environment variables** (prefixed with the uppercased project name)
3. **Configuration file** (`conf.yml` or set via `--conf` / `APPNAME_CONF` env var)
4. **Defaults** set on the flag definitions

```bash
./{{ .name }} --help           # list all flags
./{{ .name }} --conf custom.yml # use a specific config file
APPNAME_LOG_LEVEL=debug ./{{ .name }}  # env var override
```
{{ end }}
{{ if .gin_otel }}
## OpenTelemetry

This project instruments **traces**, **metrics**, and **logs** via OpenTelemetry OTLP (gRPC).
Set `--server.instrument=false` to disable.

### Local development

A local [otel-desktop-viewer](https://github.com/CtrlSpice/otel-desktop-viewer) is provided
for development — traces, metrics, and logs in a single web UI.

```bash
# Start the viewer (Docker required)
docker compose -f docker-compose.otel.yml up

# In another terminal, run the app
OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:4317 ./{{ .name }}
```
Both the service name (`OTEL_SERVICE_NAME`) and OTLP endpoint
(`OTEL_EXPORTER_OTLP_ENDPOINT`) are configured by default
(`{{ .name }}` → `http://localhost:4317`). Override them to point elsewhere.

Open [http://localhost:8000](http://localhost:8000) to explore traces, metrics, and logs.

### Production

Point the app at your OTLP collector or observability backend:

```bash
OTEL_EXPORTER_OTLP_ENDPOINT=https://your-collector:4317 ./{{ .name }}
```

All configuration is via standard [OTEL environment variables](https://opentelemetry.io/docs/concepts/sdk-configuration/).

| Signal | What's instrumented |
|--------|---------------------|
| Traces | One span per HTTP request |
| Metrics | HTTP duration/throughput + Go runtime (goroutines, memory, GC) |
| Logs | All `slog.*` calls forwarded to OTLP |
{{ end }}
{{ if .gin }}
## API Endpoints

| Path | Description |
|------|-------------|
| `GET /health` | Health check, returns `{"status":"ok"}` |
{{ end }}

## Makefile Targets

```
make help       # list all targets
make build      # compile the binary
make test       # run the test suite
make lint       # run golangci-lint
{{ if .gin_otel }}make dev        # run with OTLP enabled (go run)
{{ end }}{{ if .docker }}make docker    # build Docker image
{{ end }}{{ if .goreleaser }}make release   # create a GitHub release via goreleaser
{{ end }}make clean      # remove binary and coverage output
```

## Build Information

Version, commit hash, and build date are injected at compile time via `-ldflags`:

```bash
make build
./{{ .name }} version
# Build: 9a3b2c1
# Version: 0.1.0-dev
# Build Date: 2026-07-05T15:53:41Z
```
