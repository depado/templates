<h1 align="center">{{.name}}</h1>

<h2 align="center">

  ![Go Version](https://img.shields.io/badge/Go%20Version-latest-brightgreen.svg)
  [![Go Report Card](https://goreportcard.com/badge/{{.gitserver}}/{{.owner}}/{{.name}})](https://goreportcard.com/report/{{.gitserver}}/{{.owner}}/{{.name}})
  [![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/{{.owner}}/{{.name}}/blob/main/LICENSE)

</h2>

<h2 align="center">{{.description}}</h2>

## Prerequisites

- **Go 1.26+** — [download](https://go.dev/dl/)
- **templ** — `go install github.com/a-h/templ/cmd/templ@latest`
- **Tailwind CSS v4 standalone CLI** — [download](https://github.com/tailwindlabs/tailwindcss/releases/latest), place `tailwindcss` in your `$PATH`
- **Task** — `go install github.com/go-task/task/v3/cmd/task@latest`
- **TemplUI CLI** — `go install github.com/templui/templui/cmd/templui@latest`

## Setup

```bash
make bootstrap   # Install TemplUI components, download HTMX, tidy deps
```

## Development

```bash
task dev           # templ + tailwind hot reload, runs go run . serve --dev
```

Or with Makefile:

```bash
make run           # build assets and run in dev mode
make dev           # shortcut alias for make run
```

## Building

```bash
make build
```

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
make bootstrap  # install TemplUI components + download HTMX + tidy deps
make build      # compile the binary
make run        # build assets and run in dev mode
make test       # run the test suite
make lint       # run golangci-lint
{{ if .docker }}make docker    # build Docker image
{{ end }}{{ if .goreleaser }}make release   # create a GitHub release via goreleaser
{{ end }}make clean      # remove binary and data
```

## Build Information

```bash
make build
./{{.name}} version
# Build: 9a3b2c1
# Version: 0.1.0-dev
# Build Date: 2026-07-05T15:53:41Z
```

## Getting Started

1. Start the server with `task dev` or `make run`
2. Open `http://localhost:8090`
3. Go to `http://localhost:8090/_/` and create a user in the `users` collection
4. Log in at `http://localhost:8090/login`

## Features

- **Login / Logout** — Cookie-based auth with PocketBase `users` collection
- **Settings** — Profile display (avatar, name, email) and password change
- **SPA Navigation** — HTMX boosted navigation with partial swaps
- **Light / Dark theme** — Toggle button in the header, persisted in localStorage
- **Responsive sidebar** — TemplUI sidebar with collapse, mobile sheet, and inset variant

## Project Structure

```
.
├── main.go                # Application entrypoint
├── assets/
│   ├── embed.go           # Embedded static files
│   └── css/
│       └── input.css      # Tailwind CSS v4 entry point
├── cmd/
│   ├── conf.go            # Viper-based configuration
│   └── version.go         # Version command
├── components/            # TemplUI components (installed via CLI)
├── models/
│   └── user.go            # User model
├── router/
│   ├── auth.go            # Auth middleware, login/logout handlers
│   ├── htmx.go            # go-htmx integration
│   ├── middleware.go      # Logger and recover middleware
│   ├── pages.go           # Dashboard and settings handlers
│   ├── render.go          # Render helper
│   └── router.go          # Route registration
├── views/
│   ├── layout/
│   │   ├── base.templ     # HTML shell
│   │   └── app.templ      # Authenticated layout with sidebar
│   └── pages/
│       ├── dashboard.templ
│       ├── login.templ
│       └── settings.templ
├── migrations/
├── Taskfile.yml           # Dev tasks
└── Makefile               # Build targets
```
