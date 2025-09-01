module {{ .gitserver }}/{{ .owner }}/{{ .name }}

go 1.22.2

require (
	{{ if .gin -}}
	github.com/Depado/ginprom v1.8.1
	github.com/gin-contrib/cors v1.7.6
	github.com/gin-gonic/gin v1.10.1
	{{ end -}}
	github.com/rs/zerolog v1.34.0
	github.com/spf13/cobra v1.10.1
	github.com/spf13/viper v1.20.1
)
