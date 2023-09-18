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
- Logging with [zerolog](https://github.com/rs/zerolog)
- Inject version information (commit hash, latest tag, build date) with Makefile
- License selector
- Optional renovate configuration
- Optional docker support
    - Dockerfile with multi-step build with [distroless](https://github.com/GoogleContainerTools/distroless)
      container image for small and secure images
    - Uses `make` while building to inject version information
- Optional CI support
    - GitHub Actions: Build, test and run [golangci-lint](https://github.com/golangci/golangci-lint)
    - Drone: Build, test and run [golangci-lint](https://github.com/golangci/golangci-lint)
- Optional Gin integration:
    - Adds the proper configuration structs and flags
    - CORS customization and support
    - Prometheus instrumentation using [ginprom](https://github.com/Depado/ginprom)
    - Adds dependencies to `go.mod`

