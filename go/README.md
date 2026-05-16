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

A sample configuration file (`conf.yml`) is provided. You can also use environment variables (prefixed and uppercased based on the project name) or command-line flags.

```bash
./{{ .name }} --help
```
{{ end }}
