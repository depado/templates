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
- **shadcn-templ CLI** — `go install github.com/axadrn/shadcn-templ/v2/cmd/shadcn-templ@latest`

## Setup

Rendering this template with `qk` runs the bootstrap steps automatically (each
command is confirmed unless `--trusted` is passed). To set up an existing
project manually:

```bash
# Initialize shadcn-templ (writes components.json, generates globals.css, installs utils)
shadcn-templ init
# Install the components
shadcn-templ add sidebar button input label field card avatar dropdown-menu icon toast separator switch command textarea
# Download Datastar and the Geist fonts
mkdir -p assets/js assets/fonts/geist
curl -sL https://shadcn-templ.com/assets/fonts/geist/geist-variable.woff2 -o assets/fonts/geist/geist-variable.woff2
curl -sL https://shadcn-templ.com/assets/fonts/geist/geist-mono-variable.woff2 -o assets/fonts/geist/geist-mono-variable.woff2
curl -sL https://cdn.jsdelivr.net/gh/starfederation/datastar@v1.0.2/bundles/datastar.js -o assets/js/datastar.js
# Build assets
templ generate ./...
tailwindcss -i ./assets/css/globals.css -o ./assets/css/output.css --minify
go mod tidy
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
- Environment variables (uppercased, prefixed with project name, e.g. `{{.name | uc}}_DEV=true`)
- Configuration file (`conf.yml` or via `{{.name | uc}}_CONF` env var)

## Makefile Targets

```
make help       # list all targets
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

- **Login / Logout** — Cookie-based auth with PocketBase `users` collection; login and logout are Datastar backend actions (SSE responses)
- **Fragment navigation** — Sidebar links fetch `text/html` fragments that Datastar morphs into `#main-content`, with `history.pushState` URL sync and `popstate` back/forward support
- **Settings** — Profile display (avatar, name, email) and password change form (inline errors/success, field-preserving morph)
- **Light / Dark theme** — Toggle button in the header, persisted in localStorage
- **Responsive sidebar** — shadcn-templ sidebar with collapse, mobile sheet, and inset variant
- **Toasts** — shadcn-templ Toaster + a global fetch-error handler
- **CSRF protection** — same-origin check on state-changing endpoints

## Project Structure

```
.
├── main.go                # Application entrypoint
├── assets/
│   ├── embed.go           # Embedded static files
│   ├── css/
│   │   └── globals.css    # Tailwind CSS v4 entry point (shadcn-templ theme)
│   └── js/
│       └── datastar.js    # Datastar v1.0.2 bundle
├── cmd/
│   ├── conf.go            # Viper-based configuration
│   └── version.go         # Version command
├── components/            # shadcn-templ components (installed via CLI)
├── models/
│   └── user.go            # User model
├── router/
│   ├── auth.go            # Auth middleware, login/logout handlers
│   ├── datastar.go        # Datastar request detection + SSE helpers
│   ├── middleware.go      # Same-origin CSRF protection
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

## License

MIT
