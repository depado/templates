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

Copy the example configuration file and adjust as needed:

```bash
cp conf.example.yml conf.yml
```

Configuration can be set via:
- Configuration file (`conf.yml`)
- Environment variables (uppercased, prefixed with project name, e.g. `MYAPP_DEV=true`)

## Project Structure

```
.
├── main.go              # Application entrypoint
├── conf/                # Configuration
│   └── conf.go          # Viper-based config loader
├── migrations/          # Database migrations
├── models/              # Data models and query helpers
├── router/              # HTTP routes and middleware
│   ├── router.go        # Route definitions
│   └── middleware.go    # Custom middleware
└── pb_data/             # PocketBase data directory (gitignored)
```

## API Endpoints

- `GET /health` - Health check
- `GET /api/v1/hello/{name}` - Hello world example
- `GET /api/v1/protected/me` - Get current user (requires auth)
- `/_/` - PocketBase Admin UI
