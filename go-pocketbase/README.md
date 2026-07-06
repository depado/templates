<h1 align="center">{{ .name }}</h1>

<h2 align="center">

  ![Go Version](https://img.shields.io/badge/Go%20Version-latest-brightgreen.svg)
  [![Go Report Card](https://goreportcard.com/badge/{{ .gitserver }}/{{ .owner }}/{{ .name }})](https://goreportcard.com/report/{{ .gitserver }}/{{ .owner }}/{{ .name }})
  {{ if .license }}[![License](https://img.shields.io/badge/license-{{ .license }}-blue.svg)](https://{{ .gitserver }}/{{ .owner }}/{{ .name }}/blob/main/LICENSE){{ end }}

</h2>

<h2 align="center">{{ .description }}</h2>

## Getting Started

After generating this project, initialize dependencies:

```bash
go mod tidy
```

### Development

Run in development mode (with automigrate):

```bash
make dev
```

Or build and serve:

```bash
make serve
```

The PocketBase admin UI is available at `http://localhost:8090/_/`.

## Configuration

```bash
cp conf.example.yml conf.yml
```

Configuration can be set via (in precedence order):
- CLI flags (`--conf path/to/config.yml`)
- Environment variables (uppercased, prefixed with project name, e.g. `MYAPP_DEV=true`)
- Configuration file (`conf.yml` or via `MYAPP_CONF` env var)

## Makefile Targets

```
make help       # list all targets
make build      # compile the binary
make dev        # run in development mode
make serve      # build and serve
make test       # run the test suite
make lint       # run golangci-lint
{{ if .docker }}make docker    # build Docker image
{{ end }}{{ if .goreleaser }}make release   # create a GitHub release via goreleaser
{{ end }}make clean      # remove binary and data
```

## Build Information

```bash
make build
./{{ .name }} version
# Build: 9a3b2c1
# Version: 0.1.0-dev
# Build Date: 2026-07-05T15:53:41Z
```

## Project Structure

```
.
├── main.go              # Application entrypoint
├── cmd/                 # Viper config loader + version command
├── migrations/          # Database migrations
├── models/              # Data models and query helpers
├── router/              # HTTP routes and middleware
│   ├── router.go        # Route definitions
│   └── middleware.go    # Logger and recover middleware
└── pb_data/             # PocketBase data directory (gitignored)
```

## API Endpoints

- `GET /health` — Health check
- `GET /api/v1/hello/{name}` — Hello world example
- `GET /api/v1/protected/me` — Get current user (requires auth)
- `/_/` — PocketBase Admin UI
