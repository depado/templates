<h1 align="center">{{ .name }}</h1>

<h2 align="center">

  ![Go Version](https://img.shields.io/badge/Go%20Version-latest-brightgreen.svg)
  [![Go Report Card](https://goreportcard.com/badge/{{ .gitserver }}/{{ .owner }}/{{ .name }})](https://goreportcard.com/report/{{ .gitserver }}/{{ .owner }}/{{ .name }})
  {{ if .license }}[![License](https://img.shields.io/badge/license-{{ .license }}-blue.svg)](https://{{ .gitserver }}/{{ .owner }}/{{ .name }}/blob/master/LICENSE){{ end }}

</h2>

<h2 align="center">{{ .description }}</h2>
