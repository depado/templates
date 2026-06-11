<h1 align="center">{{.name}}</h1>

<h2 align="center">

  ![Go Version](https://img.shields.io/badge/Go%20Version-latest-brightgreen.svg)
  [![Go Report Card](https://goreportcard.com/badge/{{.gitserver}}/{{.owner}}/{{.name}})](https://goreportcard.com/report/{{.gitserver}}/{{.owner}}/{{.name}})
  [![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/{{.owner}}/{{.name}}/blob/main/LICENSE)

</h2>

<h2 align="center">{{.description}}</h2>

## Prerequisites

- **Go 1.23+** — [download](https://go.dev/dl/)
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
task dev
```

## Building

```bash
make build
```

## Configuration

```bash
cp conf.example.yml conf.yml
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
│   └── conf.go            # Viper-based configuration
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
├── utils/
│   └── templui.go         # TemplUI utilities
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
