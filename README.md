# templates

[Quokka](https://github.com/Depado/quokka) templates

## Go

### Usage

```sh
$ # Using SSH
$ qk git@github.com:Depado/templates.git -p go myproject
$ # Using HTTPS
$ qk https://github.com/Depado/templates.git -p go myproject
```

### Description

This template can be used to generate a base Go application.

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

