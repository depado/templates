# templates

A collection of [Quokka](https://github.com/Depado/quokka) templates for Go projects — from a minimal CLI app to a full-stack PocketBase + HTMX + TemplUI web application.

## Go

### Usage

```sh
$ # Using SSH
$ qk git@github.com:Depado/templates.git -p go myproject
$ # Using HTTPS
$ qk https://github.com/Depado/templates.git -p go myproject
```

### Description

This template generates a base Go application.

- CLI
    - [Cobra](https://github.com/spf13/cobra) and [Viper](https://github.com/spf13/viper)
      for CLI and configuration
    - Configuration file, flags and environment variables support
    - `version` command to display injected variables at build time
- Configurable [zerolog](https://github.com/rs/zerolog) logger (format, level, caller)
- Inject version information (commit hash, latest tag, build date) with Makefile
- License selector
- Optional [renovate](https://github.com/renovatebot/renovate) configuration
- Optional example configuration file with default values for the app
- Optional docker support
    - Dockerfile with multi-step build with [distroless](https://github.com/GoogleContainerTools/distroless)
      container image for small and secure images
    - Uses `make` while building to inject version information
    - Adds Makefile rules (tag with `latest` and commit hash)
- Optional CI support
    - GitHub Actions: Build, test and run [golangci-lint](https://github.com/golangci/golangci-lint)
    - Drone: Build, test and run [golangci-lint](https://github.com/golangci/golangci-lint)
- Optional Gin integration:
    - Adds the proper configuration structs and flags
    - Graceful shutdown support
    - CORS customization and support
    - Prometheus instrumentation using [ginprom](https://github.com/Depado/ginprom)
    - Optional unified logger instead of gin's default one
- Optional [goreleaser](https://goreleaser.com/) config

## Go PocketBase

### Usage

```sh
$ # Using SSH
$ qk git@github.com:Depado/templates.git -p go-pocketbase myproject
$ # Using HTTPS
$ qk https://github.com/Depado/templates.git -p go-pocketbase myproject
```

### Description

This template generates a Go project using [PocketBase](https://pocketbase.io/) as the backend framework.

- [PocketBase](https://pocketbase.io/) integration with built-in admin UI, database, authentication, and file storage
- Cobra CLI with `version` command and PocketBase's built-in `serve` command
- Viper-based configuration (dev mode toggle, env var prefix)
- API routes: health check, hello endpoint, protected route with `RequireAuth`
- Logger and recover middlewares
- Migrations support via PocketBase `migratecmd` plugin (auto-migration in dev mode)
- Generic model helpers (FindById, FindAllByField, FindFirst)
- Inject version information (commit hash, latest tag, build date) with Makefile
- License selector
- Optional [renovate](https://github.com/renovatebot/renovate) configuration
- Optional example configuration file with default values for the app
- Optional docker support
    - Dockerfile with multi-step build using [distroless](https://github.com/GoogleContainerTools/distroless)
    - Exposes port 8090
    - Uses `make` while building to inject version information
- Optional CI via GitHub Actions (build, test, lint)
- Optional [goreleaser](https://goreleaser.com/) config with GitHub Actions release workflow

## Go PocketBase + TemplUI + HTMX

### Usage

```sh
$ # Using SSH
$ qk git@github.com:Depado/templates.git -p go-pocketbase-htmx myproject
$ # Using HTTPS
$ qk https://github.com/Depado/templates.git -p go-pocketbase-htmx myproject
```

### Description

This template generates a full-stack Go web application using PocketBase, [Templ](https://templ.guide/), [TemplUI](https://templui.io/), [HTMX](https://htmx.org/), and [Tailwind CSS v4](https://tailwindcss.com/).

- [PocketBase](https://pocketbase.io/) backend with built-in admin UI, database, and authentication
- [Templ](https://templ.guide/) for type-safe server-side HTML templating
- [TemplUI](https://templui.io/) component library: sidebar, cards, buttons, inputs, forms, toasts, avatars, dropdowns, icons
- [HTMX](https://htmx.org/) for SPA-like navigation with partial page swaps
- [Tailwind CSS v4](https://tailwindcss.com/) with dark mode support and oklch color tokens
- Cookie-based authentication with login/logout flows
- Responsive sidebar with collapse, mobile sheet, theme toggle, and user dropdown
- Dashboard and Settings pages with form validation
- Real-time hot reload: `task dev` runs templ proxy + tailwind watch in parallel via [Taskfile](https://taskfile.dev/)
- Asset embedding with Go's `embed` package
- User model with avatar (initials fallback) and password change support
- Cobra CLI with `version` command integrated into PocketBase's root command
- Viper-based configuration (dev mode toggle, env var prefix)
- Inject version information (commit hash, latest tag, build date) with Makefile
- License selector
- Optional [renovate](https://github.com/renovatebot/renovate) configuration
- Optional example configuration file with default values for the app
- Optional docker support with [distroless](https://github.com/GoogleContainerTools/distroless) on port 8090
- Optional CI via GitHub Actions (build, test, lint)
- Optional [goreleaser](https://goreleaser.com/) config with GitHub Actions release workflow

## Common Include Templates

The `common/` directory contains reusable sub-templates shared across all templates via Quokka's [includes](https://github.com/Depado/quokka#template-creation) directive. They can be used standalone or composed into the main templates above.

### Usage

```sh
$ # Using a common template standalone
$ qk git@github.com:Depado/templates.git -p common/license myproject
$ # Common templates are also composed into other templates via includes
```

### License (`common/license`)

Adds a `LICENSE` file with configurable license type.

- Supports MIT, Apache License 2.0, GPL, LGPL, and WTFPL
- Prompts for license type, copyright owner, and copyright year

### CI (`common/ci`)

Adds a GitHub Actions CI workflow.

- Triggers on push and PR to `main`
- Uses [depado/github-actions](https://github.com/Depado/github-actions) for Go build, test, and lint

### Renovate (`common/renovate`)

Adds a [Renovate](https://github.com/renovatebot/renovate) dependency update configuration.

- Extends `config:best-practices`
- Automerges minor and digest updates

### GoReleaser (`common/goreleaser`)

Adds a [GoReleaser](https://goreleaser.com/) configuration and GitHub Actions release workflow.

- Builds for linux, windows, and darwin (amd64 + arm64)
- Generates tar.gz/zip archives, checksums, and changelogs
- Release workflow triggers on git tags via [depado/github-actions](https://github.com/Depado/github-actions)
